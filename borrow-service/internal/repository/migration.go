package repository

import "database/sql"

func RunMigration(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS borrows(
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		book_id UUID NOT NULL,
		borrow_date TIMESTAMP NOT NULL,
		due_date TIMESTAMP NOT NULL,
		return_date TIMESTAMP,
		status VARCHAR(30),
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	)
	`

	_, err := db.Exec(query)
	return err
}
