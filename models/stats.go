package models

type Stats struct {
	Total         int     `json:"total"`
	Completed     int     `json:"completed"`
	CompletedRate float64 `json:"completed_rate"`
	LastTask      Task    `json:"last_task"`
}
