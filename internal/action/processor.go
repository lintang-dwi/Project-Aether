package action

import (
	"fmt"
	"sync"

	"aether/internal/eventbus"
)

// ActionResult records the result of an executed action operation.
type ActionResult struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
	Diff    string `json:"diff"`
	Error   string `json:"error,omitempty"`
}

// Processor executes Planned File Operations safely against the project workspace.
type Processor struct {
	mu       sync.Mutex
	rootPath string
	eventBus *eventbus.EventBus
}

// NewProcessor creates an Action Processor for a workspace root path.
func NewProcessor(rootPath string, eb *eventbus.EventBus) *Processor {
	return &Processor{
		rootPath: rootPath,
		eventBus: eb,
	}
}

// ExecuteOperation executes a single FileOperation.
func (p *Processor) ExecuteOperation(op FileOperation) ActionResult {
	p.mu.Lock()
	defer p.mu.Unlock()

	diff := op.Diff()
	err := op.Execute(p.rootPath)
	res := ActionResult{
		Success: err == nil,
		Diff:    diff,
	}

	if err != nil {
		res.Error = err.Error()
	} else if p.eventBus != nil {
		p.eventBus.Publish(eventbus.NewEvent("action.executed", "action_processor", map[string]interface{}{
			"diff": diff,
		}))
	}

	return res
}

// ExecuteBatch executes multiple operations sequentially.
func (p *Processor) ExecuteBatch(ops []FileOperation) ([]ActionResult, error) {
	var results []ActionResult
	for _, op := range ops {
		res := p.ExecuteOperation(op)
		results = append(results, res)
		if !res.Success {
			return results, fmt.Errorf("action operation failed: %s", res.Error)
		}
	}
	return results, nil
}
