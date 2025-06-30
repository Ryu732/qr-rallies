package infra

import (
	"log"
	"os"
)

func SettingEnv() {
	// Docker環境では環境変数が既に設定されているため、.envファイルは不要
	log.Println("環境変数を確認中...")

	// 必要な環境変数が設定されているか確認
	requiredEnvs := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}

	for _, env := range requiredEnvs {
		value := os.Getenv(env)
		if value == "" {
			log.Printf("警告: %s が設定されていません", env)
		} else {
			log.Printf("%s: %s", env, value)
		}
	}
}
