package response

import (
	"erp/internal/app/models"
	"time"
)

type RepeatRequestResponse struct {
	ID          uint          `json:"id"`
	OrderID     uint          `json:"order_id"`
	Order       OrderResponse `json:"order"`
	Description string        `json:"description"`
	RequestedAt time.Time     `json:"requested_at"`
	ScheduledAt time.Time     `json:"scheduled_at"`
	Confirmed   bool          `json:"confirmed"`
	ConfirmedAt *time.Time    `json:"confirmed_at,omitempty"`
	ConfirmedBy *int64        `json:"confirmed_by,omitempty"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"CreatedAt"`
	UpdatedAt   time.Time     `json:"UpdatedAt"`
}

// FromRepeatRequestModel преобразует модель в response DTO
func FromRepeatRequestModel(req *models.RepeatRequest) RepeatRequestResponse {
	return RepeatRequestResponse{
		ID:          req.ID,
		OrderID:     req.OrderID,
		Order:       FromOrderModel(&req.Order),
		Description: req.Description,
		RequestedAt: req.RequestedAt,
		ScheduledAt: req.ScheduledAt,
		Confirmed:   req.Confirmed,
		ConfirmedAt: req.ConfirmedAt,
		ConfirmedBy: req.ConfirmedBy,
		Status:      req.Status,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}
}

func FromRepeatRequestList(list []models.RepeatRequest) []RepeatRequestResponse {
	resp := make([]RepeatRequestResponse, len(list))
	for i, req := range list {
		resp[i] = FromRepeatRequestModel(&req)
	}
	return resp
}
