package controllers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pilar_test_rest_api/request"
	"github.com/pilar_test_rest_api/response"
	"github.com/pilar_test_rest_api/services"
	"io"
	"os"
	"strconv"
	"time"
)

type userController struct {
	userService services.UserService
	authService services.AuthService
}

func NewUserController(userService services.UserService, authService services.AuthService) *userController {
	return &userController{
		userService: userService,
		authService: authService,
	}
}

func (controller *userController) Register(ctx echo.Context) error {
	userRegisterRequest := request.UserRegisterRequest{}

	err := ctx.Bind(&userRegisterRequest)
	if err != nil {
		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "Failed Register",
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

func (controller *userController) GetProfile(ctx echo.Context) error {
	userLogin := ctx.Get("user")
	user, ok := userLogin.(response.UserResponse)
	if !ok {
		return ctx.JSON(200, response.APIResponse{
			Status:  "error",
			Message: "failed get profile",
			Data:    "data intervace invalid",
		})
	}

	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success get profile",
		Data:    user,
	})
}

func (controller *userController) Login(ctx echo.Context) error {
	userLoginRequest := request.UserLoginRequest{}

	err := ctx.Bind(&userLoginRequest)
	if err != nil {
		panic(err)
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed login",
			Data:    err.Error(),
		})
	}

	err = ctx.Validate(&userLoginRequest)

	if err != nil {
		panic(err)
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

		fmt.Println(errorMessage)

		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "failed login",
			Data:    errorMessage,
		})
	}

	userResponse, err := controller.userService.Login(userLoginRequest)
	if err != nil {
		panic(err)
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed login",
			Data:    err.Error(),
		})
	}
	token, err := controller.authService.GenerateToken(userResponse.Id)

	if err != nil {
		panic(err.Error())
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed login",
			Data:    err.Error(),
		})
	}

	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success login",
		Data:    token,
	})
}

func (controller *userController) UploadImage(ctx echo.Context) error {
	img_file, err := ctx.FormFile("profile_image")
	if err != nil {
		fmt.Println("164")
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed upload image",
			Data:    err.Error(),
		})
	}

	// Source
	src, err := img_file.Open()
	if err != nil {
		fmt.Println("174")
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed upload image",
			Data:    err.Error(),
		})
	}

	defer src.Close()

	// Destination
	timeString := strconv.Itoa(int(time.Now().Unix()))
	staticPublicPath := "public/"
	filePath := "asset/profile_image/" + timeString + ".jpg"
	path := staticPublicPath + filePath
	dst, err := os.Create(path)
	if err != nil {
		fmt.Println("187")
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed upload image",
			Data:    err.Error(),
		})
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("197")
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed upload image",
			Data:    err.Error(),
		})
	}

	userLogin := ctx.Get("user")
	userResponse, ok := userLogin.(response.UserResponse)
	if !ok {
		fmt.Println("207")
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed upload image",
			Data:    errors.New("interface invalid"),
		})
	}

	urlImage := os.Getenv("hostname") + "/" + filePath
	newUserResponse, err := controller.userService.UploadProfileImage(urlImage, userResponse.Id)

	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success upload Image",
		Data:    newUserResponse,
	})

}
