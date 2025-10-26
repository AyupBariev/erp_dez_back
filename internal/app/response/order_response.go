package response

import "erp/internal/app/models"

type OrderResponse struct {
	ID          int               `json:"id"`
	ERPNumber   int64             `json:"erp_number"`
	ClientName  string            `json:"client_name"`
	Address     string            `json:"address"`
	Problem     string            `json:"problem"`
	ScheduledAt string            `json:"scheduled_at"`
	Status      string            `json:"status"`
	Engineer    *EngineerResponse `json:"engineer"`
}

func FromOrderModel(e *models.Order) OrderResponse {
	order := OrderResponse{
		ID:          e.ID,
		ERPNumber:   e.ERPNumber,
		ClientName:  e.ClientName,
		Address:     e.Address,
		Problem:     e.Problem,
		ScheduledAt: e.ScheduledAt.Format("2006-01-02 15:04"),
		Status:      e.Status,
		Engineer:    nil,
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
