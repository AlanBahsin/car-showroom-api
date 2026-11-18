package service

import (
	"car-showroom-api/internal/entity"
	"car-showroom-api/internal/repository"
	"errors"
)

type CreateCarInput struct {
	BrandID           uint    `json:"brand_id" binding:"required"`
	CategoryID        uint    `json:"category_id" binding:"required"`
	ModelName         string  `json:"model_name" binding:"required"`
	Year              int     `json:"year" binding:"required"`
	Price             float64 `json:"price" binding:"required"`
	Condition         string  `json:"condition"`
	Mileage           int     `json:"mileage"`
	Transmission      string  `json:"transmission" binding:"required"`
	FuelType          string  `json:"fuel_type" binding:"required"`
	Color             string  `json:"color"`
	Status            string  `json:"status"`
	IsFeatured        bool    `json:"is_featured"`
	Description       string  `json:"description"`
	EngineCapacity    int     `json:"engine_capacity"`
	Horsepower        int     `json:"horsepower"`
	SeatCapacity      int     `json:"seat_capacity"`
	FuelConsumption   float64 `json:"fuel_consumption"`
	TaxExpiryDate     string  `json:"tax_expiry_date"`
	ServiceBookStatus bool    `json:"service_book_status"`
}

type CarService interface {
	CreateBrand(name, logoURL string) (*entity.Brand, error)
	GetAllBrands() ([]entity.Brand, error)
	CreateCategory(name, description string) (*entity.Category, error)
	GetAllCategories() ([]entity.Category, error)

	CreateCar(input CreateCarInput) (*entity.Car, error)
	GetAllCars(filter repository.CarFilter) ([]entity.Car, int64, error)
	GetCarByID(id uint) (*entity.Car, error)
	UpdateCar(id uint, input CreateCarInput) (*entity.Car, error)
	DeleteCar(id uint) error
	SaveCarImages(carID uint, urls []string) ([]entity.CarImage, error)
}

type carService struct {
	repo repository.CarRepository
}

func NewCarService(repo repository.CarRepository) CarService {
	return &carService{repo}
}

func (s *carService) CreateBrand(name, logoURL string) (*entity.Brand, error) {
	if name == "" {
		return nil, errors.New("nama merek wajib diisi")
	}
	brand := &entity.Brand{Name: name, LogoURL: logoURL}
	err := s.repo.CreateBrand(brand)
	return brand, err
}

func (s *carService) GetAllBrands() ([]entity.Brand, error) {
	return s.repo.GetAllBrands()
}

func (s *carService) CreateCategory(name, description string) (*entity.Category, error) {
	if name == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}
	category := &entity.Category{Name: name, Description: description}
	err := s.repo.CreateCategory(category)
	return category, err
}

func (s *carService) GetAllCategories() ([]entity.Category, error) {
	return s.repo.GetAllCategories()
}

func (s *carService) CreateCar(input CreateCarInput) (*entity.Car, error) {
	condition := "used"
	if input.Condition == "new" {
		condition = "new"
	}
	status := "available"
	if input.Status != "" {
		status = input.Status
	}

	car := &entity.Car{
		BrandID:      input.BrandID,
		CategoryID:   input.CategoryID,
		ModelName:    input.ModelName,
		Year:         input.Year,
		Price:        input.Price,
		Condition:    condition,
		Mileage:      input.Mileage,
		Transmission: input.Transmission,
		FuelType:     input.FuelType,
		Color:        input.Color,
		Status:       status,
		IsFeatured:   input.IsFeatured,
		Description:  input.Description,
		Specification: entity.Specification{
			EngineCapacity:    input.EngineCapacity,
			Horsepower:        input.Horsepower,
			SeatCapacity:      input.SeatCapacity,
			FuelConsumption:   input.FuelConsumption,
			TaxExpiryDate:     input.TaxExpiryDate,
			ServiceBookStatus: input.ServiceBookStatus,
		},
	}

	if err := s.repo.CreateCar(car); err != nil {
		return nil, err
	}
	return s.repo.GetCarByID(car.ID)
}

func (s *carService) GetAllCars(filter repository.CarFilter) ([]entity.Car, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	return s.repo.GetAllCars(filter)
}

func (s *carService) GetCarByID(id uint) (*entity.Car, error) {
	return s.repo.GetCarByID(id)
}

func (s *carService) UpdateCar(id uint, input CreateCarInput) (*entity.Car, error) {
	car, err := s.repo.GetCarByID(id)
	if err != nil {
		return nil, errors.New("mobil tidak ditemukan")
	}

	car.BrandID = input.BrandID
	car.CategoryID = input.CategoryID
	car.ModelName = input.ModelName
	car.Year = input.Year
	car.Price = input.Price
	if input.Condition != "" {
		car.Condition = input.Condition
	}
	car.Mileage = input.Mileage
	car.Transmission = input.Transmission
	car.FuelType = input.FuelType
	car.Color = input.Color
	if input.Status != "" {
		car.Status = input.Status
	}
	car.IsFeatured = input.IsFeatured
	car.Description = input.Description

	car.Specification.EngineCapacity = input.EngineCapacity
	car.Specification.Horsepower = input.Horsepower
	car.Specification.SeatCapacity = input.SeatCapacity
	car.Specification.FuelConsumption = input.FuelConsumption
	car.Specification.TaxExpiryDate = input.TaxExpiryDate
	car.Specification.ServiceBookStatus = input.ServiceBookStatus

	if err := s.repo.UpdateCar(car); err != nil {
		return nil, err
	}
	return s.repo.GetCarByID(id)
}

func (s *carService) DeleteCar(id uint) error {
	return s.repo.DeleteCar(id)
}

func (s *carService) SaveCarImages(carID uint, urls []string) ([]entity.CarImage, error) {
	car, err := s.repo.GetCarByID(carID)
	if err != nil {
		return nil, errors.New("mobil tidak ditemukan")
	}

	hasPrimary := len(car.Images) > 0
	var images []entity.CarImage
	for i, url := range urls {
		isPrimary := false
		if !hasPrimary && i == 0 {
			isPrimary = true
		}
		images = append(images, entity.CarImage{
			CarID:     carID,
			ImageURL:  url,
			IsPrimary: isPrimary,
		})
	}

	if err := s.repo.AddCarImages(images); err != nil {
		return nil, err
	}
	return images, nil
}
