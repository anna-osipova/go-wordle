package models

type Session struct {
	ID        string `json:"id"`
	Word      string `json:"word"`
	Attempts  int    `json:"attempts"`
	CreatedAt string `json:"created_at"`
}
