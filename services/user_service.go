package services

import (
	"github.com/pilar_test_rest_api/models"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/response"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(request request.UserRegisterRequest) (response.UserResponse, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) Register(request request.UserRegisterRequest) (response.UserResponse, error) {
	user := models.User{}
	userResponse := response.UserResponse{}

	passWordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		return userResponse, err
	}

	user.Email = request.Email
	user.Password = string(passWordHash)
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Address = request.Address
	user.City = request.City
	user.Province = request.Province
	user.Country = request.Country
	user.Telephone = request.Telephone

	err = service.userRepository.Create(&user)
	if err != nil {
		return userResponse, err
	}

	userResponse.ParseToUserResponse(user)

	return userResponse, nil
}
