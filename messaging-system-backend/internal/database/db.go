package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

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
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				fmt.Println("✅ Connected to DB successfully!")

				// Create Users table if not exists
				createTable := `
					CREATE TABLE IF NOT EXISTS Users (
						id SERIAL PRIMARY KEY,
						username VARCHAR(100) UNIQUE NOT NULL,
						password TEXT NOT NULL,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
					);
				`
				if _, err := DB.Exec(createTable); err != nil {
					return fmt.Errorf("❌ Failed to create Users table: %w", err)
				}
				fmt.Println("📦 Users table ensured.")
				return nil
			}
		}
		fmt.Println("⏳ Waiting for DB to be ready...")
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("❌ Error pinging DB: %w", err)
}
