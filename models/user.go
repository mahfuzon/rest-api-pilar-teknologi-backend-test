package models

import (
	"time"
)

type User struct {
	Id           int
	Email        string
	Password     string
	FirstName    string
	LastName     string
	Telephone    string
	ProfileImage string
	Address      string
	City         string
	Province     string
	Country      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (user *User) TableName() string {
	return "users"
}
