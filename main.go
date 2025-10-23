package main

import (
	"net/http"

	_ "github.com/junicochandra/golang-api-service/docs"
	"github.com/labstack/echo/v4"
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
	e := echo.New()

	g := e.Group("/api/v1")

	g.GET("/users", GetUser)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":9000"))
}

type UserRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// @Tags Users
// @Summary Get Users
// @Description Get list of users
// @Router /users [get]
// @Accept json
// @Produce json
// @Success 200 {array} UserRequest
// @Failure 400
func GetUser(c echo.Context) error {
	data := UserRequest{
		ID:    1,
		Name:  "Junico Dwi Chandra",
		Email: "junicodwi.chandra@gmail.com",
	}
	return c.JSON(http.StatusOK, data)
}