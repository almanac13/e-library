package repository

import "database/sql"

func RunMigration(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS books(
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		category TEXT,
		available TEXT DEFAULT 'Yes',
		created_at TIMESTAMP DEFAULT NOW()
	)
	`

	_, err := db.Exec(query)
	return err
}
