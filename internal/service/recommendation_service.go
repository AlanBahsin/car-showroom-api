package service

import (
	"car-showroom-api/internal/entity"
	"car-showroom-api/internal/repository"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

type SAWInput struct {
	BudgetRange         string  `json:"budget_range"`
	MinBudget           float64 `json:"min_budget"`
	MaxBudget           float64 `json:"max_budget" binding:"required"`
	PreferredCategoryID uint    `json:"preferred_category_id"`
	PreferredTrans      string  `json:"preferred_transmission"`
	PreferredFuel       string  `json:"preferred_fuel"`
	Priority            string  `json:"priority"`
}

type RecommendationService interface {
	CalculateSAW(userID uint, input SAWInput) (*entity.Recommendation, error)
	GetHistory(userID uint) ([]entity.Recommendation, error)
	CompareCars(userID uint, carIDs []uint) ([]entity.Car, error)
}

type recommendationService struct {
	repo repository.RecommendationRepository
}

func NewRecommendationService(repo repository.RecommendationRepository) RecommendationService {
	return &recommendationService{repo}
}

type rawCarScore struct {
	Car           entity.Car
	PriceScore    float64
	CategoryScore float64
	TransScore    float64
	FuelScore     float64
	YearScore     float64
}

func (s *recommendationService) CalculateSAW(userID uint, input SAWInput) (*entity.Recommendation, error) {
	cars, err := s.repo.GetAvailableCars()
	if err != nil {
		return nil, err
	}
	if len(cars) == 0 {
		return nil, errors.New("belum ada unit mobil yang tersedia untuk direkomendasikan")
	}

	wPrice := 0.35
	wCategory := 0.25
	wTrans := 0.15
	wFuel := 0.10
	wYear := 0.15

	switch strings.ToLower(input.Priority) {
	case "harga_terjangkau":
		wPrice, wYear = 0.40, 0.10
	case "kapasitas":
		wCategory, wFuel = 0.30, 0.05
	case "irit_bbm":
		wFuel, wTrans = 0.15, 0.10
	}

	var scores []rawCarScore
	maxPriceS, maxCatS, maxTransS, maxFuelS, maxYearS := 0.0, 0.0, 0.0, 0.0, 0.0

	for _, car := range cars {
		priceS := 5.0
		if car.Price > input.MaxBudget {
			overRatio := (car.Price - input.MaxBudget) / input.MaxBudget
			priceS = math.Max(1.0, 5.0-(overRatio*10.0))
		} else if input.MinBudget > 0 && car.Price < input.MinBudget {
			priceS = 4.5
		}

		catS := 3.0
		if input.PreferredCategoryID == 0 || car.CategoryID == input.PreferredCategoryID {
			catS = 5.0
		}
		if strings.ToLower(input.Priority) == "kapasitas" && car.Specification.SeatCapacity >= 7 {
			catS = math.Min(5.0, catS+0.5)
		}

		transS := 5.0
		if input.PreferredTrans != "" && !strings.EqualFold(input.PreferredTrans, "Bebas") {
			if !strings.EqualFold(car.Transmission, input.PreferredTrans) {
				transS = 1.0
			}
		}

		fuelS := 5.0
		if input.PreferredFuel != "" && !strings.EqualFold(input.PreferredFuel, "Bebas") {
			if !strings.EqualFold(car.FuelType, input.PreferredFuel) {
				fuelS = 1.0
			}
		}
		if strings.ToLower(input.Priority) == "irit_bbm" && car.Specification.FuelConsumption >= 14.0 {
			fuelS = math.Min(5.0, fuelS+0.5)
		}

		yearS := 1.0
		if car.Year >= 2023 {
			yearS = 5.0
		} else if car.Year >= 2020 {
			yearS = 4.0
		} else if car.Year >= 2017 {
			yearS = 3.0
		} else if car.Year >= 2014 {
			yearS = 2.0
		}
		if strings.ToLower(input.Priority) == "performa" && car.Specification.EngineCapacity >= 1500 {
			yearS = math.Min(5.0, yearS+0.5)
		}

		if priceS > maxPriceS {
			maxPriceS = priceS
		}
		if catS > maxCatS {
			maxCatS = catS
		}
		if transS > maxTransS {
			maxTransS = transS
		}
		if fuelS > maxFuelS {
			maxFuelS = fuelS
		}
		if yearS > maxYearS {
			maxYearS = yearS
		}

		scores = append(scores, rawCarScore{
			Car:           car,
			PriceScore:    priceS,
			CategoryScore: catS,
			TransScore:    transS,
			FuelScore:     fuelS,
			YearScore:     yearS,
		})
	}

	var details []entity.RecommendationDetail
	for _, sc := range scores {
		rPrice := sc.PriceScore / maxPriceS
		rCat := sc.CategoryScore / maxCatS
		rTrans := sc.TransScore / maxTransS
		rFuel := sc.FuelScore / maxFuelS
		rYear := sc.YearScore / maxYearS

		finalScore := (wPrice * rPrice) + (wCategory * rCat) + (wTrans * rTrans) + (wFuel * rFuel) + (wYear * rYear)
		matchPct := math.Round(finalScore*1000) / 10

		reason := fmt.Sprintf("Cocok untuk kebutuhan %s dengan transmisi %s, bahan bakar %s, dan kapasitas %d penumpang.",
			sc.Car.Category.Name, sc.Car.Transmission, sc.Car.FuelType, sc.Car.Specification.SeatCapacity)

		details = append(details, entity.RecommendationDetail{
			CarID:           sc.Car.ID,
			Car:             sc.Car,
			PriceScore:      math.Round(rPrice*100) / 100,
			CategoryScore:   math.Round(rCat*100) / 100,
			TransScore:      math.Round(rTrans*100) / 100,
			FuelScore:       math.Round(rFuel*100) / 100,
			YearScore:       math.Round(rYear*100) / 100,
			FinalScore:      math.Round(finalScore*10000) / 10000,
			MatchPercentage: matchPct,
			Reason:          reason,
		})
	}

	sort.Slice(details, func(i, j int) bool {
		return details[i].FinalScore > details[j].FinalScore
	})

	if len(details) > 5 {
		details = details[:5]
	}
	for i := range details {
		details[i].Rank = i + 1
	}

	rec := &entity.Recommendation{
		UserID:              userID,
		BudgetRange:         input.BudgetRange,
		MinBudget:           input.MinBudget,
		MaxBudget:           input.MaxBudget,
		PreferredCategoryID: input.PreferredCategoryID,
		PreferredTrans:      input.PreferredTrans,
		PreferredFuel:       input.PreferredFuel,
		Priority:            input.Priority,
		Details:             details,
	}

	if err := s.repo.SaveRecommendation(rec); err != nil {
		return nil, err
	}

	return rec, nil
}

func (s *recommendationService) GetHistory(userID uint) ([]entity.Recommendation, error) {
	return s.repo.GetUserRecommendations(userID)
}

func (s *recommendationService) CompareCars(userID uint, carIDs []uint) ([]entity.Car, error) {
	if len(carIDs) < 2 || len(carIDs) > 3 {
		return nil, errors.New("pilih minimal 2 dan maksimal 3 mobil untuk dibandingkan")
	}

	cars, err := s.repo.GetCarsByIDs(carIDs)
	if err != nil {
		return nil, err
	}

	var idStrings []string
	for _, id := range carIDs {
		idStrings = append(idStrings, fmt.Sprintf("%d", id))
	}

	_ = s.repo.SaveComparison(&entity.Comparison{
		UserID: userID,
		CarIDs: strings.Join(idStrings, ","),
	})

	return cars, nil
}
