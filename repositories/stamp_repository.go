package repositories

import (
	"github.com/Ryu732/qr-rallies/models"
	"gorm.io/gorm"
)

type IStampRepository interface {
	CreateStamp(stamp *models.Stamp) (*models.Stamp, error)
}

// データベース用のRepository
type StampRepository struct {
	database *gorm.DB
}

func NewStampRepository(database *gorm.DB) IStampRepository {
	return &StampRepository{database: database}
}

// レスポンス用の構造体
type StampNameResponse struct {
	ID        uint   `json:"id"`
	StampName string `json:"stamp_name"`
}

// StampRepository用のCreateStampメソッドを実装
func (r *StampRepository) CreateStamp(stamp *models.Stamp) (*models.Stamp, error) {
	if err := r.database.Create(stamp).Error; err != nil {
		return nil, err
	}

	return stamp, nil
}
