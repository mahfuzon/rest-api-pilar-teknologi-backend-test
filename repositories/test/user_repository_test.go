package test

import (
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	db := database.SetDbTest()
	userRepository := repositories.NewUserRepository(db)
	err := db.Exec("TRUNCATE TABLE USERS").Error
	if err != nil {
		assert.Equal(t, nil, err, err.Error())
	}

	user, err := libraries.CreateExampleUserObject()

	if err != nil {
		assert.Equal(t, nil, err, err.Error())
	}

	_, err = userRepository.Create(&user)

	if err != nil {
		assert.Equal(t, nil, err, err.Error())
	}

}
