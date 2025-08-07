package controllers

import (
	"database/sql"
	"errors"
	"time"

	"messaging-system-backend/internal/database"
)

// GetUserStatus fetches the status of a user by ID.
func GetUserStatus(userID int) (string, error) {
	var status string
	err := database.DB.QueryRow("SELECT status FROM users WHERE id = $1", userID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return status, nil
}

// UpdateUserStatus attempts to update the user’s status using optimistic concurrency.
func UpdateUserStatus(userID int, newStatus string, lastUpdatedAt time.Time) error {
	result, err := database.DB.Exec(`
		UPDATE users
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND updated_at = $3
	`, newStatus, userID, lastUpdatedAt)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("status update conflict: data was modified by another process")
	}

	return nil
}
