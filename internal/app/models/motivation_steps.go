package models

import "time"

type MotivationStep struct {
	ID               uint    `gorm:"primaryKey"`
	Name             string  // первичка/повтор/бонус
	MinAmount        float64 // минимальная сумма заказа для перехода текущий процент
	Percent          uint
	PercentIncrement uint
	Sort             uint
	Type             string `gorm:"type:enum('primary','repeat','bonus');not null;default:'primary'"` // Тип заказа
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
