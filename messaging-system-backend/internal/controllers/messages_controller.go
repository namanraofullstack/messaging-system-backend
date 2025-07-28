package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"messaging-system-backend/internal/database"
	"messaging-system-backend/internal/models"
	"messaging-system-backend/pkg/utils"
)

// SendMessage handles sending a message to a user or group
func SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Set sender ID from JWT
	msg.SenderID = userID
	msg.CreatedAt = time.Now()

	_, err = database.DB.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content, created_at) VALUES ($1, $2, $3, $4)",
		msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent"})
}

// SendGroupMessage handles sending a message to a group
func SendGroupMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.GroupMessageInput
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is a member of the group
	var count int
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM group_members 
        WHERE group_id=$1 AND user_id=$2
    `, msg.GroupID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "You are not a member of this group", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec(`
        INSERT INTO group_messages (group_id, sender_id, content) 
        VALUES ($1, $2, $3)
    `, msg.GroupID, userID, msg.Content)
	if err != nil {
		http.Error(w, "Could not send message", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent"})
}
