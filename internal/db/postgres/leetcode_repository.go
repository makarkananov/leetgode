package postgres

import (
	"database/sql"
	"fmt"
)

type LeetcodeRepository struct {
	db *sql.DB
}

func NewLeetcodeRepository(db *sql.DB) *LeetcodeRepository {
	return &LeetcodeRepository{db: db}
}

func (r *LeetcodeRepository) AddLeetodeUser(userID int64, LeetcodeUsername string) error {
	query := "INSERT INTO Leetcode (user_id, Leetcode_username) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET Leetcode_username = $2"
	_, err := r.db.Exec(query, userID, LeetcodeUsername)
	if err != nil {
		return fmt.Errorf("failed to add Leetcode user in the database: %w", err)
	}

	return nil
}

func (r *LeetcodeRepository) GetLeetcodeUser(userID int64) (string, error) {
	query := "SELECT Leetcode_username FROM Leetcode WHERE user_id = $1"
	var LeetcodeUsername string
	err := r.db.QueryRow(query, userID).Scan(&LeetcodeUsername)
	if err != nil {
		return "", fmt.Errorf("failed to get Leetcode username: %w", err)
	}

	return LeetcodeUsername, nil
}

func (r *LeetcodeRepository) RemoveLeetcodeUser(userID int64) error {
	query := "DELETE FROM Leetcode WHERE user_id = $1"
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to remove Leetcode user: %w", err)
	}

	return nil
}
