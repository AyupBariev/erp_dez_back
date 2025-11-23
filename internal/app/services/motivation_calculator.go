package services

import (
	"erp/internal/app/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
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
	totalOrderProfit, err := strconv.ParseFloat(finishPrice, 64)
	grossProfit := math.Round(totalOrderProfit * orderPercent / 100)
	monthly.AggregatorPayout += math.Round(totalOrderProfit - grossProfit)

	if err != nil {
		return fmt.Errorf("invalid finishPrice: %w", err)
	}

	monthly.ReportsCount++
	if isRepeat {
		monthly.RepeatOrdersCount++
		monthly.RepeatOrdersAmount += grossProfit
	} else {
		monthly.PrimaryOrdersCount++
		monthly.OrdersTotalAmount += grossProfit
	}

	monthly.GrossProfit += grossProfit

	totalAmount := monthly.OrdersTotalAmount + monthly.RepeatOrdersAmount
	if monthly.PrimaryOrdersCount > 0 {
		monthly.AverageCheck = math.Round(totalAmount / float64(monthly.PrimaryOrdersCount))
	}

	// 2️⃣ Загружаем шаги мотивации
	var steps []models.MotivationStep
	if err := mc.DB.Order("percent ASC").Find(&steps).Error; err != nil {
		return err
	}

	// 3️⃣ Определяем базовую и бонусную мотивацию
	var basePercent float64
	var bonusPercent float64
	orderType := "primary"
	if isRepeat {
		orderType = "repeat"
	}

	// Найдём бонусный шаг
	var bonusStep *models.MotivationStep
	for _, step := range steps {
		if step.Type == "bonus" {
			bonusStep = &step
			break
		}
	}

	// Проверяем: активируется ли бонус
	hasBonus := bonusStep != nil && totalOrderProfit >= bonusStep.MinAmount && monthly.BonusPercent == 0

	if hasBonus {
		// ✅ Бонус считается отдельно
		bonusPercent = bonusStep.Percent
	} else {
		// ✅ Работаем с сеткой progression

		currentBase := monthly.BaseMotivationPercent
		nextBase := currentBase
		log.Printf("STEP 1: loaded BaseMotivationPercent=%+v", currentBase)
		log.Printf("STEP 2: loaded nextBase=%+v", nextBase)

		for _, step := range steps {
			log.Printf("STEP 3: loaded orderType=%+v", orderType)
			log.Printf("STEP 3.1: loaded step.Type=%+v", step.Type)
			if step.Type == orderType && step.Type != "bonus" {
				log.Printf("STEP 4: loaded totalOrderProfit=%+v", totalOrderProfit)
				log.Printf("STEP 5: loaded  step.MinAmount =%+v", step.MinAmount)
				log.Printf("STEP 6: loaded  step.Percent =%+v", step.Percent)
				log.Printf("STEP 7: loaded  currentBase =%+v", currentBase)
				if totalOrderProfit >= step.MinAmount && step.Percent > currentBase {
					nextBase = step.Percent
					break // только одно продвижение
				}
			}
		}
		log.Printf("STEP 8: loaded  basePercent =%+v", basePercent)
		basePercent = nextBase
	}

	// ✅ Сохраняем только актуальные проценты
	monthly.BaseMotivationPercent = basePercent
	log.Printf("STEP 9: loaded  basePercent =%+v var = %+v", monthly.BaseMotivationPercent, basePercent)

	monthly.BonusPercent = bonusPercent

	// Итоговый процент (без двойного бонусирования)
	motivationPercent := min(basePercent+bonusPercent, 30)

	monthly.MotivationPercent = motivationPercent

	// 💰 Общая сумма мотивации
	monthly.TotalMotivationAmount = math.Round((monthly.OrdersTotalAmount + monthly.RepeatOrdersAmount) * motivationPercent / 100)
	return mc.DB.Save(&monthly).Error
}
