package handlers

import (
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/app/response"
	"erp/internal/app/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EngineerResponse struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Username   string  `json:"username"`
	Phone      *string `json:"phone,omitempty"`
	TelegramID int64   `json:"telegram_id"`
	IsApproved bool    `json:"is_approved"`
	IsWorking  bool    `json:"is_working"`
}

type EngineerHandler struct {
	engineerService *services.EngineerService
}

func NewEngineerHandler(engineerService *services.EngineerService) *EngineerHandler {
	return &EngineerHandler{
		engineerService: engineerService,
	}
}

// CreateEngineerRequest --- Запрос на создание инженера через HTTP ---
type CreateEngineerRequest struct {
	FirstName  string `json:"first_name" binding:"required"`
	SecondName string `json:"second_name"`
	Username   string `json:"username" binding:"required"`
	Phone      string `json:"phone"`
	TelegramID int64  `json:"telegram_id"` // необязательный
}

// CreateEngineer через HTTP
func (h *EngineerHandler) CreateEngineer(c *gin.Context) {
	var req CreateEngineerRequest

	// Парсим JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаём модель инженера
	engineer := &models.Engineer{
		FirstName:  sql.NullString{String: req.FirstName, Valid: req.FirstName != ""},
		SecondName: sql.NullString{String: req.SecondName, Valid: req.SecondName != ""},
		Username:   req.Username,
		Phone:      sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		TelegramID: sql.NullInt64{Int64: req.TelegramID, Valid: req.TelegramID != 0},
		IsApproved: false, // по умолчанию
	}

	// Сохраняем через сервис
	createdEngineer, err := h.engineerService.CreateEngineer(engineer)
	if err != nil {
		if errors.Is(err, services.ErrEngineerAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Инженер с таким username или Telegram ID уже существует"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать инженера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Инженер создан",
		"engineer": createdEngineer,
	})
}

func (h *EngineerHandler) ListEngineers(c *gin.Context) {
	dateParam := c.Query("date")
	if dateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}

	engineers, err := h.engineerService.ListWorkingEngineers(dateParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.FromEngineerList(engineers))
}
