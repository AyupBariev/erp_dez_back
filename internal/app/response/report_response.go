package response

import (
	"erp/internal/domain"
	"time"
)

type CashReportResp struct {
	ID                int64      `json:"id"`
	OrderID           int64      `json:"order_id"`
	EngineerID        int64      `json:"engineer_id"`
	HasRepeat         bool       `json:"has_repeat"`
	RepeatDate        *time.Time `json:"repeat_date"`
	RepeatNote        string     `json:"repeat_note"`
	Description       string     `json:"description"`
	CreatedAt         string     `json:"created_at"`
	ToCash            float64    `json:"to_cash"`
	GaveCash          float64    `json:"gave_cash"`
	IssuedMoney       float64    `json:"issued_money"`
	MotivationPercent int64      `json:"motivation_percent"`
	Prepayment        float64    `json:"prepayment"`
	Salary            float64    `json:"salary"`
	OrderStatus       string     `json:"order_status"`
}

func FromCashReportList(list []*domain.CashReport) []CashReportResp {
	out := make([]CashReportResp, 0, len(list))
	for _, r := range list {
		out = append(out, CashReportResp{
			ID: r.ID, OrderID: r.OrderID, EngineerID: r.EngineerID,
			HasRepeat: r.HasRepeat, RepeatDate: r.RepeatDate, RepeatNote: r.RepeatNote,
			Description: r.Description, CreatedAt: r.CreatedAt.Format("2006-01-02 15:04"),
			ToCash: r.ToCash, GaveCash: r.GaveCash, IssuedMoney: r.IssuedMoney, MotivationPercent: int64(r.MotivationPercent),
			Prepayment:  r.Prepayment,
			Salary:      r.Salary,
			OrderStatus: r.OrderStatus,
		})
	}
	return out
}
