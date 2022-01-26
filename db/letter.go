package db

import (
	"github.com/anna-osipova/go-wordle/models"
)

func (db Database) CreateLetter(letter *models.Letter) error {
	var id int
	var createdAt string
	query := `INSERT INTO letters (attempt_id, letter, color) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, letter.AttemptId, letter.Letter, letter.Color).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	letter.ID = id
	letter.CreatedAt = createdAt
	return nil
}
