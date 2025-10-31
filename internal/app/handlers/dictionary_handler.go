package handlers

import (
	"erp/internal/app/models"
	"erp/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DictionaryHandler struct {
	dictService *services.DictionaryService
}

func NewDictionaryHandler(dictService *services.DictionaryService) *DictionaryHandler {
	return &DictionaryHandler{
		dictService: dictService,
	}
}

// HandleDictionary Универсальный обработчик для всех словарей
func (h *DictionaryHandler) HandleDictionary(tableName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case "GET":
			if c.Param("id") != "" {
				h.getDictionaryItem(tableName, c)
			} else {
				h.getDictionaryList(tableName, c)
			}
		case "POST":
			h.createDictionaryItem(tableName, c)
		case "PUT":
			h.updateDictionaryItem(tableName, c)
		case "DELETE":
			h.deleteDictionaryItem(tableName, c)
		}
	}
}

// GET /dictionaries/:type
func (h *DictionaryHandler) getDictionaryList(tableName string, c *gin.Context) {
	items, err := h.dictService.GetAll(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GET /dictionaries/:type/:id
func (h *DictionaryHandler) getDictionaryItem(tableName string, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	item, err := h.dictService.GetByID(tableName, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// POST /dictionaries/:type
func (h *DictionaryHandler) createDictionaryItem(tableName string, c *gin.Context) {
	var req models.CreateDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.dictService.Create(tableName, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

// PUT /dictionaries/:type/:id
func (h *DictionaryHandler) updateDictionaryItem(tableName string, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.UpdateDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.dictService.Update(tableName, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// DELETE /dictionaries/:type/:id
func (h *DictionaryHandler) deleteDictionaryItem(tableName string, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.dictService.Delete(tableName, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
