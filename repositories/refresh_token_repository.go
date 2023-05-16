package repositories

import (
	"github.com/pilar_test_rest_api/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(refreshToken *models.RefreshToken) error
	Find(refreshToken string) (models.RefreshToken, error)
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

func (refreshTokenRepository *refreshTokenRepository) Create(refreshToken *models.RefreshToken) error {
	err := refreshTokenRepository.db.Create(refreshToken).Error
	if err != nil {
		return err
	}

	return nil
}

func (refreshTokenRepository *refreshTokenRepository) Find(refreshToken string) (models.RefreshToken, error) {
	refreshTokenModel := models.RefreshToken{}
	err := refreshTokenRepository.db.Where("token = ?", refreshToken).First(&refreshTokenModel).Error
	if err != nil {
		return refreshTokenModel, err
	}

	return refreshTokenModel, nil

}
