package models

import (
	"database/sql"
	"time"
)

type OrderStatus string

type Order struct {
	ID          int           `db:"id"`
	ERPNumber   int64         `db:"erp_number"`
	SourceID    int           `db:"source_id"`
	OurPercent  float64       `db:"our_percent"`
	Price       string        `db:"price"`
	ClientName  string        `db:"client_name"`
	EngineerID  sql.NullInt64 `db:"engineer_id"`
	Engineer    *Engineer     `json:"engineer,omitempty"`
	AdminID     int64         `db:"admin_id"`
	Phones      []string      `db:"phones"` // JSON в базе
	Address     string        `db:"address"`
	Title       string        `db:"title"`
	ProblemID   sql.NullInt64 `db:"engineer_id"`
	Problem     string        `db:"problem"`
	ScheduledAt time.Time     `db:"scheduled_at"`
	Status      string        `db:"status"`
	ConfirmedAt sql.NullTime  `db:"confirmed_at"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}
