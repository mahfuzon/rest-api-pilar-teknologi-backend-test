package models

import "time"

type RefreshToken struct {
	Id        int
	Token     string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *RefreshToken) TableName() string {
	return "refresh_tokens"
}
