package db

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Ryu732/qr-rallies/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	var dsn string

	// 本番環境
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		// HerokuのDATABASE_URLをパース
		parsedURL, err := url.Parse(databaseURL)
		if err != nil {
			log.Fatal("DATABASE_URLのパースエラー:", err)
		}

		password, _ := parsedURL.User.Password()
		host := parsedURL.Host
		dbname := parsedURL.Path[1:] // 先頭の'/'を除去

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=require TimeZone=UTC",
			host, parsedURL.User.Username(), password, dbname)
	} else {
		// ローカル開発環境
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Tokyo",
			host, port, user, password, dbname)
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("データベース接続エラー:", err)
	}

	// マイグレーション実行
	if err := database.AutoMigrate(&models.Rally{}, &models.Stamp{}); err != nil {
		log.Fatal("マイグレーションエラー:", err)
	}

	return database
}
