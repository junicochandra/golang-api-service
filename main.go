package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/junicochandra/golang-api-service/docs"
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

	g.POST("get-product", GetProduct)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":9000"))
}

type ProductRequest struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"Sample Product"`
	Price int    `json:"price" example:"100"`
}

// @Tags Product
// @Summary Get Product
// @Description Get list of products
// @Router /get-product [post]
// @Param request body ProductRequest true "Payload Body [RAW]"
// @Accept json
// @Produce json
// @Success 200 {array} ProductRequest
// @Failure 400
func GetProduct(c echo.Context) error {
	data := ProductRequest{
		ID:    2,
		Name:  "Sample Product 2",
		Price: 100,
	}
	return c.JSON(http.StatusOK, data)
}