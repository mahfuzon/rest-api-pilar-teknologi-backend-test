package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"github.com/labstack/echo/v4"
	"github.com/pilar_test_rest_api/controllers"
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/middleware"
	"github.com/pilar_test_rest_api/models"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupUserController(db *gorm.DB) *controllers.UserController {
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	userService := services.NewUserService(userRepository, refreshTokenRepository)
	authService := services.NewAuthService(refreshTokenRepository)
	userController := controllers.NewUserController(userService, authService)
	return userController
}

func TestGetProfileUser(t *testing.T) {
	router := libraries.SetRouter() // setup router
	db := database.SetDbTest()      // setup db test

	// clean data
	db.Exec("TRUNCATE TABLE USERS")
	db.Exec("TRUNCATE TABLE REFRESH_TOKENS")
	// end

	userController := setupUserController(db)
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	userService := services.NewUserService(userRepository, refreshTokenRepository)
	authService := services.NewAuthService(refreshTokenRepository)

	// create data user
	user, err := libraries.CreateExampleUserObject()
	if err != nil {
		panic(err.Error())
	}
	err = db.Create(&user).Error
	if err != nil {
		panic(err.Error())
	}
	// end

	// create token authorization
	tokendEncoded, err := libraries.GenerateNewToken(user.Id, "access")
	if err != nil {
		panic(err.Error())
	}
	// end

	router.GET("/api/profile", userController.GetProfile)
	router.Use(middleware.AuthMiddleware(userService, authService))
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/profile", nil)
	rec := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer "+tokendEncoded)
	router.ServeHTTP(rec, req)

	res := rec.Result()
	body, _ := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "success get profile", responseBody["message"], responseBody["data"])
}

func TestRegister(t *testing.T) {
	db := database.SetDbTest()
	db.Exec("TRUNCATE TABLE USERS")
	db.Exec("TRUNCATE TABLE REFRESH_TOKENS")

	router := libraries.SetRouter()
	userController := setupUserController(db)
	router.POST("api/register", userController.Register)

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
	router.ServeHTTP(rec, req)

	res := rec.Result()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(responseBody)

	assert.Equal(t, "ok", responseBody["status"], responseBody["data"])
}

func TestLogin(t *testing.T) {
	db := database.SetDbTest()
	db.Exec("TRUNCATE TABLE USERS")
	db.Exec("TRUNCATE TABLE REFRESH_TOKENS")

	// create new user
	newUser, err := libraries.CreateExampleUserObject()
	if err != nil {
		panic(err.Error())
	}
	db.Create(&newUser)
	// end

	router := libraries.SetRouter()
	f := make(url.Values)
	f.Set("email", "email@example.com")
	f.Set("password", "12345678")

	userController := setupUserController(db)
	router.POST("api/login", userController.Login)

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	body, err := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, "ok", responseBody["status"], responseBody["data"])
}

func TestCreateNewAccessToken(t *testing.T) {
	// setup router dan db
	router := libraries.SetRouter()
	db := database.SetDbTest()
	// end

	db.Exec("TRUNCATE TABLE USERS")
	db.Exec("TRUNCATE TABLE REFRESH_TOKENS")

	userController := setupUserController(db)

	user, err := libraries.CreateExampleUserObject()
	if err != nil {
		panic(err.Error())
	}

	err = db.Create(&user).Error
	if err != nil {
		panic(err.Error())
	}

	tokenEncoded, err := libraries.GenerateNewToken(user.Id, "refresh")
	if err != nil {
		panic(err.Error())
	}

	refreshToken := models.RefreshToken{}
	refreshToken.Token = tokenEncoded
	refreshToken.UserId = user.Id
	err = db.Create(&refreshToken).Error
	if err != nil {
		panic(err.Error())
	}

	createAccessTokenRequest := request.CreateNewAccessTokenRequest{}
	createAccessTokenRequest.RefreshToken = refreshToken.Token
	createAccessTokenRequest.UserId = user.Id

	router.POST("/api/create-new-access-token", userController.CreateNewAccessToken)

	byte, err := json.Marshal(&createAccessTokenRequest)
	if err != nil {
		panic(err.Error())
	}
	req := httptest.NewRequest(http.MethodPost,
		"http://localhost:8000/api/create-new-access-token",
		strings.NewReader(string(byte)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(responseBody)
	assert.Equal(t, "ok", responseBody["status"], responseBody["data"])
}

func TestFindRefreshToken(t *testing.T) {
	// setup router dan db
	router := libraries.SetRouter()
	db := database.SetDbTest()
	// end

	db.Exec("TRUNCATE TABLE USERS")
	db.Exec("TRUNCATE TABLE REFRESH_TOKENS")

	userController := setupUserController(db)

	user, err := libraries.CreateExampleUserObject()
	if err != nil {
		panic(err.Error())
	}

	err = db.Create(&user).Error
	if err != nil {
		panic(err.Error())
	}

	tokenEncoded, err := libraries.GenerateNewToken(user.Id, "refresh")
	if err != nil {
		panic(err.Error())
	}

	refreshToken := models.RefreshToken{}
	refreshToken.Token = tokenEncoded
	refreshToken.UserId = user.Id
	err = db.Create(&refreshToken).Error
	if err != nil {
		panic(err.Error())
	}

	router.GET("/api/get-refresh-token/:refresh_token", userController.FindRefreshToken)

	req := httptest.NewRequest(http.MethodGet,
		"http://localhost:8000/api/get-refresh-token/"+tokenEncoded,
		nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(responseBody)
	assert.Equal(t, "ok", responseBody["status"], responseBody["data"])
}
