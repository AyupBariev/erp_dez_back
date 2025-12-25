package models

import (
	"time"
)

type OrderStatus string

type Order struct {
	ID                uint        `gorm:"primaryKey;autoIncrement"`
	ERPNumber         int64       `gorm:"column:erp_number;not null"`
	AggregatorID      int64       `gorm:"column:aggregator_id;not null"`
	Aggregator        *Aggregator `gorm:"foreignKey:AggregatorID"`
	OurPercent        float64     `gorm:"column:our_percent"`
	AggregatorPayout  float64     `gorm:"column:aggregator_payout"`
	Price             string      `gorm:"column:price"`
	FinishPrice       string      `gorm:"column:finish_price"`
	ClientName        string      `gorm:"column:client_name"`
	EngineerID        *int64      `gorm:"column:engineer_id"`
	Engineer          *Engineer   `gorm:"foreignKey:EngineerID" json:"engineer,omitempty"`
	AdminID           int64       `gorm:"column:admin_id"`
	Phones            StringArray `gorm:"type:json"`
	Address           string      `gorm:"column:address"`
	WorkVolume        string      `gorm:"column:work_volume"`
	ProblemID         int64       `gorm:"column:problem_id"`
	Problem           *Problem    `gorm:"foreignKey:ProblemID"`
	Note              string      `gorm:"column:note"`
	ScheduledAt       time.Time   `gorm:"column:scheduled_at"`
	Status            string      `gorm:"column:status"`
	ConfirmedAt       *time.Time  `gorm:"column:confirmed_at"`
	RepeatID          *uint64     `gorm:"column:repeat_id"`
	ParentOrder       *Order      `gorm:"foreignKey:RepeatID" json:"parent_order,omitempty"`
	RepeatDescription string      `gorm:"column:repeat_description"`
	RepeatedBy        string      `gorm:"column:repeated_by"`
	CreatedAt         time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time   `gorm:"column:updated_at;autoUpdateTime"`
}
