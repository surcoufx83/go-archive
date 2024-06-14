package utils

import (
	"log"
	"net/http"
)

// LogErrorAndReturnCode logs the error message and sends an HTTP error response
func LogErrorAndReturnCode(w http.ResponseWriter, message string, err error, code int) {
	log.Printf("HTTP %d: %s - %v", code, message, err)
	http.Error(w, "", code)
}
