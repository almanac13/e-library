package repository

import (
	"database/sql"
	"user-service/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user model.User) error {
	query := `
		INSERT INTO users (id, name, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Password, user.Role)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at
		FROM users
		WHERE id = $1
	`

	var user model.User

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	query := `
		SELECT id, name, email, role, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
