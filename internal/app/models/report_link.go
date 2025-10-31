package models

import "time"

type ReportLink struct {
	ID         int64     `db:"id"`
	OrderID    int64     `db:"order_id"`
	EngineerID int64     `db:"engineer_id"`
	Token      string    `db:"token"`
	ExpiresAt  time.Time `db:"expires_at"`
	CreatedAt  time.Time `db:"created_at"`
}
