// services/engineer_motivation_service.go
package services

import (
	"erp/internal/app/dto"
	"erp/internal/app/repositories"
	"erp/internal/app/response"
	"fmt"
	"time"
)

type EngineerMotivationService struct {
	repo *repositories.EngineerMotivationRepository
}

func NewEngineerMotivationService(repo *repositories.EngineerMotivationRepository) *EngineerMotivationService {
	return &EngineerMotivationService{repo: repo}
}

var ruMonths = []string{
	"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
	"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
}

func formatMonthYear(t time.Time) string {
	month := ruMonths[int(t.Month())-1]
	return fmt.Sprintf("%s %d", month, t.Year())
}

func (s *EngineerMotivationService) GetMotivationByDateRange(fromStr, toStr string) ([]dto.EngineerMotivationView, error) {
	// Парсим даты
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты 'from': %v", err)
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты 'to': %v", err)
	}

	// Валидация диапазона
	if from.After(to) {
		return nil, fmt.Errorf("дата 'from' должна быть раньше даты 'to'")
	}

	return s.repo.GetByDateRange(from, to)
}

func (s *EngineerMotivationService) GetProfitByDateRange(
	fromStr, toStr string,
) ([]response.ProfitResponse, error) {

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return nil, err
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return nil, err
	}

	if from.After(to) {
		return nil, fmt.Errorf("from > to")
	}

	rows, err := s.repo.GetProfitByDateRange(from, to)
	if err != nil {
		return nil, err
	}

	out := make([]response.ProfitResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, response.ProfitResponse{
			Period:    formatMonthYear(r.Date),
			NetProfit: r.NetProfit,
		})
	}

	return out, nil
}

func (s *EngineerMotivationService) GetMonthlyMotivation(monthStr string) ([]dto.EngineerMotivationView, error) {
	var month time.Time
	var err error
	if monthStr == "" {
		now := time.Now()
		month = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	} else {
		month, err = time.Parse("2006-01", monthStr)
		if err != nil {
			return nil, err
		}
	}
	return s.repo.GetByMonth(month)
}
