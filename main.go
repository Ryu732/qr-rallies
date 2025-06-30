// main.go
package main

import (
	"github.com/Ryu732/qr-rallies/db"
	"github.com/Ryu732/qr-rallies/infra"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// ルーティングの設定
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected", // データベース接続確認
		})
	})

	// データベースを使用するエンドポイントの例
	router.GET("/db-status", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(500, gin.H{"error": "database connection failed"})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"error": "database ping failed"})
			return
		}

		c.JSON(200, gin.H{"status": "database connected"})
	})

	return router
}

func main() {
	infra.SettingEnv()
	database := db.SetupDB()

	router := setupRouter(database)

	// ポート 8080 で起動
	router.Run(":8080")
}
