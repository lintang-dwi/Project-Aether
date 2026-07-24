package memory

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryRecord represents a single stored session decision or context memory.
type MemoryRecord struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// Engine manages short-term and persistent AI session memory.
type Engine struct {
	mu      sync.RWMutex
	records map[string]MemoryRecord
}

// NewEngine initializes a Memory Engine.
func NewEngine() *Engine {
	return &Engine{
		records: make(map[string]MemoryRecord),
	}
}

// Store saves a key-value record in memory.
func (e *Engine) Store(category, key, value string) MemoryRecord {
	e.mu.Lock()
	defer e.mu.Unlock()

	rec := MemoryRecord{
		ID:        uuid.New().String(),
		Key:       key,
		Value:     value,
		Category:  category,
		CreatedAt: time.Now().UTC(),
	}

	e.records[key] = rec
	return rec
}

// Retrieve fetches a memory record by key.
func (e *Engine) Retrieve(key string) (MemoryRecord, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rec, exists := e.records[key]
	return rec, exists
}
