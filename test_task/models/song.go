package models

import "time"

// Модель песни
type Song struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Group       string    `json:"group"`
	Song        string    `json:"song"`
	ReleaseDate string    `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
