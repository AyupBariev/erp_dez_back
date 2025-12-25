package models

import "time"

type EngineerMonthlyMotivation struct {
	ID uint `gorm:"primaryKey"`

	EngineerID int64     `gorm:"not null;index"`
	Month      time.Time `gorm:"not null;type:date"` // первый день месяца

	ReportsCount       int     `gorm:"not null;default:0"` // всего отчетов
	PrimaryOrdersCount int     `gorm:"not null;default:0"` // количество первичных заказов
	RepeatOrdersCount  int     `gorm:"not null;default:0"` // количество повторных заказов
	OrdersTotalAmount  float64 `gorm:"not null;default:0"` // прибыль с первичных заказов
	RepeatOrdersAmount float64 `gorm:"not null;default:0"` // прибыль с повторных заказов

	GrossProfit  float64 `gorm:"not null;default:0"` // валовая прибыль
	AverageCheck float64 `gorm:"not null;default:0"` // средний чек

	// 🔹 Разделение мотивации
	BaseMotivationPercent float64 `gorm:"not null;default:0"` // текущий процент по мотивационной сетке
	BonusPercent          float64 `gorm:"not null;default:0"` // бонусный процент, если достигнут бонус
	MotivationPercent     uint    `gorm:"not null;default:0"` // итоговый % мотивации (активный для расчета)
	TotalMotivationAmount float64 `gorm:"not null;default:0"` // общая сумма мотивации за месяц

	// 🔹 Учет сторонних выплат
	AggregatorPayout float64 `gorm:"not null;default:0"` // сумма выплат агрегаторам (если применимо)

	CreatedAt time.Time
	UpdatedAt time.Time
}
