package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/junicochandra/golang-api-service/docs"

	"github.com/junicochandra/golang-api-service/internal/config/database"
	"github.com/junicochandra/golang-api-service/internal/entity"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Golang API Service
// @version 1.0
// @description This is a RESTful API service built with Golang for managing data and handling requests efficiently.
// @contact.name Junico Dwi Chandra
// @contact.url https://junicochandra.com/
// @contact.email junicodwi.chandra@gmail.com
// @host localhost:9000
// @BasePath /api/v1
func main() {
	// DB Connection
	database.Connect()
	defer database.DB.Close()

	// Echo init
	e := echo.New()

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes
	g := e.Group("/api/v1")
	g.GET("/users", getUsers)

	e.Logger.Fatal(e.Start(":9000"))
}

// @Tags Users
// @Summary Get all users
// @Description Get all users from database
// @Router /users [get]
// @Accept json
// @Produce json
// @Success 200 {array} entity.User
// @Failure 400
func getUsers(c echo.Context) error {
	rows, err := database.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		users = append(users, u)
	}

	return c.JSON(http.StatusOK, users)
}
