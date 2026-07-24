package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type MetadataRepository struct {
	db *DB
}

func NewMetadataRepository(db *DB) *MetadataRepository {
	return &MetadataRepository{db: db}
}

func (r *MetadataRepository) Set(ctx context.Context, key, value string) error {
	query := `
	INSERT INTO workspace_metadata (key, value, updated_at)
	VALUES (?, ?, ?)
	ON CONFLICT(key) DO UPDATE SET
		value = excluded.value,
		updated_at = excluded.updated_at
	`
	_, err := r.db.ExecContext(ctx, query, key, value, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to set metadata key %s: %w", key, err)
	}
	return nil
}

func (r *MetadataRepository) Get(ctx context.Context, key string) (string, error) {
	query := `SELECT value FROM workspace_metadata WHERE key = ?`
	var val string
	err := r.db.QueryRowContext(ctx, query, key).Scan(&val)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("metadata key %s not found", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get metadata key %s: %w", key, err)
	}
	return val, nil
}
