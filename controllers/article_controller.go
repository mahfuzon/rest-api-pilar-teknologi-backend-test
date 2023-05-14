package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/response"
	"github.com/pilar_test_rest_api/services"
)

type articleController struct {
	articleService services.ArticleService
}

func NewArticleController(articleService services.ArticleService) *articleController {
	return &articleController{
		articleService: articleService,
	}
}

func (articleController *articleController) Get(ctx echo.Context) error {
	detailArticleRequest := request.GetDetailArticleRequest{}

	err := ctx.Bind(&detailArticleRequest)
	if err != nil {
		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "failed get data article",
			Data:    err.Error(),
		})
	}

	err = ctx.Validate(&detailArticleRequest)
	if err != nil {
		var errorMessage []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, e.Error())
		}

		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "failed get data article",
			Data:    errorMessage,
		})
	}

	articleResponse, err := articleController.articleService.Get(detailArticleRequest.Id)

	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed get article",
			Data:    err.Error(),
		})
	}

	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success get article",
		Data:    articleResponse,
	})
}

func (articleController *articleController) GetAll(ctx echo.Context) error {
	articleRequest := request.ArticleRequest{}

	err := ctx.Bind(&articleRequest)
	if err != nil {
		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "failed get data article",
			Data:    err.Error(),
		})
	}

	err = ctx.Validate(&articleRequest)
	if err != nil {
		var errorMessage []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, e.Error())
		}

		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "failed get data article",
			Data:    errorMessage,
		})
	}

	articleResponse, err := articleController.articleService.GetAll(articleRequest)

	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed get article",
			Data:    err.Error(),
		})
	}

	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success get article",
		Data:    articleResponse,
	})
}
