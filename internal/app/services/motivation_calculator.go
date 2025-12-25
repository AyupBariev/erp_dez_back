package services

import (
	"erp/internal/app/models"
	"errors"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"
)

type MotivationCalculator struct {
	DB *gorm.DB
}

//// UpdateEngineerMonthlyMotivation обновляет месячную мотивацию инженера после отправки отчета
//func (mc *MotivationCalculator) UpdateEngineerMonthlyMotivation(
//	engineerID int64,
//	totalOrderPrice string,
//	OurPercent float64,
//	isRepeat bool,
//) error {
//	now := time.Now()
//	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
//
//	var monthly models.EngineerMonthlyMotivation
//	err := mc.DB.Where("engineer_id = ? AND month = ?", engineerID, monthStart.Format("2006-01-02")).First(&monthly).Error
//	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
//		return err
//	}
//
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		monthly = models.EngineerMonthlyMotivation{
//			EngineerID: engineerID,
//			Month:      monthStart,
//		}
//	}
//
//	// 1️⃣ Разбираем входные данные
//	totalOrderProfit, err := strconv.ParseFloat(totalOrderPrice, 64)
//	grossProfit := math.Round(totalOrderProfit * OurPercent / 100)
//	monthly.AggregatorPayout += math.Round(totalOrderProfit - grossProfit)
//
//	if err != nil {
//		return fmt.Errorf("invalid totalOrderPrice: %w", err)
//	}
//
//	monthly.ReportsCount++
//	if isRepeat {
//		monthly.RepeatOrdersCount++
//		monthly.RepeatOrdersAmount += grossProfit
//	} else {
//		monthly.PrimaryOrdersCount++
//		monthly.OrdersTotalAmount += grossProfit
//	}
//
//	monthly.GrossProfit += grossProfit
//
//	totalAmount := monthly.OrdersTotalAmount + monthly.RepeatOrdersAmount
//	if monthly.PrimaryOrdersCount > 0 {
//		monthly.AverageCheck = math.Round(totalAmount / float64(monthly.PrimaryOrdersCount))
//	}
//
//	// 2️⃣ Загружаем шаги мотивации
//	var steps []models.MotivationStep
//	if err := mc.DB.Order("percent ASC").Find(&steps).Error; err != nil {
//		return err
//	}
//
//	// 3️⃣ Определяем базовую и бонусную мотивацию
//	var basePercent float64
//	var bonusPercent float64
//
//	// Найдём бонусный шаг
//	var bonusStep *models.MotivationStep
//	for _, step := range steps {
//		if step.Type == "bonus" {
//			bonusStep = &step
//			break
//		}
//	}
//
//	// Проверяем: активируется ли бонус
//	hasBonus := bonusStep != nil && totalOrderProfit >= bonusStep.MinAmount && monthly.BonusPercent == 0
//
//	if hasBonus {
//		// ✅ Бонус считается отдельно
//		bonusPercent = bonusStep.Percent
//	} else {
//		// ✅ Работаем с сеткой progression
//		currentBase := monthly.BaseMotivationPercent
//		nextBase := currentBase
//		for _, step := range steps {
//			if step.Percent > nextBase && nextBase == currentBase && step.Type != "bonus" {
//				nextBase = step.Percent
//			}
//			if totalOrderProfit >= step.MinAmount {
//				break // только одно продвижение
//			}
//		}
//		basePercent = nextBase
//	}
//
//	// ✅ Сохраняем только актуальные проценты
//	monthly.BaseMotivationPercent = basePercent
//	monthly.BonusPercent = bonusPercent
//
//	// Итоговый процент (без двойного бонусирования)
//	motivationPercent := min(basePercent+bonusPercent, 30)
//	monthly.MotivationPercent = motivationPercent
//
//	// 💰 Общая сумма мотивации
//	monthly.TotalMotivationAmount = math.Round((monthly.AggregatorPayout + monthly.GrossProfit) * motivationPercent / 100)
//	return mc.DB.Save(&monthly).Error
//}

func (mc *MotivationCalculator) CalculateReportMotivation(
	report *models.Report,
	totalOrderPrice float64,
	isOrderRepeat bool,
) error {
	var steps []models.MotivationStep
	if err := mc.DB.Order("sort ASC").Find(&steps).Error; err != nil {
		return err
	}

	monthStart := time.Date(report.CreatedAt.Year(), report.CreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)

	var monthReports []models.Report
	if err := mc.DB.
		Where("engineer_id = ? AND created_at >= ? AND created_at < ?", report.EngineerID, monthStart, monthEnd).
		Find(&monthReports).Error; err != nil {
		return err
	}

	// Справочник шагов по ID
	stepByID := make(map[uint]models.MotivationStep, len(steps))
	for _, s := range steps {
		stepByID[s.ID] = s
	}

	usedSteps := map[uint]int64{} // stepID -> reportID
	hasAnyBaseStep := false
	for _, r := range monthReports {
		if r.ID == report.ID {
			continue // ⬅️ ВАЖНО: текущий отчет не учитываем
		}
		if r.MotivationStepID != nil {
			usedSteps[*r.MotivationStepID] = r.ID
			step := stepByID[*r.MotivationStepID]
			if step.Type == "primary" || step.Type == "repeat" {
				hasAnyBaseStep = true
			}
		}
	}

	var percent uint
	var stepID *uint

	checkStep := func(step models.MotivationStep, stepValue uint) bool {
		lastReportID, used := usedSteps[step.ID]
		if used && lastReportID != report.ID {
			return false // шаг уже использован другим отчетом
		}
		if totalOrderPrice >= step.MinAmount {
			percent = stepValue
			stepID = &step.ID
			return true
		}
		return false
	}

	// 1️⃣ БОНУС
	for _, step := range steps {
		if step.Type == "bonus" && checkStep(step, step.Percent) {
			goto SAVE
		}
	}

	// 2️⃣ ПЕРВЫЙ ОТЧЁТ МЕСЯЦА
	if !hasAnyBaseStep {
		for _, step := range steps {
			if (step.Type == "primary" && !isOrderRepeat) || (step.Type == "repeat" && isOrderRepeat) {
				if checkStep(step, step.Percent) {
					goto SAVE
				}
			}
		}
	}

	// 3️⃣ ИНКРЕМЕНТ
	for _, step := range steps {
		if (step.Type == "primary" && !isOrderRepeat) || (step.Type == "repeat" && isOrderRepeat) {
			if checkStep(step, step.PercentIncrement) {
				goto SAVE
			}
		}
	}

SAVE:
	report.MotivationPercent = percent
	report.MotivationStepID = stepID
	return nil
}
func (mc *MotivationCalculator) RecalculateEngineerMonthlyMotivation(report *models.Report) error {
	monthStart := time.Date(
		report.CreatedAt.Year(),
		report.CreatedAt.Month(),
		1, 0, 0, 0, 0,
		time.UTC,
	)
	monthEnd := monthStart.AddDate(0, 1, 0)

	var monthly models.EngineerMonthlyMotivation
	err := mc.DB.
		Where("engineer_id = ? AND month = ?", report.EngineerID, monthStart).
		First(&monthly).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		monthly = models.EngineerMonthlyMotivation{
			EngineerID: report.EngineerID,
			Month:      monthStart,
		}
	} else if err != nil {
		return err
	}

	// 🔹 получаем все отчеты месяца с предзагрузкой заказов
	var reports []models.Report
	if err := mc.DB.
		Preload("Order").
		Where("engineer_id = ? AND created_at >= ? AND created_at < ?", report.EngineerID, monthStart, monthEnd).
		Find(&reports).Error; err != nil {
		return err
	}

	// 🔄 сброс агрегатов
	monthly.ReportsCount = 0
	monthly.PrimaryOrdersCount = 0
	monthly.RepeatOrdersCount = 0
	monthly.OrdersTotalAmount = 0
	monthly.RepeatOrdersAmount = 0
	monthly.GrossProfit = 0
	monthly.AggregatorPayout = 0
	monthly.AverageCheck = 0
	monthly.MotivationPercent = 0
	monthly.TotalMotivationAmount = 0

	var totalPercent uint
	var totalAmount float64

	for _, r := range reports {
		o := r.Order

		monthly.ReportsCount++
		totalPercent += r.MotivationPercent

		finishPrice, _ := strconv.ParseFloat(o.FinishPrice, 64)
		grossProfit := math.Round(finishPrice * o.OurPercent / 100)
		payout := finishPrice - grossProfit

		if o.RepeatID != nil {
			monthly.RepeatOrdersCount++
			monthly.RepeatOrdersAmount += grossProfit
		} else {
			monthly.PrimaryOrdersCount++
			monthly.OrdersTotalAmount += grossProfit
		}

		monthly.GrossProfit += grossProfit
		monthly.AggregatorPayout += payout
	}

	// 🔒 лимит 30%
	if totalPercent > 30 {
		totalPercent = 30
	}
	monthly.MotivationPercent = totalPercent

	totalAmount = monthly.OrdersTotalAmount + monthly.RepeatOrdersAmount
	if monthly.PrimaryOrdersCount > 0 {
		monthly.AverageCheck = math.Round(totalAmount / float64(monthly.PrimaryOrdersCount))
	}

	monthly.TotalMotivationAmount = math.Round((monthly.AggregatorPayout + monthly.GrossProfit) * float64(monthly.PrimaryOrdersCount) / 100)

	return mc.DB.Save(&monthly).Error
}
