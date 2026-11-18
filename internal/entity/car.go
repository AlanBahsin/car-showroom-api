package entity

import (
	"time"

	"gorm.io/gorm"
)

type Car struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	BrandID       uint           `gorm:"not null" json:"brand_id"`
	Brand         Brand          `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	CategoryID    uint           `gorm:"not null" json:"category_id"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ModelName     string         `gorm:"type:varchar(150);not null" json:"model_name"` // Contoh: Avanza 1.5 G CVT
	Year          int            `gorm:"not null" json:"year"`
	Price         float64        `gorm:"not null" json:"price"`
	Condition     string         `gorm:"type:varchar(20);default:'used'" json:"condition"` // new, used
	Mileage       int            `gorm:"default:0" json:"mileage"`                         // Kilometer tempuh
	Transmission  string         `gorm:"type:varchar(30);not null" json:"transmission"`    // Automatic, Manual
	FuelType      string         `gorm:"type:varchar(30);not null" json:"fuel_type"`       // Bensin, Diesel, Hybrid, Listrik
	Color         string         `gorm:"type:varchar(50)" json:"color"`
	Status        string         `gorm:"type:varchar(20);default:'available'" json:"status"` // available, booked, sold
	IsFeatured    bool           `gorm:"default:false" json:"is_featured"`                   // Mobil unggulan di homepage
	Description   string         `gorm:"type:text" json:"description"`
	Specification Specification  `gorm:"foreignKey:CarID;constraint:OnDelete:CASCADE" json:"specification"`
	Images        []CarImage     `gorm:"foreignKey:CarID;constraint:OnDelete:CASCADE" json:"images"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Specification struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	CarID             uint    `gorm:"uniqueIndex;not null" json:"car_id"`
	EngineCapacity    int     `json:"engine_capacity"`  // cc, contoh: 1496
	Horsepower        int     `json:"horsepower"`       // hp, contoh: 106
	SeatCapacity      int     `json:"seat_capacity"`    // jumlah penumpang, contoh: 7
	FuelConsumption   float64 `json:"fuel_consumption"` // km/liter, contoh: 14.5
	TaxExpiryDate     string  `gorm:"type:varchar(50)" json:"tax_expiry_date"`
	ServiceBookStatus bool    `gorm:"default:true" json:"service_book_status"`
}

type CarImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CarID     uint      `gorm:"not null" json:"car_id"`
	ImageURL  string    `gorm:"type:varchar(255);not null" json:"image_url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}
