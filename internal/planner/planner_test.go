package planner_test

import (
	"context"
	"testing"

	"aether/internal/aiprovider"
	actxt "aether/internal/context"
	"aether/internal/graph"
	"aether/internal/planner"
	"aether/internal/task"
)

type MockAIProvider struct{}

func (m *MockAIProvider) Name() string { return "mock" }
func (m *MockAIProvider) Complete(ctx context.Context, req aiprovider.CompletionRequest) (*aiprovider.CompletionResponse, error) {
	return &aiprovider.CompletionResponse{
		Content:      "Step 1: Edit file\nStep 2: Run tests",
		FinishReason: "stop",
	}, nil
}
func (m *MockAIProvider) StreamComplete(ctx context.Context, req aiprovider.CompletionRequest, handler func(chunk string) error) error {
	return handler("Step 1: Edit file")
}

func TestPlanner_CreatePlan(t *testing.T) {
	ge := graph.NewEngine()
	ctxEngine, err := actxt.NewEngine(ge, "gpt-4o")
	if err != nil {
		t.Fatalf("failed to create context engine: %v", err)
	}

	mockAI := &MockAIProvider{}
	pl := planner.NewPlanner(mockAI, ctxEngine)

	taskEngine := task.NewEngine()
	testTask := taskEngine.CreateTask("Refactor Config", "internal/config/config.go")

	plan, err := pl.CreatePlan(context.Background(), testTask)
	if err != nil {
		t.Fatalf("failed to create plan: %v", err)
	}

	if plan.TaskID != testTask.ID {
		t.Errorf("expected plan task ID %s, got %s", testTask.ID, plan.TaskID)
	}
	if len(plan.Steps) != 1 {
		t.Errorf("expected 1 step in plan, got %d", len(plan.Steps))
	}
	if plan.Steps[0].Type != planner.StepFileEdit {
		t.Errorf("expected StepFileEdit, got %v", plan.Steps[0].Type)
	}
}
