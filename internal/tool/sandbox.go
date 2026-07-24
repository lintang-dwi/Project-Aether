package tool

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ExecutionResult captures the output of a sandboxed command execution.
type ExecutionResult struct {
	Stdout   string        `json:"stdout"`
	Stderr   string        `json:"stderr"`
	ExitCode int           `json:"exit_code"`
	Duration time.Duration `json:"duration"`
}

// Sandbox executes command line operations safely bounded within workspace.
type Sandbox struct {
	workspaceRoot  string
	defaultTimeout time.Duration
}

// NewSandbox initializes Tool Sandbox.
func NewSandbox(workspaceRoot string, defaultTimeout time.Duration) (*Sandbox, error) {
	absPath, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return nil, err
	}

	if defaultTimeout <= 0 {
		defaultTimeout = 30 * time.Second
	}

	return &Sandbox{
		workspaceRoot:  absPath,
		defaultTimeout: defaultTimeout,
	}, nil
}

// RunCommand executes a command within workspace bounds with a timeout.
func (s *Sandbox) RunCommand(ctx context.Context, command string, args ...string) (*ExecutionResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.defaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = s.workspaceRoot

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return &ExecutionResult{
		Stdout:   strings.TrimSpace(stdoutBuf.String()),
		Stderr:   strings.TrimSpace(stderrBuf.String()),
		ExitCode: exitCode,
		Duration: duration,
	}, nil
}
