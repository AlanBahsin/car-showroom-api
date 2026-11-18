package entity

import (
	"time"

	"gorm.io/gorm"
)

type Wishlist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CarID     uint      `gorm:"not null" json:"car_id"`
	Car       Car       `gorm:"foreignKey:CarID" json:"car,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TestDrive struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CarID         uint           `gorm:"not null" json:"car_id"`
	Car           Car            `gorm:"foreignKey:CarID" json:"car,omitempty"`
	ScheduledDate string         `gorm:"type:varchar(30);not null" json:"scheduled_date"`  // YYYY-MM-DD
	ScheduledTime string         `gorm:"type:varchar(20);not null" json:"scheduled_time"`  // Contoh: "10:00 WIB"
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, completed, cancelled
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Consultation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CarID     *uint          `json:"car_id"` // Opsional jika menanyakan unit tertentu
	Car       *Car           `gorm:"foreignKey:CarID" json:"car,omitempty"`
	Subject   string         `gorm:"type:varchar(150);not null" json:"subject"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Reply     string         `gorm:"type:text" json:"reply"`
	Status    string         `gorm:"type:varchar(20);default:'open'" json:"status"` // open, answered, closed
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Promotion struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Code           string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Title          string         `gorm:"type:varchar(150);not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	DiscountAmount float64        `gorm:"not null" json:"discount_amount"` // Nominal potongan harga (Rp)
	ValidUntil     string         `gorm:"type:varchar(30);not null" json:"valid_until"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	InvoiceNumber  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"invoice_number"`
	UserID         uint           `gorm:"not null" json:"user_id"`
	User           User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CarID          uint           `gorm:"not null" json:"car_id"`
	Car            Car            `gorm:"foreignKey:CarID" json:"car,omitempty"`
	SalesID        *uint          `json:"sales_id"` // Sales yang menangani transaksi
	PromotionID    *uint          `json:"promotion_id"`
	PaymentMethod  string         `gorm:"type:varchar(20);not null" json:"payment_method"` // cash, credit
	CarPrice       float64        `gorm:"not null" json:"car_price"`
	DiscountAmount float64        `gorm:"default:0" json:"discount_amount"`
	FinalPrice     float64        `gorm:"not null" json:"final_price"`
	DownPayment    float64        `gorm:"default:0" json:"down_payment"`                    // DP jika kredit
	TenorMonths    int            `gorm:"default:0" json:"tenor_months"`                    // 12, 24, 36, 48, 60 bulan
	MonthlyInstall float64        `gorm:"default:0" json:"monthly_installment"`             // Estimasi cicilan per bulan
	Status         string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, confirmed, completed, cancelled
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
