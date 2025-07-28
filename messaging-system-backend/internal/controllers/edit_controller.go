package controllers

import (
	"fmt"
	"time"

	"messaging-system-backend/internal/database"
	"messaging-system-backend/internal/models"
)

// EditDirectMessage allows a user to edit a direct message
func EditDirectMessage(input models.EditMessageInput, userID int) error {
	var existing models.Message

	err := database.DB.QueryRow(`
		SELECT id, sender_id, content, updated_at, created_at
		FROM messages
		WHERE id = $1`, input.MessageID).Scan(
		&existing.ID,
		&existing.SenderID,
		&existing.Content,
		&existing.UpdatedAt,
		&existing.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("message not found")
	}

	if existing.SenderID != userID {
		return fmt.Errorf("you can only edit your own messages")
	}

	if time.Since(existing.CreatedAt) > time.Hour {
		return fmt.Errorf("message can no longer be edited")
	}

	inputTime := input.LastUpdatedAt.UTC().Truncate(time.Second)
	dbTime := existing.UpdatedAt.UTC().Truncate(time.Second)

	if !inputTime.Equal(dbTime) {
		return fmt.Errorf("conflict detected, please refresh the message")
	}

	_, err = database.DB.Exec(`
		UPDATE messages
		SET content = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, input.NewContent, input.MessageID)

	return err
}

// EditGroupMessage allows a user to edit a message in a group chat
func EditGroupMessage(input models.EditMessageInput, userID int) error {
	var existing struct {
		ID        int
		SenderID  int
		Content   string
		UpdatedAt time.Time
		CreatedAt time.Time
	}

	err := database.DB.QueryRow(`
		SELECT id, sender_id, content, updated_at, created_at
		FROM group_messages
		WHERE id = $1`, input.MessageID).Scan(
		&existing.ID,
		&existing.SenderID,
		&existing.Content,
		&existing.UpdatedAt,
		&existing.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("group message not found")
	}

	if existing.SenderID != userID {
		return fmt.Errorf("you can only edit your own messages")
	}

	if time.Since(existing.CreatedAt) > time.Hour {
		return fmt.Errorf("group message can no longer be edited")
	}

	inputTime := input.LastUpdatedAt.UTC().Truncate(time.Second)
	dbTime := existing.UpdatedAt.UTC().Truncate(time.Second)

	if !inputTime.Equal(dbTime) {
		return fmt.Errorf("conflict detected, please refresh the message")
	}

	_, err = database.DB.Exec(`
		UPDATE group_messages
		SET content = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, input.NewContent, input.MessageID)

	return err
}
