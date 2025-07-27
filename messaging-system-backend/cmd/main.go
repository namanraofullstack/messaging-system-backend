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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
