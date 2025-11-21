package router

import (
	"github.com/gin-gonic/gin"

	"github.com/junicochandra/golang-api-service/internal/app/auth"
	"github.com/junicochandra/golang-api-service/internal/app/payment"
	"github.com/junicochandra/golang-api-service/internal/app/user"
	"github.com/junicochandra/golang-api-service/internal/handler"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/config/database"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/repository"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/service/rabbitmq"
	"github.com/junicochandra/golang-api-service/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(rabbitSvc *rabbitmq.RabbitMQService) *gin.Engine {
	r := gin.Default()

	// Swagger
	r.GET("/api/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Dependency Injection
	userRepository := repository.NewUserRepository(database.DB)
	accountRepository := repository.NewAccountRepository(database.DB)
	transactionRepository := repository.NewTransactionRepository(database.DB)

	userUC := user.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUC)

	authUC := auth.NewAuthUseCase(userRepository)
	authHandler := handler.NewAuthHandler(authUC)

	topUpUC := payment.NewTopUpUseCase(accountRepository, transactionRepository, rabbitSvc)
	topUpHandler := handler.NewPaymentHandler(topUpUC)

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

		// Payment
		pay := api.Group("/payments")
		{
			pay.POST("/topup", topUpHandler.CreateTopUp)
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
