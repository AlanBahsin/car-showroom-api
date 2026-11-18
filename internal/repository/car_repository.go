package repository

import (
	"car-showroom-api/internal/entity"

	"gorm.io/gorm"
)

type CarFilter struct {
	Search       string
	BrandID      string
	CategoryID   string
	Transmission string
	FuelType     string
	Condition    string
	MinPrice     float64
	MaxPrice     float64
	MinYear      int
	MaxYear      int
	Status       string
	Page         int
	Limit        int
}

type CarRepository interface {
	// Brand & Category
	CreateBrand(brand *entity.Brand) error
	GetAllBrands() ([]entity.Brand, error)
	CreateCategory(category *entity.Category) error
	GetAllCategories() ([]entity.Category, error)

	// Car & Specification
	CreateCar(car *entity.Car) error
	GetAllCars(filter CarFilter) ([]entity.Car, int64, error)
	GetCarByID(id uint) (*entity.Car, error)
	UpdateCar(car *entity.Car) error
	DeleteCar(id uint) error

	// Car Images
	AddCarImages(images []entity.CarImage) error
}

type carRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &carRepository{db}
}

func (r *carRepository) CreateBrand(brand *entity.Brand) error {
	return r.db.Create(brand).Error
}

func (r *carRepository) GetAllBrands() ([]entity.Brand, error) {
	var brands []entity.Brand
	err := r.db.Order("name asc").Find(&brands).Error
	return brands, err
}

func (r *carRepository) CreateCategory(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *carRepository) GetAllCategories() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Order("id asc").Find(&categories).Error
	return categories, err
}

func (r *carRepository) CreateCar(car *entity.Car) error {
	return r.db.Create(car).Error
}

func (r *carRepository) GetAllCars(filter CarFilter) ([]entity.Car, int64, error) {
	var cars []entity.Car
	var total int64

	query := r.db.Model(&entity.Car{}).
		Preload("Brand").
		Preload("Category").
		Preload("Specification").
		Preload("Images")

	if filter.Search != "" {
		query = query.Where("model_name ILIKE ?", "%"+filter.Search+"%")
	}
	if filter.BrandID != "" {
		query = query.Where("brand_id = ?", filter.BrandID)
	}
	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.Transmission != "" {
		query = query.Where("transmission ILIKE ?", filter.Transmission)
	}
	if filter.FuelType != "" {
		query = query.Where("fuel_type ILIKE ?", filter.FuelType)
	}
	if filter.Condition != "" {
		query = query.Where("condition = ?", filter.Condition)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}
	if filter.MinYear > 0 {
		query = query.Where("year >= ?", filter.MinYear)
	}
	if filter.MaxYear > 0 {
		query = query.Where("year <= ?", filter.MaxYear)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.Limit
	err := query.Order("created_at desc").Limit(filter.Limit).Offset(offset).Find(&cars).Error

	return cars, total, err
}

func (r *carRepository) GetCarByID(id uint) (*entity.Car, error) {
	var car entity.Car
	err := r.db.Preload("Brand").
		Preload("Category").
		Preload("Specification").
		Preload("Images").
		First(&car, id).Error
	return &car, err
}

func (r *carRepository) UpdateCar(car *entity.Car) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(car).Error
}

func (r *carRepository) DeleteCar(id uint) error {
	return r.db.Delete(&entity.Car{}, id).Error
}

func (r *carRepository) AddCarImages(images []entity.CarImage) error {
	return r.db.Create(&images).Error
}
