package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pilar_test_rest_api/response"
	"github.com/pilar_test_rest_api/services"
	"net/http"
	"strings"
)

func AuthMiddleware(userService services.UserService, authService services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if !strings.Contains(authHeader, "Bearer") {
				response := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    nil,
				}
				return c.JSON(http.StatusUnauthorized, response)
			}

			tokenString := ""
			arrayToken := strings.Split(authHeader, " ")
			if len(arrayToken) == 2 {
				tokenString = arrayToken[1]
			}

			token, err := authService.ValidateToken(tokenString)
			if err != nil {
				response := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    err.Error(),
				}
				return c.JSON(http.StatusUnauthorized, response)
			}

			claim, ok := token.Claims.(jwt.MapClaims)

			if !ok || !token.Valid {
				response := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    err.Error(),
				}
				return c.JSON(http.StatusUnauthorized, response)
			}

			userID := int(claim["user_id"].(float64))

			userResponse, err := userService.FindById(userID)
			if err != nil {
				response := response.APIResponse{
					Status:  "error",
					Message: "Unauthorized",
					Data:    err.Error(),
				}
				return c.JSON(http.StatusUnauthorized, response)
			}

			c.Set("user", userResponse)
			return next(c)
		}
	}
}
