package handler

import (
	"car-showroom-api/internal/service"
	"car-showroom-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	service service.RecommendationService
}

func NewRecommendationHandler(service service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service}
}

func (h *RecommendationHandler) CalculateSAW(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input service.SAWInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input preferensi tidak valid", err.Error())
		return
	}

	result, err := h.service.CalculateSAW(userID.(uint), input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perhitungan rekomendasi SAW berhasil", result)
}

func (h *RecommendationHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")

	history, err := h.service.GetHistory(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil riwayat rekomendasi", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Riwayat rekomendasi berhasil diambil", history)
}

func (h *RecommendationHandler) CompareCars(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		CarIDs []uint `json:"car_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID mobil tidak valid", err.Error())
		return
	}

	cars, err := h.service.CompareCars(userID.(uint), input.CarIDs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Perbandingan spesifikasi mobil berhasil diambil", cars)
}
