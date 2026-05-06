package repository

import (
	"database/sql"
	"errors"

	"github.com/almanac13/e-library/book-service/internal/model"
)

// BookRepository is responsible only for database operations with books.
// Handler should not write SQL directly.
type BookRepository struct {
	db *sql.DB
}

// NewBookRepository is a constructor.
// It creates BookRepository and gives it database connection.
func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

// GetAll returns all books from database.
func (r *BookRepository) GetAll() ([]model.Book, error) {
	// Query is used when SQL returns multiple rows.
	rows, err := r.db.Query(`
		SELECT id, title, author, category, available, created_at
		FROM books
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}

	// Important: always close rows after using them.
	// This prevents memory/resource leaks.
	defer rows.Close()

	books := make([]model.Book, 0)

	// rows.Next moves through every row from SQL result.
	for rows.Next() {
		var book model.Book

		// Scan copies SQL columns into Go struct fields.
		// Order must be the same as in SELECT.
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Category,
			&book.Available,
			&book.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	// rows.Err checks if error happened during iteration.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// GetByID returns one book by ID.
func (r *BookRepository) GetByID(id int) (*model.Book, error) {
	var book model.Book

	// QueryRow is used when SQL should return only one row.
	err := r.db.QueryRow(`
		SELECT id, title, author, category, available, created_at
		FROM books
		WHERE id = $1
	`, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Category,
		&book.Available,
		&book.CreatedAt,
	)

	// sql.ErrNoRows means book with this ID does not exist.
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &book, nil
}

// Create inserts new book into database.
func (r *BookRepository) Create(req model.CreateBookRequest) (*model.Book, error) {
	var book model.Book

	// RETURNING allows PostgreSQL to return created row immediately.
	// Without RETURNING, we would need another SELECT query.
	err := r.db.QueryRow(`
		INSERT INTO books (title, author, category, available)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, author, category, available, created_at
	`,
		req.Title,
		req.Author,
		req.Category,
		req.Available,
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Category,
		&book.Available,
		&book.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

// Update changes existing book by ID.
func (r *BookRepository) Update(id int, req model.UpdateBookRequest) (*model.Book, error) {
	var book model.Book

	// UPDATE changes row.
	// RETURNING returns updated row.
	err := r.db.QueryRow(`
		UPDATE books
		SET title = $1,
			author = $2,
			category = $3,
			available = $4
		WHERE id = $5
		RETURNING id, title, author, category, available, created_at
	`,
		req.Title,
		req.Author,
		req.Category,
		req.Available,
		id,
	).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Category,
		&book.Available,
		&book.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &book, nil
}

// Delete removes book from database by ID.
func (r *BookRepository) Delete(id int) error {
	// Exec is used for SQL commands that don't need returned rows.
	result, err := r.db.Exec(`
		DELETE FROM books
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	// RowsAffected shows how many rows were deleted.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If 0 rows deleted, it means book was not found.
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindByAuthor searches books by author name.
func (r *BookRepository) FindByAuthor(author string) ([]model.Book, error) {
	rows, err := r.db.Query(`
		SELECT id, title, author, category, available, created_at
		FROM books
		WHERE author ILIKE '%' || $1 || '%'
		ORDER BY id
	`, author)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := make([]model.Book, 0)

	for rows.Next() {
		var book model.Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Category,
			&book.Available,
			&book.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
