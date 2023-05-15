package request

type CreateNewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	UserId       int    `json:"user_id" validate:"required"`
}
