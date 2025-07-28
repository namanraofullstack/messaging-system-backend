package models

import "time"

// Message models a message in the messaging system
type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// GroupMessageInput models the input for sending a message to a group
type GroupMessageInput struct {
	GroupID int    `json:"group_id"`
	Content string `json:"content"`
}

// ChatMessage models a message in a chat, which can be sent to a user or a group
type ChatMessage struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id,omitempty"`
	ReceiverID int       `json:"receiver_id,omitempty"`
	GroupID    int       `json:"group_id,omitempty"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
