package repository

import "database/sql"

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := db.Exec(query)
	return err
}
