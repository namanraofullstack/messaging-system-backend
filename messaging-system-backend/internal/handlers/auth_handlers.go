package handlers

import (
	"net/http"

	"messaging-system-backend/internal/controllers"
)

// RegisterHandler handles POST /register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.Register(w, r)
}

// LoginHandler handles POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.Login(w, r)
}

// LogoutHandler handles GET /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.Logout(w, r)
}

// ProtectedHandler handles GET /protected
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.Protected(w, r)
}
