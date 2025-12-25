package handlers

import (
	"erp/internal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProfitHandler struct {
	EngineerMotivationService *services.EngineerMotivationService
}

func NewProfitHandler(engineerMotivationService *services.EngineerMotivationService) *ProfitHandler {
	return &ProfitHandler{
		EngineerMotivationService: engineerMotivationService,
	}
}

func (h *ProfitHandler) GetProfit(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to required"})
		return
	}

	rows, err := h.EngineerMotivationService.GetProfitByDateRange(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rows)
}
