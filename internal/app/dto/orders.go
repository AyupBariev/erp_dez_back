package dto

type UpdateOrderRequest struct {
	Status       string   `json:"status"`
	ERPNumber    int64    `json:"erp_number"`
	AggregatorID int64    `json:"aggregator_id"`
	ProblemID    int64    `json:"problem_id"`
	Price        string   `json:"price"`
	FinishPrice  string   `json:"finish_price"`
	OurPercent   float64  `json:"our_percent"`
	ClientName   string   `json:"client_name"`
	Phones       []string `json:"phones"`
	Address      string   `json:"address"`
	WorkVolume   string   `json:"work_volume"`
	Note         string   `json:"note"`
	ScheduledAt  string   `json:"scheduled_at"` // ISO8601
	EngineerID   *int     `json:"engineer_id,omitempty"`
	AdminID      *int     `json:"admin_id,omitempty"`
}

type CreateOrderRequest struct {
	Status       string   `json:"status"`
	AggregatorID int      `json:"aggregator_id"`
	ProblemID    int64    `json:"problem_id" binding:"required"`
	Price        string   `json:"price"`
	OurPercent   float64  `json:"our_percent"`
	ClientName   string   `json:"client_name"`
	Phones       []string `json:"phones"`
	Address      string   `json:"address"`
	WorkVolume   string   `json:"work_volume"`
	Note         string   `json:"note"`
	ScheduledAt  string   `json:"scheduled_at"` // ISO8601
	EngineerID   *int64   `json:"engineer_id,omitempty"`
	AdminID      *int     `json:"admin_id,omitempty"`
}
