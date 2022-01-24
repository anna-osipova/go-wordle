package models

type Session struct {
	ID        int    `json:"id"`
	Word      string `json:"word"`
	Attempts  int    `json:"attempts"`
	CreatedAt string `json:"created_at"`
}
