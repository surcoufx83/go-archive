package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-archive/api"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	api.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Defer closing the database connection until the main function exits
	defer func() {
		if err := api.DB.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/api/helloworld", api.HelloWorld).Methods("GET")
	r.HandleFunc("/api/checkdb", api.CheckDBConnection).Methods("GET")
	r.HandleFunc("/api/file/{id:[0-9]+}", api.GetFile).Methods("GET")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
