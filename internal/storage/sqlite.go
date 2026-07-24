package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DB wraps sql.DB for SQLite storage.
type DB struct {
	*sql.DB
}

// OpenDB opens a SQLite database using modernc.org/sqlite (pure Go).
// If dbPath is ":memory:", an in-memory database is created.
func OpenDB(dbPath string) (*DB, error) {
	if dbPath == "" {
		dbPath = ":memory:"
	}

	dsn := dbPath
	if dbPath != ":memory:" {
		dsn = fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", dbPath)
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database at %s: %w", dbPath, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	s := &DB{DB: db}
	if err := s.Migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return s, nil
}

// Migrate executes DDL schemas to ensure tables exist.
func (db *DB) Migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS workspace_metadata (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS nodes (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		path TEXT NOT NULL,
		language TEXT,
		line_start INTEGER,
		line_end INTEGER,
		metadata_json TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_nodes_type ON nodes(type);
	CREATE INDEX IF NOT EXISTS idx_nodes_path ON nodes(path);

	CREATE TABLE IF NOT EXISTS edges (
		id TEXT PRIMARY KEY,
		from_id TEXT NOT NULL,
		to_id TEXT NOT NULL,
		type TEXT NOT NULL,
		weight REAL DEFAULT 1.0,
		metadata_json TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (from_id) REFERENCES nodes(id) ON DELETE CASCADE,
		FOREIGN KEY (to_id) REFERENCES nodes(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_edges_from ON edges(from_id);
	CREATE INDEX IF NOT EXISTS idx_edges_to ON edges(to_id);
	CREATE INDEX IF NOT EXISTS idx_edges_type ON edges(type);

	CREATE TABLE IF NOT EXISTS checkpoints (
		id TEXT PRIMARY KEY,
		commit_hash TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(schema)
	return err
}
