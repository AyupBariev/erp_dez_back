package models

import "database/sql"

type Engineer struct {
	ID         int            `db:"id" json:"id"`
	FirstName  sql.NullString `db:"first_name" json:"first_name"`
	SecondName sql.NullString `db:"second_name" json:"second_name"`
	Username   string         `db:"username" json:"username"`
	Phone      sql.NullString `db:"phone" json:"phone"`
	TelegramID int64          `db:"telegram_id" json:"telegram_id"`
	IsApproved bool           `db:"is_approved" json:"is_approved"`
}
