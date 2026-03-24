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

type Response struct {
	Total int           `json:"total"`
	Tasks []models.Task `json:"tasks"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (server *Server) TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()

		var newTask models.NewTask
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid JSON"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

		responseTask, addErr := server.Store.AddTask(newTask)
		if addErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid JSON: title can't be empty"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

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

		var response Response
		var err error
		response.Tasks, err = server.Store.GetAllTasks(title, text, complete)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid Query params"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}
		response.Total = len(response.Tasks)

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

		case "time", "-time":
			order := strings.HasPrefix(sortingType, "-")
			sort.Slice(response.Tasks, func(i, j int) bool {
				if !order {
					return response.Tasks[i].Time.Before(response.Tasks[j].Time)
				}
				return response.Tasks[i].Time.After(response.Tasks[j].Time)
			})
		}

		offset := 0
		limit := response.Total

		offsetString := query.Get("offset")
		if offsetString != "" {
			o, err := strconv.Atoi(offsetString)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid offset: must be integer"}); err != nil {
					fmt.Println("Encoding error", err)
					return
				}
				return
			}
			offset = o
		}

		limitString := query.Get("limit")
		if limitString != "" {
			l, err := strconv.Atoi(limitString)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid limit: must be integer"}); err != nil {
					fmt.Println("Encoding error", err)
					return
				}
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid method"}); err != nil {
			fmt.Println("Encoding error", err)
			return
		}
		return
	}
}

func (server *Server) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, idErr := parseID(r)
		if idErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid id: must be integer"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

		responseTask, exist := server.Store.GetTask(id)
		if !exist {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Task not found"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid id: must be integer"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

		var newTask models.NewTask
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid JSON"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

		responseTask, changingErr := server.Store.ChangeTask(id, newTask)
		if changingErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Task not found"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(responseTask); err != nil {
			fmt.Println("Encoding error: ", err)
			return
		}

	case http.MethodPatch:
		defer r.Body.Close()

		id, idErr := parseID(r)
		if idErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid id: must be integer"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

		var patchTask models.PatchTask
		if err := json.NewDecoder(r.Body).Decode(&patchTask); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid JSON"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

		responseTask, changingErr := server.Store.PartialChangeTask(id, patchTask)
		if changingErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Task not found"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid id: must be integer"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
			return
		}

		if try := server.Store.DeleteTask(id); !try {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Task not found"}); err != nil {
				fmt.Println("Encoding error", err)
				return
			}
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
