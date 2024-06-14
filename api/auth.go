package api

import (
	"database/sql"
	"encoding/json"
	"go-archive/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func NextcloudOAuth2(w http.ResponseWriter, r *http.Request) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	clientID := os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	redirectURI := os.Getenv("OAUTH2_REDIRECT_URI")
	tokenURL := os.Getenv("OAUTH2_TOKEN_URL")

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		utils.LogErrorAndReturnCode(w, "Invalid form data", err, http.StatusBadRequest)
		return
	}

	code := r.PostFormValue("archauth_oauth2_code")
	state := r.PostFormValue("archauth_oauth2_state")

	if code == "" || state == "" {
		utils.LogErrorAndReturnCode(w, "Missing code or state", err, http.StatusBadRequest)
		return
	}

	// Request token from Nextcloud
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", redirectURI)
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)

	resp, err := http.PostForm(tokenURL, form)
	if err != nil {
		utils.LogErrorAndReturnCode(w, "Error requesting Nextcloud token", err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.LogErrorAndReturnCode(w, "Error requesting Nextcloud token", err, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.LogErrorAndReturnCode(w, "Error reading Nextcloud response", err, http.StatusInternalServerError)
		return
	}

	var tokenData map[string]interface{}
	err = json.Unmarshal(body, &tokenData)
	if err != nil {
		utils.LogErrorAndReturnCode(w, "Error parsing Nextcloud response", err, http.StatusInternalServerError)
		return
	}

	// Extract token data
	accessToken, _ := tokenData["access_token"].(string)
	refreshToken, _ := tokenData["refresh_token"].(string)
	expiresIn, _ := tokenData["expires_in"].(float64)

	// Dummy user info extraction
	userInfo := map[string]interface{}{"username": "dummyuser"} // Replace with actual user info extraction

	username, _ := userInfo["username"].(string)

	// Ensure user exists
	var userID int
	err = DB.QueryRow("SELECT id FROM users WHERE loginname = ?", username).Scan(&userID)
	if err == sql.ErrNoRows {
		// Create user if not exists
		result, err := DB.Exec("INSERT INTO users (loginname, password, email, enabled) VALUES (?, ?, ?, ?)", username, "", "", 1)
		if err != nil {
			utils.LogErrorAndReturnCode(w, "Error creating user", err, http.StatusInternalServerError)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			utils.LogErrorAndReturnCode(w, "Error getting last insert id", err, http.StatusInternalServerError)
			return
		}
		userID = int(id)
	} else if err != nil {
		utils.LogErrorAndReturnCode(w, "Error querying user", err, http.StatusInternalServerError)
		return
	}

	// Store token data
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	_, err = DB.Exec("INSERT INTO user_tokens (userid, page, access_token, refresh_token, expires_in, valid) VALUES (?, ?, ?, ?, ?, ?)", userID, "oauth2", accessToken, refreshToken, int(expiresIn), expiresAt)
	if err != nil {
		utils.LogErrorAndReturnCode(w, "Error storing token", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Authentication successful",
		"user_id": userID,
	})
}
