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
			writeError(w, http.StatusBadRequest, "Invalid json")
			return
		}

		responseTask, addErr := server.Store.AddTask(newTask)
		if addErr != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON: title can't be empty")
			return
		}

		writeJSON(w, http.StatusCreated, responseTask)

	case http.MethodGet:
		query := r.URL.Query()
		title := query.Get("title")
		text := query.Get("text")
		complete := query.Get("complete")

		var response Response
		var err error
		response.Tasks, err = server.Store.GetAllTasks(title, text, complete)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid Query params")
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
				writeError(w, http.StatusBadRequest, "Invalid offset: must be integer")
				return
			}
			offset = o
		}

		limitString := query.Get("limit")
		if limitString != "" {
			l, err := strconv.Atoi(limitString)
			if err != nil {
				writeError(w, http.StatusBadRequest, "Invalid limit: must be integer")
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

		writeJSON(w, http.StatusOK, response)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Invalid method")
		return
	}
}

func (server *Server) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, idErr := parseID(r)
		if idErr != nil {
			writeError(w, http.StatusBadRequest, "Invalid id: must be integer")
			return
		}

		responseTask, exist := server.Store.GetTask(id)
		if !exist {
			writeError(w, http.StatusNotFound, "Task not found")
			return
		}

		writeJSON(w, http.StatusOK, responseTask)

	case http.MethodPut:
		defer r.Body.Close()

		id, idErr := parseID(r)
		if idErr != nil {
			writeError(w, http.StatusBadRequest, "Invalid id: must be integer")
			return
		}

		var newTask models.NewTask
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		responseTask, changingErr := server.Store.ChangeTask(id, newTask)
		if changingErr != nil {
			writeError(w, http.StatusNotFound, "Task not found")
			return
		}

		writeJSON(w, http.StatusOK, responseTask)

	case http.MethodPatch:
		defer r.Body.Close()

		id, idErr := parseID(r)
		if idErr != nil {
			writeError(w, http.StatusBadRequest, "Invalid id: must be integer")
			return
		}

		var patchTask models.PatchTask
		if err := json.NewDecoder(r.Body).Decode(&patchTask); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		responseTask, changingErr := server.Store.PartialChangeTask(id, patchTask)
		if changingErr != nil {
			writeError(w, http.StatusNotFound, "Task not found")
			return
		}

		writeJSON(w, http.StatusOK, responseTask)

	case http.MethodDelete:
		id, idErr := parseID(r)
		if idErr != nil {
			writeError(w, http.StatusBadRequest, "Invalid id: must be integer")
			return
		}

		if try := server.Store.DeleteTask(id); !try {
			writeError(w, http.StatusNotFound, "Task not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Invalid method")
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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println("Encoding error: ", err)
		return
	}
	return
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Message: message}); err != nil {
		fmt.Println("Encoding error: ", err)
		return
	}
	return
}
