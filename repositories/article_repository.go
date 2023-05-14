package repositories

import (
	"fmt"
	"github.com/pilar_test_rest_api/models"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetAll(limit int, offset int) ([]models.Article, error)
	Get(id int) (models.Article, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

func (articleRepository *articleRepository) GetAll(limit int, offset int) ([]models.Article, error) {
	var articles []models.Article
	fmt.Println(limit)
	fmt.Println(offset)

	query := articleRepository.db
	if limit != 0 {
		query = query.Limit(limit)
	}

	if offset != 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&articles).Error
	if err != nil {
		return articles, err
	}

	return articles, nil

}

func (articleRepository *articleRepository) Get(id int) (models.Article, error) {
	article := models.Article{}

	err := articleRepository.db.First(&article, id).Error
	if err != nil {
		return article, err
	}

	return article, nil
}
