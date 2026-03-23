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
	tasks map[int]models.Task
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
	ts.tasks[int(id)] = task
	ts.mtx.Unlock()

	return task
}

func (ts *TaskStore) DeleteTask(id int) bool {
	ts.mtx.Lock()

	_, exist := ts.tasks[id]
	if !exist {
		ts.mtx.Unlock()
		return false
	}

	delete(ts.tasks, id)

	ts.mtx.Unlock()

	return true
}

func (ts *TaskStore) ChangeTask(id int, newTask models.NewTask) (models.Task, bool) {
	ts.mtx.Lock()

	_, exist := ts.tasks[id]
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

	ts.tasks[id] = task

	ts.mtx.Unlock()

	return task, true
}

func (ts *TaskStore) GetAllTasks() Response {
	ts.mtx.RLock()

	response := Response{
		Total: len(ts.tasks),
		Tasks: make([]models.Task, 0, len(ts.tasks)),
	}

	for _, task := range ts.tasks {
		response.Tasks = append(response.Tasks, task)
	}

	ts.mtx.RUnlock()

	return response
}

func (ts *TaskStore) GetAllTasksFiltered(filter string) Response {
	ts.mtx.RLock()

	response := Response{
		Total: len(ts.tasks),
		Tasks: make([]models.Task, 0, len(ts.tasks)),
	}

	for _, task := range ts.tasks {
		if strings.Contains(task.Title, filter) {
			response.Tasks = append(response.Tasks, task)
		}
	}

	ts.mtx.RUnlock()

	return response
}

func (ts *TaskStore) GetTask(id int) (models.Task, bool) {
	ts.mtx.RLock()

	responseTask, exist := ts.tasks[id]
	if !exist {
		ts.mtx.RUnlock()
		return models.Task{}, false
	}

	ts.mtx.RUnlock()

	return responseTask, true
}
