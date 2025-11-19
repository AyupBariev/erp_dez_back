package domain

import "time"

type CashReport struct {
	ID                int64
	OrderID           int64
	EngineerID        int64
	HasRepeat         bool
	RepeatDate        *time.Time
	RepeatNote        string
	Description       string
	CreatedAt         time.Time
	ToCash            float64 // orders.finish_price * orders.our_percent / 100
	GaveCash          float64 // reports.gave_cach
	IssuedMoney       float64 // reports.issued_money (only if motivation>=20 and prepayment<=salary/2)
	MotivationPercent float64 // engineer_monthly_motivations.motivation_percent
	Prepayment        float64 // engineer_monthly_motivations.total_motivation_amount * 0.5
	Salary            float64 // engineer_monthly_motivations.total_motivation_amount
	OrderStatus       string  // orders.status
}
