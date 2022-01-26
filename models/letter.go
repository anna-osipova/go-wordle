package models

type Letter struct {
	ID        int    `json:"-"`
	AttemptId int    `json:"-"`
	Letter    string `json:"letter"`
	Color     string `json:"color"`
	CreatedAt string `json:"-"`
}
