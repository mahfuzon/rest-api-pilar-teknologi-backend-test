package request

type UserRegisterRequest struct {
	Email     string `json:"email" validate:"required,email" form:"email"`
	Password  string `json:"password" validate:"required" form:"password"`
	FirstName string `json:"first_name" validate:"required" form:"first_name"`
	LastName  string `json:"last_name" validate:"required" form:"last_name"`
	Telephone string `json:"telephone" validate:"required" form:"telephone"`
	Address   string `json:"address" validate:"required" form:"address"`
	City      string `json:"city" validate:"required" form:"city"`
	Province  string `json:"province" validate:"required" form:"province"`
	Country   string `json:"country" validate:"required" form:"country"`
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
