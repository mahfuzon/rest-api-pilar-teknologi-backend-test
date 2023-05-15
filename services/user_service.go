package services

import (
	"errors"
	"fmt"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/models"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/response"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(request request.UserRegisterRequest) (response.UserRegisterResponse, error)
	FindById(id int) (response.UserResponse, error)
	Login(loginRequest request.UserLoginRequest) (response.LoginResponse, error)
	UploadProfileImage(url string, userId int) (response.UserResponse, error)
}

type userService struct {
	userRepository repositories.UserRepository
	repositories.RefreshTokenRepository
}

func NewUserService(userRepository repositories.UserRepository, refreshTokenRepository repositories.RefreshTokenRepository) UserService {
	return &userService{
		userRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (service *userService) Register(userRegisterRequest request.UserRegisterRequest) (response.UserRegisterResponse, error) {
	user := models.User{}
	userRegisterResponse := response.UserRegisterResponse{}

	// cek user is exists
	checUserIsExist, err := service.userRepository.FindByEmail(userRegisterRequest.Email)
	if err != nil {
		fmt.Println(checUserIsExist)
		if checUserIsExist.Id != 0 {
			return userRegisterResponse, errors.New("user already exists")
		}
	}
	// end

	// hash password
	passWordHash, err := bcrypt.GenerateFromPassword([]byte(userRegisterRequest.Password), bcrypt.MinCost)
	if err != nil {
		return userRegisterResponse, err
	}
	// end

	// fill all field object user model
	user.Email = userRegisterRequest.Email
	user.Password = string(passWordHash)
	user.FirstName = userRegisterRequest.FirstName
	user.LastName = userRegisterRequest.LastName
	user.Address = userRegisterRequest.Address
	user.City = userRegisterRequest.City
	user.Province = userRegisterRequest.Province
	user.Country = userRegisterRequest.Country
	user.Telephone = userRegisterRequest.Telephone
	// end

	// create new user
	err = service.userRepository.Create(&user)
	if err != nil {
		return userRegisterResponse, err
	}
	// end

	// generate access token and refresh token
	accessTokenString, err := libraries.GenerateNewToken(user.Id, "access")
	refreshTokenString, err := libraries.GenerateNewToken(user.Id, "refresh")
	if err != nil {
		return userRegisterResponse, err
	}
	// end

	// create refresh token
	refreshToken := models.RefreshToken{
		Token:  refreshTokenString,
		UserId: user.Id,
	}
	err = service.RefreshTokenRepository.Create(&refreshToken)
	if err != nil {
		return userRegisterResponse, err
	}
	// end

	// fill all field userResponse
	userResponse := response.UserResponse{
		Id:           user.Id,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Telephone:    user.Telephone,
		ProfileImage: user.ProfileImage,
		Address:      user.Address,
		City:         user.City,
		Province:     user.Province,
		Country:      user.Country,
	}
	// end

	// fill all user login response
	userLoginResponse := response.LoginResponse{
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
	}
	// end

	userRegisterResponse.User = userResponse
	userRegisterResponse.Token = userLoginResponse

	return userRegisterResponse, nil
}

func (service *userService) FindById(id int) (response.UserResponse, error) {
	userResponse := response.UserResponse{}

	// find user
	user, err := service.userRepository.FindById(id)
	if err != nil {
		return userResponse, err
	}
	// end

	userResponse.ParseToUserResponse(user)

	return userResponse, nil
}

func (service *userService) Login(loginRequest request.UserLoginRequest) (response.LoginResponse, error) {
	// buat object kosong login response
	loginResponse := response.LoginResponse{}
	// end

	// get user by email
	user, err := service.userRepository.FindByEmail(loginRequest.Email)
	if err != nil {
		return loginResponse, err
	}
	// end

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		return loginResponse, err
	}
	// end

	// generate access token string
	accessTokenString, err := libraries.GenerateNewToken(user.Id, "access")
	if err != nil {
		return loginResponse, err
	}

	// generate refresh token string
	refreshTokenString, err := libraries.GenerateNewToken(user.Id, "refresh")
	if err != nil {
		return loginResponse, err
	}
	// end

	// create refreshToken
	refreshToken := models.RefreshToken{
		Token:  refreshTokenString,
		UserId: user.Id,
	}
	err = service.RefreshTokenRepository.Create(&refreshToken)
	if err != nil {
		return loginResponse, err
	}
	// end create refresh token

	loginResponse.Token = accessTokenString
	loginResponse.RefreshToken = refreshTokenString

	return loginResponse, nil
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
