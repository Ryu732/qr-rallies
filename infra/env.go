package infra

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SettingEnv() {
	if os.Getenv("DYNO") == "" {
		if err := godotenv.Load(); err != nil {
			log.Printf(".envファイルが見つかりません: %v", err)
		}
	}

	// 開発環境時のみ確認
	if os.Getenv("DYNO") == "" {
		requiredEnvs := []string{
			"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		}

		for _, env := range requiredEnvs {
			if os.Getenv(env) == "" {
				log.Printf("警告: %s 環境変数が設定されていません", env)
			}
		}
	}
}
