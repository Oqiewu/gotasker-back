package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gotasker/gotasker-back/src/internal/models"
)

// ErrTaskNotFound возвращается, когда задача с указанным ID не существует
var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, req *models.CreateTaskRequest) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tasks (title, description)
		 VALUES ($1, $2)
		 RETURNING id, title, description, completed_at, created_at, updated_at`,
		req.Title, req.Description,
	).Scan(&task.ID, &task.Title, &task.Description, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, description, completed_at, created_at, updated_at
		 FROM tasks
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}
	return tasks, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, completed_at, created_at, updated_at
		 FROM tasks WHERE id = $1`,
		id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func (r *TaskRepository) Update(ctx context.Context, id int, req *models.UpdateTaskRequest) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.QueryRowContext(ctx,
		`UPDATE tasks SET
			title       = COALESCE($2, title),
			description = COALESCE($3, description),
			completed_at = CASE
				WHEN $4::boolean IS NULL THEN completed_at
				WHEN $4 THEN COALESCE(completed_at, CURRENT_TIMESTAMP)
				ELSE NULL
			END,
			updated_at = CURRENT_TIMESTAMP
		 WHERE id = $1
		 RETURNING id, title, description, completed_at, created_at, updated_at`,
		id, req.Title, req.Description, req.Completed,
	).Scan(&task.ID, &task.Title, &task.Description, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if affected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
