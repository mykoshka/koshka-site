package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	LoadEnv() // Load environment variables

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_USER", "postgres"),
		GetEnv("DB_PASSWORD", "password"),
		GetEnv("DB_NAME", "koshka"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	log.Println("Connected to database successfully")
}
