package dto

import "time"

type ProfitRaw struct {
	Date      time.Time `gorm:"column:date"`
	NetProfit float64   `gorm:"column:net_profit"`
}
