package models

import "time"

type TaskModel struct {
	ID          int
	Title       string
	Text        string
	Completed   bool
	CreatedTime time.Time
	UpdatedTime time.Time
}
