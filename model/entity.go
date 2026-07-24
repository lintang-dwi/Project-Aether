package model

import "time"

// EntityType defines the classification of a Graph Node.
type EntityType string

const (
	EntityFile     EntityType = "FILE"
	EntityFolder   EntityType = "FOLDER"
	EntityModule   EntityType = "MODULE"
	EntityPackage  EntityType = "PACKAGE"
	EntityFunction EntityType = "FUNCTION"
	EntityStruct   EntityType = "STRUCT"
	EntityClass    EntityType = "CLASS"
	EntityInterface EntityType = "INTERFACE"
	EntityVariable EntityType = "VARIABLE"
)

// Node represents an entity in the Knowledge Graph.
type Node struct {
	ID         string            `json:"id"`
	Type       EntityType        `json:"type"`
	Name       string            `json:"name"`
	Path       string            `json:"path"`
	Language   string            `json:"language,omitempty"`
	LineStart  int               `json:"line_start,omitempty"`
	LineEnd    int               `json:"line_end,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// FileMetadata represents detailed metadata for a file in the workspace.
type FileMetadata struct {
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	ModTime    time.Time `json:"mod_time"`
	Hash       string    `json:"hash,omitempty"`
	IsIgnored  bool      `json:"is_ignored"`
}
