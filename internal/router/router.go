package router

import (
	"github.com/gin-gonic/gin"
	"github.com/junicochandra/golang-api-service/internal/config/database"
	"github.com/junicochandra/golang-api-service/internal/handler"
	"github.com/junicochandra/golang-api-service/internal/repository"
	authUseCase "github.com/junicochandra/golang-api-service/internal/usecase/auth"
	userUseCase "github.com/junicochandra/golang-api-service/internal/usecase/user"
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
		// Users
		api.GET("/users", userHandler.GetUsers)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUserByID)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)

		// Auth
		api.POST("/auth/register", authHandler.Register)
	}

	return r
}
