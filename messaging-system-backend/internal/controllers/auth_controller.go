package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"messaging-system-backend/internal/database"
	"messaging-system-backend/internal/models"
	"messaging-system-backend/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	u.Password = string(hash)

	// Insert into DB
	_, err = database.DB.Exec("INSERT INTO Users (username, password) VALUES ($1, $2)", u.Username, u.Password)
	if err != nil {
		http.Error(w, "Error saving user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "registered"})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get user from DB
	var userID int
	var hashedPwd string
	err := database.DB.QueryRow("SELECT id, password FROM Users WHERE username = $1", creds.Username).Scan(&userID, &hashedPwd)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT including user ID
	token, err := utils.GenerateJWT(userID, creds.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Logout handles user logout and blacklists the token
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from Authorization header
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		http.Error(w, "No token provided", http.StatusBadRequest)
		return
	}

	tokenString := strings.TrimPrefix(auth, "Bearer ")

	// Blacklist the token
	err := utils.BlacklistToken(tokenString)
	if err != nil {
		http.Error(w, "Error blacklisting token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully. Token has been blacklisted.",
	})
}

// Protected handles access to protected route
func Protected(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Access granted to protected route.",
	})
}
