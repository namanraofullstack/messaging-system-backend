package handlers

import (
	"encoding/json"
	"net/http"

	"messaging-system-backend/internal/controllers"
	"messaging-system-backend/internal/models"
	"messaging-system-backend/pkg/utils"
)

// EditDirectMessageHandler handles the request to edit a direct message
func EditDirectMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input models.EditMessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = controllers.EditDirectMessage(input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Message edited successfully"})
}

// EditGroupMessageHandler handles the request to edit a group message
func EditGroupMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input models.EditMessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = controllers.EditGroupMessage(input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Group message edited successfully"})
}
