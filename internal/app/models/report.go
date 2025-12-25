package models

import "time"

type Report struct {
	ID                int64      `gorm:"primaryKey" json:"id"`
	OrderID           int64      `json:"order_id"`
	Order             Order      `gorm:"foreignKey:OrderID"`
	EngineerID        int64      `json:"engineer_id"`
	HasRepeat         bool       `json:"has_repeat"`
	RepeatDate        *time.Time `json:"repeat_date"`
	RepeatNote        string     `json:"repeat_note"`
	Description       string     `json:"description"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	IssuedMoney       float64    `json:"issued_money"`
	GaveCash          float64    `json:"gave_cash"`
	MotivationPercent uint       `json:"motivation_percent"`
	MotivationStepID  *uint      `json:"motivation_step_id"`
}
