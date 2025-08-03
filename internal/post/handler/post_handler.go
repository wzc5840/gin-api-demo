package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wzc5840/gin-api-demo/internal/post/service"
	"github.com/wzc5840/gin-api-demo/pkg/logger"
	"github.com/wzc5840/gin-api-demo/pkg/util"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, err := h.postService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "認証が必要です")
		return
	}

	var req service.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Create post bind error:", err)
		util.BadRequestResponse(c, "無効なリクエスト形式です")
		return
	}

	post, err := h.postService.CreatePost(userID, &req)
	if err != nil {
		logger.Error("Create post error:", err)
		util.BadRequestResponse(c, err.Error())
		return
	}

	logger.Info("Post created:", post.ID)
	util.CreatedResponse(c, "投稿を作成しました", post)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "無効な投稿IDです")
		return
	}

	incrementView := c.Query("view") == "true"
	post, err := h.postService.GetPostByID(uint(postID), incrementView)
	if err != nil {
		logger.Error("Get post error:", err)
		util.NotFoundResponse(c, "投稿が見つかりません")
		return
	}

	util.SuccessResponse(c, "投稿を取得しました", post)
}

func (h *PostHandler) GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.DefaultQuery("status", "published")

	resp, err := h.postService.GetPostList(page, limit, status)
	if err != nil {
		logger.Error("Get post list error:", err)
		util.InternalServerErrorResponse(c, "投稿リストの取得に失敗しました")
		return
	}

	util.SuccessResponse(c, "投稿リストを取得しました", resp)
}

func (h *PostHandler) GetMyPosts(c *gin.Context) {
	userID, err := h.postService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "認証が必要です")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.postService.GetMyPosts(userID, page, limit)
	if err != nil {
		logger.Error("Get my posts error:", err)
		util.InternalServerErrorResponse(c, "投稿の取得に失敗しました")
		return
	}

	util.SuccessResponse(c, "投稿を取得しました", resp)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, err := h.postService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "認証が必要です")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "無効な投稿IDです")
		return
	}

	var req service.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Update post bind error:", err)
		util.BadRequestResponse(c, "無効なリクエスト形式です")
		return
	}

	post, err := h.postService.UpdatePost(userID, uint(postID), &req)
	if err != nil {
		logger.Error("Update post error:", err)
		if err.Error() == "自分の投稿のみ更新できます" {
			util.UnauthorizedResponse(c, err.Error())
		} else {
			util.BadRequestResponse(c, err.Error())
		}
		return
	}

	logger.Info("Post updated:", post.ID)
	util.SuccessResponse(c, "投稿を更新しました", post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, err := h.postService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "認証が必要です")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "無効な投稿IDです")
		return
	}

	err = h.postService.DeletePost(userID, uint(postID))
	if err != nil {
		logger.Error("Delete post error:", err)
		if err.Error() == "自分の投稿のみ削除できます" {
			util.UnauthorizedResponse(c, err.Error())
		} else if err.Error() == "record not found" {
			util.NotFoundResponse(c, "投稿が見つかりません")
		} else {
			util.InternalServerErrorResponse(c, "投稿の削除に失敗しました")
		}
		return
	}

	logger.Info("Post deleted:", postID)
	util.SuccessResponse(c, "投稿を削除しました", map[string]interface{}{})
}