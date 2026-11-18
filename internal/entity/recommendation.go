package entity

import (
	"time"
)

type Recommendation struct {
	ID                  uint                   `gorm:"primaryKey" json:"id"`
	UserID              uint                   `gorm:"not null" json:"user_id"`
	BudgetRange         string                 `gorm:"type:varchar(50)" json:"budget_range"`
	MinBudget           float64                `json:"min_budget"`
	MaxBudget           float64                `json:"max_budget"`
	PreferredCategoryID uint                   `json:"preferred_category_id"`
	PreferredTrans      string                 `gorm:"type:varchar(30)" json:"preferred_transmission"`
	PreferredFuel       string                 `gorm:"type:varchar(30)" json:"preferred_fuel"`
	Priority            string                 `gorm:"type:varchar(50)" json:"priority"` // irit_bbm, kapasitas, kenyamanan, performa, harga_terjangkau
	Details             []RecommendationDetail `gorm:"foreignKey:RecommendationID;constraint:OnDelete:CASCADE" json:"details"`
	CreatedAt           time.Time              `json:"created_at"`
}

type RecommendationDetail struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	RecommendationID uint    `gorm:"not null" json:"recommendation_id"`
	CarID            uint    `gorm:"not null" json:"car_id"`
	Car              Car     `gorm:"foreignKey:CarID" json:"car"`
	PriceScore       float64 `json:"price_score"`
	CategoryScore    float64 `json:"category_score"`
	TransScore       float64 `json:"trans_score"`
	FuelScore        float64 `json:"fuel_score"`
	YearScore        float64 `json:"year_score"`
	FinalScore       float64 `json:"final_score"`      // Skala 0.00 - 1.00
	MatchPercentage  float64 `json:"match_percentage"` // Skala 0 - 100%
	Rank             int     `json:"rank"`
	Reason           string  `gorm:"type:text" json:"reason"`
}

type Comparison struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CarIDs    string    `gorm:"type:varchar(100);not null" json:"car_ids"` // Contoh: "1,2,3"
	CreatedAt time.Time `json:"created_at"`
}
