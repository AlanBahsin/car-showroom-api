package config

import (
	"car-showroom-api/internal/entity"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database Supabase: ", err)
	}

	log.Println("Berhasil terhubung ke Supabase PostgreSQL!")

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Brand{},
		&entity.Category{},
		&entity.Car{},
		&entity.Specification{},
		&entity.CarImage{},
		&entity.Recommendation{},
		&entity.RecommendationDetail{},
		&entity.Comparison{},
		&entity.Wishlist{},
		&entity.TestDrive{},
		&entity.Consultation{},
		&entity.Promotion{},
		&entity.Transaction{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi tabel: ", err)
	}

	return db
}
