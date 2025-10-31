package handlers

import (
	"erp/internal/app/models"
	"erp/internal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type EngineerTargetHandler struct {
	service *services.EngineerTargetService
}

func NewEngineerTargetHandler(service *services.EngineerTargetService) *EngineerTargetHandler {
	return &EngineerTargetHandler{service: service}
}

func (h *EngineerTargetHandler) GetTargets(c *gin.Context) {
	engineerID, _ := strconv.Atoi(c.Param("engineer_id"))
	monthStart := c.Query("month_start")
	monthEnd := c.Query("month_end")
	targets, err := h.service.GetTargets(uint(engineerID), monthStart, monthEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, targets)
}

func (h *EngineerTargetHandler) CreateTarget(c *gin.Context) {
	var target models.EngineerMotivationTarget
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateTarget(&target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, target)
}

func (h *EngineerTargetHandler) UpdateTarget(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var target models.EngineerMotivationTarget
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	target.ID = uint(id)
	if err := h.service.UpdateTarget(&target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, target)
}
