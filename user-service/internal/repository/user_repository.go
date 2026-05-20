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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUserName(id, name string) error {
	query := `
		UPDATE users
		SET name = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, name, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(result)
}

func (r *UserRepository) UpdateUserRole(id, role string) error {
	query := `
		UPDATE users
		SET role = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, role, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(result)
}

func (r *UserRepository) ChangePassword(id, password string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, password, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(result)
}

func (r *UserRepository) DeleteUser(id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(result)
}

func (r *UserRepository) UserExists(id string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE id = $1
		)
	`

	var exists bool

	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) CountUsers() (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
	`

	var count int

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) GetUsersByRole(role string) ([]model.User, error) {
	query := `
		SELECT id, name, email, role, created_at
		FROM users
		WHERE role = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, role)
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func checkRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
