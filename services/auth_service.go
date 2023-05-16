package services

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/models"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/response"
	"os"
	"time"
)

type AuthService interface {
	GenerateAccessToken(userID int) (string, error)
	ValidateToken(token string, typeToken string) (*jwt.Token, error)
	GenerateNewAccessToken(tokenRequest request.CreateNewAccessTokenRequest) (response.CreateNewAccessTokenResponse, error)
	FindRefreshToken(findRefreshTokenRequest request.FindRefreshTokenRequest) (models.RefreshToken, error)
}

type authService struct {
	refreshTokenRepository repositories.RefreshTokenRepository
}

func NewAuthService(refreshTokenRepository repositories.RefreshTokenRepository) AuthService {
	return &authService{
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (s *authService) FindRefreshToken(findRefreshTokenRequest request.FindRefreshTokenRequest) (models.RefreshToken, error) {
	refreshToken, err := s.refreshTokenRepository.Find(findRefreshTokenRequest.RefreshToken)
	if err != nil {
		return refreshToken, err
	}

	return refreshToken, nil
}

func (s *authService) GenerateAccessToken(userID int) (string, error) {
	token, err := libraries.GenerateNewToken(userID, "access")

	if err != nil {
		return token, err
	}

	return token, nil
}

func (s *authService) ValidateToken(encodedToken string, typeToken string) (*jwt.Token, error) {
	secretKey := os.Getenv("jwt_secret_key")
	if typeToken == "refresh" {
		secretKey = os.Getenv("jwt_secret_refresh_token_key")
	}
	token, err := libraries.VerifyTokenBySecretKey(encodedToken, secretKey)

	if err != nil {
		return token, err
	}

	return token, nil
}

func (s *authService) GenerateNewAccessToken(request request.CreateNewAccessTokenRequest) (response.CreateNewAccessTokenResponse, error) {
	createNewAccessTokenResponse := response.CreateNewAccessTokenResponse{}

	oldRefreshToken, err := s.refreshTokenRepository.Find(request.RefreshToken)
	if err != nil {
		return createNewAccessTokenResponse, err
	}

	// verify token base on secret key
	tokenJwt, err := libraries.VerifyTokenBySecretKey(oldRefreshToken.Token, os.Getenv("jwt_secret_refresh_token_key"))
	if err != nil {
		return createNewAccessTokenResponse, err
	}

	claim, err := libraries.DecodeEncodedTokenToMapClaim(tokenJwt)

	//verify user id
	if oldRefreshToken.UserId != claim.UserId {
		return createNewAccessTokenResponse, errors.New("invalid user id")
	}

	//verify expired token
	if time.Now().Unix() > claim.ExpiredAt {
		return createNewAccessTokenResponse, errors.New("token expired")
	}

	newRefreshTokenEncode, err := libraries.GenerateNewToken(oldRefreshToken.UserId, "refresh")
	if err != nil {
		return createNewAccessTokenResponse, errors.New("failed generate refresh token")
	}

	newRefreshToken := models.RefreshToken{
		Token:  newRefreshTokenEncode,
		UserId: request.UserId,
	}

	err = s.refreshTokenRepository.Create(&newRefreshToken)
	if err != nil {
		return createNewAccessTokenResponse, errors.New("failed create new refresh token")
	}

	// generate new access token
	newAccessTokenEncode, err := libraries.GenerateNewToken(request.UserId, "access")
	if err != nil {
		return createNewAccessTokenResponse, err
	}
	// end

	createNewAccessTokenResponse.Token = newAccessTokenEncode

	return createNewAccessTokenResponse, nil
}
