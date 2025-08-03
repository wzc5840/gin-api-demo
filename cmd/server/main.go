package main

import (
	"log"
	"os"

	"github.com/wzc5840/gin-api-demo/pkg/logger"
	"github.com/wzc5840/gin-api-demo/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger.Init()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=gin_demo port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := router.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}