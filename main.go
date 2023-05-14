package main

import (
	"github.com/joho/godotenv"
	"github.com/pilar_test_rest_api/controllers"
	"github.com/pilar_test_rest_api/database"
	"github.com/pilar_test_rest_api/libraries"
	"github.com/pilar_test_rest_api/middleware"
	"github.com/pilar_test_rest_api/repositories"
	"github.com/pilar_test_rest_api/services"
)

func main() {
	router := libraries.SetRouter()
	db := database.SetDb()

	err := godotenv.Load()
	if err != nil {
		router.Logger.Fatal("Error loading .env file")
	}
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService()
	userController := controllers.NewUserController(userService, authService)

	articleRepository := repositories.NewArticleRepository(db)
	articleService := services.NewArticleService(articleRepository)
	articleController := controllers.NewArticleController(articleService)

	api := router.Group("/api")

	api.GET("/article/:id", articleController.Get)
	api.GET("/article", articleController.GetAll)
	api.POST("/register", userController.Register)
	api.POST("/login", userController.Login)
	api.GET("/profile", userController.GetProfile, middleware.AuthMiddleware(userService, authService))
	api.POST("/upload", userController.UploadImage, middleware.AuthMiddleware(userService, authService))
	router.Logger.Fatal(router.Start(":8000"))
}
