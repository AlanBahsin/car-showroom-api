package handler

import (
	"car-showroom-api/internal/service"
	"car-showroom-api/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) AddWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input struct {
		CarID uint `json:"car_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID mobil wajib diisi", err.Error())
		return
	}
	if err := h.service.AddWishlist(userID.(uint), input.CarID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambah wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Mobil berhasil disimpan ke wishlist", nil)
}

func (h *TransactionHandler) RemoveWishlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	carID, _ := strconv.Atoi(c.Param("car_id"))
	if err := h.service.RemoveWishlist(userID.(uint), uint(carID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Mobil dihapus dari wishlist", nil)
}

func (h *TransactionHandler) GetWishlists(c *gin.Context) {
	userID, _ := c.Get("user_id")
	list, err := h.service.GetWishlists(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Daftar wishlist berhasil diambil", list)
}

func (h *TransactionHandler) BookTestDrive(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input struct {
		CarID         uint   `json:"car_id" binding:"required"`
		ScheduledDate string `json:"scheduled_date" binding:"required"`
		ScheduledTime string `json:"scheduled_time" binding:"required"`
		Notes         string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Data booking test drive tidak lengkap", err.Error())
		return
	}
	td, err := h.service.BookTestDrive(userID.(uint), input.CarID, input.ScheduledDate, input.ScheduledTime, input.Notes)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat jadwal test drive", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Jadwal test drive berhasil diajukan", td)
}

func (h *TransactionHandler) GetTestDrives(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	list, err := h.service.GetTestDrives(userID.(uint), role.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data test drive", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Daftar test drive berhasil diambil", list)
}

func (h *TransactionHandler) UpdateTestDriveStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Status wajib diisi", err.Error())
		return
	}
	if err := h.service.UpdateTestDriveStatus(uint(id), input.Status); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal update status test drive", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Status test drive berhasil diperbarui", nil)
}

func (h *TransactionHandler) CreateConsultation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input struct {
		CarID   *uint  `json:"car_id"`
		Subject string `json:"subject" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Subjek dan pesan wajib diisi", err.Error())
		return
	}
	cons, err := h.service.CreateConsultation(userID.(uint), input.CarID, input.Subject, input.Message)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengirim konsultasi", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Pertanyaan konsultasi berhasil dikirim", cons)
}

func (h *TransactionHandler) GetConsultations(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	list, err := h.service.GetConsultations(userID.(uint), role.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil konsultasi", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Daftar konsultasi berhasil diambil", list)
}

func (h *TransactionHandler) ReplyConsultation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Reply string `json:"reply" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Balasan wajib diisi", err.Error())
		return
	}
	if err := h.service.ReplyConsultation(uint(id), input.Reply); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membalas konsultasi", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Balasan konsultasi berhasil dikirim", nil)
}

func (h *TransactionHandler) CreatePromotion(c *gin.Context) {
	var input struct {
		Code           string  `json:"code" binding:"required"`
		Title          string  `json:"title" binding:"required"`
		Description    string  `json:"description"`
		DiscountAmount float64 `json:"discount_amount" binding:"required"`
		ValidUntil     string  `json:"valid_until" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Data promo tidak lengkap", err.Error())
		return
	}
	p, err := h.service.CreatePromotion(input.Code, input.Title, input.Description, input.ValidUntil, input.DiscountAmount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Gagal membuat promo", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Promo berhasil ditambahkan", p)
}

func (h *TransactionHandler) GetActivePromotions(c *gin.Context) {
	list, err := h.service.GetActivePromotions()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil promo", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Daftar promo aktif berhasil diambil", list)
}

func (h *TransactionHandler) SimulateCredit(c *gin.Context) {
	var input service.CreditSimulationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input simulasi kredit tidak valid", err.Error())
		return
	}
	res, err := h.service.SimulateCredit(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Simulasi kredit berhasil dihitung", res)
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input service.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Data pemesanan tidak valid", err.Error())
		return
	}
	tx, err := h.service.CreateTransaction(userID.(uint), input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Pemesanan unit mobil berhasil dibuat", tx)
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	list, err := h.service.GetTransactions(userID.(uint), role.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data transaksi", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Daftar transaksi berhasil diambil", list)
}

func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	salesID, _ := c.Get("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Status string `json:"status" binding:"required"` // confirmed, completed, cancelled
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Status transaksi wajib diisi", err.Error())
		return
	}
	if err := h.service.UpdateTransactionStatus(uint(id), salesID.(uint), input.Status); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Status transaksi berhasil diperbarui", nil)
}
