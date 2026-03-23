package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Saik0-0/TaskManager/models"
	"github.com/Saik0-0/TaskManager/storage"
	"net/http"
	"strconv"
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

	case http.MethodGet:
		query := r.URL.Query()
		filter := query.Get("filter")
		response := storage.Response{}
		if filter != "" {
			response = server.Store.GetAllTasksFiltered(filter)
		} else {
			response = server.Store.GetAllTasks()
		}

		offset := 0
		limit := response.Total

		offsetString := query.Get("offset")
		if offsetString != "" {
			o, err := strconv.Atoi(offsetString)
			if err != nil {
				http.Error(w, "Invalid offset: must be integer", http.StatusBadRequest)
			}
			offset = o
		}

		limitString := query.Get("limit")
		if limitString != "" {
			l, err := strconv.Atoi(limitString)
			if err != nil {
				http.Error(w, "Invalid limit: must be integer", http.StatusBadRequest)
			}
			limit = l
		}

		if offset > response.Total {
			offset = response.Total
		}
		end := offset + limit
		if end > response.Total {
			end = response.Total
		}

		response.Tasks = response.Tasks[offset:end]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Println("Encoding error: ", err)
			return
		}
		
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
}
