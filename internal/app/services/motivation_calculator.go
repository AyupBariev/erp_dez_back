package services

import (
	"erp/internal/app/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type MotivationCalculator struct {
	DB *gorm.DB
}

// UpdateEngineerMonthlyMotivation обновляет месячную мотивацию инженера после отправки отчета
func (mc *MotivationCalculator) UpdateEngineerMonthlyMotivation(
	engineerID int64,
	finishPrice string,
	orderPercent float64,
	isRepeat bool,
) error {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	var monthly models.EngineerMonthlyMotivation
	err := mc.DB.Where("engineer_id = ? AND month = ?", engineerID, monthStart.Format("2006-01-02")).First(&monthly).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		monthly = models.EngineerMonthlyMotivation{
			EngineerID: engineerID,
			Month:      monthStart,
		}
	}

	// 1️⃣ Разбираем входные данные
	amount, err := strconv.ParseFloat(finishPrice, 64)
	if err != nil {
		return fmt.Errorf("invalid finishPrice: %w", err)
	}

	monthly.ReportsCount++
	if isRepeat {
		monthly.RepeatOrdersCount++
		monthly.RepeatOrdersAmount += amount
	} else {
		monthly.PrimaryOrdersCount++
		monthly.OrdersTotalAmount += amount
	}

	monthly.GrossProfit += amount * orderPercent / 100

	totalAmount := monthly.OrdersTotalAmount + monthly.RepeatOrdersAmount
	if monthly.PrimaryOrdersCount > 0 {
		monthly.AverageCheck = totalAmount / float64(monthly.PrimaryOrdersCount)
	}

	// 2️⃣ Загружаем шаги мотивации
	var steps []models.MotivationStep
	if err := mc.DB.Order("sort ASC").Find(&steps).Error; err != nil {
		return err
	}

	// 3️⃣ Определяем процент по текущему заказу
	var motivationPercent float64
	orderType := "primary"
	if isRepeat {
		orderType = "repeat"
	}

	for _, step := range steps {
		if step.OrderType == orderType && amount >= step.MinAmount {
			motivationPercent = step.Percent
		}
	}

	// 4️⃣ Проверяем бонус
	hasBonus := totalAmount >= 100000 // если хотя бы один заказ на 100к+
	if hasBonus {
		motivationPercent += 5
	}

	// Ограничение по максимуму
	if motivationPercent > 30 {
		motivationPercent = 30
	}

	monthly.MotivationPercent = motivationPercent

	// 5️⃣ Расчёт общей мотивации
	monthly.TotalMotivationAmount = totalAmount * motivationPercent / 100

	return mc.DB.Save(&monthly).Error
}
