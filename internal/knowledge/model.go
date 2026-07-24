package knowledge

import (
	"aether/model"
)

// KnowledgeModel is the top-level structured representation of project understanding.
type KnowledgeModel struct {
	Entities      map[string]model.Node `json:"entities"`
	Relationships map[string]model.Edge `json:"relationships"`
}

// NewKnowledgeModel initializes an empty KnowledgeModel.
func NewKnowledgeModel() *KnowledgeModel {
	return &KnowledgeModel{
		Entities:      make(map[string]model.Node),
		Relationships: make(map[string]model.Edge),
	}
}

// AddNode inserts or updates a Node entity in the model.
func (km *KnowledgeModel) AddNode(node model.Node) {
	km.Entities[node.ID] = node
}

// AddEdge inserts or updates a Relationship edge in the model.
func (km *KnowledgeModel) AddEdge(edge model.Edge) {
	km.Relationships[edge.ID] = edge
}

// GetNodes returns all entities as a slice.
func (km *KnowledgeModel) GetNodes() []model.Node {
	nodes := make([]model.Node, 0, len(km.Entities))
	for _, n := range km.Entities {
		nodes = append(nodes, n)
	}
	return nodes
}

// GetEdges returns all relationships as a slice.
func (km *KnowledgeModel) GetEdges() []model.Edge {
	edges := make([]model.Edge, 0, len(km.Relationships))
	for _, e := range km.Relationships {
		edges = append(edges, e)
	}
	return edges
}
