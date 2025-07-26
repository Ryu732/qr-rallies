package repositories

import (
	"github.com/Ryu732/qr-rallies/models"
	"gorm.io/gorm"
)

type IRallyRepository interface {
	FindAllRallies() (*[]RallyNameResponse, error)
	FindRallyByID(id uint) (*models.Rally, error)
	CreateRally(rally *models.Rally) (*models.Rally, error)
	DeleteRally(id uint) error
	CheckRallyNameExists(name string) (bool, error)
}

// データベース用のRepository
type RallyRepository struct {
	database *gorm.DB
}

func NewRallyRepository(database *gorm.DB) IRallyRepository {
	return &RallyRepository{database: database}
}

// レスポンス用の構造体
type RallyNameResponse struct {
	ID        uint   `json:"id"`
	RallyName string `json:"rally_name"`
}

// RallyRepository用のFindAllRalliesメソッドを実装
func (r *RallyRepository) FindAllRallies() (*[]RallyNameResponse, error) {
	var rallies []RallyNameResponse

	if err := r.database.Model(&models.Rally{}).
		Select("id", "rally_name").
		Find(&rallies).Error; err != nil {
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

// RallyRepository用のCheckRallyNameExistsメソッドを実装
func (r *RallyRepository) CheckRallyNameExists(name string) (bool, error) {
	var count int64
	if err := r.database.Model(&models.Rally{}).Where("rally_name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
