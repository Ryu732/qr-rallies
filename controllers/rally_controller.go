package controller

import (
	"github.com/Ryu732/qr-rallies/repositories"
	"github.com/gin-gonic/gin"
)

type IRallyController interface {
	FindAllRallies(ctx *gin.Context)
}

type RallyController struct {
	repository repositories.IRallyRepository
}

func NewRallyController(repository repositories.IRallyRepository) IRallyController {
	return &RallyController{repository: repository}
}

// メソッドレシーバーを追加
func (c *RallyController) FindAllRallies(ctx *gin.Context) {
	rallies, err := c.repository.FindAllRallies()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "ラリーの取得に失敗しました"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "全てのラリーを取得しました",
		"data":    rallies,
	})
}
