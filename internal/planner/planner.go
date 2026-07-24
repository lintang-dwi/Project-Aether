package planner

import (
	"context"
	"fmt"
	"time"

	"aether/internal/aiprovider"
	actxt "aether/internal/context"
	"aether/internal/task"
	"github.com/google/uuid"
)

type StepType string

const (
	StepFileEdit   StepType = "FILE_EDIT"
	StepFileCreate StepType = "FILE_CREATE"
	StepFileDelete StepType = "FILE_DELETE"
	StepRunCommand StepType = "RUN_COMMAND"
)

// PlanStep represents a single discrete action step within an execution plan.
type PlanStep struct {
	ID          string   `json:"id"`
	Sequence    int      `json:"sequence"`
	Type        StepType `json:"type"`
	TargetPath  string   `json:"target_path"`
	Instruction string   `json:"instruction"`
	Description string   `json:"description"`
}

// Plan represents an execution plan generated for a task.
type Plan struct {
	ID        string     `json:"id"`
	TaskID    string     `json:"task_id"`
	Goal      string     `json:"goal"`
	Steps     []PlanStep `json:"steps"`
	CreatedAt time.Time  `json:"created_at"`
}

// Planner generates execution plans combining AI Reasoning and Context.
type Planner struct {
	provider aiprovider.Provider
	ctxEngine *actxt.Engine
}

// NewPlanner creates a new Planner instance.
func NewPlanner(provider aiprovider.Provider, ctxEngine *actxt.Engine) *Planner {
	return &Planner{
		provider:  provider,
		ctxEngine: ctxEngine,
	}
}

// CreatePlan generates a structured execution plan for the given task.
func (p *Planner) CreatePlan(ctx context.Context, t *task.Task) (*Plan, error) {
	if t == nil {
		return nil, fmt.Errorf("task cannot be nil")
	}

	prepared := p.ctxEngine.PrepareContextForTarget(t.Goal, t.TargetPath, 4000)

	prompt := fmt.Sprintf("Generate an execution plan for Goal: %s\nTarget Path: %s\nContext:\n%s", t.Goal, t.TargetPath, prepared.UserPrompt)

	req := aiprovider.CompletionRequest{
		System: prepared.SystemPrompt,
		Prompt: prompt,
	}

	resp, err := p.provider.Complete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ai planner completion failed: %w", err)
	}

	// Create structured plan
	plan := &Plan{
		ID:        uuid.New().String(),
		TaskID:    t.ID,
		Goal:      t.Goal,
		CreatedAt: time.Now().UTC(),
		Steps: []PlanStep{
			{
				ID:          uuid.New().String(),
				Sequence:    1,
				Type:        StepFileEdit,
				TargetPath:  t.TargetPath,
				Instruction: resp.Content,
				Description: fmt.Sprintf("Apply changes for goal: %s", t.Goal),
			},
		},
	}

	return plan, nil
}
