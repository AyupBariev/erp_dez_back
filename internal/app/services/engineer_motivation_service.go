// services/engineer_motivation_service.go
package services

import (
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"time"
)

type EngineerMotivationService struct {
	repo *repositories.EngineerMotivationRepository
}

func NewEngineerMotivationService(repo *repositories.EngineerMotivationRepository) *EngineerMotivationService {
	return &EngineerMotivationService{repo: repo}
}

func (s *EngineerMotivationService) GetMonthlyMotivation(monthStr string) ([]models.EngineerMonthlyMotivation, error) {
	var month time.Time
	var err error
	if monthStr == "" {
		now := time.Now()
		month = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		month, err = time.Parse("2006-01", monthStr)
		if err != nil {
			return nil, err
		}
	}

	return s.repo.GetByMonth(month)
}
