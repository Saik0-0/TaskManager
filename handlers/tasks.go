package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Saik0-0/TaskManager/models"
	"github.com/Saik0-0/TaskManager/storage"
	"net/http"
)

type Server struct {
	Store *storage.TaskStore
}

func (server *Server) TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()

		var newTask models.NewTask
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		responseTask := server.Store.AddTask(newTask)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responseTask); err != nil {
			fmt.Println("Encoding error: ", err)
			return
		}

	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
}
