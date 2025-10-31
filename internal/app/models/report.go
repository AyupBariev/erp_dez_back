package models

import "time"

type Report struct {
	ID          int64      `db:"id"`
	OrderID     int64      `db:"order_id"`
	EngineerID  int64      `db:"engineer_id"`
	HasRepeat   bool       `db:"has_repeat"`
	RepeatDate  *time.Time `db:"repeat_date"`
	RepeatNote  string     `db:"repeat_note"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
}
