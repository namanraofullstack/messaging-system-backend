package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"messaging-system-backend/internal/database"
	"messaging-system-backend/internal/models"
	"messaging-system-backend/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	u.Password = string(hash)

	// Insert into DB
	_, err = database.DB.Exec("INSERT INTO Users (username, password) VALUES ($1, $2)", u.Username, u.Password)
	if err != nil {
		http.Error(w, "Error saving user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "registered"})
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get user from DB
	var userID int
	var hashedPwd string
	err := database.DB.QueryRow("SELECT id, password FROM Users WHERE username = $1", creds.Username).Scan(&userID, &hashedPwd)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT including user ID
	token, err := utils.GenerateJWT(userID, creds.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Logout handles user logout and blacklists the token
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract token from Authorization header
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		http.Error(w, "No token provided", http.StatusBadRequest)
		return
	}

	tokenString := strings.TrimPrefix(auth, "Bearer ")

	// Blacklist the token
	err := utils.BlacklistToken(tokenString)
	if err != nil {
		http.Error(w, "Error blacklisting token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully. Token has been blacklisted.",
	})
}

// Protected handles access to protected route
func Protected(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Access granted to protected route.",
	})
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Set sender ID from JWT
	msg.SenderID = userID
	msg.CreatedAt = time.Now()

	_, err = database.DB.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content, created_at) VALUES ($1, $2, $3, $4)",
		msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to send message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent"})
}

func SendGroupMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.GroupMessageInput
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is a member of the group
	var count int
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM group_members 
        WHERE group_id=$1 AND user_id=$2
    `, msg.GroupID, userID).Scan(&count)
	if err != nil || count == 0 {
		http.Error(w, "You are not a member of this group", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec(`
        INSERT INTO group_messages (group_id, sender_id, content) 
        VALUES ($1, $2, $3)
    `, msg.GroupID, userID, msg.Content)
	if err != nil {
		http.Error(w, "Could not send message", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent"})
}

func CreateGroup(input models.CreateGroupInput, creatorID int) (map[string]interface{}, error) {
	// Validate member count
	if len(input.Members) > 25 {
		return nil, errors.New("Group cannot have more than 25 members")
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
		return nil, errors.New("Group creator must be a member of the group")
	}

	// Create group
	var groupID int
	err := database.DB.QueryRow(`
		INSERT INTO groups (name) VALUES ($1) RETURNING id
	`, input.Name).Scan(&groupID)
	if err != nil {
		return nil, errors.New("Error creating group")
	}

	// Add members
	for _, memberID := range input.Members {
		isAdmin := (memberID == creatorID)

		_, err := database.DB.Exec(`
			INSERT INTO group_members (group_id, user_id, is_admin)
			VALUES ($1, $2, $3)
		`, groupID, memberID, isAdmin)
		if err != nil {
			return nil, errors.New("Error adding members to group")
		}
	}

	return map[string]interface{}{
		"message":  "Group created",
		"group_id": groupID,
	}, nil
}

func AddMemberToGroup(input models.AddGroupMemberInput, requesterID int) (map[string]string, error) {
	var isAdmin bool
	err := database.DB.QueryRow(`
        SELECT is_admin FROM group_members
        WHERE group_id = $1 AND user_id = $2
    `, input.GroupID, requesterID).Scan(&isAdmin)

	if err == sql.ErrNoRows || !isAdmin {
		return nil, errors.New("Only admins can add members")
	} else if err != nil {
		return nil, errors.New("Database error checking admin status")
	}

	var count int
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM group_members WHERE group_id = $1
    `, input.GroupID).Scan(&count)

	if err != nil {
		return nil, errors.New("Failed to count members")
	}
	if count >= 25 {
		return nil, errors.New("Group already has 25 members")
	}

	_, err = database.DB.Exec(`
        INSERT INTO group_members (group_id, user_id, is_admin)
        VALUES ($1, $2, false)
    `, input.GroupID, input.UserID)
	if err != nil {
		return nil, errors.New("Error adding member")
	}

	return map[string]string{
		"message": "Member added to group",
	}, nil
}

func PromoteMemberToAdmin(input models.PromoteMemberInput, requesterID int) (map[string]string, error) {
	// Check if requester is an admin
	var isAdmin bool
	err := database.DB.QueryRow(`
        SELECT is_admin FROM group_members 
        WHERE group_id = $1 AND user_id = $2
    `, input.GroupID, requesterID).Scan(&isAdmin)
	if err != nil || !isAdmin {
		return nil, errors.New("Only admins can promote members")
	}

	// Count current admins
	var adminCount int
	err = database.DB.QueryRow(`
        SELECT COUNT(*) FROM group_members 
        WHERE group_id = $1 AND is_admin = true
    `, input.GroupID).Scan(&adminCount)
	if err != nil {
		return nil, errors.New("Failed to check admin count")
	}
	if adminCount >= 2 {
		return nil, errors.New("Group already has 2 admins")
	}

	// Promote the member
	res, err := database.DB.Exec(`
        UPDATE group_members SET is_admin = true 
        WHERE group_id = $1 AND user_id = $2
    `, input.GroupID, input.UserID)
	if err != nil {
		return nil, errors.New("Failed to promote member")
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return nil, errors.New("Member not found in group")
	}

	return map[string]string{
		"message": "Member promoted to admin",
	}, nil
}

func DemoteAdminToMember(input models.PromoteOrDemoteInput, requesterID int) (map[string]string, error) {
	// Step 1: Verify requester is admin
	var isAdmin bool
	err := database.DB.QueryRow(`
		SELECT is_admin FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, requesterID).Scan(&isAdmin)
	if err != nil || !isAdmin {
		return nil, fmt.Errorf("Only admins can demote")
	}

	// Step 2: Check if target user is admin
	var isTargetAdmin bool
	err = database.DB.QueryRow(`
		SELECT is_admin FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, input.UserID).Scan(&isTargetAdmin)
	if err != nil {
		return nil, fmt.Errorf("Target user not found in group")
	}
	if !isTargetAdmin {
		return nil, fmt.Errorf("User is not an admin")
	}

	// Step 3: Prevent self-demotion
	if input.UserID == requesterID {
		return nil, fmt.Errorf("You cannot demote yourself")
	}

	// Step 4: Perform demotion
	_, err = database.DB.Exec(`
		UPDATE group_members 
		SET is_admin = false 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("Failed to demote user")
	}

	return map[string]string{
		"message": "Admin demoted to member",
	}, nil
}

func RemoveMemberFromGroup(input models.PromoteOrDemoteInput, requesterID int) (map[string]string, error) {
	// Step 1: Check if requester is admin
	var isAdmin bool
	err := database.DB.QueryRow(`
		SELECT is_admin FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, requesterID).Scan(&isAdmin)
	if err != nil || !isAdmin {
		return nil, fmt.Errorf("Only admins can remove members")
	}

	// Step 2: Prevent removing self
	if input.UserID == requesterID {
		return nil, fmt.Errorf("You cannot remove yourself")
	}

	// Step 3: Optional - Prevent removing another admin
	var targetIsAdmin bool
	err = database.DB.QueryRow(`
		SELECT is_admin FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, input.UserID).Scan(&targetIsAdmin)
	if err != nil {
		return nil, fmt.Errorf("Target user not found in group")
	}
	if targetIsAdmin {
		return nil, fmt.Errorf("Cannot remove another admin")
	}

	// Step 4: Remove member
	_, err = database.DB.Exec(`
		DELETE FROM group_members 
		WHERE group_id = $1 AND user_id = $2
	`, input.GroupID, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("Failed to remove member")
	}

	return map[string]string{
		"message": "Member removed from group",
	}, nil
}

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
		return nil, errors.New("Invalid chat type")
	}
}

func CallHuggingFaceSummarizer(messages []string) (string, error) {
	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Hugging Face API key not set")
	}

	input := strings.Join(messages, "\n")
	requestBody, err := json.Marshal(map[string]string{
		"inputs": input,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/philschmid/bart-large-cnn-samsum", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Hugging Face API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("Hugging Face API error: %d - %s", resp.StatusCode, string(body))
	}

	var result []map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode Hugging Face response: %v", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("empty response from Hugging Face")
	}

	summary, ok := result[0]["summary_text"]
	if !ok {
		return "", fmt.Errorf("missing summary_text in response")
	}

	return summary, nil
}

func SummarizeGroupMessages(groupID int) (map[string]interface{}, error) {
	// Add context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := database.DB.QueryContext(ctx, `
        SELECT u.username, gm.content
        FROM group_messages gm
        JOIN users u ON gm.sender_id = u.id
        WHERE gm.group_id = $1
        ORDER BY gm.created_at DESC
        LIMIT 20
    `, groupID)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}
	defer rows.Close()

	var messages []string
	var userList []string
	for rows.Next() {
		var username, msg string
		if err := rows.Scan(&username, &msg); err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		messages = append(messages, fmt.Sprintf("%s: %s", username, msg))
		userList = append(userList, username)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	summary, err := CallHuggingFaceSummarizer(messages)
	if err != nil {
		return nil, fmt.Errorf("summarization failed: %v", err)
	}

	return map[string]interface{}{
		"summary": summary,
		"users":   unique(userList),
	}, nil
}

func unique(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, val := range input {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}
