package controllers

import (
	"database/sql"
	"errors"
	"fmt"

	"messaging-system-backend/internal/database"
	"messaging-system-backend/internal/models"
)

// CreateGroup creates a new group with the given input and creator ID
func CreateGroup(input models.CreateGroupInput, creatorID int) (map[string]interface{}, error) {
	// Validate member count
	if len(input.Members) > 25 {
		return nil, errors.New("group cannot have more than 25 members")
	}

	// Ensure creatorID is in members list
	isCreatorPresent := false
	for _, memberID := range input.Members {
		if memberID == creatorID {
			isCreatorPresent = true
			break
		}
	}
	if !isCreatorPresent {
		return nil, errors.New("group creator must be a member of the group")
	}

	// Create group
	var groupID int
	err := database.DB.QueryRow(`
		INSERT INTO groups (name) VALUES ($1) RETURNING id
	`, input.Name).Scan(&groupID)
	if err != nil {
		return nil, errors.New("error creating group")
	}

	// Add members
	for _, memberID := range input.Members {
		isAdmin := (memberID == creatorID)

		_, err := database.DB.Exec(`
			INSERT INTO group_members (group_id, user_id, is_admin)
			VALUES ($1, $2, $3)
		`, groupID, memberID, isAdmin)
		if err != nil {
			return nil, errors.New("error adding members to group")
		}
	}

	return map[string]interface{}{
		"message":  "Group created",
		"group_id": groupID,
	}, nil
}

// AddMemberToGroup adds a member to an existing group
func AddMemberToGroup(input models.AddGroupMemberInput, requesterID int) (map[string]string, error) {
	var isAdmin bool
	err := database.DB.QueryRow(`
        SELECT is_admin FROM group_members
        WHERE group_id = $1 AND user_id = $2
    `, input.GroupID, requesterID).Scan(&isAdmin)

	if err == sql.ErrNoRows || !isAdmin {
		return nil, errors.New("only admins can add members")
	} else if err != nil {
		return nil, errors.New("database error checking admin status")
	}

	var count int
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM group_members WHERE group_id = $1
    `, input.GroupID).Scan(&count)

	if err != nil {
		return nil, errors.New("failed to count members")
	}
	if count >= 25 {
		return nil, errors.New("group already has 25 members")
	}

	_, err = database.DB.Exec(`
        INSERT INTO group_members (group_id, user_id, is_admin)
        VALUES ($1, $2, false)
    `, input.GroupID, input.UserID)
	if err != nil {
		return nil, errors.New("error adding member")
	}

	return map[string]string{
		"message": "Member added to group",
	}, nil
}

// PromoteMemberToAdmin promotes a member to admin in a group
func PromoteMemberToAdmin(input models.PromoteMemberInput, requesterID int) (map[string]string, error) {
	// Check if requester is an admin
	var isAdmin bool
	err := database.DB.QueryRow(`
        SELECT is_admin FROM group_members 
        WHERE group_id = $1 AND user_id = $2
    `, input.GroupID, requesterID).Scan(&isAdmin)
	if err != nil || !isAdmin {
		return nil, errors.New("only admins can promote members")
	}

	// Count current admins
	var adminCount int
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM group_members 
        WHERE group_id = $1 AND is_admin = true
    `, input.GroupID).Scan(&adminCount)
	if err != nil {
		return nil, errors.New("failed to check admin count")
	}
	if adminCount >= 2 {
		return nil, errors.New("group already has 2 admins")
	}

	// Promote the member
	res, err := database.DB.Exec(`
        UPDATE group_members SET is_admin = true 
        WHERE group_id = $1 AND user_id = $2
    `, input.GroupID, input.UserID)
	if err != nil {
		return nil, errors.New("failed to promote member")
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, errors.New("member not found in group")
	}

	return map[string]string{
		"message": "Member promoted to admin",
	}, nil
}

// DemoteAdminToMember demotes an admin back to a regular member
func DemoteAdminToMember(input models.PromoteOrDemoteInput, requesterID int) (map[string]string, error) {
	// Step 1: Verify requester is admin
	var isAdmin bool
	err := database.DB.QueryRow(`
		SELECT is_admin FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, requesterID).Scan(&isAdmin)
	if err != nil || !isAdmin {
		return nil, fmt.Errorf("only admins can demote")
	}

	// Step 2: Check if target user is admin
	var isTargetAdmin bool
	err = database.DB.QueryRow(`
		SELECT is_admin FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, input.UserID).Scan(&isTargetAdmin)
	if err != nil {
		return nil, fmt.Errorf("target user not found in group")
	}
	if !isTargetAdmin {
		return nil, fmt.Errorf("user is not an admin")
	}

	// Step 3: Prevent self-demotion
	if input.UserID == requesterID {
		return nil, fmt.Errorf("you cannot demote yourself")
	}

	// Step 4: Perform demotion
	_, err = database.DB.Exec(`
		UPDATE group_members 
		SET is_admin = false 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to demote user")
	}

	return map[string]string{
		"message": "Admin demoted to member",
	}, nil
}

// RemoveMemberFromGroup removes a member from a group
func RemoveMemberFromGroup(input models.PromoteOrDemoteInput, requesterID int) (map[string]string, error) {
	// Step 1: Check if requester is admin
	var isAdmin bool
	err := database.DB.QueryRow(`
		SELECT is_admin FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, requesterID).Scan(&isAdmin)
	if err != nil || !isAdmin {
		return nil, fmt.Errorf("only admins can remove members")
	}

	// Step 2: Prevent removing self
	if input.UserID == requesterID {
		return nil, fmt.Errorf("you cannot remove yourself")
	}

	// Step 3: Optional - Prevent removing another admin
	var targetIsAdmin bool
	err = database.DB.QueryRow(`
		SELECT is_admin FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, input.UserID).Scan(&targetIsAdmin)
	if err != nil {
		return nil, fmt.Errorf("target user not found in group")
	}
	if targetIsAdmin {
		return nil, fmt.Errorf("cannot remove another admin")
	}

	// Step 4: Remove member
	_, err = database.DB.Exec(`
		DELETE FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove member")
	}

	return map[string]string{
		"message": "Member removed from group",
	}, nil
}
