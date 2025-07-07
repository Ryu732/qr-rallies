// main.go
package main

import (
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

	router := gin.Default()
	router.Use(cors.Default())

	// ルーティンググループ
	rallyRouter := router.Group("/rallies")

	rallyRouter.GET("", RallyController.FindAllRallies)
	rallyRouter.GET("/:id", RallyController.FindRallyByID)
	rallyRouter.POST("", RallyController.CreateRally)
	rallyRouter.DELETE(":id", RallyController.DeleteRally)
	rallyRouter.POST("/login/:id", RallyController.LoginRally)

	return router
}

func main() {
	infra.SettingEnv()
	db := db.SetupDB()
	router := setupRouter(db)

	// ポート 8080 で起動
	router.Run(":8080")
}
