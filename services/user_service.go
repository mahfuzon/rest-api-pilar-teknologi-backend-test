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
	FindById(id int) (response.UserResponse, error)
	Login(loginRequest request.UserLoginRequest) (response.UserResponse, error)
	UploadProfileImage(url string, userId int) (response.UserResponse, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) Register(userRegisterRequest request.UserRegisterRequest) (response.UserResponse, error) {
	user := models.User{}
	userResponse := response.UserResponse{}

	passWordHash, err := bcrypt.GenerateFromPassword([]byte(userRegisterRequest.Password), bcrypt.MinCost)
	if err != nil {
		return userResponse, err
	}

	user.Email = userRegisterRequest.Email
	user.Password = string(passWordHash)
	user.FirstName = userRegisterRequest.FirstName
	user.LastName = userRegisterRequest.LastName
	user.Address = userRegisterRequest.Address
	user.City = userRegisterRequest.City
	user.Province = userRegisterRequest.Province
	user.Country = userRegisterRequest.Country
	user.Telephone = userRegisterRequest.Telephone

	err = service.userRepository.Create(&user)
	if err != nil {
		return userResponse, err
	}

	userResponse.ParseToUserResponse(user)

	return userResponse, nil
}

func (service *userService) FindById(id int) (response.UserResponse, error) {
	userResponse := response.UserResponse{}

	user, err := service.userRepository.FindById(id)
	if err != nil {
		return userResponse, err
	}

	userResponse.ParseToUserResponse(user)

	return userResponse, nil
}

func (service *userService) Login(loginRequest request.UserLoginRequest) (response.UserResponse, error) {
	userResponse := response.UserResponse{}

	user, err := service.userRepository.FindByEmail(loginRequest.Email)
	if err != nil {
		return userResponse, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return userResponse, err
	}

	userResponse.ParseToUserResponse(user)

	return userResponse, nil
}

func (service *userService) UploadProfileImage(url string, userId int) (response.UserResponse, error) {
	userResponse := response.UserResponse{}
	user, err := service.userRepository.FindById(userId)
	if err != nil {
		return userResponse, err
	}

	user.ProfileImage = url

	err = service.userRepository.Update(&user)
	if err != nil {
		return userResponse, err
	}

	userResponse.Id = user.Id
	userResponse.Email = user.Email
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.Telephone = user.Telephone
	userResponse.Address = user.Address
	userResponse.City = user.City
	userResponse.Province = user.Province
	userResponse.Country = user.Country
	userResponse.ProfileImage = user.ProfileImage

	return userResponse, nil
}
