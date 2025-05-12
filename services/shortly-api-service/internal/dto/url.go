package dto

import "time"

type CreateUrlResponseDTO struct {
	ID          uint   `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortKey    string `json:"short_url"`
	Title       string `json:"title"`
}

type GetUrlResponseDTO struct {
	ID          uint      `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortKey    string    `json:"short_url"`
	Title       string    `json:"title"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}
