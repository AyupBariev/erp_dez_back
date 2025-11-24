package models

// EngineerMotivationView используется только для выборки объединённых данных
// из таблиц engineers и engineer_monthly_motivations.
type EngineerMotivationView struct {
	EngineerID            int64   `json:"engineer_id"`
	EngineerName          string  `json:"engineer_name"`
	ReportsCount          int     `json:"reports_count"`
	PrimaryOrdersCount    int     `json:"primary_orders_count"`
	RepeatOrdersCount     int     `json:"repeat_orders_count"`
	OrdersTotalAmount     float64 `json:"orders_total_amount"`
	RepeatOrdersAmount    float64 `json:"repeat_orders_amount"`
	GrossProfit           float64 `json:"gross_profit"`
	NetProfit             float64 `json:"net_profit"`
	AverageCheck          float64 `json:"average_check"`
	MotivationPercent     float64 `json:"motivation_percent"`
	TotalMotivationAmount float64 `json:"total_motivation_amount"`
	AggregatorPayout      float64 `json:"aggregator_payout"`
	TotalAmount           float64 `json:"total_amount"`
}
