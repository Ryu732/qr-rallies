// main.go
package main

import (
	"log"
	"os"

	controller "github.com/Ryu732/qr-rallies/controllers"
	"github.com/Ryu732/qr-rallies/db"
	"github.com/Ryu732/qr-rallies/infra"
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
