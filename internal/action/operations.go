package action

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// FileOperation defines atomic file mutation.
type FileOperation interface {
	Execute(rootPath string) error
	Diff() string
}

type CreateFileOp struct {
	Path    string
	Content string
}

func (op *CreateFileOp) Execute(rootPath string) error {
	fullPath := filepath.Join(rootPath, op.Path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}
	return atomicWrite(fullPath, []byte(op.Content))
}

func (op *CreateFileOp) Diff() string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain("", op.Content, false)
	return dmp.DiffPrettyText(diffs)
}

type EditFileOp struct {
	Path       string
	NewContent string
}

func (op *EditFileOp) Execute(rootPath string) error {
	fullPath := filepath.Join(rootPath, op.Path)
	return atomicWrite(fullPath, []byte(op.NewContent))
}

func (op *EditFileOp) Diff() string {
	fullPath := op.Path
	oldContentBytes, _ := os.ReadFile(fullPath)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(oldContentBytes), op.NewContent, false)
	return dmp.DiffPrettyText(diffs)
}

type DeleteFileOp struct {
	Path string
}

func (op *DeleteFileOp) Execute(rootPath string) error {
	fullPath := filepath.Join(rootPath, op.Path)
	return os.Remove(fullPath)
}

func (op *DeleteFileOp) Diff() string {
	return fmt.Sprintf("Delete file: %s", op.Path)
}

func atomicWrite(filePath string, data []byte) error {
	tmpPath := filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file %s: %w", tmpPath, err)
	}
	if err := os.Rename(tmpPath, filePath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file to %s: %w", filePath, err)
	}
	return nil
}
