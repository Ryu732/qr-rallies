package controller

import (
	"strconv"

	"github.com/Ryu732/qr-rallies/models"
	"github.com/Ryu732/qr-rallies/repositories"
	"github.com/gin-gonic/gin"
)

type IRallyController interface {
	FindAllRallies(ctx *gin.Context)
	FindRallyByID(ctx *gin.Context)
	CreateRally(ctx *gin.Context)
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

func (c *RallyController) FindRallyByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "無効なIDです"})
		return
	}

	rally, err := c.repository.FindRallyByID(uint(id))
	if err != nil {
		ctx.JSON(404, gin.H{"error": "ラリーが見つかりません"})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "ラリーを取得しました",
		"data":    rally,
	})
}

func (c *RallyController) CreateRally(ctx *gin.Context) {
	var rally models.Rally
	if err := ctx.ShouldBindJSON(&rally); err != nil {
		ctx.JSON(400, gin.H{"error": "無効なリクエストボディです"})
		return
	}

	createdRally, err := c.repository.CreateRally(&rally)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "ラリーの作成に失敗しました"})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "ラリーを作成しました",
		"data":    createdRally,
	})
}
