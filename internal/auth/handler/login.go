package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wzc5840/gin-api-demo-01/internal/auth/service"
	"github.com/wzc5840/gin-api-demo-01/pkg/logger"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Login bind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		logger.Error("Login error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Login failed",
			"message": err.Error(),
		})
		return
	}

	logger.Info("User logged in:", resp.User.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    resp,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Register bind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		logger.Error("Register error:", err)
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Registration failed",
			"message": err.Error(),
		})
		return
	}

	logger.Info("User registered:", resp.User.Username)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"data":    resp,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		logger.Error("Get profile error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}