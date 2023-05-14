package models

import "time"

type Article struct {
	Id        int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *Article) TableName() string {
	return "articles"
}
