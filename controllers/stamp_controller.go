package controller

import (
	"fmt"

	"github.com/Ryu732/qr-rallies/repositories"
	"github.com/gin-gonic/gin"
)

type IStampController interface {
	CreateStamp(ctx *gin.Context)
}

type StampController struct {
	repository repositories.IStampRepository
}

func NewStampController(repository repositories.IStampRepository) IStampController {
	return &StampController{repository: repository}
}

type ImageRequest struct {
	Themas  string   `json:"themas" binding:"required"`
	Concept []string `json:"concept" binding:"required"`
}

func (c *StampController) CreateStamp(ctx *gin.Context) {
	var jsonReq ImageRequest
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(400, gin.H{"error": "テーマ選択とデザインコンセプトが必要です"})
		return
	}

	// プロンプトを構築
	prompt := c.buildPrompt(jsonReq.Themas, jsonReq.Concept)

	// AIにスタンプの画像を生成させる処理
	// TODO: AI画像生成APIを呼び出す

	ctx.JSON(200, gin.H{
		"message": "スタンプを作成しました",
		"prompt":  prompt, // デバッグ用
	})
}

// プロンプト構築
func (c *StampController) buildPrompt(themas string, concepts []string) string {
	conceptStr := ""
	for _, concept := range concepts {
		conceptStr += fmt.Sprintf("  - %s\n", concept)
	}

	promptTemplate := `# 命令書:
		あなたはプロのデザイナーです。以下の制約条件に基づいて、スタンプラリーのスタンプをデザインしてください。

		# 制約条件:
			- イラスト風で制作する
			- スタンプのテーマ: %s
			- スタンプのコンセプト:%s
	`

	return fmt.Sprintf(promptTemplate, themas, conceptStr)
}
