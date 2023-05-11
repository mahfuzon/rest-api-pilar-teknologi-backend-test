package response

import "github.com/pilar_test_rest_api/models"

type UserResponse struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"last_name"`
	Telephone    string `json:"telephone"`
	ProfileImage string `json:"profile_image"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Province     string `json:"province"`
	Country      string `json:"country"`
}

func (userResponse *UserResponse) ParseToUserResponse(user models.User) {
	userResponse.Id = user.Id
	userResponse.Email = user.Email
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.Telephone = user.Telephone
	userResponse.ProfileImage = user.ProfileImage
	userResponse.Address = user.Address
	userResponse.City = user.City
	userResponse.Province = user.Province
	userResponse.Country = user.Country
}
