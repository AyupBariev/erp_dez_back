package handlers

import (
	"erp/internal/app/services"
	"net/http"
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
		"price":       order.Price,
		"engineer_id": link.EngineerID,
	}})
}

func (h *ReportHandler) SubmitReport(c *gin.Context) {
	var req services.SubmitReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ReportService.SubmitReport(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Отчет успешно отправлен"})
}
