package services

import (
	"erp/internal/app/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"
)

type MotivationCalculator struct {
	DB *gorm.DB
}

// UpdateEngineerMonthlyMotivation обновляет месячную мотивацию инженера после отправки отчета
func (mc *MotivationCalculator) UpdateEngineerMonthlyMotivation(
	engineerID int64,
	orderAmount string,
	orderPercent float64, // теперь передаем OurPercent
	isRepeat bool,
) error {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	var monthly models.EngineerMonthlyMotivation
	err := mc.DB.Where("engineer_id = ? AND month = ?", engineerID, monthStart).First(&monthly).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Создаем новую запись
		monthly = models.EngineerMonthlyMotivation{
			EngineerID: engineerID,
			Month:      monthStart,
		}
	}

	// Обновляем счетчики
	monthly.ReportsCount += 1
	amount, err := strconv.ParseFloat(orderAmount, 64)
	if err != nil {
		return fmt.Errorf("invalid orderAmount: %w", err)
	}
	if !isRepeat {
		monthly.PrimaryOrdersCount += 1
		monthly.OrdersTotalAmount += amount
		monthly.GrossProfit += amount * orderPercent / 100
	} else {
		monthly.RepeatOrdersCount += 1
		monthly.RepeatOrdersAmount += amount
	}

	// Средний чек по первичным заказам
	if monthly.PrimaryOrdersCount > 0 {
		monthly.AverageCheck = monthly.OrdersTotalAmount / float64(monthly.PrimaryOrdersCount)
	} else {
		monthly.AverageCheck = 0
	}

	// Рассчитываем процент мотивации
	var steps []models.MotivationStep
	if err := mc.DB.Order("sort ASC").Find(&steps).Error; err != nil {
		return err
	}

	motivationPercent := 0.0
	for _, step := range steps {
		switch step.OrderType {
		case "primary":
			if monthly.OrdersTotalAmount >= step.MinAmount {
				motivationPercent += step.Percent
			}
		case "repeat":
			if monthly.RepeatOrdersAmount >= step.MinAmount {
				motivationPercent += step.Percent
			}
		case "bonus":
			if monthly.OrdersTotalAmount+monthly.RepeatOrdersAmount >= step.MinAmount {
				motivationPercent += step.Percent
			}
		}
	}

	// Ограничение 30%
	monthly.MotivationPercent = math.Min(motivationPercent, 30)

	// Сумма мотивации (берем все заказы)
	monthly.TotalMotivation = (monthly.OrdersTotalAmount + monthly.RepeatOrdersAmount) * monthly.MotivationPercent / 100

	// Сохраняем
	return mc.DB.Save(&monthly).Error
}
