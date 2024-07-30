package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func DbConnect() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalf("DATABASE_URL not set in .env file")
	}

	var openErr error
	if db, openErr = sql.Open("postgres", databaseUrl); openErr != nil {
		log.Fatalf("Error opening database connection: %v", openErr)
	}

	if pingErr := db.Ping(); pingErr != nil {
		log.Fatalf("Error connecting to database: %v", pingErr)
	}

	fmt.Println("Connected!")

	return db
}

func RegUser(Username, Email string) error {
	conn := DbConnect()
	defer conn.Close()

	query := `INSERT INTO users (username, email ) VALUES ($1, $2)`
	if _, execErr := conn.Exec(query, Username, Email); execErr != nil {
		return execErr
	}

	return nil
}
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
