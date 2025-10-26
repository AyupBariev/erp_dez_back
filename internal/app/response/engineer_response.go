package response

import "erp/internal/app/models"

type EngineerResponse struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"first_name"`
	SecondName string  `json:"second_name"`
	Name       string  `json:"name"`
	Username   string  `json:"username"`
	Phone      *string `json:"phone,omitempty"`
	TelegramID int64   `json:"telegram_id"`
	IsApproved bool    `json:"is_approved"`
	IsWorking  bool    `json:"is_working"`
}

func FromEngineerModel(e *models.Engineer) EngineerResponse {
	name := e.Username
	if e.FirstName.Valid || e.SecondName.Valid {
		name = e.FirstName.String + " " + e.SecondName.String
	}

	var phone *string
	if e.Phone.Valid {
		phone = &e.Phone.String
	}

	return EngineerResponse{
		ID:         e.ID,
		FirstName:  e.FirstName.String,
		SecondName: e.SecondName.String,
		Name:       name,
		Username:   e.Username,
		Phone:      phone,
		TelegramID: e.TelegramID,
		IsApproved: e.IsApproved,
		IsWorking:  true, //e.IsWorking,
	}
}

func FromEngineerList(list []*models.Engineer) []EngineerResponse {
	resp := make([]EngineerResponse, 0, len(list))
	for _, e := range list {
		resp = append(resp, FromEngineerModel(e))
	}
	return resp
}
