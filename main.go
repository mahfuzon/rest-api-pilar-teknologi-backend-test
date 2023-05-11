package main

import (
	"github.com/pilar_test_rest_api/controllers"
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/services"
)

func main() {
	router := libraries.SetRouter()
	db := database.SetDb()

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	api := router.Group("/api")

	api.POST("/register", userController.Register)

	router.Logger.Fatal(router.Start(":8000"))
}
