package handlers

import (
	"erp/internal/app/services"
	"erp/internal/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminHandler struct {
	engineerService *services.EngineerService
	telegramHandler *TelegramHandler
}

func NewAdminHandler(engineerService *services.EngineerService, telegramHandler *TelegramHandler) *AdminHandler {
	return &AdminHandler{engineerService: engineerService, telegramHandler: telegramHandler}
}

func (h *AdminHandler) ApproveEngineer(c *gin.Context) {
	var input struct {
		EngineerID int64 `json:"engineer_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	engineer, err := h.engineerService.ApproveEngineer(input.EngineerID)
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to approve engineer %d", input.EngineerID), err)
		c.JSON(500, gin.H{"error": "Не удалось подтвердить инженера"})
		return
	}

	if engineer.TelegramID != 0 {
		h.telegramHandler.sendMessage(engineer.TelegramID, "✅ Бот активирован\nДля перехода в главное меню нажмите кнопку ниже", "init")
	}

	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}
