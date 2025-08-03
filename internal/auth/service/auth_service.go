package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wzc5840/gin-api-demo-01/internal/auth/model"
	"github.com/wzc5840/gin-api-demo-01/internal/auth/repository"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	if !s.verifyPassword(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	token := s.generateToken(user)

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	existingUser, err := s.userRepo.GetUserByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, err = s.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword := s.hashPassword(req.Password)

	user := &model.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	token := s.generateToken(user)

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func (s *AuthService) verifyPassword(password, hashedPassword string) bool {
	return s.hashPassword(password) == hashedPassword
}

func (s *AuthService) generateToken(user *model.User) string {
	data := fmt.Sprintf("%s:%d:%d", user.Username, user.ID, time.Now().Unix())
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func (s *AuthService) GetCurrentUser(c *gin.Context) (*model.User, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, errors.New("user not authenticated")
	}

	id, ok := userID.(uint)
	if !ok {
		return nil, errors.New("invalid user ID")
	}

	return s.userRepo.GetUserByID(id)
}