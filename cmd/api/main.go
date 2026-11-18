package main

import (
	"car-showroom-api/internal/config"
	"car-showroom-api/internal/handler"
	"car-showroom-api/internal/middleware"
	"car-showroom-api/internal/repository"
	"car-showroom-api/internal/service"
	"car-showroom-api/pkg/utils"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: file .env tidak ditemukan, menggunakan environment sistem")
	}

	db := config.ConnectDB()

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	carRepo := repository.NewCarRepository(db)
	carService := service.NewCarService(carRepo)
	carHandler := handler.NewCarHandler(carService)

	recRepo := repository.NewRecommendationRepository(db)
	recService := service.NewRecommendationService(recRepo)
	recHandler := handler.NewRecommendationHandler(recService)

	txRepo := repository.NewTransactionRepository(db)
	txService := service.NewTransactionService(txRepo)
	txHandler := handler.NewTransactionHandler(txService)

	r := gin.Default()

	// CORS Middleware agar frontend Next.js (localhost:3000) bisa mengakses API ini
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Menyajikan gambar mobil yang diunggah
	r.Static("/uploads", "./uploads")

	// Health Check Endpoint
	r.GET("/api/v1/health", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, "AutoDrive Showroom API berjalan normal", nil)
	})

	// Routing API v1
	v1 := r.Group("/api/v1")
	{
		// Auth Routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
		}

		// Public Routes (Pengunjung tanpa login)
		v1.GET("/brands", carHandler.GetAllBrands)
		v1.GET("/categories", carHandler.GetAllCategories)
		v1.GET("/cars", carHandler.GetAllCars)
		v1.GET("/cars/:id", carHandler.GetCarByID)
		v1.GET("/promotions", txHandler.GetActivePromotions)
		v1.POST("/credit-simulation", txHandler.SimulateCredit)

		// Protected Buyer Routes (Semua user terautentikasi)
		buyer := v1.Group("/buyer")
		buyer.Use(middleware.AuthMiddleware())
		{
			// SAW & Compare
			buyer.POST("/recommendations/calculate", recHandler.CalculateSAW)
			buyer.GET("/recommendations/history", recHandler.GetHistory)
			buyer.POST("/cars/compare", recHandler.CompareCars)

			// Wishlist
			buyer.POST("/wishlists", txHandler.AddWishlist)
			buyer.GET("/wishlists", txHandler.GetWishlists)
			buyer.DELETE("/wishlists/:car_id", txHandler.RemoveWishlist)

			// Test Drive & Konsultasi
			buyer.POST("/test-drives", txHandler.BookTestDrive)
			buyer.GET("/test-drives", txHandler.GetTestDrives)
			buyer.POST("/consultations", txHandler.CreateConsultation)
			buyer.GET("/consultations", txHandler.GetConsultations)

			// Pemesanan Unit
			buyer.POST("/transactions", txHandler.CreateTransaction)
			buyer.GET("/transactions", txHandler.GetTransactions)
		}

		// Protected Sales & Admin Routes (Operasional Showroom)
		sales := v1.Group("/sales")
		sales.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("sales", "admin"))
		{
			sales.GET("/test-drives", txHandler.GetTestDrives)
			sales.PATCH("/test-drives/:id/status", txHandler.UpdateTestDriveStatus)
			sales.GET("/consultations", txHandler.GetConsultations)
			sales.PATCH("/consultations/:id/reply", txHandler.ReplyConsultation)
			sales.GET("/transactions", txHandler.GetTransactions)
			sales.PATCH("/transactions/:id/status", txHandler.UpdateTransactionStatus)
		}

		// Protected Admin Routes (Manajemen Master & Stok)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			admin.GET("/dashboard-check", func(c *gin.Context) {
				utils.SuccessResponse(c, http.StatusOK, "Selamat datang di Dashboard Admin AutoDrive", nil)
			})
			admin.POST("/brands", carHandler.CreateBrand)
			admin.POST("/categories", carHandler.CreateCategory)
			admin.POST("/cars", carHandler.CreateCar)
			admin.PUT("/cars/:id", carHandler.UpdateCar)
			admin.DELETE("/cars/:id", carHandler.DeleteCar)
			admin.POST("/cars/:id/images", carHandler.UploadCarImages)
			admin.POST("/promotions", txHandler.CreatePromotion)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server AutoDrive Showroom berjalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}
