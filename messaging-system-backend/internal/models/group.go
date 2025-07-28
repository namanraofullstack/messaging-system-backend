package models

// Preview models for groups to be used in the messaging system
type Group struct {
	ID          int
	Name        string
	Description string
}

// GroupMember models a member of a group with their user ID and admin status
type GroupMember struct {
	UserID  int
	GroupID int
	IsAdmin bool
}

// GroupMessageInput models the input for sending a message to a group
type CreateGroupInput struct {
	Name    string `json:"name"`
	Members []int  `json:"members"`
}

// AddGroupMemberInput models the input for adding a member to a group
type AddGroupMemberInput struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}

// PromoteMemberInput models the input for promoting a member to admin in a group
type PromoteMemberInput struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}

// DemoteMemberInput models the input for demoting an admin to a member in a group
type PromoteOrDemoteInput struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}
