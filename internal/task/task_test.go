package task_test

import (
	"testing"

	"aether/internal/task"
)

func TestTaskEngine_Lifecycle(t *testing.T) {
	engine := task.NewEngine()

	rootTask := engine.CreateTask("Refactor Logger", "internal/observability")
	if rootTask.Status != task.StatusPending {
		t.Errorf("expected PENDING, got %s", rootTask.Status)
	}

	subTask, err := engine.AddSubTask(rootTask.ID, "Migrate to Zap", "internal/observability/logger.go")
	if err != nil {
		t.Fatalf("failed to add subtask: %v", err)
	}
	if subTask.ParentID != rootTask.ID {
		t.Errorf("expected parent ID %s, got %s", rootTask.ID, subTask.ParentID)
	}

	if err := engine.UpdateStatus(subTask.ID, task.StatusCompleted, ""); err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	fetched, _ := engine.GetTask(subTask.ID)
	if fetched.Status != task.StatusCompleted {
		t.Errorf("expected COMPLETED, got %s", fetched.Status)
	}
	if fetched.CompletedAt == nil {
		t.Error("expected CompletedAt timestamp to be set")
	}
}
