package services

import (
	"erp/internal/app/dto"
	"erp/internal/app/models"
	"erp/internal/utils"
	"time"
)

type RepeatRequestService struct {
	requestRepo  RepeatRequestRepository
	orderService *OrderService
}

type RepeatRequestRepository interface {
	Create(req *models.RepeatRequest) error
	GetByID(id uint) (*models.RepeatRequest, error)
	Update(req *models.RepeatRequest) error

	Confirm(
		requestID uint,
		curatorID int64,
		scheduledAt time.Time,
		repeatOrderID uint64,
	) error

	List(status string) ([]models.RepeatRequest, error)
}

func NewRepeatRequestService(requestRepo RepeatRequestRepository, orderService *OrderService) *RepeatRequestService {
	return &RepeatRequestService{
		requestRepo:  requestRepo,
		orderService: orderService,
	}
}

func (s *RepeatRequestService) Create(
	orderID uint,
	engineerID uint,
	description string,
	scheduledAt time.Time,
) error {

	req := &models.RepeatRequest{
		OrderID:     orderID,
		EngineerID:  engineerID,
		Description: description,
		RequestedAt: time.Now(),
		ScheduledAt: scheduledAt,
		Status:      "pending",
		Confirmed:   false,
	}

	return s.requestRepo.Create(req)
}

func (s *RepeatRequestService) Confirm(
	requestID uint,
	managerID int64,
	payload dto.CreateOrderRequest,
) error {
	scheduledAt := time.Now()
	if payload.ScheduledAt != "" {
		var err error
		scheduledAt, err = utils.ParseScheduledAt(payload.ScheduledAt)
		if err != nil {
			return err
		}
		scheduledAt = scheduledAt.UTC()
	}
	parentOrder, err := s.orderService.GetOrderByErpNumber(*payload.RepeatErpNumber)
	if err != nil {
		return err
	}
	parentOrderID := uint64(parentOrder.ID)

	// Преобразуем CreateOrderRequest в models.Order
	newRepeatOrder := &models.Order{
		AggregatorID:      int64(payload.AggregatorID),
		ProblemID:         payload.ProblemID,
		Price:             payload.Price,
		OurPercent:        payload.OurPercent,
		ClientName:        payload.ClientName,
		Phones:            payload.Phones,
		Address:           payload.Address,
		AdminID:           managerID,
		WorkVolume:        payload.WorkVolume,
		Note:              payload.Note,
		ScheduledAt:       scheduledAt,
		EngineerID:        payload.EngineerID,
		Status:            "working",
		RepeatID:          &parentOrderID,
		RepeatDescription: payload.Note,
		RepeatedBy:        "curator",
	}

	err = s.orderService.CreateOrder(newRepeatOrder)
	if err != nil {
		return err
	}

	return s.requestRepo.Confirm(requestID, managerID, scheduledAt, uint64(newRepeatOrder.ID))
}

func (s *RepeatRequestService) List(status string) ([]models.RepeatRequest, error) {
	return s.requestRepo.List(status)
}
