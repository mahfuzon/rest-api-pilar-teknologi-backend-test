package request

type FindRefreshTokenRequest struct {
	RefreshToken string `param:"refresh_token" validate:"required"`
}
