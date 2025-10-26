package handlers

import (
	"erp/internal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var creds struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	if err := c.ShouldBind(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.AuthService.Authenticate(creds.Username, creds.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Генерация токена
	token, err := h.AuthService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Возвращаем токен и минимальные данные
	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"expires_in":   86400,
	})
}

func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	tokenString, exists := c.Get("tokenString")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token not found"})
		return
	}

	if err := h.AuthService.InvalidateToken(tokenString.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to invalidate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
