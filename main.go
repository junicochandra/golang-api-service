package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/junicochandra/golang-api-service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type User struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	CreatedAt sql.NullString `json:"created_at"`
}

var db *sql.DB

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
	_ = godotenv.Load()
	initDB()

	// Swagger
	e := echo.New()
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
// @Success 200 {array} User
// @Failure 400
func getUsers(c echo.Context) error {
	rows, err := db.Query("SELECT id, name, email, created_at FROM users")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		users = append(users, u)
	}

	return c.JSON(http.StatusOK, users)
}

func initDB() {
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error :", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB ping error :", err)
	}

	log.Println("DB connected")
}
