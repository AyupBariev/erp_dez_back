package services

import (
	"erp/internal/app/repositories"
	"erp/internal/pkg/telphin"
)

type CallService struct {
	Telphin  *telphin.TelphinClient
	CallRepo *repositories.CallRepository // работа с таблицей calls
}

func NewCallService(t *telphin.TelphinClient, repo *repositories.CallRepository) *CallService {
	return &CallService{Telphin: t, CallRepo: repo}
}
