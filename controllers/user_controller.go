package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/response"
	"github.com/pilar_test_rest_api/services"
)

type userController struct {
	userService services.UserService
}

func NewUserController(service services.UserService) *userController {
	return &userController{
		userService: service,
	}
}

func (controller userController) Register(ctx echo.Context) error {
	userRegisterRequest := request.UserRegisterRequest{}

	err := ctx.Bind(&userRegisterRequest)

	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed register",
			Data:    err.Error(),
		})
	}

	err = ctx.Validate(&userRegisterRequest)
	if err != nil {
		fmt.Println(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return ctx.JSON(422, response.APIResponse{
				Status:  "error",
				Message: "Failed Register",
				Data:    err.Error(),
			})
		}
		var errorMessage []string

		for _, e := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, e.Error())
		}

		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "failed register",
			Data:    errorMessage,
		})
	}

	userRegisterResponse, err := controller.userService.Register(userRegisterRequest)
	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed register",
			Data:    err.Error(),
		})
	}

	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success create",
		Data:    userRegisterResponse,
	})
}
