package handlers

import (
	"erp/internal/app/dto"
	"erp/internal/app/response"
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

	order, err := h.OrderService.GetOrderByID(link.OrderID)
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

func (h *ReportHandler) ListCashReports(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	reps, err := h.ReportService.GetCashReports(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.FromCashReportList(reps))
}

func (h *ReportHandler) ReceiveCash(c *gin.Context) {
	type req struct {
		GaveCash    float64 `json:"gave_cash" binding:"required"`
		IssuedMoney float64 `json:"issued_money"`
		Comment     string  `json:"comment"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	erpNumber, err := strconv.ParseInt(c.Param("erp_number"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}
	order, err := h.OrderService.GetOrderByErpNumber(erpNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("cant load order for payout: %w", err)})
		return
	}

	if err := h.ReportService.ReceiveCash(int64(order.ID), r.GaveCash, r.IssuedMoney, r.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. если инженеру выдали деньги – фиксируем выплату
	if r.IssuedMoney > 0 {

		if order.EngineerID == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order has no engineer"})
			return
		}
		monthStr := order.ScheduledAt.Format("2006-01")
		if err := h.EngineerService.PayAdvance(*order.EngineerID, monthStr, r.IssuedMoney); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("advance payout failed: %w", err)})
			return
		}
	}

	if _, err := h.OrderService.Update(order, dto.UpdateOrderRequest{
		Status: "closed_finally",
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
