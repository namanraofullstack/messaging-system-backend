package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"messaging-system-backend/internal/controllers"
	"messaging-system-backend/internal/database"
	"messaging-system-backend/pkg/utils"
)

// GetGroupSummary handles GET /groups/summary
func GetGroupSummary(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")

	groupIDStr := r.URL.Query().Get("group_id")
	if groupIDStr == "" {
		http.Error(w, `{"error": "group_id parameter is required"}`, http.StatusBadRequest)
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid group_id format"}`, http.StatusBadRequest)
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Check membership with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err = database.DB.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM group_members 
            WHERE group_id = $1 AND user_id = $2
        )
    `, groupID, userID).Scan(&exists)

	if err != nil {
		log.Printf("Database error checking membership: %v", err)
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, `{"error": "You are not a member of this group"}`, http.StatusForbidden)
		return
	}

	summaryData, err := controllers.SummarizeGroupMessages(groupID)
	if err != nil {
		log.Printf("Summary generation error: %v", err)
		http.Error(w, `{"error": "Failed to generate summary"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(summaryData); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}
