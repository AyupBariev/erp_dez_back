package services

import (
	"erp/internal/app/models"
	"erp/internal/app/repositories"
)

type MotivationService struct {
	repo *repositories.MotivationRepository
}

func NewMotivationService(repo *repositories.MotivationRepository) *MotivationService {
	return &MotivationService{repo: repo}
}

func (s *MotivationService) CreateStep(step *models.MotivationStep) error {
	return s.repo.Create(step)
}

func (s *MotivationService) UpdateStep(step *models.MotivationStep) error {
	return s.repo.Update(step)
}

func (s *MotivationService) DeleteStep(id uint) error {
	return s.repo.Delete(id)
}

func (s *MotivationService) ListSteps() ([]models.MotivationStep, error) {
	return s.repo.GetAll()
}
