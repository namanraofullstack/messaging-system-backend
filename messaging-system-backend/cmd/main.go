package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"messaging-system-backend/internal/database"
	"messaging-system-backend/internal/handlers"
	"messaging-system-backend/internal/middleware"
)

func main() {
	// Initialize database
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	if err := database.EnsureTables(); err != nil {
		log.Fatal(err)
	}

	// Initialize Redis
	err = database.InitRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Aunthentication routes
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Protected route using JWT middleware
	http.Handle("/protected", middleware.JWTMiddleware(http.HandlerFunc(handlers.ProtectedHandler)))

	// Message sending routes
	http.HandleFunc("/send", handlers.SendMessageHandler)
	http.HandleFunc("/group/message", handlers.GroupMessageHandler)

	//Group management routes
	http.HandleFunc("/group/create", handlers.CreateGroup)
	http.HandleFunc("/group/add-member", handlers.AddMemberToGroup)
	http.HandleFunc("/group/remove-member", handlers.RemoveMemberFromGroup)
	http.HandleFunc("/group/promote", handlers.PromoteMemberToAdmin)
	http.HandleFunc("/group/demote", handlers.DemoteAdminToMember)

	//Chat previews and messages
	http.HandleFunc("/chats/latest-dm-previews", handlers.ViewLatestUserChats)
	http.HandleFunc("/chats/latest-group-previews", handlers.ViewLatestGroups)
	http.HandleFunc("/chats/messages", handlers.ViewChatMessages)

	//Group messages summary
	http.HandleFunc("/groups/summary", handlers.GetGroupSummary)

	// Edit message routes
	http.HandleFunc("/edit/direct", handlers.EditDirectMessageHandler)
	http.HandleFunc("/edit/group", handlers.EditGroupMessageHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
