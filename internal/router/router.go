package router

import (
	"github.com/gin-gonic/gin"

	authUseCase "github.com/junicochandra/golang-api-service/internal/app/auth"
	userUseCase "github.com/junicochandra/golang-api-service/internal/app/user"
	"github.com/junicochandra/golang-api-service/internal/handler"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/config/database"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/repository"
	"github.com/junicochandra/golang-api-service/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger
	r.GET("/api/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Dependency Injection
	userRepository := repository.NewUserRepository(database.DB)
	userUC := userUseCase.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUC)

	authUC := authUseCase.NewAuthUseCase(userRepository)
	authHandler := handler.NewAuthHandler(authUC)

	// Routes
	api := r.Group("/api/v1")
	{
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Users
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Protected Routes (JWT Required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", handler.Profile)
			protected.POST("/auth/logout", authHandler.Logout)
		}
	}

	return r
}
