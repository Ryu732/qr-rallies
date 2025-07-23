package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

	// AIで画像を生成
	imageURL, err := GenerateImage(prompt)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "画像生成に失敗しました: " + err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"message":   "スタンプを作成しました",
		"image_url": imageURL,
		"prompt":    prompt, // デバッグ用
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
		- スタンプのコンセプト:
		%s`

	return fmt.Sprintf(promptTemplate, themas, conceptStr)
}

// Azure OpenAI用のAI画像生成関数
func GenerateImage(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	apiURL := os.Getenv("OPENAI_API_URL")

	if apiKey == "" || apiURL == "" {
		return "", fmt.Errorf("azure OpenAI API key or URL not set")
	}

	// Azure OpenAI用のリクエストボディ
	requestBody := map[string]interface{}{
		"prompt":  prompt,
		"size":    "1024x1024",
		"style":   "vivid",
		"quality": "standard",
		"n":       1,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// HTTPリクエストを作成
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Azure OpenAI用のヘッダー
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey) // Bearer認証ではなく、api-keyを使用

	// リクエストを実行
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスを読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// レスポンス構造
	var response struct {
		Created int64 `json:"created"`
		Data    []struct {
			URL           string `json:"url"`
			RevisedPrompt string `json:"revised_prompt,omitempty"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(response.Data) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	return response.Data[0].URL, nil
}
