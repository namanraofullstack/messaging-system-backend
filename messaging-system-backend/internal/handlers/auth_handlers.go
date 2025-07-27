package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"messaging-system-backend/internal/controllers"
	"messaging-system-backend/internal/models"
	"messaging-system-backend/pkg/utils"
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

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.SendMessage(w, r)
}

func GroupMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.SendGroupMessage(w, r)
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var input models.CreateGroupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := controllers.CreateGroup(input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func AddMemberToGroup(w http.ResponseWriter, r *http.Request) {
	var input models.AddGroupMemberInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := controllers.AddMemberToGroup(input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func PromoteMemberToAdmin(w http.ResponseWriter, r *http.Request) {
	var input models.PromoteMemberInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	requesterID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := controllers.PromoteMemberToAdmin(input, requesterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func DemoteAdminToMember(w http.ResponseWriter, r *http.Request) {
	var input models.PromoteOrDemoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	requesterID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := controllers.DemoteAdminToMember(input, requesterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func RemoveMemberFromGroup(w http.ResponseWriter, r *http.Request) {
	var input models.PromoteOrDemoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	requesterID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	resp, err := controllers.RemoveMemberFromGroup(input, requesterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

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
