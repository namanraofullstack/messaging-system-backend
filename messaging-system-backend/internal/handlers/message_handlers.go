package handlers

import (
	"net/http"

	"messaging-system-backend/internal/controllers"
)

// SendMessageHandler handles POST /send
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.SendMessage(w, r)
}

// GroupMessageHandler handles POST /groups/messages
func GroupMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.SendGroupMessage(w, r)
}
