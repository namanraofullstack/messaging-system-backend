package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"messaging-system-backend/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// BlacklistToken adds a token to the Redis blacklist
func BlacklistToken(tokenString string) error {
	// Parse token to get expiration time
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	// Get expiration time
	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("invalid expiration time")
	}

	expirationTime := time.Unix(int64(exp), 0)

	// Calculate TTL (time until token expires)
	ttl := time.Until(expirationTime)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	// Store token in Redis with TTL
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", tokenString)

	return database.RedisClient.Set(ctx, key, "blacklisted", ttl).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted
func IsTokenBlacklisted(tokenString string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", tokenString)

	result := database.RedisClient.Get(ctx, key)
	return result.Err() == nil // If no error, key exists (token is blacklisted)
}
