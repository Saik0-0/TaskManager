package storage

import (
	"fmt"
	"github.com/Saik0-0/TaskManager/models"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type TaskStore struct {
	Tasks  map[int]models.Task
	NextID atomic.Int64
	mtx    sync.RWMutex
}

func (ts *TaskStore) AddTask(newTask models.NewTask) (models.Task, error) {
	if newTask.Title == "" {
		return models.Task{}, fmt.Errorf("title can't be empty")
	}

	id := ts.NextID.Add(1)

	task := models.Task{
		ID:          int(id),
		Title:       newTask.Title,
		Text:        newTask.Text,
		Completed:   newTask.Completed,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	ts.mtx.Lock()
	ts.Tasks[int(id)] = task
	ts.mtx.Unlock()

	return task, nil
}

func (ts *TaskStore) DeleteTask(id int) bool {
	ts.mtx.Lock()

	if _, exist := ts.Tasks[id]; !exist {
		ts.mtx.Unlock()
		return false
	}

	delete(ts.Tasks, id)

	ts.mtx.Unlock()

	return true
}

func (ts *TaskStore) ChangeTask(id int, newTask models.NewTask) (models.Task, error) {
	ts.mtx.Lock()

	currTask, exist := ts.Tasks[id]
	if !exist {
		ts.mtx.Unlock()
		return models.Task{}, fmt.Errorf("task not found")
	}

	if newTask.Title == "" {
		ts.mtx.Unlock()
		return models.Task{}, fmt.Errorf("title can't be empty")
	}

	task := models.Task{
		ID:          id,
		Title:       newTask.Title,
		Text:        newTask.Text,
		Completed:   newTask.Completed,
		CreatedTime: currTask.CreatedTime,
		UpdatedTime: time.Now(),
	}

	ts.Tasks[id] = task

	ts.mtx.Unlock()

	return task, nil
}

func (ts *TaskStore) PartialChangeTask(id int, patchTask models.PatchTask) (models.Task, error) {
	ts.mtx.Lock()

	currentTask, exist := ts.Tasks[id]
	if !exist {
		ts.mtx.Unlock()
		return models.Task{}, fmt.Errorf("task not found")
	}

	if patchTask.Title != nil {
		if *patchTask.Title == "" {
			ts.mtx.Unlock()
			return models.Task{}, fmt.Errorf("title can't be empty")
		}
		currentTask.Title = *patchTask.Title
	}
	if patchTask.Text != nil {
		currentTask.Text = *patchTask.Text
	}
	if patchTask.Completed != nil {
		currentTask.Completed = *patchTask.Completed
	}
	currentTask.UpdatedTime = time.Now()

	ts.Tasks[id] = currentTask

	ts.mtx.Unlock()

	return currentTask, nil
}

func (ts *TaskStore) GetAllTasks(titleFilter string, textFilter string, completeFilter string) ([]models.Task, error) {
	ts.mtx.RLock()

	response := make([]models.Task, 0, len(ts.Tasks))

	flag := true
	var err error
	if completeFilter != "" {
		flag, err = strconv.ParseBool(completeFilter)
		if err != nil {
			ts.mtx.RUnlock()
			return response, err
		}
	}

	for _, task := range ts.Tasks {
		if titleFilter == "" || strings.Contains(task.Title, titleFilter) {
			if textFilter == "" || strings.Contains(task.Text, textFilter) {
				if completeFilter != "" {
					if task.Completed == flag {
						response = append(response, task)
					}
				} else {
					response = append(response, task)
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

func (ts *TaskStore) GetStats() models.Stats {
	var stats models.Stats
	var lastTime time.Time
	flag := true

	ts.mtx.RLock()
	for _, task := range ts.Tasks {
		if flag {
			lastTime = task.CreatedTime
			flag = false
		}
		stats.Total++

		if task.Completed {
			stats.Completed++
		}

		if task.CreatedTime.After(lastTime) {
			stats.LastTask = task
			lastTime = task.CreatedTime
		}
	}

	if stats.Total != 0 {
		stats.CompletedRate = float64(stats.Completed) / float64(stats.Total)
	}

	ts.mtx.RUnlock()

	return stats
}
