package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Saik0-0/TaskManager/models"
	"github.com/Saik0-0/TaskManager/storage"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
		title := query.Get("title")
		text := query.Get("text")
		complete := query.Get("complete")

		response, err := server.Store.GetAllTasks(title, text, complete)
		if err != nil {
			http.Error(w, "Invalid Query params", http.StatusBadRequest)
			return
		}

		sortingType := query.Get("sort")
		switch sortingType {
		case "title", "-title":
			order := strings.HasPrefix(sortingType, "-")
			sort.Slice(response.Tasks, func(i, j int) bool {
				if !order {
					return response.Tasks[i].Title < response.Tasks[j].Title
				}
				return response.Tasks[i].Title > response.Tasks[j].Title
			})

		case "completed", "-completed":
			order := strings.HasPrefix(sortingType, "-")
			sort.Slice(response.Tasks, func(i, j int) bool {
				if !order {
					return fromBoolToInt(response.Tasks[i].Completed) > fromBoolToInt(response.Tasks[j].Completed)
				}
				return fromBoolToInt(response.Tasks[i].Completed) < fromBoolToInt(response.Tasks[j].Completed)
			})
		}

		offset := 0
		limit := response.Total

		offsetString := query.Get("offset")
		if offsetString != "" {
			o, err := strconv.Atoi(offsetString)
			if err != nil {
				http.Error(w, "Invalid offset: must be integer", http.StatusBadRequest)
				return
			}
			offset = o
		}

		limitString := query.Get("limit")
		if limitString != "" {
			l, err := strconv.Atoi(limitString)
			if err != nil {
				http.Error(w, "Invalid limit: must be integer", http.StatusBadRequest)
				return
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

func (server *Server) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, idErr := parseID(r)
		if idErr != nil {
			http.Error(w, "Invalid id: must be integer", http.StatusBadRequest)
			return
		}

		responseTask, exist := server.Store.GetTask(id)
		if !exist {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(responseTask); err != nil {
			fmt.Println("Encoding error: ", err)
			return
		}

	case http.MethodPut:
		defer r.Body.Close()

		id, idErr := parseID(r)
		if idErr != nil {
			http.Error(w, "Invalid id: must be integer", http.StatusBadRequest)
			return
		}

		var newTask models.NewTask
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		responseTask, exist := server.Store.ChangeTask(id, newTask)
		if !exist {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(responseTask); err != nil {
			fmt.Println("Encoding error: ", err)
			return
		}

	case http.MethodDelete:
		id, idErr := parseID(r)
		if idErr != nil {
			http.Error(w, "Invalid id: must be integer", http.StatusBadRequest)
			return
		}

		if try := server.Store.DeleteTask(id); !try {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:

	}
}

func fromBoolToInt(flag bool) int {
	if flag {
		return 1
	}
	return 0
}

func parseID(r *http.Request) (int, error) {
	idString := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, idErr := strconv.Atoi(idString)
	if idErr != nil {
		return -1, fmt.Errorf("id parsing error %w", idErr)
	}
	return id, nil
}
