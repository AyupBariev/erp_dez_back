package response

type ProfitResponse struct {
	Period    string  `json:"period"` // "Декабрь 2025"
	NetProfit float64 `json:"net_profit"`
}
