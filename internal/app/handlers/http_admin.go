package handlers

import (
	"erp/internal/app/services"
	"erp/internal/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	if telegramID := engineer.GetTelegramID(); telegramID != nil {
		h.telegramHandler.sendMessage(*telegramID, "✅ Бот активирован\nДля перехода в главное меню нажмите кнопку ниже", "init")
	}

	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

func (h *AdminHandler) ToggleWorkingStatusEngineer(c *gin.Context) {
	engineerIDStr := c.Param("engineer_id")
	engineerID, err := strconv.ParseInt(engineerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid engineer ID"})
		return
	}

	var input struct {
		IsWorking bool `json:"is_working"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-94827366643": err.Error()})
		return
	}

	engineer, err := h.engineerService.UpdateEngineerWorkingStatus(engineerID, input.IsWorking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, engineer)
}
