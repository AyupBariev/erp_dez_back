package services

import (
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"erp/internal/pkg/logger"
	"errors"
	"fmt"
)

var ErrEngineerAlreadyExists = errors.New("engineer already exists")

type EngineerService struct {
	engineerRepo *repositories.EngineerRepository
}

func NewEngineerService(engineerRepo *repositories.EngineerRepository) *EngineerService {
	return &EngineerService{engineerRepo: engineerRepo}
}

// GetEngineerByTelegramID Найти инженера по Telegram ID
func (s *EngineerService) GetEngineerByTelegramID(telegramID int64) (*models.Engineer, error) {
	engineer, err := s.engineerRepo.FindByTelegramID(telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // не найден
		}
		logger.LogError(fmt.Sprintf("Failed to find engineer by telegramID=%d", telegramID), err)
		return nil, err
	}
	return engineer, nil
}

// GetEngineerByID Найти инженера по ID
func (s *EngineerService) GetEngineerByID(ID int64) (*models.Engineer, error) {
	engineer, err := s.engineerRepo.FindApprovedByID(ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // не найден
		}
		logger.LogError(fmt.Sprintf("Failed to find engineer by ID=%d", ID), err)
		return nil, err
	}
	return engineer, nil
}

// CreateEngineer Создать инженера
func (s *EngineerService) CreateEngineer(engineer *models.Engineer) (*models.Engineer, error) {
	existingEngineer, err := s.engineerRepo.FindByTelegramID(engineer.TelegramID)
	if err == nil && existingEngineer != nil {
		return nil, ErrEngineerAlreadyExists
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check engineer existence: %w", err)
	}

	// Всегда ставим false при создании
	engineer.IsApproved = false

	if err := s.engineerRepo.Create(engineer); err != nil {
		logger.LogError("Failed to create engineer", err)
		return nil, err
	}

	return engineer, nil
}

func (s *EngineerService) ListEngineers(date string) ([]*models.Engineer, error) {
	engineers, err := s.engineerRepo.ListWorking(date)
	if err != nil {
		return nil, fmt.Errorf("failed to get engineers: %w", err)
	}
	return engineers, nil
}

// ApproveEngineer Активировать учетку инженера
func (s *EngineerService) ApproveEngineer(engineerID int64) (*models.Engineer, error) {
	return s.engineerRepo.ApproveByID(engineerID)
}
