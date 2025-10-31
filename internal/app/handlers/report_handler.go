package handlers

import (
	"erp/internal/app/services"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	ReportService   *services.ReportService
	OrderService    *services.OrderService
	EngineerService *services.EngineerService
}
type ReportFormResponse struct {
	OrderID        string `json:"order_id"`
	ClientName     string `json:"client_name"`
	Address        string `json:"address"`
	Problem        string `json:"problem"`
	AggregatorName string `json:"aggregator_name"`
	EngineerName   string `json:"engineer_name"`
}
type SubmitReportRequest struct {
	Token       string  `json:"token" binding:"required"`
	HasRepeat   bool    `json:"has_repeat"`
	RepeatDate  *string `json:"repeat_date,omitempty"` // формат "2006-01-02T15:04"
	RepeatNote  string  `json:"repeat_note,omitempty"`
	Description string  `json:"description,omitempty"`
}

func NewReportHandler(reportService *services.ReportService, orderService *services.OrderService, engineerService *services.EngineerService) *ReportHandler {
	return &ReportHandler{reportService, orderService, engineerService}
}

// Проверка токена (валидация)
func (h *ReportHandler) GetReportByToken(c *gin.Context) {
	token := c.Param("token")
	link, err := h.ReportService.GetByToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная ссылка"})
		return
	}

	if link.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Срок действия ссылки истёк"})
		return
	}

	order, err := h.OrderService.GetOrderByErpNumber(link.OrderID)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"erp_number":  order.ERPNumber,
		"client_name": order.ClientName,
		"address":     order.Address,
		"engineer_id": link.EngineerID,
	}})
}

func (h *ReportHandler) GetReportForm(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}

	link, err := h.ReportService.GetByToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid token"})
		return
	}

	if link.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token expired"})
		return
	}

	order, err := h.OrderService.GetOrderByErpNumber(link.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	engineer, err := h.EngineerService.GetEngineerByID(link.EngineerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "engineer not found"})
		return
	}
	engineerName := fmt.Sprintf("%s %s", engineer.FirstName, engineer.SecondName)
	resp := ReportFormResponse{
		OrderID:        strconv.FormatInt(order.ERPNumber, 10),
		ClientName:     order.ClientName,
		Address:        order.Address,
		Problem:        order.Problem.Name,
		AggregatorName: order.Aggregator.Name,
		EngineerName:   engineerName,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ReportHandler) SubmitReport(c *gin.Context) {
	var req SubmitReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ReportService.SubmitReport(services.SubmitReportRequest(req)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Отчет успешно отправлен"})
}
