package controllers_test

import (
	"sync"
	"testing"
	"time"

	"messaging-system-backend/internal/controllers"
	"messaging-system-backend/internal/database"
)

func TestConcurrentUserStatusUpdate(t *testing.T) {
	// Setup: Insert test user
	_, err := database.DB.Exec(`INSERT INTO users (id, username, password) 
		VALUES (9999, 'testuser_concurrency', 'hashed') 
		ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	userID := 9999

	// Get initial updated_at timestamp
	var lastUpdatedAt time.Time
	err = database.DB.QueryRow(`SELECT updated_at FROM users WHERE id = $1`, userID).Scan(&lastUpdatedAt)
	if err != nil {
		t.Fatalf("Failed to get initial updated_at: %v", err)
	}

	statuses := []string{"Available", "Busy", "Away", "Offline", "Sleeping"}
	successCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, status := range statuses {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			err := controllers.UpdateUserStatus(userID, s, lastUpdatedAt)
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				successCount++
				t.Logf("Success: status updated to %s", s)
			} else {
				t.Logf(" Conflict: failed to update to %s (%v)", s, err)
			}
		}(status)
	}

	wg.Wait()

	if successCount != 1 {
		t.Errorf("Expected only 1 successful update, but got %d", successCount)
	}
}
