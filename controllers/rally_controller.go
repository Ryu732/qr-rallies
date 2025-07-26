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
	DeleteRally(ctx *gin.Context)
	LoginRally(ctx *gin.Context)
	CheckRallyName(ctx *gin.Context)
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

func (c *RallyController) DeleteRally(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "無効なIDです"})
		return
	}

	err = c.repository.DeleteRally(uint(id))
	if err != nil {
		ctx.JSON(404, gin.H{"error": "ラリーが見つかりません"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "ラリーを削除しました",
	})
}

func (c *RallyController) LoginRally(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "無効なIDです"})
		return
	}

	// パスワードを含むリクエストボディを構造体にバインド
	var loginRequest struct {
		Password string `json:"Password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(400, gin.H{"error": "パスワードが必要です"})
		return
	}

	// ラリーIDからパスワードを取得
	rally, err := c.repository.FindRallyByID(uint(id))
	if err != nil {
		ctx.JSON(404, gin.H{"error": "ラリーが見つかりません"})
		return
	}

	// パスワード照合
	if rally.Password != loginRequest.Password {
		ctx.JSON(401, gin.H{"error": "パスワードが間違っています"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "ログインしました",
		"data":    rally,
	})
}

func (c *RallyController) CheckRallyName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(400, gin.H{"error": "名前パラメータが必要です"})
		return
	}

	exists, err := c.repository.CheckRallyNameExists(name)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "名前チェックに失敗しました"})
		return
	}

	available := !exists
	message := "この名前は使用可能です"
	if !available {
		message = "この名前は既に使用されています"
	}

	ctx.JSON(200, gin.H{
		"name":      name,
		"available": available,
		"message":   message,
	})
}
