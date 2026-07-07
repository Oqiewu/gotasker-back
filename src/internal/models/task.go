package models

import (
	"encoding/json"
	"time"
)

const responseTimeLayout = "15:04:05 02.01.2006"

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (t Task) MarshalJSON() ([]byte, error) {
	var completedAt *string
	if t.CompletedAt != nil {
		s := t.CompletedAt.Format(responseTimeLayout)
		completedAt = &s
	}
	return json.Marshal(struct {
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		CompletedAt *string `json:"completed_at"`
		CreatedAt   string  `json:"created_at"`
		UpdatedAt   string  `json:"updated_at"`
	}{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		CompletedAt: completedAt,
		CreatedAt:   t.CreatedAt.Format(responseTimeLayout),
		UpdatedAt:   t.UpdatedAt.Format(responseTimeLayout),
	})
}

// CreateTaskRequest — тело запроса на создание задачи
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description"`
}

// UpdateTaskRequest — тело запроса на обновление задачи (частичное)
type UpdateTaskRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=255"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
