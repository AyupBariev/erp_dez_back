package handlers

import (
	"erp/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EngineerMotivationHandler struct {
	service *services.EngineerMotivationService
}

func NewEngineerMotivationHandler(service *services.EngineerMotivationService) *EngineerMotivationHandler {
	return &EngineerMotivationHandler{service: service}
}

func (h *EngineerMotivationHandler) GetMonthlyMotivation(c *gin.Context) {
	month := c.Query("month") // yyyy-mm
	result, err := h.service.GetMonthlyMotivation(month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
