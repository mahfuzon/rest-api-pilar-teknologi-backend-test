package services

import (
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/response"
)

type ArticleService interface {
	Get(id int) (response.ArticleResponse, error)
	GetAll(articleRequest request.ArticleRequest) ([]response.ArticleResponse, error)
}

type articleService struct {
	articleRepository repositories.ArticleRepository
}

func NewArticleService(articleRepository repositories.ArticleRepository) ArticleService {
	return &articleService{
		articleRepository: articleRepository,
	}
}

func (articleService *articleService) Get(id int) (response.ArticleResponse, error) {
	articleResponse := response.ArticleResponse{}
	article, err := articleService.articleRepository.Get(id)
	if err != nil {
		return articleResponse, err
	}

	articleResponse.Id = article.Id
	articleResponse.Title = article.Title
	articleResponse.Body = article.Body

	return articleResponse, nil

}

func (articleService *articleService) GetAll(articleRequest request.ArticleRequest) ([]response.ArticleResponse, error) {
	var listArticleResponse []response.ArticleResponse
	listArticle, err := articleService.articleRepository.GetAll(articleRequest.Limit, articleRequest.Offset)
	if err != nil {
		return listArticleResponse, nil
	}

	if len(listArticle) > 0 {
		for _, article := range listArticle {
			articleResponse := response.ArticleResponse{
				Id:    article.Id,
				Title: article.Title,
				Body:  article.Body,
			}

			listArticleResponse = append(listArticleResponse, articleResponse)
		}
	}

	return listArticleResponse, nil
}
