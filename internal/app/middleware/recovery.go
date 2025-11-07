package middleware

import (
	"erp/internal/pkg/logger"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := string(debug.Stack())
				logger.LogError("[PANIC]", fmt.Errorf("%v\n%s", rec, stack))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		}()

		c.Next()

		// Логируем все gin.Errors
		for _, e := range c.Errors {
			logger.LogError("[GIN ERROR]", e.Err)
		}
	}
}
