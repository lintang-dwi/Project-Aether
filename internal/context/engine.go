package context

import (
	"aether/internal/graph"
	"aether/model"
)

// Engine extracts context from Knowledge Graph and prepares LLM prompts.
type Engine struct {
	graphEngine *graph.Engine
	builder     *Builder
}

// NewEngine initializes Context Engine.
func NewEngine(ge *graph.Engine, modelName string) (*Engine, error) {
	builder, err := NewBuilder(modelName)
	if err != nil {
		return nil, err
	}
	return &Engine{
		graphEngine: ge,
		builder:     builder,
	}, nil
}

// PrepareContextForTarget extracts graph nodes relevant to a target path and builds token-budgeted context.
func (e *Engine) PrepareContextForTarget(goal string, targetPath string, tokenBudget int) PreparedContext {
	systemRole := "You are Project Aether AI Runtime. System reasoning is guided by structured Knowledge Graph context."

	var nodes []model.Node
	if e.graphEngine != nil {
		nodes = e.graphEngine.QueryNodes(graph.FilterOptions{
			PathPrefix: targetPath,
		})
	}

	return e.builder.BuildStructuredContext(systemRole, goal, nodes, tokenBudget)
}

// GetBuilder returns the internal context builder.
func (e *Engine) GetBuilder() *Builder {
	return e.builder
}
