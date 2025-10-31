package handlers

import (
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/app/response"
	"erp/internal/app/services"
	"erp/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type OrderHandler struct {
	OrderService    *services.OrderService
	EngineerService *services.EngineerService
}

type OrderAssignResponse struct {
	ERPNumber   int64  `json:"erp_number"`
	ClientName  string `json:"client_name"`
	Address     string `json:"address"`
	Note        string `json:"note"`
	ScheduledAt string `json:"scheduled_at"`
	Status      string `json:"status"`
	Engineer    string `json:"engineer"`
}

func NewOrderHandler(orderService *services.OrderService, engineerService *services.EngineerService) *OrderHandler {
	return &OrderHandler{
		OrderService:    orderService,
		EngineerService: engineerService,
	}
}

// CreateOrderRequest --- Запрос на создание инженера через HTTP ---
type CreateOrderRequest struct {
	AggregatorID int      `json:"aggregator_id"`
	ProblemID    int      `json:"problem_id" binding:"required"`
	Price        string   `json:"price"`
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

func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Добавь ProblemID
	order := &models.Order{
		AggregatorID: int64(req.AggregatorID),
		ProblemID:    utils.Int64ToNullInt64(int64(req.ProblemID)),
		Price:        req.Price,
		OurPercent:   req.OurPercent,
		ClientName:   req.ClientName,
		AdminID:      userID.(int64),
		Phones:       req.Phones,
		Address:      req.Address,
		WorkVolume:   req.WorkVolume,
		Note:         req.Note,
	}

	// Обработка scheduled_at
	if req.ScheduledAt != "" {
		scheduledAt, err := time.Parse("2006-01-02T15:04", req.ScheduledAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled_at format, use YYYY-MM-DDTHH:MM"})
			return
		}
		order.ScheduledAt = scheduledAt
	}

	// Обработка engineer_id и статуса
	if req.EngineerID != nil {
		order.EngineerID = sql.NullInt64{Int64: int64(*req.EngineerID), Valid: true}
		order.Status = "in_proccess"
	} else {
		order.EngineerID = sql.NullInt64{Valid: false}
		order.Status = "new"
	}

	if err := h.OrderService.CreateOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	dateParam := c.Query("date") // может быть пустым
	var date *string
	if dateParam != "" {
		date = &dateParam
	}

	orders, err := h.OrderService.GetOrders(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := response.FromOrderList(orders)

	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) AssignOrderHandler(c *gin.Context) {
	var input struct {
		ErpNumber  int64 `json:"order_number"`
		EngineerID int64 `json:"engineer_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	engineer, err := h.EngineerService.GetEngineerByID(input.EngineerID)
	if err != nil || engineer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "engineer not found or not approved"})
		return
	}

	order, err := h.OrderService.GetOrderForAssign(input.ErpNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// 🚫 Проверяем, что заказ уже принят кем-то
	if order.EngineerID.Valid && order.Status == "working" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order already confirmed by engineer"})
		return
	}

	// Обновляем инженера и статус
	order.EngineerID = sql.NullInt64{Int64: int64(engineer.ID), Valid: true}
	order.Engineer = engineer
	order.Status = "in_proccess"

	if err := h.OrderService.UpdateEngineerAndStatus(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := response.FromOrderModel(order)

	c.JSON(http.StatusOK, resp)
}
