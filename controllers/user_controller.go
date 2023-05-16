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

type UserController struct {
	userService services.UserService
	authService services.AuthService
}

func NewUserController(userService services.UserService, authService services.AuthService) *UserController {
	return &UserController{
		userService: userService,
		authService: authService,
	}
}

func (controller *UserController) Register(ctx echo.Context) error {
	// binding data inputan user ke struct register
	userRegisterRequest := request.UserRegisterRequest{}
	err := ctx.Bind(&userRegisterRequest)
	if err != nil {
		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "Failed Register",
			Data:    err.Error(),
		})
	}
	// end

	// validasi input
	err = ctx.Validate(&userRegisterRequest)
	if err != nil {
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
	// end

	// panggil service register
	userRegisterResponse, err := controller.userService.Register(userRegisterRequest)
	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed register",
			Data:    err.Error(),
		})
	}
	// end

	// response
	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success create",
		Data:    userRegisterResponse,
	})
}

func (controller *UserController) GetProfile(ctx echo.Context) error {
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

func (controller *UserController) Login(ctx echo.Context) error {
	// binding inputan user ke struct userLogin request
	userLoginRequest := request.UserLoginRequest{}
	err := ctx.Bind(&userLoginRequest)
	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed login",
			Data:    err.Error(),
		})
	}
	//  end

	// validasi input
	err = ctx.Validate(&userLoginRequest)
	if err != nil {
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
			Message: "failed login",
			Data:    errorMessage,
		})
	}
	// end

	// panggil service login
	loginResponse, err := controller.userService.Login(userLoginRequest)
	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed login",
			Data:    err.Error(),
		})
	}
	// end

	// response
	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success login",
		Data:    loginResponse,
	})
}

func (controller *UserController) UploadImage(ctx echo.Context) error {
	imgFile, err := ctx.FormFile("profile_image")
	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed upload image",
			Data:    err.Error(),
		})
	}

	// Source
	src, err := imgFile.Open()
	if err != nil {
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

func (controller *UserController) CreateNewAccessToken(ctx echo.Context) error {
	createNewAccessTokenRequest := request.CreateNewAccessTokenRequest{}

	// binding data inputan user ke struct input
	err := ctx.Bind(&createNewAccessTokenRequest)
	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed create new access token",
			Data:    err.Error(),
		})
	}
	// end

	// validasi input
	err = ctx.Validate(&createNewAccessTokenRequest)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return ctx.JSON(422, response.APIResponse{
				Status:  "error",
				Message: "failed create access token",
				Data:    err.Error(),
			})
		}
		var errorMessage []string

		for _, e := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, e.Error())
		}

		return ctx.JSON(422, response.APIResponse{
			Status:  "error",
			Message: "failed create new access token",
			Data:    errorMessage,
		})
	}
	// end

	// panggil service get new acces token
	createNewAccessTokenResponse, err := controller.authService.GenerateNewAccessToken(createNewAccessTokenRequest)
	if err != nil {
		return ctx.JSON(400, response.APIResponse{
			Status:  "error",
			Message: "failed create new access token",
			Data:    err.Error(),
		})
	}
	// end

	// response
	return ctx.JSON(200, response.APIResponse{
		Status:  "ok",
		Message: "success create new token access",
		Data:    createNewAccessTokenResponse,
	})

}

func (controller *UserController) FindRefreshToken(ctx echo.Context) error {
	findrefreshTokenRequest := request.FindRefreshTokenRequest{}

	// binding data request
	err := ctx.Bind(&findrefreshTokenRequest)
	if err != nil {
		apiResponse := response.APIResponse{
			Status:  "error",
			Message: "failed get refresh token",
			Data:    err.Error(),
		}
		return ctx.JSON(400, apiResponse)
	}
	// end

	// validate data request
	err = ctx.Validate(&findrefreshTokenRequest)
	if err != nil {
		var errorMessage []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, e.Error())
		}

		apiResponse := response.APIResponse{
			Status:  "error",
			Message: "failed get refresh token",
			Data:    errorMessage,
		}
		return ctx.JSON(400, apiResponse)
	}
	//	end

	// dapatkan user login

	// service find refresh token
	refreshToken, err := controller.authService.FindRefreshToken(findrefreshTokenRequest)
	if err != nil {
		apiResponse := response.APIResponse{
			Status:  "error",
			Message: "failed get refresh token",
			Data:    err.Error(),
		}
		return ctx.JSON(400, apiResponse)
	}
	// end

	apiResponse := response.APIResponse{
		Status:  "ok",
		Message: "success get refresh token",
		Data:    refreshToken,
	}
	return ctx.JSON(400, apiResponse)
}
