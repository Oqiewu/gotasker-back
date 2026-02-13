package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
