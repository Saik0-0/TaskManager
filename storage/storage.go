package storage

import (
	"github.com/Saik0-0/TaskManager/models"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	nextID atomic.Int64
)

type TaskStore struct {
	Tasks map[int]models.Task
	mtx   sync.RWMutex
}

type Response struct {
	Total int           `json:"total"`
	Tasks []models.Task `json:"tasks"`
}

func (ts *TaskStore) AddTask(newTask models.NewTask) models.Task {
	id := nextID.Add(1)

	task := models.Task{
		ID:        int(id),
		Title:     newTask.Title,
		Text:      newTask.Text,
		Completed: newTask.Completed,
	}

	ts.mtx.Lock()
	ts.Tasks[int(id)] = task
	ts.mtx.Unlock()

	return task
}

func (ts *TaskStore) DeleteTask(id int) bool {
	ts.mtx.Lock()

	_, exist := ts.Tasks[id]
	if !exist {
		ts.mtx.Unlock()
		return false
	}

	delete(ts.Tasks, id)

	ts.mtx.Unlock()

	return true
}

func (ts *TaskStore) ChangeTask(id int, newTask models.NewTask) (models.Task, bool) {
	ts.mtx.Lock()

	_, exist := ts.Tasks[id]
	if !exist {
		ts.mtx.Unlock()
		return models.Task{}, false
	}

	task := models.Task{
		ID:        id,
		Title:     newTask.Title,
		Text:      newTask.Text,
		Completed: newTask.Completed,
	}

	ts.Tasks[id] = task

	ts.mtx.Unlock()

	return task, true
}

func (ts *TaskStore) GetAllTasks() Response {
	ts.mtx.RLock()

	response := Response{
		Total: len(ts.Tasks),
		Tasks: make([]models.Task, 0, len(ts.Tasks)),
	}

	for _, task := range ts.Tasks {
		response.Tasks = append(response.Tasks, task)
	}

	ts.mtx.RUnlock()

	return response
}

func (ts *TaskStore) GetAllTasksFiltered(filter string) Response {
	ts.mtx.RLock()

	response := Response{
		Total: 0,
		Tasks: make([]models.Task, 0, len(ts.Tasks)),
	}

	for _, task := range ts.Tasks {
		if strings.Contains(task.Title, filter) {
			response.Tasks = append(response.Tasks, task)
			response.Total++
		}
	}

	ts.mtx.RUnlock()

	return response
}

func (ts *TaskStore) GetTask(id int) (models.Task, bool) {
	ts.mtx.RLock()

	responseTask, exist := ts.Tasks[id]
	if !exist {
		ts.mtx.RUnlock()
		return models.Task{}, false
	}

	ts.mtx.RUnlock()

	return responseTask, true
}
