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

	// Heroku環境の検出（DYNOまたはHEROKU環境変数の存在で判定）
	isHeroku := os.Getenv("DYNO") != "" || os.Getenv("HEROKU_APP_NAME") != ""

	// 本番環境（Heroku）
	if isHeroku {
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			log.Fatal("Heroku環境でDATABASE_URLが設定されていません")
		}

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

		log.Printf("Heroku環境でデータベース接続: %s", host)
	} else if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		// DATABASE_URLが設定されている場合（Heroku以外）
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

		log.Printf("ローカル環境でデータベース接続: %s:%s", host, port)
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("データベース接続エラー: %v", err)
		log.Printf("データベース接続文字列: %s", dsn)
		log.Fatal("データベース接続エラー:", err)
	}

	log.Printf("データベース接続成功!")

	// マイグレーション実行
	if err := database.AutoMigrate(&models.Rally{}, &models.Stamp{}); err != nil {
		log.Fatal("マイグレーションエラー:", err)
	}

	return database
}
