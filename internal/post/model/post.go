package model

import (
	"time"

	"gorm.io/gorm"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type Post struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Title       string         `json:"title" gorm:"not null;size:255"`
	Content     string         `json:"content" gorm:"type:text"`
	Summary     string         `json:"summary" gorm:"size:500"`
	Status      PostStatus     `json:"status" gorm:"default:'draft'"`
	AuthorID    uint           `json:"author_id" gorm:"not null;index"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	Tags        string         `json:"tags" gorm:"size:500"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Post) TableName() string {
	return "posts"
}