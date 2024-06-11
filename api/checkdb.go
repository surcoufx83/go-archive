package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var DB *sql.DB

func CheckDBConnection(w http.ResponseWriter, r *http.Request) {
	err := DB.Ping()
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Database connection successful!")
}
