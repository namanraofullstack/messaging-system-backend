package database

import "fmt"

// EnsureTables checks and creates necessary tables in the database.
func EnsureTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		receiver_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	CREATE TABLE IF NOT EXISTS group_members (
		id SERIAL PRIMARY KEY,
		group_id INT REFERENCES groups(id) ON DELETE CASCADE,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		is_admin BOOLEAN DEFAULT FALSE,
		UNIQUE(group_id, user_id)
	);
	CREATE TABLE IF NOT EXISTS group_messages (
		id SERIAL PRIMARY KEY,
		group_id INT REFERENCES groups(id) ON DELETE CASCADE,
		sender_id INT REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("❌ Failed to run schema migrations: %w", err)
	}
	fmt.Println("📦 DB schema ensured.")
	return nil
}
