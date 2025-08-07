package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"messaging-system-backend/internal/controllers"
	"messaging-system-backend/pkg/utils"
)

type SetStatusRequest struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GET /status?id=123 — Requires valid token
func GetUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	_, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	status, err := controllers.GetUserStatus(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	resp := map[string]string{"status": status}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// POST /status — Only sets your own status
func SetUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SetStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, "Status cannot be empty", http.StatusBadRequest)
		return
	}

	if len(req.Status) > 1000 {
		http.Error(w, "Status cannot exceed 1000 characters", http.StatusBadRequest)
		return
	}

	// 💡 Update with optimistic concurrency check
	if err := controllers.UpdateUserStatus(userID, req.Status, req.UpdatedAt); err != nil {
		if err.Error() == "status update conflict: data was modified by another process" {
			http.Error(w, "Conflict: status was recently updated", http.StatusConflict)
		} else {
			http.Error(w, "Failed to update status", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"status updated"}`))
}
