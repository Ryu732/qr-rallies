// main.go（自動マイグレーション＆airなし対応版）
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
			"env": gin.H{
				"port":     os.Getenv("PORT"),
				"dyno":     os.Getenv("DYNO"),
				"database": os.Getenv("DATABASE_URL") != "",
				"openai":   os.Getenv("OPENAI_API_KEY") != "",
			},
		})
	})

	// ルーティンググループ
	rallyRouter := router.Group("/rallies")
	rallyRouter.GET("", RallyController.FindAllRallies)
	rallyRouter.GET("/check-name", RallyController.CheckRallyName)
	rallyRouter.GET("/:id", RallyController.FindRallyByID)
	rallyRouter.POST("", RallyController.CreateRally)
	rallyRouter.DELETE("/:id", RallyController.DeleteRally)
	rallyRouter.POST("/login/:id", RallyController.LoginRally)

	stampRouter := router.Group("/stamps")
	stampRouter.POST("", StampController.CreateStamp)

	return router
}

// 自動マイグレーション関数
func runAutoMigration(database *gorm.DB) {
	log.Printf("=== 自動マイグレーション開始 ===")

	if err := database.AutoMigrate(&models.Rally{}, &models.Stamp{}); err != nil {
		log.Fatalf("マイグレーションエラー: %v", err)
	}

	log.Printf("=== 自動マイグレーション完了 ===")
}

func main() {
	log.Printf("=== QR Rally アプリケーション開始 ===")

	// 環境変数を読み込み
	infra.SettingEnv()

	// ポート設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("ポート設定: %s", port)

	// 環境判定
	if dyno := os.Getenv("DYNO"); dyno != "" {
		log.Printf("Heroku環境で実行中: %s", dyno)
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Printf("ローカル環境で実行中")
	}

	// データベース接続
	log.Printf("データベース接続開始...")
	database := db.SetupDB()
	log.Printf("データベース接続完了")

	// 自動マイグレーション実行
	runAutoMigration(database)

	// ルーター設定
	router := setupRouter(database)

	log.Printf("=== サーバー起動中... ポート:%s ===", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("サーバー起動失敗:", err)
	}
}
