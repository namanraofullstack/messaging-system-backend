package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"messaging-system-backend/internal/controllers"
	"messaging-system-backend/pkg/utils"
)

// ViewLatestUserChats handles GET /users/chats
func ViewLatestUserChats(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chats, err := controllers.GetLatestMessagesFromUsers(userID)
	if err != nil {
		http.Error(w, "Could not fetch chats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chats)
}

// ViewLatestGroups handles GET /groups/latest
func ViewLatestGroups(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groups, err := controllers.GetLatestGroupsWithMessages(userID)
	if err != nil {
		http.Error(w, "Could not fetch groups", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

// ViewChatMessages handles GET /chats/messages
func ViewChatMessages(w http.ResponseWriter, r *http.Request) {
	chatType := r.URL.Query().Get("type") // "dm" or "group"
	chatIDStr := r.URL.Query().Get("id")

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	msgs, err := controllers.GetLatestMessages(chatType, chatID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(msgs)
}
