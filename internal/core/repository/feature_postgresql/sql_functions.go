package feature_postgresql

import (
	"context"
	"fmt"
	"github.com/Saik0-0/TaskManager/internal/core/transport/dto"
	"github.com/jackc/pgx/v5"
	"time"
)

func CreateTask(ctx context.Context, conn pgx.Conn, newTask dto.NewTask) error {
	if newTask.Title == "" {
		return fmt.Errorf("title can't be empty")
	}

	task := dto.Task{
		Title:       newTask.Title,
		Text:        newTask.Text,
		Completed:   newTask.Completed,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	sqlQuery := `
	INSERT INTO tasks (title, description, completed, created_time, updated_time)
	VALUES ($1, $2, $3, $4, $5);
	`

	_, err := conn.Exec(ctx, sqlQuery, task.Title, task.Text, task.Completed, task.CreatedTime, task.UpdatedTime)

	return err
}

//func UpdateTask(ctx context.Context, conn pgx.Conn, patchTask dto.PatchTask) error {
//
//}
