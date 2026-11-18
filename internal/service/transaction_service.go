package service

import (
	"car-showroom-api/internal/entity"
	"car-showroom-api/internal/repository"
	"errors"
	"fmt"
	"math"
	"time"
)

type CreditSimulationInput struct {
	CarPrice       float64 `json:"car_price" binding:"required"`
	DownPayment    float64 `json:"down_payment" binding:"required"`
	TenorMonths    int     `json:"tenor_months" binding:"required"` // 12, 24, 36, 48, 60
	AnnualInterest float64 `json:"annual_interest"`                 // Opsional, default 6.0%
}

type CreditSimulationResult struct {
	CarPrice           float64 `json:"car_price"`
	DownPayment        float64 `json:"down_payment"`
	PrincipalLoan      float64 `json:"principal_loan"`
	TenorMonths        int     `json:"tenor_months"`
	AnnualInterestRate float64 `json:"annual_interest_rate"`
	MonthlyInstallment float64 `json:"monthly_installment"`
	TotalPayment       float64 `json:"total_payment"`
}

type CreateOrderInput struct {
	CarID         uint    `json:"car_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"` // cash, credit
	PromoCode     string  `json:"promo_code"`
	DownPayment   float64 `json:"down_payment"`
	TenorMonths   int     `json:"tenor_months"`
}

type TransactionService interface {
	// Wishlist
	AddWishlist(userID, carID uint) error
	RemoveWishlist(userID, carID uint) error
	GetWishlists(userID uint) ([]entity.Wishlist, error)

	// Test Drive
	BookTestDrive(userID, carID uint, date, timeSlot, notes string) (*entity.TestDrive, error)
	GetTestDrives(userID uint, role string) ([]entity.TestDrive, error)
	UpdateTestDriveStatus(id uint, status string) error

	// Consultation
	CreateConsultation(userID uint, carID *uint, subject, message string) (*entity.Consultation, error)
	GetConsultations(userID uint, role string) ([]entity.Consultation, error)
	ReplyConsultation(id uint, reply string) error

	// Promotion & Credit Simulation
	CreatePromotion(code, title, desc, validUntil string, discount float64) (*entity.Promotion, error)
	GetActivePromotions() ([]entity.Promotion, error)
	SimulateCredit(input CreditSimulationInput) (*CreditSimulationResult, error)

	// Order / Transaction
	CreateTransaction(userID uint, input CreateOrderInput) (*entity.Transaction, error)
	GetTransactions(userID uint, role string) ([]entity.Transaction, error)
	UpdateTransactionStatus(txID, salesID uint, status string) error
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo}
}

func (s *transactionService) AddWishlist(userID, carID uint) error {
	return s.repo.AddWishlist(&entity.Wishlist{UserID: userID, CarID: carID})
}

func (s *transactionService) RemoveWishlist(userID, carID uint) error {
	return s.repo.RemoveWishlist(userID, carID)
}

func (s *transactionService) GetWishlists(userID uint) ([]entity.Wishlist, error) {
	return s.repo.GetUserWishlists(userID)
}

func (s *transactionService) BookTestDrive(userID, carID uint, date, timeSlot, notes string) (*entity.TestDrive, error) {
	td := &entity.TestDrive{
		UserID:        userID,
		CarID:         carID,
		ScheduledDate: date,
		ScheduledTime: timeSlot,
		Status:        "pending",
		Notes:         notes,
	}
	err := s.repo.CreateTestDrive(td)
	return td, err
}

func (s *transactionService) GetTestDrives(userID uint, role string) ([]entity.TestDrive, error) {
	return s.repo.GetTestDrives(userID, role)
}

func (s *transactionService) UpdateTestDriveStatus(id uint, status string) error {
	return s.repo.UpdateTestDriveStatus(id, status)
}

func (s *transactionService) CreateConsultation(userID uint, carID *uint, subject, message string) (*entity.Consultation, error) {
	c := &entity.Consultation{
		UserID:  userID,
		CarID:   carID,
		Subject: subject,
		Message: message,
		Status:  "open",
	}
	err := s.repo.CreateConsultation(c)
	return c, err
}

func (s *transactionService) GetConsultations(userID uint, role string) ([]entity.Consultation, error) {
	return s.repo.GetConsultations(userID, role)
}

func (s *transactionService) ReplyConsultation(id uint, reply string) error {
	return s.repo.ReplyConsultation(id, reply)
}

func (s *transactionService) CreatePromotion(code, title, desc, validUntil string, discount float64) (*entity.Promotion, error) {
	p := &entity.Promotion{
		Code:           code,
		Title:          title,
		Description:    desc,
		DiscountAmount: discount,
		ValidUntil:     validUntil,
		IsActive:       true,
	}
	err := s.repo.CreatePromotion(p)
	return p, err
}

func (s *transactionService) GetActivePromotions() ([]entity.Promotion, error) {
	return s.repo.GetActivePromotions()
}

func (s *transactionService) SimulateCredit(input CreditSimulationInput) (*CreditSimulationResult, error) {
	if input.DownPayment >= input.CarPrice {
		return nil, errors.New("uang muka (DP) harus lebih kecil dari harga mobil")
	}
	if input.TenorMonths <= 0 {
		return nil, errors.New("tenor cicilan tidak valid")
	}

	rate := input.AnnualInterest
	if rate <= 0 {
		rate = 6.0 // Default bunga flat 6% per tahun
	}

	principal := input.CarPrice - input.DownPayment
	years := float64(input.TenorMonths) / 12.0
	totalInterest := principal * (rate / 100.0) * years
	totalLoan := principal + totalInterest
	monthly := math.Round(totalLoan / float64(input.TenorMonths))

	return &CreditSimulationResult{
		CarPrice:           input.CarPrice,
		DownPayment:        input.DownPayment,
		PrincipalLoan:      principal,
		TenorMonths:        input.TenorMonths,
		AnnualInterestRate: rate,
		MonthlyInstallment: monthly,
		TotalPayment:       input.DownPayment + totalLoan,
	}, nil
}

func (s *transactionService) CreateTransaction(userID uint, input CreateOrderInput) (*entity.Transaction, error) {
	car, err := s.repo.GetCarPriceByID(input.CarID)
	if err != nil {
		return nil, errors.New("mobil tidak ditemukan")
	}
	if car.Status != "available" {
		return nil, errors.New("maaf, mobil ini sedang dibooking atau sudah terjual")
	}

	var discount float64
	var promoID *uint
	if input.PromoCode != "" {
		promo, err := s.repo.FindPromotionByCode(input.PromoCode)
		if err != nil {
			return nil, errors.New("kode promo tidak valid atau sudah tidak aktif")
		}
		discount = promo.DiscountAmount
		promoID = &promo.ID
	}

	finalPrice := math.Max(0, car.Price-discount)
	var monthlyInstall float64

	if input.PaymentMethod == "credit" {
		sim, err := s.SimulateCredit(CreditSimulationInput{
			CarPrice:    finalPrice,
			DownPayment: input.DownPayment,
			TenorMonths: input.TenorMonths,
		})
		if err != nil {
			return nil, err
		}
		monthlyInstall = sim.MonthlyInstallment
	}

	invoice := fmt.Sprintf("INV-AD-%d-%d", time.Now().Year(), time.Now().Unix()%100000)

	tx := &entity.Transaction{
		InvoiceNumber:  invoice,
		UserID:         userID,
		CarID:          car.ID,
		PromotionID:    promoID,
		PaymentMethod:  input.PaymentMethod,
		CarPrice:       car.Price,
		DiscountAmount: discount,
		FinalPrice:     finalPrice,
		DownPayment:    input.DownPayment,
		TenorMonths:    input.TenorMonths,
		MonthlyInstall: monthlyInstall,
		Status:         "pending",
	}

	if err := s.repo.CreateTransaction(tx); err != nil {
		return nil, err
	}

	// Otomatis kunci status mobil menjadi 'booked'
	_ = s.repo.UpdateCarStatus(car.ID, "booked")

	return tx, nil
}

func (s *transactionService) GetTransactions(userID uint, role string) ([]entity.Transaction, error) {
	return s.repo.GetTransactions(userID, role)
}

func (s *transactionService) UpdateTransactionStatus(txID, salesID uint, status string) error {
	tx, err := s.repo.GetTransactionByID(txID)
	if err != nil {
		return errors.New("transaksi tidak ditemukan")
	}

	if err := s.repo.UpdateTransactionStatus(txID, salesID, status); err != nil {
		return err
	}

	// Sinkronisasi status stok unit mobil
	switch status {
	case "completed":
		_ = s.repo.UpdateCarStatus(tx.CarID, "sold")
	case "cancelled":
		_ = s.repo.UpdateCarStatus(tx.CarID, "available")
	}

	return nil
}
