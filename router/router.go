package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wzc5840/gin-api-demo/internal/auth/handler"
	"github.com/wzc5840/gin-api-demo/internal/auth/repository"
	"github.com/wzc5840/gin-api-demo/internal/auth/service"
	"github.com/wzc5840/gin-api-demo/pkg/middleware"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r.GET("/hello", func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>欢迎页面</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            font-family: 'Arial', sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            overflow: hidden;
        }
        
        .welcome-container {
            text-align: center;
            color: white;
            position: relative;
        }
        
        .welcome-title {
            font-size: 4rem;
            font-weight: bold;
            margin-bottom: 2rem;
            opacity: 0;
            animation: fadeInUp 2s ease-out forwards;
        }
        
        .welcome-subtitle {
            font-size: 1.5rem;
            opacity: 0;
            animation: fadeInUp 2s ease-out 0.5s forwards;
        }
        
        .particles {
            position: absolute;
            width: 100%;
            height: 100%;
            pointer-events: none;
        }
        
        .particle {
            position: absolute;
            width: 4px;
            height: 4px;
            background: rgba(255, 255, 255, 0.8);
            border-radius: 50%;
            animation: float 6s infinite ease-in-out;
        }
        
        @keyframes fadeInUp {
            from {
                opacity: 0;
                transform: translateY(30px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
        
        @keyframes float {
            0%, 100% {
                transform: translateY(0px) rotate(0deg);
                opacity: 0.8;
            }
            50% {
                transform: translateY(-20px) rotate(180deg);
                opacity: 1;
            }
        }
        
        .particle:nth-child(1) { left: 10%; animation-delay: 0s; }
        .particle:nth-child(2) { left: 20%; animation-delay: 0.5s; }
        .particle:nth-child(3) { left: 30%; animation-delay: 1s; }
        .particle:nth-child(4) { left: 40%; animation-delay: 1.5s; }
        .particle:nth-child(5) { left: 50%; animation-delay: 2s; }
        .particle:nth-child(6) { left: 60%; animation-delay: 2.5s; }
        .particle:nth-child(7) { left: 70%; animation-delay: 3s; }
        .particle:nth-child(8) { left: 80%; animation-delay: 3.5s; }
        .particle:nth-child(9) { left: 90%; animation-delay: 4s; }
        .particle:nth-child(10) { left: 15%; animation-delay: 4.5s; }
    </style>
</head>
<body>
    <div class="particles">
        <div class="particle"></div>
        <div class="particle"></div>
        <div class="particle"></div>
        <div class="particle"></div>
        <div class="particle"></div>
        <div class="particle"></div>
        <div class="particle"></div>
        <div class="particle"></div>
        <div class="particle"></div>
        <div class="particle"></div>
    </div>
    
    <div class="welcome-container">
        <h1 class="welcome-title">欢迎来到我的网站</h1>
        <p class="welcome-subtitle">Hello, World! 🌟</p>
    </div>
</body>
</html>
        `
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		protected := api.Group("/user")
		protected.Use(middleware.AuthMiddleware(userRepo))
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.GET("/list", authHandler.GetUserList)
			protected.GET("/:id", authHandler.GetUserDetail)
			protected.PUT("/:id", authHandler.UpdateUser)
			protected.DELETE("/:id", authHandler.DeleteUser)
		}
	}

	return r
}