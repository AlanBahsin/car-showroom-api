package repository

import (
	"car-showroom-api/internal/entity"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	// Wishlist
	AddWishlist(w *entity.Wishlist) error
	RemoveWishlist(userID, carID uint) error
	GetUserWishlists(userID uint) ([]entity.Wishlist, error)

	// Test Drive
	CreateTestDrive(td *entity.TestDrive) error
	GetTestDrives(userID uint, role string) ([]entity.TestDrive, error)
	UpdateTestDriveStatus(id uint, status string) error

	// Consultation
	CreateConsultation(c *entity.Consultation) error
	GetConsultations(userID uint, role string) ([]entity.Consultation, error)
	ReplyConsultation(id uint, reply string) error

	// Promotion
	CreatePromotion(p *entity.Promotion) error
	GetActivePromotions() ([]entity.Promotion, error)
	FindPromotionByCode(code string) (*entity.Promotion, error)

	// Transaction & Car Status Sync
	CreateTransaction(tx *entity.Transaction) error
	GetTransactions(userID uint, role string) ([]entity.Transaction, error)
	GetTransactionByID(id uint) (*entity.Transaction, error)
	UpdateTransactionStatus(id uint, salesID uint, status string) error
	UpdateCarStatus(carID uint, status string) error
	GetCarPriceByID(carID uint) (*entity.Car, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) AddWishlist(w *entity.Wishlist) error {
	var existing entity.Wishlist
	if err := r.db.Where("user_id = ? AND car_id = ?", w.UserID, w.CarID).First(&existing).Error; err == nil {
		return nil // Sudah ada di wishlist
	}
	return r.db.Create(w).Error
}

func (r *transactionRepository) RemoveWishlist(userID, carID uint) error {
	return r.db.Where("user_id = ? AND car_id = ?", userID, carID).Delete(&entity.Wishlist{}).Error
}

func (r *transactionRepository) GetUserWishlists(userID uint) ([]entity.Wishlist, error) {
	var list []entity.Wishlist
	err := r.db.Preload("Car").Preload("Car.Brand").Preload("Car.Images").
		Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *transactionRepository) CreateTestDrive(td *entity.TestDrive) error {
	return r.db.Create(td).Error
}

func (r *transactionRepository) GetTestDrives(userID uint, role string) ([]entity.TestDrive, error) {
	var list []entity.TestDrive
	query := r.db.Preload("User").Preload("Car").Preload("Car.Brand")
	if role == "buyer" {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *transactionRepository) UpdateTestDriveStatus(id uint, status string) error {
	return r.db.Model(&entity.TestDrive{}).Where("id = ?", id).Update("status", status).Error
}

func (r *transactionRepository) CreateConsultation(c *entity.Consultation) error {
	return r.db.Create(c).Error
}

func (r *transactionRepository) GetConsultations(userID uint, role string) ([]entity.Consultation, error) {
	var list []entity.Consultation
	query := r.db.Preload("User").Preload("Car")
	if role == "buyer" {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *transactionRepository) ReplyConsultation(id uint, reply string) error {
	return r.db.Model(&entity.Consultation{}).Where("id = ?", id).
		Updates(map[string]interface{}{"reply": reply, "status": "answered"}).Error
}

func (r *transactionRepository) CreatePromotion(p *entity.Promotion) error {
	return r.db.Create(p).Error
}

func (r *transactionRepository) GetActivePromotions() ([]entity.Promotion, error) {
	var list []entity.Promotion
	err := r.db.Where("is_active = ?", true).Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *transactionRepository) FindPromotionByCode(code string) (*entity.Promotion, error) {
	var p entity.Promotion
	err := r.db.Where("code = ? AND is_active = ?", code, true).First(&p).Error
	return &p, err
}

func (r *transactionRepository) CreateTransaction(tx *entity.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) GetTransactions(userID uint, role string) ([]entity.Transaction, error) {
	var list []entity.Transaction
	query := r.db.Preload("User").Preload("Car").Preload("Car.Brand")
	if role == "buyer" {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *transactionRepository) GetTransactionByID(id uint) (*entity.Transaction, error) {
	var tx entity.Transaction
	err := r.db.Preload("Car").First(&tx, id).Error
	return &tx, err
}

func (r *transactionRepository) UpdateTransactionStatus(id uint, salesID uint, status string) error {
	return r.db.Model(&entity.Transaction{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "sales_id": salesID}).Error
}

func (r *transactionRepository) UpdateCarStatus(carID uint, status string) error {
	return r.db.Model(&entity.Car{}).Where("id = ?", carID).Update("status", status).Error
}

func (r *transactionRepository) GetCarPriceByID(carID uint) (*entity.Car, error) {
	var car entity.Car
	err := r.db.First(&car, carID).Error
	return &car, err
}
