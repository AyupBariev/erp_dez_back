package services

import (
	"erp/internal/app/dto"
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"fmt"
	"time"
)

type OrderService struct {
	orderRepo       *repositories.OrderRepository
	notification    *NotificationService
	engineerService *EngineerService
}

func NewOrderService(orderRepo *repositories.OrderRepository, notification *NotificationService, engineerService *EngineerService) *OrderService {
	return &OrderService{orderRepo: orderRepo, notification: notification, engineerService: engineerService}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	nextErpNumber, err := s.orderRepo.GetNextERPNumber()
	if err != nil {
		return err
	}

	order.ERPNumber = nextErpNumber
	if order.Status == "" {
		if order.EngineerID != nil {
			order.Status = "in_proccess"
		} else {
			order.Status = "new"
		}
	}

	if order.EngineerID != nil {
		go s.notification.NotifyEngineerNewOrder(order, order.Engineer)
	}

	return s.orderRepo.Create(order)
}

func (s *OrderService) Update(order *models.Order, req dto.UpdateOrderRequest) (*models.Order, error) {
	if req.AggregatorID != 0 {
		order.AggregatorID = req.AggregatorID
	}
	if req.ProblemID != 0 {
		order.ProblemID = req.ProblemID
	}
	if req.OurPercent != 0 {
		order.OurPercent = req.OurPercent
	}
	if req.ClientName != "" {
		order.ClientName = req.ClientName
	}
	if req.Address != "" {
		order.Address = req.Address
	}
	if req.WorkVolume != "" {
		order.WorkVolume = req.WorkVolume
	}
	if req.Price != "" {
		order.Price = req.Price
	}
	if req.Note != "" {
		order.Note = req.Note
	}
	if req.ScheduledAt != "" {
		scheduledAt, _ := time.ParseInLocation("2006-01-02T15:04", req.ScheduledAt, time.Local)
		order.ScheduledAt = scheduledAt.UTC()
	}
	if req.Phones != nil {
		order.Phones = req.Phones
	}
	if req.Status != "" {
		order.Status = req.Status
	}

	if err := s.orderRepo.UpdateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrders(date *string) ([]*models.Order, error) {
	if date != nil {
		_, err := time.Parse("2006-01-02", *date)
		if err != nil {
			return nil, err
		}
	}
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

func (s *OrderService) UpdateEngineerAndStatus(engineerID int64, erpNumber int64, status string) (*models.Order, error) {
	engineer, err := s.engineerService.GetEngineerByID(engineerID)
	if err != nil || engineer == nil {
		return nil, fmt.Errorf("engineer not found or not approved")
	}

	order, err := s.GetOrderForAssign(erpNumber)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	// 🚫 Проверяем, что заказ уже принят кем-то
	if order.EngineerID != nil && order.Status == status {
		return nil, fmt.Errorf("order already confirmed by engineer")
	}

	// Обновляем инженера и статус
	order.EngineerID = &engineerID
	order.Engineer = engineer
	order.Status = status

	if order.EngineerID != nil {
		go s.notification.NotifyEngineerNewOrder(order, order.Engineer)
	}

	err = s.orderRepo.UpdateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

func (s *OrderService) UnassignOrder(erpNumber int64) (*models.Order, error) {
	order, err := s.GetOrderForAssign(erpNumber)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	// Сохраняем текущего инженера для уведомления
	previousEngineer := order.Engineer

	// Обнуляем инженера и меняем статус
	order.EngineerID = nil
	order.ConfirmedAt = nil
	order.Engineer = nil
	order.Status = "new"

	err = s.orderRepo.UpdateOrder(order, "EngineerID", "ConfirmedAt", "Status")
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Отправляем уведомление предыдущему инженеру
	if previousEngineer != nil {
		go s.notification.NotifyEngineerOrderUnassigned(order, previousEngineer)
	}

	return order, nil
}

func (s *OrderService) EngineerAcceptOrderByErpNumber(order *models.Order) error {
	// Обновляем статус
	order.Status = "working"
	now := time.Now()
	order.ConfirmedAt = &now

	// Сохраняем изменения
	return s.orderRepo.UpdateOrder(order)
}

func (s *OrderService) GetOrderByID(orderID int64) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetOrderByErpNumber(erpNumber int64) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderByErpNumber(erpNumber)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) Delete(erpOrderNumber int64) error {
	return s.orderRepo.Delete(erpOrderNumber)
}
