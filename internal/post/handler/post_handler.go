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
		util.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	var req service.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Create post bind error:", err)
		util.BadRequestResponse(c, "Invalid request format")
		return
	}

	post, err := h.postService.CreatePost(userID, &req)
	if err != nil {
		logger.Error("Create post error:", err)
		util.BadRequestResponse(c, err.Error())
		return
	}

	logger.Info("Post created:", post.ID)
	util.CreatedResponse(c, "Post created successfully", post)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid post ID")
		return
	}

	incrementView := c.Query("view") == "true"
	post, err := h.postService.GetPostByID(uint(postID), incrementView)
	if err != nil {
		logger.Error("Get post error:", err)
		util.NotFoundResponse(c, "Post not found")
		return
	}

	util.SuccessResponse(c, "Post retrieved successfully", post)
}

func (h *PostHandler) GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.DefaultQuery("status", "published")

	resp, err := h.postService.GetPostList(page, limit, status)
	if err != nil {
		logger.Error("Get post list error:", err)
		util.InternalServerErrorResponse(c, "Failed to get post list")
		return
	}

	util.SuccessResponse(c, "Post list retrieved successfully", resp)
}

func (h *PostHandler) GetMyPosts(c *gin.Context) {
	userID, err := h.postService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.postService.GetMyPosts(userID, page, limit)
	if err != nil {
		logger.Error("Get my posts error:", err)
		util.InternalServerErrorResponse(c, "Failed to get posts")
		return
	}

	util.SuccessResponse(c, "Posts retrieved successfully", resp)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, err := h.postService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid post ID")
		return
	}

	var req service.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Update post bind error:", err)
		util.BadRequestResponse(c, "Invalid request format")
		return
	}

	post, err := h.postService.UpdatePost(userID, uint(postID), &req)
	if err != nil {
		logger.Error("Update post error:", err)
		if err.Error() == "you can only update your own posts" {
			util.UnauthorizedResponse(c, err.Error())
		} else {
			util.BadRequestResponse(c, err.Error())
		}
		return
	}

	logger.Info("Post updated:", post.ID)
	util.SuccessResponse(c, "Post updated successfully", post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, err := h.postService.GetCurrentUserID(c)
	if err != nil {
		util.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		util.BadRequestResponse(c, "Invalid post ID")
		return
	}

	err = h.postService.DeletePost(userID, uint(postID))
	if err != nil {
		logger.Error("Delete post error:", err)
		if err.Error() == "you can only delete your own posts" {
			util.UnauthorizedResponse(c, err.Error())
		} else if err.Error() == "record not found" {
			util.NotFoundResponse(c, "Post not found")
		} else {
			util.InternalServerErrorResponse(c, "Failed to delete post")
		}
		return
	}

	logger.Info("Post deleted:", postID)
	util.SuccessResponse(c, "Post deleted successfully", map[string]interface{}{})
}