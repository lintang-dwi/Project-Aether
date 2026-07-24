package git

import (
	"fmt"
	"path/filepath"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Checkpoint represents a Git snapshot state.
type Checkpoint struct {
	ID          string    `json:"id"`
	Hash        string    `json:"hash"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
}

// Engine wraps go-git operations for checkpointing and rollback guarantees.
type Engine struct {
	repoPath string
	repo     *gogit.Repository
}

// NewEngine initializes or opens a Git repository at repoPath.
func NewEngine(repoPath string) (*Engine, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}

	repo, err := gogit.PlainOpen(absPath)
	if err == gogit.ErrRepositoryNotExists {
		repo, err = gogit.PlainInit(absPath, false)
		if err != nil {
			return nil, fmt.Errorf("failed to init git repository at %s: %w", absPath, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to open git repository at %s: %w", absPath, err)
	}

	return &Engine{
		repoPath: absPath,
		repo:     repo,
	}, nil
}

// CreateCheckpoint stages all changes and creates a Git commit snapshot.
func (e *Engine) CreateCheckpoint(message string) (*Checkpoint, error) {
	wt, err := e.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	if err := wt.AddWithOptions(&gogit.AddOptions{All: true}); err != nil {
		return nil, fmt.Errorf("failed to stage changes: %w", err)
	}

	commitHash, err := wt.Commit(message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "Project Aether Autonomous Engine",
			Email: "aether@runtime.local",
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create commit checkpoint: %w", err)
	}

	cp := &Checkpoint{
		ID:        fmt.Sprintf("checkpoint-%s", commitHash.String()[:7]),
		Hash:      commitHash.String(),
		Message:   message,
		CreatedAt: time.Now().UTC(),
	}

	return cp, nil
}

// Rollback resets the workspace back to a specific commit hash (100% rollback guarantee).
func (e *Engine) Rollback(commitHash string) error {
	wt, err := e.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	hash := plumbing.NewHash(commitHash)
	err = wt.Reset(&gogit.ResetOptions{
		Commit: hash,
		Mode:   gogit.HardReset,
	})
	if err != nil {
		return fmt.Errorf("failed to hard reset to commit %s: %w", commitHash, err)
	}

	return nil
}
