// main.go
package main

import (
	"log"
	"os"

	controller "github.com/Ryu732/qr-rallies/controllers"
	"github.com/Ryu732/qr-rallies/db"
	"github.com/Ryu732/qr-rallies/infra"
	"github.com/Ryu732/qr-rallies/models"
	"github.com/Ryu732/qr-rallies/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	rallyRepository := repositories.NewRallyRepository(db)
	RallyController := controller.NewRallyController(rallyRepository)
	stampRepository := repositories.NewStampRepository(db)
	StampController := controller.NewStampController(stampRepository)

	router := gin.Default()
	router.Use(cors.Default())

	// ヘルスチェック用エンドポイント
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "QR Rally API is running",
			"status":  "healthy",
		})
	})

	// 手動マイグレーション用エンドポイント
	router.POST("/migrate", func(c *gin.Context) {
		log.Printf("手動マイグレーション開始...")

		if err := db.AutoMigrate(&models.Rally{}, &models.Stamp{}); err != nil {
			log.Printf("マイグレーションエラー: %v", err)
			c.JSON(500, gin.H{"error": "マイグレーション失敗", "details": err.Error()})
			return
		}

		log.Printf("手動マイグレーション完了")
		c.JSON(200, gin.H{"message": "マイグレーション完了"})
	})

	// ルーティンググループ
	rallyRouter := router.Group("/rallies")

	rallyRouter.GET("", RallyController.FindAllRallies)
	rallyRouter.GET("/:id", RallyController.FindRallyByID)
	rallyRouter.POST("", RallyController.CreateRally)
	rallyRouter.DELETE("/:id", RallyController.DeleteRally)
	rallyRouter.POST("/login/:id", RallyController.LoginRally)

	stampRouter := router.Group("/stamps")
	stampRouter.POST("", StampController.CreateStamp)

	return router
}

func main() {
	// 環境変数を読み込み
	infra.SettingEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// データベース接続
	database := db.SetupDB()
	router := setupRouter(database)

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
