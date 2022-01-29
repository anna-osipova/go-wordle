package game

import (
	"github.com/anna-osipova/go-wordle/db"
)

type Attempt struct {
	ID        uint     `gorm:"primaryKey" json:"-"`
	SessionId string   `gorm:"column:session_id" json:"-"`
	WordGuess string   `gorm:"column:word_guess" json:"word_guess"`
	Letters   []Letter `json:"letters"`
}

type Letter struct {
	ID        uint    `gorm:"primaryKey" json:"-"`
	Letter    string  `json:"letter"`
	Color     string  `json:"color"`
	AttemptID uint    `json:"-"`
	Attempt   Attempt `gorm:"foreignKey:AttemptID;references:ID" json:"-"`
}

func AutoMigrate() {
	dbInstance := db.GetDB()

	dbInstance.AutoMigrate(&Attempt{})
	dbInstance.AutoMigrate(&Letter{})
}

func CreateAttempt(attempt *Attempt) error {
	dbInstance := db.GetDB()

	result := dbInstance.Create(&attempt)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAttempts(sessionId string) ([]Attempt, error) {
	var attempts []Attempt
	dbInstance := db.GetDB()

	result := dbInstance.Where(&Attempt{SessionId: sessionId}).Find(&attempts)
	return attempts, result.Error
}
