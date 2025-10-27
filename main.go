package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/junicochandra/golang-api-service/docs"

	"github.com/junicochandra/golang-api-service/internal/config/database"
	"github.com/junicochandra/golang-api-service/internal/dto/user"
)

// @Title Golang API Service
// @Version 1.0
// @Description This RESTful API service is developed in Golang using the Gin framework. It provides structured endpoints for efficient data management and high-performance request handling.
// @Contact.name Junico Dwi Chandra
// @Contact.url https://junicochandra.com/
// @Contact.email junicodwi.chandra@gmail.com
// @Host localhost:9000
// @BasePath /api/v1
func main() {
	// DB Connection
	database.Connect()
	defer database.DB.Close()

	// Gin init
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", getUsers)
	}

	r.Run(":9000")
}

// @Tags Users
// @Summary Get all users
// @Description Get all users from database
// @Router /users [get]
// @Accept json
// @Produce json
// @Success 200 {array} user.UserListResponse
// @Failure 400
func getUsers(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []user.UserListResponse
	for rows.Next() {
		var row user.UserListResponse
		if err := rows.Scan(&row.ID, &row.Name, &row.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, row)
	}

	c.JSON(http.StatusOK, users)
}
