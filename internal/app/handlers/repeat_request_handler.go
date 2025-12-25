package handlers

import (
	"erp/internal/app/dto"
	"erp/internal/app/response"
	"erp/internal/app/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RepeatRequestHandler struct {
	service *services.RepeatRequestService
}

func NewRepeatRequestHandler(repeatRequestService *services.RepeatRequestService) *RepeatRequestHandler {
	return &RepeatRequestHandler{repeatRequestService}
}

func (h *RepeatRequestHandler) List(c *gin.Context) {
	status := c.Query("status")

	reqs, err := h.service.List(status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	responseList := response.FromRepeatRequestList(reqs)
	c.JSON(http.StatusOK, responseList)
}

func (h *RepeatRequestHandler) Confirm(c *gin.Context) {
	// Получаем ID повторного запроса из URL
	idParam := c.Param("id") // например, /repeat_requests/:id/confirm
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing request id"})
		return
	}

	requestID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	// Получаем managerID из контекста (например, через middleware авторизации)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Десериализуем payload из JSON
	var payload dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем сервис
	if err := h.service.Confirm(uint(requestID), userID.(int64), payload); err != nil {
		// Логируем ошибку
		fmt.Printf("RepeatRequest Confirm error: %+v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
