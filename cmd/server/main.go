package main

import (
	"log"
	"os"
	"time"

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

	var db *gorm.DB
	var err error
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		
		logger.Infof("データベース接続に失敗しました (試行 %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("データベース接続に失敗しました（リトライ後）:", err)
	}

	logger.Info("データベースに正常に接続しました")

	r := router.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("サーバーがポート " + port + " で開始されました")
	if err := r.Run(":" + port); err != nil {
		log.Fatal("サーバーの開始に失敗しました:", err)
	}
}