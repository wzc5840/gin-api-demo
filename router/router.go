package router

import (
	"github.com/gin-gonic/gin"
	authHandler "github.com/wzc5840/gin-api-demo/internal/auth/handler"
	"github.com/wzc5840/gin-api-demo/internal/auth/repository"
	authService "github.com/wzc5840/gin-api-demo/internal/auth/service"
	postHandler "github.com/wzc5840/gin-api-demo/internal/post/handler"
	postRepository "github.com/wzc5840/gin-api-demo/internal/post/repository"
	postService "github.com/wzc5840/gin-api-demo/internal/post/service"
	"github.com/wzc5840/gin-api-demo/pkg/middleware"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	authServiceInstance := authService.NewAuthService(userRepo)
	authHandlerInstance := authHandler.NewAuthHandler(authServiceInstance)

	postRepo := postRepository.NewPostRepository(db)
	postServiceInstance := postService.NewPostService(postRepo)
	postHandlerInstance := postHandler.NewPostHandler(postServiceInstance)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandlerInstance.Login)
			auth.POST("/register", authHandlerInstance.Register)
		}

		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware(userRepo))
		{
			user.GET("/profile", authHandlerInstance.GetProfile)
			user.GET("/list", authHandlerInstance.GetUserList)
			user.GET("/:id", authHandlerInstance.GetUserDetail)
			user.PUT("/:id", authHandlerInstance.UpdateUser)
			user.DELETE("/:id", authHandlerInstance.DeleteUser)
		}

		posts := api.Group("/posts")
		{
			posts.GET("", postHandlerInstance.GetPostList)
			posts.GET("/:id", postHandlerInstance.GetPost)
		}

		protectedPosts := api.Group("/posts")
		protectedPosts.Use(middleware.AuthMiddleware(userRepo))
		{
			protectedPosts.POST("", postHandlerInstance.CreatePost)
			protectedPosts.PUT("/:id", postHandlerInstance.UpdatePost)
			protectedPosts.DELETE("/:id", postHandlerInstance.DeletePost)
			protectedPosts.GET("/my", postHandlerInstance.GetMyPosts)
		}
	}

	return r
}