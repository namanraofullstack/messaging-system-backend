package models

import "time"

// Preview models for messages and groups to be used in the messaging system
type MessagePreview struct {
	ID             int       `json:"id"`
	SenderID       int       `json:"sender_id"`
	ReceiverID     int       `json:"receiver_id"`
	SenderStatus   string    `json:"sender_status"`   // Add this
	ReceiverStatus string    `json:"receiver_status"` // Add this
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}


// GroupPreview models a preview of a group with the last message and its timestamp
type GroupPreview struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	LastMessage    string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
	SenderStatus   string    `json:"sender_status"`
	ReceiverStatus string    `json:"receiver_status"`
}

