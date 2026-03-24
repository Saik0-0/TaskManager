package main

import (
	"fmt"
	"github.com/Saik0-0/TaskManager/handlers"
	"github.com/Saik0-0/TaskManager/models"
	"github.com/Saik0-0/TaskManager/storage"
	"net/http"
)

func main() {
	taskStore := storage.TaskStore{
		Tasks: make(map[int]models.Task),
	}

	server := handlers.Server{
		Store: &taskStore,
	}

	http.HandleFunc("/tasks", server.TasksHandler)
	http.HandleFunc("/tasks/", server.TaskHandler)
	http.HandleFunc("/stats", server.StatsHandler)

	if err := http.ListenAndServe(":9092", nil); err != nil {
		fmt.Println("Listening error: ", err)
	}
}
