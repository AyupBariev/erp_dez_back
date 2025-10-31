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

// GetEngineerByTelegramUsername Найти инженера по Username
func (s *EngineerService) GetEngineerByUsername(username string) (*models.Engineer, error) {
	engineer, err := s.engineerRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.LogError(fmt.Sprintf("Failed to find engineer by username=%d", username), err)
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
	// Проверяем дубликат по Telegram ID
	if telegramID := engineer.GetTelegramID(); telegramID != nil {
		existingByTelegramId, err := s.engineerRepo.FindByTelegramID(*telegramID)
		if err == nil && existingByTelegramId != nil {
			return nil, fmt.Errorf("telegram conflict: %w", ErrEngineerAlreadyExists)
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to check engineer by telegram id: %w", err)
		}
	}
	// Проверяем дубликат по Username
	existingByUsername, err := s.engineerRepo.FindByUsername(engineer.Username)
	if err == nil && existingByUsername != nil {
		logger.LogInfo(fmt.Sprintf("Conflict: username '%s' already exists (id=%d)", existingByUsername.Username, existingByUsername.ID))
		return nil, fmt.Errorf("username conflict: %w", ErrEngineerAlreadyExists)
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check engineer by username: %w", err)
	}

	// Всегда ставим false при создании
	engineer.IsApproved = false

	if err := s.engineerRepo.Create(engineer); err != nil {
		logger.LogError("Failed to create engineer", err)
		return nil, err
	}

	return engineer, nil
}

func (s *EngineerService) UpdateTelegramID(id int64, telegramID int64) error {
	return s.engineerRepo.UpdateTelegramID(id, telegramID)
}

func (s *EngineerService) ListWorkingEngineers(date string) ([]*models.Engineer, error) {
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
