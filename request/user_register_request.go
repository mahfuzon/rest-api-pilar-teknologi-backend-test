package request

type UserRegisterRequest struct {
	Email     string `form:"email" validate:"required,email"`
	Password  string `form:"password" validate:"required"`
	FirstName string `form:"first_name" validate:"required"`
	LastName  string `form:"last_name" validate:"required"`
	Telephone string `form:"telephone" validate:"required"`
	Address   string `form:"address" validate:"required"`
	City      string `form:"city" validate:"required"`
	Province  string `form:"province" validate:"required"`
	Country   string `form:"country" validate:"required"`
}

func CreateUserRegisterRequestExample() UserRegisterRequest {
	return UserRegisterRequest{
		Email:     "email@example.com",
		Password:  "12345678",
		FirstName: "first name",
		LastName:  "last name",
		Telephone: "0812334vv",
		Address:   "address",
		City:      "city",
		Province:  "province",
		Country:   "country",
	}
}
