package handlers

import (
	"encoding/json"
	"net/http"

	"messaging-system-backend/internal/controllers"
	"messaging-system-backend/internal/models"
	"messaging-system-backend/pkg/utils"
)

// CreateGroup handles POST /groups
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

// AddMemberToGroup handles POST /groups/add-member
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

// PromoteMemberToAdmin handles POST /groups/promote-member
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

// DemoteAdminToMember handles POST /groups/demote-admin
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

// RemoveMemberFromGroup handles POST /groups/remove-member
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
