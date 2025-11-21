package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/junicochandra/golang-api-service/docs"

	"github.com/junicochandra/golang-api-service/internal/bootstrap"
	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/config/database"
)

// @Title Golang API Service
// @Version 1.0
// @Description This RESTful API service is developed in Golang using the Gin framework. It provides structured endpoints for efficient data management and high-performance request handling.
// @Contact.name Junico Dwi Chandra
// @Contact.url https://junicochandra.com/
// @Contact.email junicodwi.chandra@gmail.com
// @Host localhost:9000
// @BasePath /api/v1

// @SecurityDefinitions.apikey BearerAuth
// @In header
// @Name Authorization
func main() {
	// DB Connection
	database.Connect()
	db := database.DB
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	bootstrap.Run()
}
