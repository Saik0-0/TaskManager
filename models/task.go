package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	Time      time.Time `json:"time"`
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
