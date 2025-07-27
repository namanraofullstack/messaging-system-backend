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
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	// Initialize Redis
	err = database.InitRedis()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	// Public routes
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Protected route using JWT middleware
	http.Handle("/protected", middleware.JWTMiddleware(http.HandlerFunc(handlers.ProtectedHandler)))

	// Message sending route
	http.HandleFunc("/send", handlers.SendMessageHandler)
	// Group message route
	http.HandleFunc("/group/message", handlers.GroupMessageHandler)

	http.HandleFunc("/group/create", handlers.CreateGroup)

	http.HandleFunc("/group/add-member", handlers.AddMemberToGroup)
	http.HandleFunc("/group/remove-member", handlers.RemoveMemberFromGroup)

	http.HandleFunc("/group/promote", handlers.PromoteMemberToAdmin)
	http.HandleFunc("/group/demote", handlers.DemoteAdminToMember)

	http.HandleFunc("/chats/latest-dm-previews", handlers.ViewLatestUserChats)
	http.HandleFunc("/chats/latest-group-previews", handlers.ViewLatestGroups)
	http.HandleFunc("/chats/messages", handlers.ViewChatMessages)

	http.HandleFunc("/groups/summary", handlers.GetGroupSummary)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
