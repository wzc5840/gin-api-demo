package repository

import (
	"github.com/wzc5840/gin-api-demo/internal/post/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	db.AutoMigrate(&model.Post{})
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) GetPostByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) GetAllPosts(limit, offset int, status string) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	query := r.db.Model(&model.Post{})
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, total, err
}

func (r *PostRepository) GetPostsByAuthor(authorID uint, limit, offset int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("author_id = ?", authorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, total, err
}

func (r *PostRepository) UpdatePost(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) DeletePost(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}

func (r *PostRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}