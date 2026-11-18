package repository

import (
	"car-showroom-api/internal/entity"

	"gorm.io/gorm"
)

type RecommendationRepository interface {
	GetAvailableCars() ([]entity.Car, error)
	SaveRecommendation(rec *entity.Recommendation) error
	GetUserRecommendations(userID uint) ([]entity.Recommendation, error)
	GetCarsByIDs(ids []uint) ([]entity.Car, error)
	SaveComparison(comp *entity.Comparison) error
}

type recommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &recommendationRepository{db}
}

func (r *recommendationRepository) GetAvailableCars() ([]entity.Car, error) {
	var cars []entity.Car
	err := r.db.Preload("Brand").
		Preload("Category").
		Preload("Specification").
		Preload("Images").
		Where("status = ?", "available").
		Find(&cars).Error
	return cars, err
}

func (r *recommendationRepository) SaveRecommendation(rec *entity.Recommendation) error {
	return r.db.Create(rec).Error
}

func (r *recommendationRepository) GetUserRecommendations(userID uint) ([]entity.Recommendation, error) {
	var recs []entity.Recommendation
	err := r.db.Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Order("rank asc")
	}).
		Preload("Details.Car").
		Preload("Details.Car.Brand").
		Preload("Details.Car.Category").
		Preload("Details.Car.Specification").
		Preload("Details.Car.Images").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&recs).Error
	return recs, err
}

func (r *recommendationRepository) GetCarsByIDs(ids []uint) ([]entity.Car, error) {
	var cars []entity.Car
	err := r.db.Preload("Brand").
		Preload("Category").
		Preload("Specification").
		Preload("Images").
		Where("id IN ?", ids).
		Find(&cars).Error
	return cars, err
}

func (r *recommendationRepository) SaveComparison(comp *entity.Comparison) error {
	return r.db.Create(comp).Error
}
