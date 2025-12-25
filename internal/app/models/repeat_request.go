package models

import (
	"gorm.io/gorm"
	"time"
)

type RepeatRequest struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OrderID       uint           `json:"order_id"`
	Order         Order          `gorm:"foreignKey:OrderID;references:ID"`
	EngineerID    uint           `json:"engineer_id"`
	Engineer      Engineer       `gorm:"foreignKey:EngineerID;references:ID"`
	Description   string         `gorm:"type:text" json:"description"`
	RequestedAt   time.Time      `json:"requested_at"`
	ScheduledAt   time.Time      `json:"scheduled_at"`
	Confirmed     bool           `json:"confirmed" gorm:"default:false"`
	ConfirmedAt   *time.Time     `json:"confirmed_at"`
	ConfirmedBy   *int64         `json:"confirmed_by"`
	ConfirmedUser *User          `gorm:"foreignKey:ConfirmedBy" json:"confirmed_user"`
	RepeatOrderID *uint64        `json:"repeat_order_id"`
	RepeatOrder   *Order         `gorm:"foreignKey:RepeatOrderID" json:"repeat_order"`
	Status        string         `json:"status" gorm:"default:pending"` // pending, confirmed, rejected
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
