package test

import (
	"encoding/json"
	"fmt"
	"github.com/pilar_test_rest_api/controllers"
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/models"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/response"
	"github.com/pilar_test_rest_api/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetDetailArticle(t *testing.T) {
	e := libraries.SetRouter()
	db := database.SetDbTest()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/article/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	article := models.Article{
		Title: "title",
		Body:  "description",
	}

	db.Exec("TRUNCATE TABLE ARTICLES")

	err := db.Create(&article).Error
	if err != nil {
		assert.Equal(t, nil, err, err.Error())
	}

	articleRepository := repositories.NewArticleRepository(db)
	articleService := services.NewArticleService(articleRepository)
	articleController := controllers.NewArticleController(articleService)

	articleController.Get(c)
	apiResponse := response.APIResponse{}
	json.Unmarshal(rec.Body.Bytes(), &apiResponse)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "ok", apiResponse.Status, apiResponse.Data)
}

func TestGetAllArticle(t *testing.T) {
	e := libraries.SetRouter()
	db := database.SetDbTest()
	q := make(url.Values)
	q.Set("limit", "5")
	q.Set("offset", "0")
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8000/api/article?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db.Exec("TRUNCATE TABLE ARTICLES")

	param := c.QueryParam("limit")
	fmt.Println(param + "hjhjhjhj")

	for i := 0; i < 10; i++ {
		article := models.Article{
			Title: "title",
			Body:  "description",
		}
		err := db.Create(&article).Error
		if err != nil {
			assert.Equal(t, nil, err, err.Error())
		}
	}

	articleRepository := repositories.NewArticleRepository(db)
	articleService := services.NewArticleService(articleRepository)
	articleController := controllers.NewArticleController(articleService)

	articleController.GetAll(c)
	apiResponse := response.APIResponse{}
	json.Unmarshal(rec.Body.Bytes(), &apiResponse)

	fmt.Println(apiResponse)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "ok", apiResponse.Status, apiResponse.Data)
}
