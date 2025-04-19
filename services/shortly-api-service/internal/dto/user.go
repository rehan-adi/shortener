package dto

import "time"

type UserDTO struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	UrlsCount int       `json:"url_count"`
	CreatedAt time.Time `json:"created"`
}
