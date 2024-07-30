package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func DbConnect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalf("DATABASE_URL not set in .env file")
	}

	db, err := sql.Open("postgres", databaseUrl)
	CheckError(err)
	defer db.Close()

	err = db.Ping()
	CheckError(err)
	fmt.Println("Connected!")

	// QueryDB(db)
}

// func QueryDB(db *sql.DB) {
// 	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`
// 	rows, err := db.Query(query)
// 	CheckError(err)
// 	defer rows.Close()
// 	fmt.Println("Tables:")
// 	for rows.Next() {
// 		var tableName string
// 		err := rows.Scan(&tableName)
// 		CheckError(err)
// 		fmt.Println(tableName)
// 	}
// 	err = rows.Err()
// 	CheckError(err)
// }

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
