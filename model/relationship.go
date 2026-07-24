package model

import "time"

// RelationType defines the nature of connection between two entities.
type RelationType string

const (
	RelContains   RelationType = "CONTAINS"
	RelImports    RelationType = "IMPORTS"
	RelCalls      RelationType = "CALLS"
	RelExtends    RelationType = "EXTENDS"
	RelImplements RelationType = "IMPLEMENTS"
	RelDependsOn  RelationType = "DEPENDS_ON"
)

// Edge represents a directed relationship between two Nodes in the Knowledge Graph.
type Edge struct {
	ID        string            `json:"id"`
	FromID    string            `json:"from_id"`
	ToID      string            `json:"to_id"`
	Type      RelationType      `json:"type"`
	Weight    float64           `json:"weight,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}
