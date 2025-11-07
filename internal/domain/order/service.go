package order

import (
	"errors"
	"time"
)

type Notifier interface {
	NotifyEngineerNewOrder(order *Order)
}

type Service struct {
	repo     Repository
	notifier Notifier
}

func NewService(repo Repository, notifier Notifier) *Service {
	return &Service{repo: repo, notifier: notifier}
}

func (s *Service) Create(order *Order) error {
	nextErp, err := s.repo.GetNextERPNumber()
	if err != nil {
		return err
	}
	order.ERPNumber = nextErp
	order.Status = "new"

	if order.EngineerID != nil {
		go s.notifier.NotifyEngineerNewOrder(order)
	}

	return s.repo.Create(order)
}

func (s *Service) Update(erpNumber int64, updateFn func(*Order) error) (*Order, error) {
	order, err := s.repo.GetByERPNumber(erpNumber)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if err := updateFn(order); err != nil {
		return nil, err
	}

	err = s.repo.Update(order)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByERPNumber(erpNumber)
}

func (s *Service) AcceptByEngineer(erpNumber int64) error {
	order, err := s.repo.GetByERPNumber(erpNumber)
	if err != nil {
		return err
	}

	now := time.Now()
	order.Status = "working"
	order.ConfirmedAt = &now

	return s.repo.Update(order)
}
