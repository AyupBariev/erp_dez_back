package models

import "time"

type EngineerPayout struct {
	ID             uint      `gorm:"primaryKey"`
	EngineerID     int64     `gorm:"index"`
	Month          time.Time `gorm:"type:date"`
	Prepayment     float64   `gorm:"not null"`
	PaidPrepayment float64   `gorm:"not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
