package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/junicochandra/golang-api-service/docs"

	"github.com/junicochandra/golang-api-service/internal/domain/entity"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/config/database"
	"github.com/junicochandra/golang-api-service/internal/infrastructure/service/rabbitmq"
	"github.com/junicochandra/golang-api-service/internal/router"
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

	// Router setup
	r := router.SetupRouter()

	// Quick RabbitMQ smoke test
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Printf("use rabbitmq default url for test")
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	broker := rabbitmq.New(rabbitURL)
	if err := broker.Connect(); err != nil {
		log.Printf("RabbitMQ connect failed: %v", err)
	} else {
		defer func() {
			if err := broker.Close(); err != nil {
				log.Printf("RabbitMQ close error: %v", err)
			}
		}()

		testQueue := "conn_test_queue"
		body := []byte(fmt.Sprintf("hello at %s", time.Now().Format(time.RFC3339)))
		if err := broker.Publish(testQueue, body); err != nil {
			log.Printf("RabbitMQ publish failed: %v", err)
		} else {
			log.Println("Published test message to RabbitMQ")
		}

		// Add delay to ensure message is ready to be consumed
		time.Sleep(200 * time.Millisecond)

		// Consume the message
		msg, err := broker.ConsumeOne(testQueue)
		if err != nil {
			log.Printf("ConsumeOne error: %v", err)
		} else {
			log.Printf("Consumed test message: %s", string(msg))
		}
	}

	// Run server
	r.Run(":9000")
}
