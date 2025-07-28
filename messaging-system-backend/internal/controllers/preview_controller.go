package controllers

import (
	"errors"

	"messaging-system-backend/internal/database"
	"messaging-system-backend/internal/models"
)

// GetLatestMessagesFromUsers retrieves the latest messages from users
func GetLatestMessagesFromUsers(userID int) ([]models.MessagePreview, error) {
	rows, err := database.DB.Query(`
		SELECT DISTINCT ON (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id)) 
			id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE sender_id = $1 OR receiver_id = $1
		ORDER BY LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id), created_at DESC
		LIMIT 10
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var previews []models.MessagePreview
	for rows.Next() {
		var msg models.MessagePreview
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		previews = append(previews, msg)
	}
	return previews, nil
}

// GetLatestGroupsWithMessages retrieves the latest groups with messages
func GetLatestGroupsWithMessages(userID int) ([]models.GroupPreview, error) {
	rows, err := database.DB.Query(`
		SELECT g.id, g.name, COALESCE(m.content, '') AS last_message, COALESCE(m.created_at, NOW()) AS last_message_time
		FROM groups g
		INNER JOIN group_members gm ON g.id = gm.group_id
		LEFT JOIN LATERAL (
			SELECT content, created_at 
			FROM group_messages 
			WHERE group_id = g.id 
			ORDER BY created_at DESC LIMIT 1
		) m ON true
		WHERE gm.user_id = $1
		ORDER BY last_message_time DESC
		LIMIT 10
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.GroupPreview
	for rows.Next() {
		var g models.GroupPreview
		if err := rows.Scan(&g.ID, &g.Name, &g.LastMessage, &g.LastMessageTime); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// GetLatestMessages retrieves the latest messages for a specific chat type (DM or group)
func GetLatestMessages(chatType string, chatID int, userID int) ([]models.ChatMessage, error) {
	switch chatType {
	case "dm":
		rows, err := database.DB.Query(`
			SELECT id, sender_id, receiver_id, content, created_at
			FROM messages
			WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
			ORDER BY created_at DESC
			LIMIT 10
		`, userID, chatID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var messages []models.ChatMessage
		for rows.Next() {
			var msg models.ChatMessage
			if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt); err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
		return messages, nil

	case "group":
		rows, err := database.DB.Query(`
			SELECT id, group_id, sender_id, content, created_at
			FROM group_messages
			WHERE group_id = $1
			ORDER BY created_at DESC
			LIMIT 10
		`, chatID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var messages []models.ChatMessage
		for rows.Next() {
			var msg models.ChatMessage
			if err := rows.Scan(&msg.ID, &msg.GroupID, &msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
		return messages, nil

	default:
		return nil, errors.New("invalid chat type")
	}
}
