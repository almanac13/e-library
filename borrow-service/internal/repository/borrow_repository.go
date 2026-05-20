package repository

import (
	"borrow-service/internal/model"
	"database/sql"
	"time"
)

type BorrowRepository struct {
	db *sql.DB
}

func NewBorrowRepository(db *sql.DB) *BorrowRepository {
	return &BorrowRepository{db: db}
}

func (r *BorrowRepository) Create(borrow *model.Borrow) error {
	query := `
		INSERT INTO borrows 
		(id, user_id, book_id, borrow_date, due_date, return_date, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := r.db.Exec(
		query,
		borrow.ID,
		borrow.UserID,
		borrow.BookID,
		borrow.BorrowDate,
		borrow.DueDate,
		borrow.ReturnDate,
		borrow.Status,
		borrow.CreatedAt,
		borrow.UpdatedAt,
	)

	return err
}

func (r *BorrowRepository) GetAll() ([]model.Borrow, error) {
	query := `
		SELECT id, user_id, book_id, borrow_date, due_date, return_date, status, created_at, updated_at
		FROM borrows
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrows []model.Borrow

	for rows.Next() {
		var borrow model.Borrow

		err := rows.Scan(
			&borrow.ID,
			&borrow.UserID,
			&borrow.BookID,
			&borrow.BorrowDate,
			&borrow.DueDate,
			&borrow.ReturnDate,
			&borrow.Status,
			&borrow.CreatedAt,
			&borrow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		borrows = append(borrows, borrow)
	}

	return borrows, rows.Err()
}
func (r *BorrowRepository) GetByID(id string) (*model.Borrow, error) {
	query := `
		SELECT id, user_id, book_id, borrow_date, due_date, return_date, status, created_at, updated_at
		FROM borrows
		WHERE id = $1
	`

	var borrow model.Borrow

	err := r.db.QueryRow(query, id).Scan(
		&borrow.ID,
		&borrow.UserID,
		&borrow.BookID,
		&borrow.BorrowDate,
		&borrow.DueDate,
		&borrow.ReturnDate,
		&borrow.Status,
		&borrow.CreatedAt,
		&borrow.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &borrow, nil
}

func (r *BorrowRepository) ReturnBorrow(id string, returnDate time.Time) (*model.Borrow, error) {
	query := `
		UPDATE borrows
		SET return_date = $1, status = $2, updated_at = $3
		WHERE id = $4 AND status = $5
		RETURNING id, user_id, book_id, borrow_date, due_date, return_date, status, created_at, updated_at
	`

	var borrow model.Borrow
	now := time.Now()

	err := r.db.QueryRow(
		query,
		returnDate,
		model.BorrowStatusReturned,
		now,
		id,
		model.BorrowStatusActive,
	).Scan(
		&borrow.ID,
		&borrow.UserID,
		&borrow.BookID,
		&borrow.BorrowDate,
		&borrow.DueDate,
		&borrow.ReturnDate,
		&borrow.Status,
		&borrow.CreatedAt,
		&borrow.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &borrow, nil
}
