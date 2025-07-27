package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Message struct {
	ID         int
	SenderID   int
	ReceiverID *int
	GroupID    *int
	Content    string
	CreatedAt  string
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
