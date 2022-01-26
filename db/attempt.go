package db

import (
	"github.com/anna-osipova/go-wordle/models"
)

func (db Database) CreateAttempt(attempt *models.Attempt) error {
	var id int
	var createdAt string
	query := `INSERT INTO attempts (word, session_id) VALUES ($1, $2) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, attempt.Word, attempt.SessionId).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	attempt.ID = id
	attempt.CreatedAt = createdAt
	return nil
}
