package storage

import (
	"fmt"
	"github.com/Saik0-0/TaskManager/models"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type TaskStore struct {
	Tasks  map[int]models.Task
	NextID atomic.Int64
	mtx    sync.RWMutex
}

type Response struct {
	Total int           `json:"total"`
	Tasks []models.Task `json:"tasks"`
}

func (ts *TaskStore) AddTask(newTask models.NewTask) (models.Task, error) {
	if newTask.Title == "" {
		return models.Task{}, fmt.Errorf("title can't be empty")
	}

	id := ts.NextID.Add(1)

	task := models.Task{
		ID:        int(id),
		Title:     newTask.Title,
		Text:      newTask.Text,
		Completed: newTask.Completed,
	}

	ts.mtx.Lock()
	ts.Tasks[int(id)] = task
	ts.mtx.Unlock()

	return task, nil
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

func (ts *TaskStore) ChangeTask(id int, newTask models.NewTask) (models.Task, error) {
	ts.mtx.Lock()

	_, exist := ts.Tasks[id]
	if !exist {
		ts.mtx.Unlock()
		return models.Task{}, fmt.Errorf("task not found")
	}

	if newTask.Title == "" {
		ts.mtx.Unlock()
		return models.Task{}, fmt.Errorf("title can't be empty")
	}

	task := models.Task{
		ID:        id,
		Title:     newTask.Title,
		Text:      newTask.Text,
		Completed: newTask.Completed,
	}

	ts.Tasks[id] = task

	ts.mtx.Unlock()

	return task, nil
}

func (ts *TaskStore) GetAllTasks(titleFilter string, textFilter string, completeFilter string) (Response, error) {
	ts.mtx.RLock()

	response := Response{
		Total: 0,
		Tasks: make([]models.Task, 0, len(ts.Tasks)),
	}

	for _, task := range ts.Tasks {
		if titleFilter == "" || strings.Contains(task.Title, titleFilter) {
			if textFilter == "" || strings.Contains(task.Text, textFilter) {
				if completeFilter != "" {
					flag, err := strconv.ParseBool(completeFilter)
					if err != nil {
						return response, err
					}
					if flag && task.Completed {
						response.Tasks = append(response.Tasks, task)
						response.Total++
					}
					if !flag && !task.Completed {
						response.Tasks = append(response.Tasks, task)
						response.Total++
					}
				} else {
					response.Tasks = append(response.Tasks, task)
					response.Total++
				}
			}
		}
	}

	ts.mtx.RUnlock()

	return response, nil
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
