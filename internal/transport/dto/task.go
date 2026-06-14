package dto

import "time"

type TaskDTO struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Text        string    `json:"text"`
	Completed   bool      `json:"completed"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

type NewTaskDTO struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type PatchTaskDTO struct {
	Title     *string `json:"title"`
	Text      *string `json:"text"`
	Completed *bool   `json:"completed"`
}
