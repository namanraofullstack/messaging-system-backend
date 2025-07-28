package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection using environment variables for configuration.
func InitDB() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	for i := 0; i < 10; i++ {
		db, err := sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				DB = db
				fmt.Println("✅ Connected to DB successfully!")
				return nil
			}
		}
		fmt.Printf("⏳ DB not ready (%v). Retrying...\n", err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("❌ Could not connect to DB: %w", err)
}
