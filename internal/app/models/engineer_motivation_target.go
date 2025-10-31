package models

import "time"

type EngineerMotivationTarget struct {
	ID                uint      `gorm:"primaryKey"`
	EngineerID        uint      `gorm:"not null;index"`
	MotivationPercent float64   `gorm:"not null"`
	TargetMonth       time.Time `gorm:"not null"`
	ConfirmedByAdmin  bool      `gorm:"default:false"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
