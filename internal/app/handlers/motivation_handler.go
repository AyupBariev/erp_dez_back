package handlers

import (
	"erp/internal/app/models"
	"erp/internal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MotivationHandler struct {
	service *services.MotivationService
}

func NewMotivationHandler(service *services.MotivationService) *MotivationHandler {
	return &MotivationHandler{service: service}
}

func (h *MotivationHandler) ListSteps(c *gin.Context) {
	steps, err := h.service.ListSteps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, steps)
}

func (h *MotivationHandler) CreateStep(c *gin.Context) {
	var step models.MotivationStep
	if err := c.ShouldBindJSON(&step); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateStep(&step); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, step)
}

func (h *MotivationHandler) UpdateStep(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var step models.MotivationStep
	if err := c.ShouldBindJSON(&step); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	step.ID = uint(id)
	if err := h.service.UpdateStep(&step); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, step)
}

func (h *MotivationHandler) DeleteStep(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteStep(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
