package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wzc5840/gin-api-demo/internal/auth/service"
	"github.com/wzc5840/gin-api-demo/pkg/logger"
	"github.com/wzc5840/gin-api-demo/pkg/util"
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
		util.BadRequestResponse(c, "Invalid request format")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		logger.Error("Login error:", err)
		util.UnauthorizedResponse(c, "Login failed")
		return
	}

	logger.Info("User logged in:", resp.User.Username)
	util.SuccessResponse(c, "Login successful", resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Register bind error:", err)
		util.BadRequestResponse(c, "Invalid request format")
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		logger.Error("Register error:", err)
		util.ConflictResponse(c, "Registration failed")
		return
	}

	logger.Info("User registered:", resp.User.Username)
	util.CreatedResponse(c, "Registration successful", resp)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		logger.Error("Get profile error:", err)
		util.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	util.SuccessResponse(c, "Profile retrieved successfully", user)
}

func (h *AuthHandler) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.authService.GetUserList(page, limit)
	if err != nil {
		logger.Error("Get user list error:", err)
		util.InternalServerErrorResponse(c, "Failed to get user list")
		return
	}

	util.SuccessResponse(c, "User list retrieved successfully", resp)
}

func (h *AuthHandler) GetUserDetail(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid user ID")
		return
	}

	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		logger.Error("Get user detail error:", err)
		util.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	if user.ID != uint(userID) {
		targetUser, err := h.authService.GetUserByID(uint(userID))
		if err != nil {
			util.NotFoundResponse(c, "User not found")
			return
		}
		user = targetUser
	}

	util.SuccessResponse(c, "User detail retrieved successfully", user)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid user ID")
		return
	}

	currentUserID, err := h.authService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	if currentUserID != uint(userID) {
		util.UnauthorizedResponse(c, "You can only update your own profile")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Update user bind error:", err)
		util.BadRequestResponse(c, "Invalid request format")
		return
	}

	user, err := h.authService.UpdateUserProfile(uint(userID), &req)
	if err != nil {
		logger.Error("Update user error:", err)
		util.ConflictResponse(c, err.Error())
		return
	}

	logger.Info("User updated:", user.Username)
	util.SuccessResponse(c, "User updated successfully", user)
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid user ID")
		return
	}

	currentUserID, err := h.authService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	err = h.authService.DeleteUser(currentUserID, uint(targetUserID))
	if err != nil {
		logger.Error("Delete user error:", err)
		if err.Error() == "cannot delete your own account" {
			util.BadRequestResponse(c, "Cannot delete your own account")
		} else if err.Error() == "user not found" {
			util.NotFoundResponse(c, "User not found")
		} else {
			util.InternalServerErrorResponse(c, "Failed to delete user")
		}
		return
	}

	logger.Info("User deleted:", targetUserID)
	util.SuccessResponse(c, "User deleted successfully", map[string]interface{}{})
}