package service

import (
	"github.com/anna-osipova/go-wordle/db"
)

type Session struct {
	ID   string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Word string `gorm:"word"`
}

func AutoMigrate() {
	dbInstance := db.GetDB()

	dbInstance.AutoMigrate(&Session{})
}

func GetSessionById(sessionId string) (Session, error) {
	var sessionModel Session
	dbInstance := db.GetDB()
	result := dbInstance.Where(Session{
		ID: sessionId,
	}).Take(&sessionModel)
	return sessionModel, result.Error
}

func CreateSession(session *Session) error {
	dbInstance := db.GetDB()
	result := dbInstance.Create(&session)
	return result.Error
}
