package router

import (
	"github.com/gin-gonic/gin"
	"github.com/junicochandra/golang-api-service/internal/config/database"
	"github.com/junicochandra/golang-api-service/internal/handler"
	"github.com/junicochandra/golang-api-service/internal/repository"
	"github.com/junicochandra/golang-api-service/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger
	r.GET("/api/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Dependency Injection
	userRepository := repository.NewUserRepository(database.DB)
	userUseCase := usecase.NewUserRepository(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)

	// Routes
	api := r.Group("/api/v1")
	{
		api.GET("/users", userHandler.GetUsers)
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUserDetail)
		api.PUT("/users/:id", userHandler.UpdateUser)
		api.DELETE("/users/:id", userHandler.DeleteUser)
	}

	return r
}
