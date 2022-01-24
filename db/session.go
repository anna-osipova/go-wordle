package db

import (
	"database/sql"

	"github.com/anna-osipova/go-wordle/models"
)

func (db Database) CreateSession(session *models.Session) error {
	var id int
	var createdAt string
	query := `INSERT INTO sessions (word) VALUES ($1) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, session.Word).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	session.ID = id
	session.CreatedAt = createdAt
	return nil
}
func (db Database) GetSessionById(sessionId int) (models.Session, error) {
	session := models.Session{}
	query := `SELECT * FROM sessions WHERE id = $1;`
	row := db.Conn.QueryRow(query, sessionId)
	switch err := row.Scan(&session.ID, &session.Word, &session.CreatedAt); err {
	case sql.ErrNoRows:
		return session, ErrNoMatch
	default:
		return session, err
	}
}

func (db Database) UpdateSessionAttemptCount(sessionId int, sessionData models.Session) (models.Session, error) {
	session := models.Session{}
	query := `UPDATE sessions SET attempts=$1 WHERE id=$2 RETURNING id, attempts, created_at;`
	err := db.Conn.QueryRow(query, sessionData.Attempts, sessionId).Scan(&session.ID, &session.Attempts, &session.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return session, ErrNoMatch
		}
		return session, err
	}
	return session, nil
}
