package test

import (
	"encoding/json"
	"fmt"
	"github.com/pilar_test_rest_api/controllers"
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/middleware"
	"github.com/pilar_test_rest_api/models"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/services"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetDetailArticle(t *testing.T) {
	db := database.SetDbTest()

	// clean data
	db.Exec("TRUNCATE TABLE USERS")
	db.Exec("TRUNCATE TABLE REFRESH_TOKENS")
	db.Exec("TRUNCATE TABLE ARTICLES")
	// end

	// create data article
	article := models.Article{
		Title: "title",
		Body:  "description",
	}
	err := db.Create(&article).Error
	if err != nil {
		panic(err.Error())
	}
	// end

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

	articleRepository := repositories.NewArticleRepository(db)
	articleService := services.NewArticleService(articleRepository)
	articleController := controllers.NewArticleController(articleService)
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	authService := services.NewAuthService(refreshTokenRepository)
	userService := services.NewUserService(userRepository, refreshTokenRepository)
	authMiddleware := middleware.AuthMiddleware(userService, authService)
	e := libraries.SetRouter()
	e.GET("/api/article/:id", articleController.Get, authMiddleware)
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/article/"+strconv.Itoa(article.Id), nil)
	rec := httptest.NewRecorder()
	req.Header.Set("Authorization", "Bearer "+tokendEncoded)
	e.ServeHTTP(rec, req)

	res := rec.Result()

	body, err := io.ReadAll(res.Body)
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, 200, res.StatusCode)
}

func TestGetAllArticle(t *testing.T) {
	db := database.SetDbTest()

	// clean data
	db.Exec("TRUNCATE TABLE USERS")
	db.Exec("TRUNCATE TABLE REFRESH_TOKENS")
	db.Exec("TRUNCATE TABLE ARTICLES")
	// end

	// create data article
	for i := 0; i < 10; i++ {
		article := models.Article{
			Title: "title",
			Body:  "description",
		}
		err := db.Create(&article).Error
		if err != nil {
			panic(err.Error())
		}
	}
	// end

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

	articleRepository := repositories.NewArticleRepository(db)
	articleService := services.NewArticleService(articleRepository)
	articleController := controllers.NewArticleController(articleService)
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	authService := services.NewAuthService(refreshTokenRepository)
	userService := services.NewUserService(userRepository, refreshTokenRepository)
	authMiddleware := middleware.AuthMiddleware(userService, authService)

	e := libraries.SetRouter()

	e.GET("api/article", articleController.GetAll, authMiddleware)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/article?limit=5&offset=5", nil)
	rec := httptest.NewRecorder()

	req.Header.Set("Authorization", "Bearer "+tokendEncoded)
	e.ServeHTTP(rec, req)

	res := rec.Result()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var responseBody map[string]interface{}

	json.Unmarshal(body, &responseBody)

	fmt.Println(responseBody)
	assert.Equal(t, 401, res.StatusCode, responseBody["message"])
}
