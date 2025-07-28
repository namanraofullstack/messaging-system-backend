package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"messaging-system-backend/internal/database"
)

// CallHuggingFaceSummarizer sends a request to the Hugging Face summarization model
func CallHuggingFaceSummarizer(messages []string) (string, error) {
	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("hugging Face API key not set")
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
		return "", fmt.Errorf("hugging Face API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("hugging Face API error: %d - %s", resp.StatusCode, string(body))
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

// SummarizeGroupMessages retrieves and summarizes the messages in a group
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

// unique removes duplicate strings from a slice
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
