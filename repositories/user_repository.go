package repositories

import (
	"github.com/pilar_test_rest_api/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindById(id int) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Update(user *models.User) error
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

func (r *userRepository) FindById(id int) (models.User, error) {
	user := models.User{}
	err := r.db.First(&user, id).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	user := models.User{}
	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (userRepository *userRepository) Update(user *models.User) error {
	err := userRepository.db.Save(user).Error
	if err != nil {
		return err
	}

	return nil
}
