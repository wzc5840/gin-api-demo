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
		util.BadRequestResponse(c, "無効なリクエスト形式です")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		logger.Error("Login error:", err)
		util.UnauthorizedResponse(c, "ログインに失敗しました")
		return
	}

	logger.Info("User logged in:", resp.User.Username)
	util.SuccessResponse(c, "ログインしました", resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Register bind error:", err)
		util.BadRequestResponse(c, "無効なリクエスト形式です")
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		logger.Error("Register error:", err)
		util.ConflictResponse(c, "登録に失敗しました")
		return
	}

	logger.Info("User registered:", resp.User.Username)
	util.CreatedResponse(c, "登録しました", resp)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		logger.Error("Get profile error:", err)
		util.UnauthorizedResponse(c, "認証が必要です")
		return
	}

	util.SuccessResponse(c, "プロフィールを取得しました", user)
}

func (h *AuthHandler) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.authService.GetUserList(page, limit)
	if err != nil {
		logger.Error("Get user list error:", err)
		util.InternalServerErrorResponse(c, "ユーザーリストの取得に失敗しました")
		return
	}

	util.SuccessResponse(c, "ユーザーリストを取得しました", resp)
}

func (h *AuthHandler) GetUserDetail(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "無効なユーザーIDです")
		return
	}

	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		logger.Error("Get user detail error:", err)
		util.UnauthorizedResponse(c, "認証が必要です")
		return
	}

	if user.ID != uint(userID) {
		targetUser, err := h.authService.GetUserByID(uint(userID))
		if err != nil {
			util.NotFoundResponse(c, "ユーザーが見つかりません")
			return
		}
		user = targetUser
	}

	util.SuccessResponse(c, "ユーザー詳細を取得しました", user)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "無効なユーザーIDです")
		return
	}

	currentUserID, err := h.authService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "認証が必要です")
		return
	}

	if currentUserID != uint(userID) {
		util.UnauthorizedResponse(c, "自分のプロフィールのみ更新できます")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Update user bind error:", err)
		util.BadRequestResponse(c, "無効なリクエスト形式です")
		return
	}

	user, err := h.authService.UpdateUserProfile(uint(userID), &req)
	if err != nil {
		logger.Error("Update user error:", err)
		util.ConflictResponse(c, err.Error())
		return
	}

	logger.Info("User updated:", user.Username)
	util.SuccessResponse(c, "ユーザー情報を更新しました", user)
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "無効なユーザーIDです")
		return
	}

	currentUserID, err := h.authService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "認証が必要です")
		return
	}

	err = h.authService.DeleteUser(currentUserID, uint(targetUserID))
	if err != nil {
		logger.Error("Delete user error:", err)
		if err.Error() == "自分のアカウントは削除できません" {
			util.BadRequestResponse(c, "自分のアカウントは削除できません")
		} else if err.Error() == "ユーザーが見つかりません" {
			util.NotFoundResponse(c, "ユーザーが見つかりません")
		} else {
			util.InternalServerErrorResponse(c, "ユーザーの削除に失敗しました")
		}
		return
	}

	logger.Info("User deleted:", targetUserID)
	util.SuccessResponse(c, "ユーザーを削除しました", map[string]interface{}{})
}