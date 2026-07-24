package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "PENDING"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)

// Task represents an executable unit of work.
type Task struct {
	ID          string     `json:"id"`
	ParentID    string     `json:"parent_id,omitempty"`
	Goal        string     `json:"goal"`
	TargetPath  string     `json:"target_path,omitempty"`
	Status      TaskStatus `json:"status"`
	SubTasks    []*Task    `json:"sub_tasks,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Engine manages creation, decomposition, and lifecycle of tasks.
type Engine struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

// NewEngine creates a new Task Engine instance.
func NewEngine() *Engine {
	return &Engine{
		tasks: make(map[string]*Task),
	}
}

// CreateTask initializes a new top-level task from user goal.
func (e *Engine) CreateTask(goal, targetPath string) *Task {
	e.mu.Lock()
	defer e.mu.Unlock()

	task := &Task{
		ID:         uuid.New().String(),
		Goal:       goal,
		TargetPath: targetPath,
		Status:     StatusPending,
		CreatedAt:  time.Now().UTC(),
	}

	e.tasks[task.ID] = task
	return task
}

// AddSubTask attaches a sub-task to a parent task.
func (e *Engine) AddSubTask(parentID, goal, targetPath string) (*Task, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	parent, exists := e.tasks[parentID]
	if !exists {
		return nil, fmt.Errorf("parent task %s not found", parentID)
	}

	subTask := &Task{
		ID:         uuid.New().String(),
		ParentID:   parentID,
		Goal:       goal,
		TargetPath: targetPath,
		Status:     StatusPending,
		CreatedAt:  time.Now().UTC(),
	}

	parent.SubTasks = append(parent.SubTasks, subTask)
	e.tasks[subTask.ID] = subTask
	return subTask, nil
}

// UpdateStatus updates the status of a task.
func (e *Engine) UpdateStatus(taskID string, status TaskStatus, errStr string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	task, exists := e.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	task.Status = status
	task.Error = errStr
	if status == StatusCompleted || status == StatusFailed {
		now := time.Now().UTC()
		task.CompletedAt = &now
	}

	return nil
}

// GetTask retrieves a task by ID.
func (e *Engine) GetTask(taskID string) (*Task, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	task, exists := e.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}
	return task, nil
}
