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

	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		receiver_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
	CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);

	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);

	CREATE TABLE IF NOT EXISTS group_members (
		id SERIAL PRIMARY KEY,
		group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		is_admin BOOLEAN DEFAULT FALSE,
		UNIQUE(group_id, user_id)
	);

	CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
	CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members(user_id);

	CREATE TABLE IF NOT EXISTS group_messages (
		id SERIAL PRIMARY KEY,
		group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		sender_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_group_messages_group_id ON group_messages(group_id);
	CREATE INDEX IF NOT EXISTS idx_group_messages_sender_id ON group_messages(sender_id);

	ALTER TABLE messages ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
	ALTER TABLE group_messages ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'Available' CHECK (char_length(status) <= 1000);
	ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;




	`

	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("❌ Failed to run schema migrations: %w", err)
	}
	fmt.Println("📦 DB schema ensured.")
	return nil
}
