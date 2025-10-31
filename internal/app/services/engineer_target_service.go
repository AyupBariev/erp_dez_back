package services

import (
	"erp/internal/app/models"
	"erp/internal/app/repositories"
)

type EngineerTargetService struct {
	repo *repositories.EngineerTargetRepository
}

func NewEngineerTargetService(repo *repositories.EngineerTargetRepository) *EngineerTargetService {
	return &EngineerTargetService{repo: repo}
}

func (s *EngineerTargetService) CreateTarget(target *models.EngineerMotivationTarget) error {
	return s.repo.Create(target)
}

func (s *EngineerTargetService) UpdateTarget(target *models.EngineerMotivationTarget) error {
	return s.repo.Update(target)
}

func (s *EngineerTargetService) GetTargets(engineerID uint, monthStart, monthEnd string) ([]models.EngineerMotivationTarget, error) {
	return s.repo.GetByEngineerAndMonth(engineerID, monthStart, monthEnd)
}
