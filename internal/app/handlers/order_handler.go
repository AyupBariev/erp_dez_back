package handlers

import (
	"erp/internal/app/dto"
	"erp/internal/app/models"
	"erp/internal/app/response"
	"erp/internal/app/services"
	"erp/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Добавь ProblemID
	order := &models.Order{
		AggregatorID: int64(req.AggregatorID),
		ProblemID:    req.ProblemID,
		Price:        req.Price,
		OurPercent:   req.OurPercent,
		ClientName:   req.ClientName,
		AdminID:      userID.(int64),
		Phones:       req.Phones,
		Address:      req.Address,
		WorkVolume:   req.WorkVolume,
		Note:         req.Note,
		Status:       req.Status,
	}

	if req.ScheduledAt != "" {
		scheduledAt, err := utils.ParseScheduledAt(req.ScheduledAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid scheduled_at format, use YYYY-MM-DDTHH:MM or YYYY-MM-DD HH:MM",
			})
			return
		}
		order.ScheduledAt = scheduledAt.UTC()
	}

	if req.RepeatID != nil {
		order.RepeatID = req.RepeatID
		order.RepeatedBy = "curator"
	} else {
		var parentOrderID *uint64 = nil
		if req.RepeatErpNumber != nil {
			parentOrder, err := h.OrderService.GetOrderByErpNumber(*req.RepeatErpNumber)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Заказ не существует"})
				return
			}
			id64 := uint64(parentOrder.ID)
			parentOrderID = &id64
		}
		order.RepeatID = parentOrderID
	}

	// Обработка engineer_id и статуса
	if req.EngineerID != nil {
		order.EngineerID = req.EngineerID
	} else {
		order.EngineerID = nil
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
		ErpNumber  int64 `json:"erp_number"`
		EngineerID int64 `json:"engineer_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем сервис с правильными параметрами
	order, err := h.OrderService.UpdateEngineerAndStatus(input.EngineerID, input.ErpNumber, "in_proccess")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := response.FromOrderModel(order)
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) UnAssignOrderHandler(c *gin.Context) {
	var input struct {
		ErpNumber int64 `json:"erp_number"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.OrderService.UnassignOrder(input.ErpNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := response.FromOrderModel(order)
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	erpNumber, err := strconv.Atoi(c.Param("erp_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order number"})
		return
	}
	var req dto.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.OrderService.GetOrderByErpNumber(int64(erpNumber))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("cant load order for payout: %w", err)})
		return
	}
	item, err := h.OrderService.Update(order, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	erpOrderNumber, err := strconv.Atoi(c.Param("erp_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order number"})
		return
	}

	err = h.OrderService.Delete(int64(erpOrderNumber))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
