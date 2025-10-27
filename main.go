package main

import (
	"database/sql"
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
	r.GET("/api/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	api := r.Group("/api/v1")
	{
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
		api.GET("/users/:id", getUserDetail)
		api.PUT("/users/:id", updateUser)
		api.DELETE("/users/:id", deleteUser)
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
	result, err := database.DB.Exec("INSERT INTO users SET name = ?, email = ?, password = ?", req.Name, req.Email, string(hashedPassword))
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

// @Tags Users
// @Summary Get user detail
// @Description Get detail information of a specific user by ID
// @Router /users/{id} [get]
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} user.UserDetailResponse
// @Failure 400
// @Failure 404
// @Failure 500
func getUserDetail(c *gin.Context) {
	id := c.Param("id")

	var res user.UserDetailResponse
	err := database.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&res.ID, &res.Name, &res.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, res)
}

// @Tags Users
// @Summary Update user
// @Description Update user data by ID
// @Router /users/{id} [put]
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body user.UserDetailRequest true "Updated user data"
// @Success 200 {object} user.UserUpdateResponse
// @Failure 400
// @Failure 404
// @Failure 500
func updateUser(c *gin.Context) {
	id := c.Param("id")

	var req user.UserDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user in database
	result, err := database.DB.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", req.Name, req.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Response with updated user
	res := user.UserUpdateResponse{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	c.JSON(http.StatusOK, res)

}

// @Tags Users
// @Summary Delete a user
// @Description Delete a user from the database by ID
// @Router /users/{id} [delete]
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
func deleteUser(c *gin.Context) {
	id := c.Param("id")

	rowsAffected, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := rowsAffected.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
