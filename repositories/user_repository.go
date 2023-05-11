package repositories

import (
	"github.com/pilar_test_rest_api/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	err := r.db.Create(user).Error

	if err != nil {
		return err
	}

	return nil
}
