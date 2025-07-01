package repositories

import (
	"github.com/Ryu732/qr-rallies/models"
	"gorm.io/gorm"
)

type IRallyRepository interface {
	FindAllRallies() (*[]models.Rally, error)
	FindRallyByID(id uint) (*models.Rally, error)
}

// データベース用のRepository
type RallyRepository struct {
	database *gorm.DB
}

func NewRallyRepository(database *gorm.DB) IRallyRepository {
	return &RallyRepository{database: database}
}

// RallyRepository用のFindAllRalliesメソッドを実装
func (r *RallyRepository) FindAllRallies() (*[]models.Rally, error) {
	var rallies []models.Rally

	if err := r.database.Find(&rallies).Error; err != nil {
		return nil, err
	}

	return &rallies, nil
}

// RallyRepository用のFindRallyByIDメソッドを実装
func (r *RallyRepository) FindRallyByID(id uint) (*models.Rally, error) {
	var rally models.Rally

	if err := r.database.First(&rally, id).Error; err != nil {
		return nil, err
	}

	return &rally, nil
}
