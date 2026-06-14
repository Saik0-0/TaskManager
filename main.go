package main

import (
	"context"
	"fmt"
	"github.com/Saik0-0/TaskManager/internal/repository/postgres"
	"github.com/Saik0-0/TaskManager/internal/service"
	"github.com/Saik0-0/TaskManager/internal/transport/handlers"
	"github.com/joho/godotenv"
	"net/http"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	conf, err := postgres.NewConfig()
	if err != nil {
		panic(err)
	}

	pool, err := postgres.NewConnectionPool(ctx, conf)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successful connection")

	repo := postgres.NewRepository(pool)
	taskService := service.NewTaskService(repo)

	server := handlers.NewServer(taskService)

	http.HandleFunc("/tasks", server.TasksHandler)
	http.HandleFunc("/tasks/", server.TaskHandler)
	http.HandleFunc("/stats", server.StatsHandler)

	if err := http.ListenAndServe(":9092", nil); err != nil {
		fmt.Println("Listening error: ", err)
	}
}
