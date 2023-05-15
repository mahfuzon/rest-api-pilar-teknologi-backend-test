package libraries

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
)

type Claim struct {
	UserId    int
	ExpiredAt int64
}

func DecodeEncodedTokenToMapClaim(tokenJwt *jwt.Token) (Claim, error) {
	claim := Claim{}
	// parsing to  jwt map claim
	mapClaim, ok := tokenJwt.Claims.(jwt.MapClaims)
	if !tokenJwt.Valid || !ok {
		return claim, errors.New("token invalid")
	}

	claim.UserId = int(mapClaim["user_id"].(float64))
	claim.ExpiredAt = int64(mapClaim["expired_at"].(float64))

	return claim, nil
}

func VerifyTokenBySecretKey(encodedToken string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}

func GenerateNewToken(userId int, typeToken string) (string, error) {
	expiredAt := 60
	if typeToken == "refresh" {
		expiredAt = 3600
	}
	claim := jwt.MapClaims{}
	claim["user_id"] = userId
	claim["expired_at"] = expiredAt

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	secretKey := os.Getenv("jwt_secret_key")
	if typeToken == "refresh" {
		secretKey = os.Getenv("jwt_secret_refresh_token_key")
	}

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}
