package repositories

import (
	"github.com/Ryu732/qr-rallies/models"
	"gorm.io/gorm"
)

type IRallyRepository interface {
	FindAllRallies() (*[]models.Rally, error)
	FindRallyByID(id uint) (*models.Rally, error)
	CreateRally(rally *models.Rally) (*models.Rally, error)
	DeleteRally(id uint) error
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

// RallyRepository用のCreateRallyメソッドを実装
func (r *RallyRepository) CreateRally(rally *models.Rally) (*models.Rally, error) {
	if err := r.database.Create(rally).Error; err != nil {
		return nil, err
	}

	return rally, nil
}

// RallyRepository用のDeleteRallyメソッドを実装
func (r *RallyRepository) DeleteRally(id uint) error {
	if err := r.database.Delete(&models.Rally{}, id).Error; err != nil {
		return err
	}

	return nil
}
