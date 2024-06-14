package api

import (
	"database/sql"
	"fmt"
	"go-archive/utils"
	"net/http"
)

var DB *sql.DB

func CheckDBConnection(w http.ResponseWriter, r *http.Request) {
	err := DB.Ping()
	if err != nil {
		utils.LogErrorAndReturnCode(w, "Database connection failed", err, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Database connection successful!")
}
