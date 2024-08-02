package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalf("DATABASE_URL not set in .env file")
	}

	var openErr error
	DB, openErr = sql.Open("postgres", databaseUrl)
	if openErr != nil {
		log.Fatalf("Error opening database connection: %v", openErr)
	}

	if pingErr := DB.Ping(); pingErr != nil {
		log.Fatalf("Error connecting to database: %v", pingErr)
	}

	log.Println("Connected to DB")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
