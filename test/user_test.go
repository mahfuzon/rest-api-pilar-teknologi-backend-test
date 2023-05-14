package test

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pilar_test_rest_api/controllers"
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/response"
	"github.com/pilar_test_rest_api/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetProfileUser(t *testing.T) {
	router := libraries.SetRouter() // setup router

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/profile", nil)
	rec := httptest.NewRecorder()
	ctx := router.NewContext(req, rec)
	db := database.SetDbTest()

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService()
	userController := controllers.NewUserController(userService, authService)

	user, err := libraries.CreateExampleUserObject()
	if err != nil {
		panic(err.Error())
	}

	db.Exec("TRUNCATE TABLE USERS")

	userRepository.Create(&user)

	token, err := authService.GenerateToken(user.Id)
	if err != nil {
		panic(err.Error())
	}

	ctx.Request().Header.Set("Authorization", "Bearer "+token)

	authToken := ctx.Request().Header.Get("Authorization")

	if !strings.Contains(authToken, "Bearer") {
		panic(errors.New("invalid token bearer"))
	}

	tokenString := ""
	arrayToken := strings.Split(authToken, " ")
	if len(arrayToken) == 2 {
		tokenString = arrayToken[1]
	}

	jwtToken, err := authService.ValidateToken(tokenString)
	if err != nil {
		panic(err.Error())
	}

	claim, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok || !jwtToken.Valid {
		panic(err.Error())
	}

	userID := int(claim["user_id"].(float64))

	userResponse, err := userService.FindById(userID)
	if err != nil {
		panic(err.Error())
	}

	ctx.Set("user", userResponse)

	userController.GetProfile(ctx)

	apiResponse := response.APIResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), &apiResponse)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, "ok", apiResponse.Status, apiResponse.Data)
}

func TestRegister(t *testing.T) {
	router := libraries.SetRouter()
	f := make(url.Values)
	f.Set("first_name", "Jon Snow")
	f.Set("email", "jon@labstack.com")
	f.Set("password", "12345678")
	f.Set("last_name", "last name")
	f.Set("telephone", "081278160")
	f.Set("address", "address")
	f.Set("city", "city")
	f.Set("province", "province")
	f.Set("country", "country")
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/register", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	ctx := router.NewContext(req, rec)
	db := database.SetDbTest()

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService()
	userController := controllers.NewUserController(userService, authService)

	db.Exec("TRUNCATE TABLE USERS")

	userController.Register(ctx)

	apiResponse := response.APIResponse{}
	err := json.Unmarshal(rec.Body.Bytes(), &apiResponse)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, "ok", apiResponse.Status, apiResponse.Data)
}

func TestLogin(t *testing.T) {
	router := libraries.SetRouter()
	f := make(url.Values)
	f.Set("email", "email@example.com")
	f.Set("password", "12345678")

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	ctx := router.NewContext(req, rec)
	db := database.SetDbTest()

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService()
	userController := controllers.NewUserController(userService, authService)

	db.Exec("TRUNCATE TABLE USERS")
	newUser, err := libraries.CreateExampleUserObject()
	if err != nil {
		panic(err.Error())
	}
	db.Create(&newUser)

	userController.Login(ctx)

	apiResponse := response.APIResponse{}
	err = json.Unmarshal(rec.Body.Bytes(), &apiResponse)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, "ok", apiResponse.Status, apiResponse.Data)
}
