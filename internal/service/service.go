package service

import (
	"context"
	"fmt"
	"github.com/Saik0-0/TaskManager/internal/repository/models"
	"github.com/Saik0-0/TaskManager/internal/transport/dto"
)

type TaskRepository interface {
	AddTask(ctx context.Context, newTask dto.NewTaskDTO) (models.TaskModel, error)
	DeleteTask(ctx context.Context, id int) error
	GetTask(ctx context.Context, id int) (models.TaskModel, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

//type TaskStore struct {
//	Tasks  map[int]dto.TaskDTO
//	NextID atomic.Int64
//	mtx    sync.RWMutex
//}

func (ts *TaskService) AddTask(ctx context.Context, newTask dto.NewTaskDTO) (dto.TaskDTO, error) {
	if newTask.Title == "" {
		return dto.TaskDTO{}, fmt.Errorf("title can't be empty")
	}

	taskModel, err := ts.repo.AddTask(ctx, newTask)
	if err != nil {
		return dto.TaskDTO{}, fmt.Errorf("repository error: %w", err)
	}

	taskDTO := dto.TaskDTO{
		ID:          taskModel.ID,
		Title:       taskModel.Title,
		Text:        taskModel.Text,
		Completed:   taskModel.Completed,
		CreatedTime: taskModel.CreatedTime,
		UpdatedTime: taskModel.UpdatedTime,
	}

	return taskDTO, nil
}

func (ts *TaskService) DeleteTask(ctx context.Context, id int) error {
	return ts.repo.DeleteTask(ctx, id)
}

func (ts *TaskService) GetTask(ctx context.Context, id int) (dto.TaskDTO, error) {
	taskModel, err := ts.repo.GetTask(ctx, id)
	if err != nil {
		return dto.TaskDTO{}, fmt.Errorf("repository error: %w", err)
	}

	taskDTO := dto.TaskDTO{
		ID:          taskModel.ID,
		Title:       taskModel.Title,
		Text:        taskModel.Text,
		Completed:   taskModel.Completed,
		CreatedTime: taskModel.CreatedTime,
		UpdatedTime: taskModel.UpdatedTime,
	}

	return taskDTO, nil
}

//func (ts *TaskStore) DeleteTask(id int) bool {
//	ts.mtx.Lock()
//
//	if _, exist := ts.Tasks[id]; !exist {
//		ts.mtx.Unlock()
//		return false
//	}
//
//	delete(ts.Tasks, id)
//
//	ts.mtx.Unlock()
//
//	return true
//}
//
//func (ts *TaskStore) ChangeTask(id int, newTask dto.NewTaskDTO) (dto.TaskDTO, error) {
//	ts.mtx.Lock()
//
//	currTask, exist := ts.Tasks[id]
//	if !exist {
//		ts.mtx.Unlock()
//		return dto.TaskDTO{}, fmt.Errorf("task not found")
//	}
//
//	if newTask.Title == "" {
//		ts.mtx.Unlock()
//		return dto.TaskDTO{}, fmt.Errorf("title can't be empty")
//	}
//
//	task := dto.TaskDTO{
//		ID:          id,
//		Title:       newTask.Title,
//		Text:        newTask.Text,
//		Completed:   newTask.Completed,
//		CreatedTime: currTask.CreatedTime,
//		UpdatedTime: time.Now(),
//	}
//
//	ts.Tasks[id] = task
//
//	ts.mtx.Unlock()
//
//	return task, nil
//}
//
//func (ts *TaskStore) PartialChangeTask(id int, patchTask dto.PatchTaskDTO) (dto.TaskDTO, error) {
//	ts.mtx.Lock()
//
//	currentTask, exist := ts.Tasks[id]
//	if !exist {
//		ts.mtx.Unlock()
//		return dto.TaskDTO{}, fmt.Errorf("task not found")
//	}
//
//	if patchTask.Title != nil {
//		if *patchTask.Title == "" {
//			ts.mtx.Unlock()
//			return dto.TaskDTO{}, fmt.Errorf("title can't be empty")
//		}
//		currentTask.Title = *patchTask.Title
//	}
//	if patchTask.Text != nil {
//		currentTask.Text = *patchTask.Text
//	}
//	if patchTask.Completed != nil {
//		currentTask.Completed = *patchTask.Completed
//	}
//	currentTask.UpdatedTime = time.Now()
//
//	ts.Tasks[id] = currentTask
//
//	ts.mtx.Unlock()
//
//	return currentTask, nil
//}
//
//func (ts *TaskStore) GetAllTasks(titleFilter string, textFilter string, completeFilter string) ([]dto.TaskDTO, error) {
//	ts.mtx.RLock()
//
//	response := make([]dto.TaskDTO, 0, len(ts.Tasks))
//
//	flag := true
//	var err error
//	if completeFilter != "" {
//		flag, err = strconv.ParseBool(completeFilter)
//		if err != nil {
//			ts.mtx.RUnlock()
//			return response, err
//		}
//	}
//
//	for _, task := range ts.Tasks {
//		if titleFilter == "" || strings.Contains(task.Title, titleFilter) {
//			if textFilter == "" || strings.Contains(task.Text, textFilter) {
//				if completeFilter != "" {
//					if task.Completed == flag {
//						response = append(response, task)
//					}
//				} else {
//					response = append(response, task)
//				}
//			}
//		}
//	}
//
//	ts.mtx.RUnlock()
//
//	return response, nil
//}
//
//func (ts *TaskStore) GetTask(id int) (dto.TaskDTO, bool) {
//	ts.mtx.RLock()
//
//	responseTask, exist := ts.Tasks[id]
//	if !exist {
//		ts.mtx.RUnlock()
//		return dto.TaskDTO{}, false
//	}
//
//	ts.mtx.RUnlock()
//
//	return responseTask, true
//}
//
//func (ts *TaskStore) GetStats() dto.Stats {
//	var stats dto.Stats
//	var lastTime time.Time
//	flag := true
//
//	ts.mtx.RLock()
//	for _, task := range ts.Tasks {
//		if flag {
//			lastTime = task.CreatedTime
//			flag = false
//		}
//		stats.Total++
//
//		if task.Completed {
//			stats.Completed++
//		}
//
//		if task.CreatedTime.After(lastTime) {
//			stats.LastTask = task
//			lastTime = task.CreatedTime
//		}
//	}
//
//	if stats.Total != 0 {
//		stats.CompletedRate = float64(stats.Completed) / float64(stats.Total)
//	}
//
//	ts.mtx.RUnlock()
//
//	return stats
//}
