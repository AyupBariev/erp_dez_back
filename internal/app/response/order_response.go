package response

import (
	"erp/internal/app/models"
)

type OrderResponse struct {
	ID           int                    `json:"id"`
	ERPNumber    int64                  `json:"erp_number"`
	ClientName   string                 `json:"client_name"`
	Address      string                 `json:"address"`
	Price        string                 `json:"price"`
	OurPercent   float64                `json:"our_percent"`
	WorkVolume   string                 `json:"work_volume"`
	AggregatorID int64                  `json:"aggregator_id"`
	ProblemID    int                    `json:"problem_id"`
	Aggregator   *models.BaseDictionary `json:"aggregator"`
	Problem      *models.BaseDictionary `json:"problem"`
	ScheduledAt  string                 `json:"scheduled_at"`
	Status       string                 `json:"status"`
	Phones       []string               `json:"phones"`
	Engineer     *EngineerResponse      `json:"engineer"`
}

func FromOrderModel(e *models.Order) OrderResponse {
	order := OrderResponse{
		ID:           e.ID,
		ERPNumber:    e.ERPNumber,
		ClientName:   e.ClientName,
		Address:      e.Address,
		Price:        e.Price,
		OurPercent:   e.OurPercent,
		WorkVolume:   e.WorkVolume,
		ProblemID:    e.Problem.ID,
		Problem:      e.Problem,
		AggregatorID: e.AggregatorID,
		Aggregator:   e.Aggregator,
		ScheduledAt:  e.ScheduledAt.Format("2006-01-02 15:04"),
		Status:       e.Status,
		Phones:       e.Phones,
		Engineer:     nil,
	}

	if e.Engineer != nil {
		engineer := FromEngineerModel(e.Engineer)
		order.Engineer = &engineer
	}

	return order
}

func FromOrderList(list []*models.Order) []OrderResponse {
	resp := make([]OrderResponse, 0, len(list))
	for _, e := range list {
		resp = append(resp, FromOrderModel(e))
	}
	return resp
}
