package handler

import (
	"car-showroom-api/internal/repository"
	"car-showroom-api/internal/service"
	"car-showroom-api/pkg/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	service service.CarService
}

func NewCarHandler(service service.CarService) *CarHandler {
	return &CarHandler{service}
}

func (h *CarHandler) CreateBrand(c *gin.Context) {
	var input struct {
		Name    string `json:"name" binding:"required"`
		LogoURL string `json:"logo_url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input merek tidak valid", err.Error())
		return
	}

	brand, err := h.service.CreateBrand(input.Name, input.LogoURL)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Merek mobil berhasil ditambahkan", brand)
}

func (h *CarHandler) GetAllBrands(c *gin.Context) {
	brands, err := h.service.GetAllBrands()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data merek", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Daftar merek berhasil diambil", brands)
}

func (h *CarHandler) CreateCategory(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input kategori tidak valid", err.Error())
		return
	}

	category, err := h.service.CreateCategory(input.Name, input.Description)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Kategori mobil berhasil ditambahkan", category)
}

func (h *CarHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data kategori", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Daftar kategori berhasil diambil", categories)
}

func (h *CarHandler) CreateCar(c *gin.Context) {
	var input service.CreateCarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Data mobil tidak valid", err.Error())
		return
	}

	car, err := h.service.CreateCar(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal menambahkan mobil", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Data mobil berhasil ditambahkan", car)
}

func (h *CarHandler) GetAllCars(c *gin.Context) {
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)
	minYear, _ := strconv.Atoi(c.Query("min_year"))
	maxYear, _ := strconv.Atoi(c.Query("max_year"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := repository.CarFilter{
		Search:       c.Query("search"),
		BrandID:      c.Query("brand_id"),
		CategoryID:   c.Query("category_id"),
		Transmission: c.Query("transmission"),
		FuelType:     c.Query("fuel_type"),
		Condition:    c.Query("condition"),
		Status:       c.Query("status"),
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
		MinYear:      minYear,
		MaxYear:      maxYear,
		Page:         page,
		Limit:        limit,
	}

	cars, total, err := h.service.GetAllCars(filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil katalog mobil", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Katalog mobil berhasil diambil", gin.H{
		"items": cars,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *CarHandler) GetCarByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID mobil tidak valid", nil)
		return
	}

	car, err := h.service.GetCarByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Mobil tidak ditemukan", nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Detail mobil berhasil diambil", car)
}

func (h *CarHandler) UpdateCar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID mobil tidak valid", nil)
		return
	}

	var input service.CreateCarInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Data update mobil tidak valid", err.Error())
		return
	}

	car, err := h.service.UpdateCar(uint(id), input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data mobil berhasil diperbarui", car)
}

func (h *CarHandler) DeleteCar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID mobil tidak valid", nil)
		return
	}

	if err := h.service.DeleteCar(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus mobil", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Data mobil berhasil dihapus", nil)
}

func (h *CarHandler) UploadCarImages(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID mobil tidak valid", nil)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format upload harus multipart/form-data", err.Error())
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Pilih minimal satu file gambar pada field 'images'", nil)
		return
	}

	uploadDir := "uploads/cars"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat folder upload", err.Error())
		return
	}

	var uploadedURLs []string
	for i, file := range files {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("car_%d_%d_%d%s", id, time.Now().Unix(), i, ext)
		dst := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan file gambar", err.Error())
			return
		}
		uploadedURLs = append(uploadedURLs, "/uploads/cars/"+filename)
	}

	savedImages, err := h.service.SaveCarImages(uint(id), uploadedURLs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Foto mobil berhasil diunggah", savedImages)
}
