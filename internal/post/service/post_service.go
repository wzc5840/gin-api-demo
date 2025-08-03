package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wzc5840/gin-api-demo/internal/post/model"
	"github.com/wzc5840/gin-api-demo/internal/post/repository"
)

type PostService struct {
	postRepo *repository.PostRepository
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Summary string `json:"summary"`
	Status  string `json:"status"`
	Tags    string `json:"tags"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Summary string `json:"summary"`
	Status  string `json:"status"`
	Tags    string `json:"tags"`
}

type PostListResponse struct {
	Posts []*model.Post `json:"posts"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) CreatePost(userID uint, req *CreatePostRequest) (*model.Post, error) {
	status := model.PostStatusDraft
	if req.Status != "" {
		switch req.Status {
		case "draft", "published", "archived":
			status = model.PostStatus(req.Status)
		default:
			return nil, errors.New("無効な投稿ステータスです")
		}
	}

	post := &model.Post{
		Title:     req.Title,
		Content:   req.Content,
		Summary:   req.Summary,
		Status:    status,
		AuthorID:  userID,
		Tags:      req.Tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if status == model.PostStatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := s.postRepo.CreatePost(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPostByID(id uint, incrementView bool) (*model.Post, error) {
	post, err := s.postRepo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	if incrementView && post.Status == model.PostStatusPublished {
		s.postRepo.IncrementViewCount(id)
		post.ViewCount++
	}

	return post, nil
}

func (s *PostService) GetPostList(page, limit int, status string) (*PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	posts, total, err := s.postRepo.GetAllPosts(limit, offset, status)
	if err != nil {
		return nil, err
	}

	return &PostListResponse{
		Posts: posts,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *PostService) GetMyPosts(userID uint, page, limit int) (*PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	posts, total, err := s.postRepo.GetPostsByAuthor(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &PostListResponse{
		Posts: posts,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *PostService) UpdatePost(userID, postID uint, req *UpdatePostRequest) (*model.Post, error) {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != userID {
		return nil, errors.New("自分の投稿のみ更新できます")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Summary != "" {
		post.Summary = req.Summary
	}
	if req.Tags != "" {
		post.Tags = req.Tags
	}

	if req.Status != "" {
		switch req.Status {
		case "draft", "published", "archived":
			oldStatus := post.Status
			post.Status = model.PostStatus(req.Status)
			
			if oldStatus != model.PostStatusPublished && post.Status == model.PostStatusPublished {
				now := time.Now()
				post.PublishedAt = &now
			}
		default:
			return nil, errors.New("無効な投稿ステータスです")
		}
	}

	post.UpdatedAt = time.Now()

	if err := s.postRepo.UpdatePost(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(userID, postID uint) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return err
	}

	if post.AuthorID != userID {
		return errors.New("自分の投稿のみ削除できます")
	}

	return s.postRepo.DeletePost(postID)
}

func (s *PostService) GetCurrentUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("ユーザー認証が必要です")
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, errors.New("無効なユーザーIDです")
	}

	return id, nil
}