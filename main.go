package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load environment variables from .env file
	// This will load variables from .env if it exists, otherwise it assumes
	// they are set directly in the environment.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it. Assuming environment variables are set directly.")
	}

	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE") // Using DB_SSLMODE to avoid conflict with common SSL_MODE

	// Check if essential variables are set
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("One or more essential database environment variables (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME) are not set.")
	}

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database!")

	// You can now use 'db' to perform database operations
	// For example, to get the current database version:
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Error querying database version: %v", err)
	}
	fmt.Printf("PostgreSQL Version: %s\n", version)
}
