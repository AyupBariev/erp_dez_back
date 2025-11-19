package domain

import "time"

type AggregatorDayPayout struct {
	Date       time.Time
	Aggregator string  // источник
	OrderCount int64   // count orders where orders.aggregator_id = aggregator.id
	OrdersSum  float64 // сумма orders.finish_price
	AvgCheck   float64 // OrdersSum / OrderCount
	LeadCost   float64 // сколько мы должны агрегатору
	//Paid       bool
	//PaidAt     *time.Time
}
