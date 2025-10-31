package models

import "time"

type EngineerMonthlyMotivation struct {
	ID uint `gorm:"primaryKey"`

	EngineerID   int64     `gorm:"not null;index"`
	EngineerName string    `json:"engineer_name"`
	Month        time.Time `gorm:"not null"` // первый день месяца

	ReportsCount       int     `gorm:"not null;default:0"` // всего отчетов
	PrimaryOrdersCount int     `gorm:"not null;default:0"` // количество первичных заказов
	RepeatOrdersCount  int     `gorm:"not null;default:0"` // количество повторных заказов
	OrdersTotalAmount  float64 `gorm:"not null;default:0"` // сумма первичных заказов
	RepeatOrdersAmount float64 `gorm:"not null;default:0"` // сумма повторов

	GrossProfit  float64 `gorm:"not null;default:0"` // валовая прибыль
	AverageCheck float64 `gorm:"not null;default:0"` // средний чек (по первичным заказам)

	MotivationPercent float64 `gorm:"not null;default:0"` // итоговый % мотивации
	TotalMotivation   float64 `gorm:"not null;default:0"` // сумма мотивации
	ConfirmedByAdmin  bool    `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
