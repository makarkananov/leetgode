package postgres

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(userID int64) (*User, error) {
	query := "SELECT id, current_state FROM users WHERE id = $1"
	row := r.db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.id, &user.CurrentState)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) AddUser(id int64, currentState string) error {
	query := `
		INSERT INTO users (id, current_state) 
		VALUES ($1, $2) 
	`
	_, err := r.db.Exec(query, id, currentState)
	if err != nil {
		return fmt.Errorf("failed to add/update user in the database: %w", err)
	}

	return nil
}

func (r *UserRepository) UpdateUserState(userID int64, currentState string) error {
	query := "UPDATE users SET current_state = $1 WHERE id = $2"
	_, err := r.db.Exec(query, currentState, userID)
	if err != nil {
		return fmt.Errorf("failed to update user state in the database: %w", err)
	}

	return nil
}

type User struct {
	id           string
	CurrentState string
}
