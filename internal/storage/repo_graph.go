package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"aether/model"
)

type GraphRepository struct {
	db *DB
}

func NewGraphRepository(db *DB) *GraphRepository {
	return &GraphRepository{db: db}
}

func (r *GraphRepository) SaveNode(ctx context.Context, node model.Node) error {
	metaJSON, _ := json.Marshal(node.Metadata)
	query := `
	INSERT INTO nodes (id, type, name, path, language, line_start, line_end, metadata_json, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		type = excluded.type,
		name = excluded.name,
		path = excluded.path,
		language = excluded.language,
		line_start = excluded.line_start,
		line_end = excluded.line_end,
		metadata_json = excluded.metadata_json,
		updated_at = excluded.updated_at
	`
	now := time.Now().UTC()
	createdAt := node.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}

	_, err := r.db.ExecContext(ctx, query,
		node.ID, string(node.Type), node.Name, node.Path,
		node.Language, node.LineStart, node.LineEnd,
		string(metaJSON), createdAt, now,
	)
	if err != nil {
		return fmt.Errorf("failed to save node %s: %w", node.ID, err)
	}
	return nil
}

func (r *GraphRepository) GetNode(ctx context.Context, id string) (*model.Node, error) {
	query := `SELECT id, type, name, path, language, line_start, line_end, metadata_json, created_at, updated_at FROM nodes WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var node model.Node
	var nodeType, metaJSON string
	err := row.Scan(&node.ID, &nodeType, &node.Name, &node.Path, &node.Language, &node.LineStart, &node.LineEnd, &metaJSON, &node.CreatedAt, &node.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("node %s not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s: %w", id, err)
	}
	node.Type = model.EntityType(nodeType)
	if metaJSON != "" {
		_ = json.Unmarshal([]byte(metaJSON), &node.Metadata)
	}
	return &node, nil
}

func (r *GraphRepository) SaveEdge(ctx context.Context, edge model.Edge) error {
	metaJSON, _ := json.Marshal(edge.Metadata)
	query := `
	INSERT INTO edges (id, from_id, to_id, type, weight, metadata_json, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		type = excluded.type,
		weight = excluded.weight,
		metadata_json = excluded.metadata_json
	`
	now := time.Now().UTC()
	createdAt := edge.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}

	_, err := r.db.ExecContext(ctx, query,
		edge.ID, edge.FromID, edge.ToID, string(edge.Type),
		edge.Weight, string(metaJSON), createdAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save edge %s: %w", edge.ID, err)
	}
	return nil
}

func (r *GraphRepository) GetEdgesFrom(ctx context.Context, fromID string) ([]model.Edge, error) {
	query := `SELECT id, from_id, to_id, type, weight, metadata_json, created_at FROM edges WHERE from_id = ?`
	rows, err := r.db.QueryContext(ctx, query, fromID)
	if err != nil {
		return nil, fmt.Errorf("failed to query edges from %s: %w", fromID, err)
	}
	defer rows.Close()

	var edges []model.Edge
	for rows.Next() {
		var edge model.Edge
		var relType, metaJSON string
		if err := rows.Scan(&edge.ID, &edge.FromID, &edge.ToID, &relType, &edge.Weight, &metaJSON, &edge.CreatedAt); err != nil {
			return nil, err
		}
		edge.Type = model.RelationType(relType)
		if metaJSON != "" {
			_ = json.Unmarshal([]byte(metaJSON), &edge.Metadata)
		}
		edges = append(edges, edge)
	}
	return edges, nil
}
