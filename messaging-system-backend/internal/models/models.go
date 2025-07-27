package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type Group struct {
	ID          int
	Name        string
	Description string
}

type GroupMember struct {
	UserID  int
	GroupID int
	IsAdmin bool
}

type GroupMessageInput struct {
	GroupID int    `json:"group_id"`
	Content string `json:"content"`
}

type CreateGroupInput struct {
	Name    string `json:"name"`
	Members []int  `json:"members"` // Include user_ids to add (must include creator)
}

type AddGroupMemberInput struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}

type PromoteMemberInput struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"` // user to promote
}

type PromoteOrDemoteInput struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}
