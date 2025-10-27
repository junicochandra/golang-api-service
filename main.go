package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"

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
		v1.POST("/users", createUser)
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

// @Tags Users
// @Summary Create a new user
// @Description Add a new user to the database
// @Router /users [post]
// @Accept json
// @Produce json
// @Param user body user.UserCreateRequest true "User data"
// @Success 201 {object} user.UserCreateResponse
// @Failure 400
// @Failure 500
func createUser(c *gin.Context) {
	var req user.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user into database
	result, err := database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", req.Name, req.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response with created user
	insertedID, _ := result.LastInsertId()
	res := user.UserCreateResponse{
		ID:    uint64(insertedID),
		Name:  req.Name,
		Email: req.Email,
	}

	c.JSON(http.StatusCreated, res)
}
