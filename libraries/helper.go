package libraries

import (
	"github.com/pilar_test_rest_api/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateExampleUserObject() (models.User, error) {
	user := models.User{}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.Email = "email@example.com"
	user.Password = string(passwordHash)
	user.FirstName = "first name"
	user.LastName = "last name"
	user.Telephone = "08153725xxx"
	user.ProfileImage = ""
	user.Address = "address"
	user.City = "city"
	user.Province = "province"
	user.Country = "country"

	return user, nil
}
