package postgres

import (
	"context"
	"fmt"
	"github.com/Saik0-0/TaskManager/internal/repository/models"
	"github.com/Saik0-0/TaskManager/internal/transport/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetTask(ctx context.Context, id int) (models.TaskModel, error) {
	query := `
	SELECT * FROM tasks
	WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var taskModel models.TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Title,
		&taskModel.Text,
		&taskModel.Completed,
		&taskModel.CreatedTime,
		&taskModel.UpdatedTime,
	)
	if err != nil {
		return models.TaskModel{}, fmt.Errorf("scanning error: %w", err)
	}

	return taskModel, nil
}

func (r *Repository) DeleteTask(ctx context.Context, id int) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1
	`
	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("executing error: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id=%d: not founded", id)
	}
	return nil
}

func (r *Repository) AddTask(ctx context.Context, newTask dto.NewTaskDTO) (models.TaskModel, error) {
	query := `
	INSERT INTO tasks (title, description, completed)
	VALUES ($1, $2, $3)
	RETURNING id, title, description, completed, created_time, updated_time
	`
	row := r.pool.QueryRow(ctx, query,
		newTask.Title,
		newTask.Text,
		newTask.Completed)

	var taskModel models.TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Title,
		&taskModel.Text,
		&taskModel.Completed,
		&taskModel.CreatedTime,
		&taskModel.UpdatedTime,
	)
	if err != nil {
		return models.TaskModel{}, fmt.Errorf("scanning error: %w", err)
	}

	return taskModel, nil
}
