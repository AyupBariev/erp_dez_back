package models

import (
	"database/sql"
	"time"
)

type OrderStatus string

type Order struct {
	ID                int             `db:"id"`
	ERPNumber         int64           `db:"erp_number"`
	AggregatorID      int64           `db:"aggregator_id"`
	Aggregator        *BaseDictionary `json:"aggregator,omitempty"`
	OurPercent        float64         `db:"our_percent"`
	Price             string          `db:"price"`
	ClientName        string          `db:"client_name"`
	EngineerID        sql.NullInt64   `db:"engineer_id"`
	Engineer          *Engineer       `json:"engineer,omitempty"`
	AdminID           int64           `db:"admin_id"`
	Phones            []string        `db:"phones"` // JSON в базе
	Address           string          `db:"address"`
	WorkVolume        string          `db:"work_volume"`
	ProblemID         sql.NullInt64   `db:"problem_id"`
	Problem           *BaseDictionary `json:"problem,omitempty"`
	Note              string          `db:"note"`
	ScheduledAt       time.Time       `db:"scheduled_at"`
	Status            string          `db:"status"`
	ConfirmedAt       sql.NullTime    `db:"confirmed_at"`
	RepeatID          int
	RepeatDescription string
	RepeatedBy        string
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
