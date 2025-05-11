package dto

import "time"

type UserDTO struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created"`
}

type UpdateUserDTO struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
}
