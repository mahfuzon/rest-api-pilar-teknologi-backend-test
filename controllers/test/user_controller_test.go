package test

import (
	"github.com/labstack/echo/v4"
	"github.com/pilar_test_rest_api/controllers"
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {

	userJson := `{
  "email": "email@example.com",
  "password": "12345678",
  "first_name": "first_name",
  "last_name": "last_name",
  "telephone": "38947384xxx",
  "address": "address",
  "city": "city",
  "province": "province",
  "country": "country"
}`

	e := libraries.SetRouter()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db := database.SetDbTest()

	err := db.Exec("TRUNCATE TABLE USERS").Error
	if err != nil {
		assert.Equal(t, nil, err)
	}
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	h := controllers.NewUserController(userService)

	err = h.Register(c)
	if err != nil {
		assert.Equal(t, nil, err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)

}
