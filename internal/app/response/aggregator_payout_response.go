package response

import (
	"erp/internal/domain"
)

type AggregatorDayPayoutResp struct {
	Date       string  `json:"date"`
	Aggregator string  `json:"aggregator"`
	OrderCount int64   `json:"order_count"`
	OrdersSum  float64 `json:"orders_sum"`
	AvgCheck   float64 `json:"avg_check"`
	LeadCost   float64 `json:"lead_cost"`
	//Paid       bool    `json:"paid"`
	//PaidAt     *string `json:"paid_at,omitempty"`
}

func FromAggregatorDayPayoutList(list []*domain.AggregatorDayPayout) []AggregatorDayPayoutResp {
	out := make([]AggregatorDayPayoutResp, 0, len(list))
	for _, v := range list {
		p := AggregatorDayPayoutResp{
			Date:       v.Date.Format("2006-01-02"),
			Aggregator: v.Aggregator,
			OrderCount: v.OrderCount,
			OrdersSum:  v.OrdersSum,
			AvgCheck:   v.AvgCheck,
			LeadCost:   v.LeadCost,
		}
		out = append(out, p)
	}
	return out
}
