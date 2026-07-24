package graph

import (
	"fmt"
	"sync"

	"aether/internal/knowledge"
	"aether/model"
	dgraph "github.com/dominikbraun/graph"
)

// Engine provides graph algorithms (traversal, shortest path, dependency lookup) over KnowledgeModel entities.
type Engine struct {
	mu     sync.RWMutex
	graph  dgraph.Graph[string, model.Node]
	nodes  map[string]model.Node
	edges  map[string]model.Edge
}

// NewEngine creates a directed Knowledge Graph Engine.
func NewEngine() *Engine {
	g := dgraph.New(func(n model.Node) string { return n.ID }, dgraph.Directed())
	return &Engine{
		graph: g,
		nodes: make(map[string]model.Node),
		edges: make(map[string]model.Edge),
	}
}

// LoadKnowledgeModel populates the graph engine from a KnowledgeModel.
func (e *Engine) LoadKnowledgeModel(km *knowledge.KnowledgeModel) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, node := range km.Entities {
		if err := e.graph.AddVertex(node); err != nil {
			// Skip if already added
		}
		e.nodes[node.ID] = node
	}

	for _, edge := range km.Relationships {
		if err := e.graph.AddEdge(edge.FromID, edge.ToID, dgraph.EdgeWeight(int(edge.Weight))); err != nil {
			// Skip duplicate edges or missing vertices
		}
		e.edges[edge.ID] = edge
	}

	return nil
}

// GetDependencies returns all direct target nodes pointed to by entityID.
func (e *Engine) GetDependencies(entityID string) ([]model.Node, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	adjMap, err := e.graph.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	targets, exists := adjMap[entityID]
	if !exists {
		return nil, fmt.Errorf("entity %s not found in graph", entityID)
	}

	deps := make([]model.Node, 0, len(targets))
	for targetID := range targets {
		if n, ok := e.nodes[targetID]; ok {
			deps = append(deps, n)
		}
	}

	return deps, nil
}

// GetDependents returns all source nodes that point to targetID (reverse dependency lookup).
func (e *Engine) GetDependents(targetID string) ([]model.Node, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	adjMap, err := e.graph.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	var dependents []model.Node
	for sourceID, targets := range adjMap {
		if _, pointsTo := targets[targetID]; pointsTo {
			if n, ok := e.nodes[sourceID]; ok {
				dependents = append(dependents, n)
			}
		}
	}

	return dependents, nil
}

// FindShortestPath calculates the shortest dependency path from sourceID to targetID.
func (e *Engine) FindShortestPath(sourceID, targetID string) ([]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return dgraph.ShortestPath(e.graph, sourceID, targetID)
}

// VertexCount returns the total number of nodes in the graph.
func (e *Engine) VertexCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return len(e.nodes)
}
