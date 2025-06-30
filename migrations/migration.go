package main

import (
	"log"

	"github.com/Ryu732/qr-rallies/db"
	"github.com/Ryu732/qr-rallies/infra"
	"github.com/Ryu732/qr-rallies/models"
)

func main() {
	infra.SettingEnv()
	database := db.SetupDB()

	if err := database.AutoMigrate(&models.Rally{}, &models.Stamp{}); err != nil {
		log.Printf("マイグレーションエラー: %v", err)
		panic("failed to migrate database")
	}

	log.Println("マイグレーション完了!")
}
