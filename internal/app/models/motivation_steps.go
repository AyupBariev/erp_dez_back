package models

import "time"

type MotivationStep struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  // первичка/повтор/бонус
	MinAmount float64 // минимальная сумма заказа для перехода текущий процент
	Percent   float64
	Sort      uint
	OrderType string `gorm:"type:enum('primary','repeat','bonus');not null;default:'primary'"` // Тип заказа
	CreatedAt time.Time
	UpdatedAt time.Time
}
