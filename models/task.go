package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Text        string    `json:"text"`
	Completed   bool      `json:"completed"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

type NewTask struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type PatchTask struct {
	Title     *string `json:"title"`
	Text      *string `json:"text"`
	Completed *bool   `json:"completed"`
}
