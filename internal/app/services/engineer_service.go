package services

import (
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"erp/internal/pkg/logger"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type PayoutRow struct {
	EngineerID int64  `json:"engineer_id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Month      string `json:"month"`

	Salary      float64 `json:"salary"`
	Advance     float64 `json:"advance"`
	PaidAdvance float64 `json:"paid_advance"`
	Left        float64 `json:"left"`
	Total       float64 `json:"total"`

	CanPay bool `json:"can_pay"`
}

var ErrEngineerAlreadyExists = errors.New("engineer already exists")

type EngineerService struct {
	engineerRepo *repositories.EngineerRepository
	db           *gorm.DB
}

func NewEngineerService(engineerRepo *repositories.EngineerRepository, db *gorm.DB) *EngineerService {
	return &EngineerService{engineerRepo: engineerRepo, db: db}
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

func (s *EngineerService) UpdateEngineerWorkingStatus(engineerID int64, isWorking bool) (*models.Engineer, error) {
	engineer, err := s.engineerRepo.FindByID(engineerID)
	if err != nil {
		return nil, fmt.Errorf("engineer not found: %w", err)
	}

	engineer.IsWorking = isWorking

	engineer, err = s.engineerRepo.UpdateWorkingStatus(engineerID, isWorking)
	if err != nil {
		return nil, fmt.Errorf("failed to update engineer working status: %w", err)
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

func (s *EngineerService) GetPayoutsByPeriod(from, to string) ([]PayoutRow, error) {
	start, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil, err
	}

	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// 1. инженеры
	var engineers []models.Engineer
	if err := s.db.Find(&engineers).Error; err != nil {
		return nil, err
	}

	// 2. мотивация за период
	var motivations []models.EngineerMonthlyMotivation
	s.db.
		Where("month >= ? AND month <= ?", start, end).
		Find(&motivations)

	// 3. выплаты за период
	var payouts []models.EngineerPayout
	s.db.
		Where("created_at BETWEEN ? AND ?", start, end).
		Find(&payouts)

	// агрегация
	motMap := make(map[int64]float64)
	for _, m := range motivations {
		motMap[m.EngineerID] += m.TotalMotivationAmount
	}

	payMap := make(map[int64]float64)
	for _, p := range payouts {
		payMap[p.EngineerID] += p.PaidPrepayment
	}

	var out []PayoutRow

	for _, e := range engineers {
		totalMot := motMap[int64(e.ID)]
		salary := totalMot / 2
		advance := salary
		paid := payMap[int64(e.ID)]

		left := advance - paid
		if left < 0 {
			left = 0
		}

		out = append(out, PayoutRow{
			EngineerID:  int64(e.ID),
			FirstName:   e.FirstName.String,
			SecondName:  e.SecondName.String,
			Salary:      salary,
			Advance:     advance,
			PaidAdvance: paid,
			Left:        left,
			Total:       salary + advance,
			CanPay:      left > 0,
		})
	}

	return out, nil
}

func (s *EngineerService) GetMonthPayouts(month string) ([]PayoutRow, error) {
	monthDate, _ := time.Parse("2006-01", month)
	monthStart := time.Date(monthDate.Year(), monthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthString := monthStart.Format("2006-01") + "-01"
	// 1. Получаем ВСЕХ инженеров
	var engineers []models.Engineer
	if err := s.db.Find(&engineers).Error; err != nil {
		return nil, err
	}

	// 2. Загружаем мотивацию за месяц
	var mot []models.EngineerMonthlyMotivation
	s.db.Where("month = ?", monthString).Find(&mot)

	motMap := make(map[int64]models.EngineerMonthlyMotivation)
	for _, m := range mot {
		motMap[m.EngineerID] = m
	}

	// 3. Загружаем выплаты за месяц
	var pays []models.EngineerPayout
	s.db.Where("month = ?", monthString).Find(&pays)

	payMap := make(map[int64]models.EngineerPayout)
	for _, p := range pays {
		payMap[p.EngineerID] = p
	}

	var out []PayoutRow

	// 4. Строим строки для каждого инженера
	for _, e := range engineers {

		// Мотивация или 0
		motivation := motMap[int64(e.ID)]
		salary := motivation.TotalMotivationAmount / 2

		// Аванс = половина ЗП
		advance := salary

		// Выплаты или 0
		payout := payMap[int64(e.ID)]
		paid := payout.PaidPrepayment

		// Осталось выплатить
		left := advance - paid
		if left < 0 {
			left = 0
		}

		canPay := motivation.MotivationPercent >= 20 && left > 0
		out = append(out, PayoutRow{
			EngineerID:  int64(e.ID),
			FirstName:   e.FirstName.String,
			SecondName:  e.SecondName.String,
			Month:       month,
			Salary:      salary,
			Advance:     advance,
			PaidAdvance: paid,
			Left:        left,
			Total:       salary + advance,
			CanPay:      canPay,
		})
	}

	return out, nil
}

func (s *EngineerService) PayAdvance(engineerID int64, month string, amount float64) error {
	monthDate, _ := time.Parse("2006-01", month)
	monthStart := time.Date(monthDate.Year(), monthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthString := monthDate.Format("2006-01") + "-01"

	var m models.EngineerMonthlyMotivation
	if err := s.db.Where("engineer_id = ? AND month = ?", engineerID, monthString).First(&m).Error; err != nil {
		return errors.New("motivation not found")
	}

	salary := m.TotalMotivationAmount
	advance := salary / 2

	var p models.EngineerPayout
	err := s.db.Where("engineer_id = ? AND month = ?", engineerID, monthString).First(&p).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		p = models.EngineerPayout{
			EngineerID:     engineerID,
			Month:          monthStart,
			Prepayment:     advance,
			PaidPrepayment: 0,
		}
		if err := s.db.Create(&p).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// проверяем лимит
	if p.PaidPrepayment+amount > advance+1e-2 {
		return fmt.Errorf("превышает доступный аванс: максимум %.2f ₽", advance-p.PaidPrepayment)
	}

	p.PaidPrepayment += amount
	return s.db.Save(&p).Error
}
