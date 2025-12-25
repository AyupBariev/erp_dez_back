package response

import (
	"erp/internal/app/models"
)

type OrderResponse struct {
	ID                int                `json:"id"`
	ERPNumber         int64              `json:"erp_number"`
	ClientName        string             `json:"client_name"`
	Address           string             `json:"address"`
	Price             string             `json:"price"`
	FinishPrice       string             `json:"finish_price"`
	OurPercent        float64            `json:"our_percent"`
	WorkVolume        string             `json:"work_volume"`
	AggregatorID      int64              `json:"aggregator_id"`
	ProblemID         int64              `json:"problem_id"`
	Aggregator        *models.Aggregator `json:"aggregator"`
	Problem           *models.Problem    `json:"problem"`
	ScheduledAt       string             `json:"scheduled_at"`
	Status            string             `json:"status"`
	Phones            []string           `json:"phones"`
	Engineer          *EngineerResponse  `json:"engineer"`
	RepeatID          *uint64            `json:"repeat_id"`
	RepeatERPNumber   *int64             `json:"repeat_erp_number,omitempty"`
	RepeatDescription string             `json:"repeat_description"`
	RepeatedBy        string             `json:"repeated_by"`
}

func FromOrderModel(e *models.Order) OrderResponse {
	var repeatERP *int64
	if e.ParentOrder != nil {
		repeatERP = &e.ParentOrder.ERPNumber
	}
	order := OrderResponse{
		ID:                int(e.ID),
		ERPNumber:         e.ERPNumber,
		ClientName:        e.ClientName,
		Address:           e.Address,
		Price:             e.Price,
		FinishPrice:       e.FinishPrice,
		OurPercent:        e.OurPercent,
		WorkVolume:        e.WorkVolume,
		ScheduledAt:       e.ScheduledAt.Format("2006-01-02 15:04"),
		Status:            e.Status,
		Phones:            e.Phones,
		RepeatID:          e.RepeatID,
		RepeatERPNumber:   repeatERP,
		RepeatDescription: e.RepeatDescription,
		RepeatedBy:        e.RepeatedBy,
		Engineer:          nil,
	}

	if e.Problem.ID != 0 {
		order.ProblemID = e.Problem.ID
		order.Problem = e.Problem
	}

	if e.Aggregator.ID != 0 {
		order.AggregatorID = e.Aggregator.ID
		order.Aggregator = e.Aggregator
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
