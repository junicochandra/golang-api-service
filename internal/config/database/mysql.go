package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning .env file not found")
	}

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?parseTime=true"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error :", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping error :", err)
	}

	log.Println("DB connected successfully")
}
