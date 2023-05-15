package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/response"
	"github.com/pilar_test_rest_api/services"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(userService services.UserService, authService services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if !strings.Contains(authHeader, "Bearer") {
				apiResponse := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    nil,
				}
				return c.JSON(http.StatusUnauthorized, apiResponse)
			}

			tokenString := ""
			arrayToken := strings.Split(authHeader, " ")
			if len(arrayToken) == 2 {
				tokenString = arrayToken[1]
			}

			token, err := authService.ValidateToken(tokenString, "access")
			if err != nil {
				apiResponse := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    err.Error(),
				}
				return c.JSON(http.StatusUnauthorized, apiResponse)
			}

			claim, err := libraries.DecodeEncodedTokenToMapClaim(token)
			if err != nil {
				apiResponse := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    err.Error(),
				}
				return c.JSON(http.StatusUnauthorized, apiResponse)
			}

			if time.Now().Unix() > claim.ExpiredAt {
				apiResponse := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    "token expired",
				}
				return c.JSON(http.StatusUnauthorized, apiResponse)
			}

			userID := claim.UserId

			userResponse, err := userService.FindById(userID)
			if err != nil {
				apiResponse := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    err.Error(),
				}
				return c.JSON(http.StatusUnauthorized, apiResponse)
			}

			c.Set("user", userResponse)
			return next(c)
		}
	}
}
