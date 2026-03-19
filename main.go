package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	db             *sql.DB
	dbStartupError error
)

func main() {
	// Load environment variables from .env file
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
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Check if essential variables are set
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		dbStartupError = fmt.Errorf("one or more essential database environment variables are not set")
		log.Println("Database Connection Warning:", dbStartupError)
	} else {
		// Construct the connection string
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

		// Open a database connection
		var err error // shadow err from go dotenv, actually err is already defined in main
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			dbStartupError = fmt.Errorf("error opening database connection: %v", err)
			log.Println("Database Connection Warning:", dbStartupError)
		} else {
			// Ping the database to verify the connection
			err = db.Ping()
			if err != nil {
				dbStartupError = fmt.Errorf("error connecting to the database: %v", err)
				log.Println("Database Connection Warning:", dbStartupError)
			} else {
				fmt.Println("Successfully connected to the PostgreSQL database!")
			}
		}
	}

	// Setup HTTP Server
	http.HandleFunc("/", welcomeHandler)
	
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type PageData struct {
	Connected bool
	Version   string
	ErrorMsg  string
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Connected: false,
	}

	if dbStartupError != nil {
		data.ErrorMsg = dbStartupError.Error()
	} else if db != nil {
		var version string
		err := db.QueryRow("SELECT version()").Scan(&version)
		if err == nil {
			data.Connected = true
			data.Version = version
		} else {
			data.ErrorMsg = fmt.Sprintf("Error querying database version: %v", err)
		}
	} else {
		data.ErrorMsg = "Database connection is not initialized."
	}

	tmpl, err := template.ParseFiles("welcome.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}
