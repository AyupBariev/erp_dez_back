package services

import (
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"fmt"
	"time"
)

type OrderService struct {
	orderRepo    *repositories.OrderRepository
	notification *NotificationService
}

func NewOrderService(orderRepo *repositories.OrderRepository, notification *NotificationService) *OrderService {
	return &OrderService{orderRepo: orderRepo, notification: notification}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	maxERP, err := s.orderRepo.GetMaxERPNumber()
	if err != nil {
		return err
	}

	order.ERPNumber = maxERP + 1
	order.Status = "new"

	if order.EngineerID.Valid {
		go s.notification.NotifyEngineerNewOrder(order, order.Engineer)
	}

	return s.orderRepo.Create(order)
}

func (s *OrderService) GetOrders(date *string) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetOrders(date)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (s *OrderService) GetTodayOrders(chatID int64) ([]models.Order, error) {
	return s.orderRepo.GetTodayOrders(chatID)
}

func (s *OrderService) GetRepeatOrders(chatID int64) ([]models.Order, error) {
	return s.orderRepo.GetRepeatOrders(chatID)
}

func (s *OrderService) GetCashOrders(chatID int64) ([]models.Order, error) {
	return s.orderRepo.GetCashOrders(chatID)
}

func (s *OrderService) GetOrderForAssign(ErpNumber int64) (*models.Order, error) {
	return s.orderRepo.GetOrderByErpNumber(ErpNumber)
}

func (s *OrderService) UpdateEngineerAndStatus(order *models.Order) error {
	if order.EngineerID.Valid {
		go s.notification.NotifyEngineerNewOrder(order, order.Engineer)
	}
	return s.orderRepo.Update(order)
}

func (s *OrderService) MarkAsAcceptedByErpNumber(erpNumber int64) error {
	order, err := s.orderRepo.GetOrderByErpNumber(erpNumber)
	if err != nil {
		return err
	}

	// Обновляем статус
	order.Status = "confirmed"
	order.ConfirmedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	// Сохраняем изменения
	return s.orderRepo.Update(order)
}
