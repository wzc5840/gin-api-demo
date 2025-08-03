package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wzc5840/gin-api-demo/internal/auth/repository"
	"github.com/wzc5840/gin-api-demo/pkg/logger"
)

func AuthMiddleware(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error("No authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "No authorization header provided",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Error("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		if token == "" {
			logger.Error("Empty token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Empty token",
			})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("user_id", uint(1))
		c.Next()
	}
}