package test

import (
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister(t *testing.T) {
	db := database.SetDbTest()
	err := db.Exec("TRUNCATE TABLE USERS").Error
	if err != nil {
		assert.Equal(t, nil, err, err.Error())
	}
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userRegisterRequest := request.CreateUserRegisterRequestExample()
	_, err = userService.Register(userRegisterRequest)
	if err != nil {
		assert.Equal(t, nil, err, err.Error())
	}
}
