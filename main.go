package main

import (
	"context"
	"fmt"
	"github.com/Saik0-0/TaskManager/iternal/core/repository/feature_postgresql"
	"github.com/Saik0-0/TaskManager/iternal/core/service/storage"
	"github.com/Saik0-0/TaskManager/iternal/core/transport/dto"
	"github.com/Saik0-0/TaskManager/iternal/core/transport/handlers"
	"net/http"
)

func main() {
	ctx := context.Background()
	_, err := feature_postgresql.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}

	taskStore := storage.TaskStore{
		Tasks: make(map[int]dto.Task),
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
