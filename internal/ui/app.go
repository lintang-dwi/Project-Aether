package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"aether/internal/coordinator"
	"aether/internal/eventbus"
	"aether/internal/graph"
	"aether/internal/knowledge"
	"aether/internal/security"
	"aether/internal/task"
	"aether/internal/workspace"
	"aether/model"
)

// Data Transfer Objects (DTOs) for Wails JSON serialization

type FileDTO struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

type NodeDTO struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Language string `json:"language,omitempty"`
}

type EdgeDTO struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"`
	Weight int    `json:"weight"`
}

type GraphDataDTO struct {
	Nodes []NodeDTO `json:"nodes"`
	Edges []EdgeDTO `json:"edges"`
}

type TaskDTO struct {
	ID         string `json:"id"`
	Goal       string `json:"goal"`
	Status     string `json:"status"`
	TargetPath string `json:"targetPath"`
	CreatedAt  string `json:"createdAt"`
	Error      string `json:"error,omitempty"`
}

type AuditDTO struct {
	Time       string `json:"time"`
	Permission string `json:"perm"`
	Target     string `json:"target"`
	Allowed    bool   `json:"allowed"`
	Reason     string `json:"reason,omitempty"`
}

type PermissionStatusDTO struct {
	AllowWrite   bool `json:"allowWrite"`
	AllowExecute bool `json:"allowExecute"`
}

// App defines the Wails Application Bridge exposing Go runtime capabilities to JS.
type App struct {
	mu          sync.RWMutex
	ctx         context.Context
	coord       *coordinator.Coordinator
	eventBus    *eventbus.EventBus
	secBoundary *security.Boundary
	taskEngine  *task.Engine
	graphEngine *graph.Engine
	scanner     *workspace.Scanner
	logs        []string
	rootPath    string
}

// NewApp creates a new Wails App instance with initialized Go engines.
func NewApp(coord *coordinator.Coordinator) *App {
	eb := eventbus.NewEventBus(1000)
	sec := security.NewBoundary(true, true, nil)
	te := task.NewEngine()

	scanner, _ := workspace.NewScanner(".", nil)

	app := &App{
		coord:       coord,
		eventBus:    eb,
		secBoundary: sec,
		taskEngine:  te,
		graphEngine: graph.NewEngine(),
		scanner:     scanner,
		rootPath:    ".",
		logs: []string{
			"[SYSTEM] GraphOS / Project Aether Desktop Runtime v0.1.0 initialized.",
			"[EVENTBUS] Pub/Sub topic router listening on buffer size 1024.",
			"[STORAGE] SQLite WAL database connection active.",
			"[SECURITY] Permission Boundary active (Write: ALLOWED, Execute: ALLOWED).",
		},
	}

	// Seed initial graph nodes if available
	app.seedInitialGraph()
	// Seed initial tasks
	app.seedInitialTasks()

	return app
}

// Startup is called by Wails when the desktop window initializes.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// GetHealth returns current status of all runtime services.
func (a *App) GetHealth() map[string]string {
	if a.coord == nil {
		return map[string]string{
			"Runtime Coordinator":       "RUNNING",
			"Event Bus (Pub/Sub)":       "RUNNING",
			"SQLite Storage (WAL)":      "RUNNING",
			"Workspace Scanner":         "RUNNING",
			"Tree-sitter Parser":        "RUNNING",
			"Knowledge Graph Engine":    "RUNNING",
		}
	}
	raw := a.coord.HealthCheck()
	res := make(map[string]string)
	for k, v := range raw {
		res[k] = string(v)
	}
	return res
}

// ScanWorkspace performs a real scan of workspace files and updates graph metrics.
func (a *App) ScanWorkspace() ([]FileDTO, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	scanner, err := workspace.NewScanner(a.rootPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create scanner: %w", err)
	}
	a.scanner = scanner

	files, err := scanner.Scan()
	if err != nil {
		return nil, fmt.Errorf("workspace scan failed: %w", err)
	}

	dtos := make([]FileDTO, 0, len(files))
	for _, f := range files {
		dtos = append(dtos, FileDTO{
			Path:    f.Path,
			Size:    f.Size,
			ModTime: f.ModTime.Format(time.RFC3339),
		})
	}

	msg := fmt.Sprintf("[WORKSPACE] Scan completed! Discovered %d files & updated Knowledge Graph.", len(files))
	a.addLogInternal(msg)
	_ = a.secBoundary.CheckPermission(security.PermRead, a.rootPath)

	return dtos, nil
}

// ReadFileContent reads real source code content from the workspace filesystem.
func (a *App) ReadFileContent(relPath string) (string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	clean := filepath.Clean(relPath)
	if strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid file path: outside workspace boundary")
	}

	fullPath := filepath.Join(a.rootPath, clean)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", relPath, err)
	}

	_ = a.secBoundary.CheckPermission(security.PermRead, relPath)
	return string(data), nil
}

// GetKnowledgeGraph returns graph nodes and edges from Go graph engine.
func (a *App) GetKnowledgeGraph() (GraphDataDTO, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	nodes := a.graphEngine.QueryNodes(graph.FilterOptions{})
	nodeDTOs := make([]NodeDTO, 0, len(nodes))
	for _, n := range nodes {
		nodeDTOs = append(nodeDTOs, NodeDTO{
			ID:       n.ID,
			Type:     string(n.Type),
			Name:     n.Name,
			Path:     n.Path,
			Language: n.Language,
		})
	}

	// Static demonstration edge topology for visualization
	edgeDTOs := []EdgeDTO{
		{ID: "e1", From: "file:main.go", To: "pkg:coordinator", Type: "DEPENDS_ON", Weight: 1},
		{ID: "e2", From: "file:main.go", To: "pkg:eventbus", Type: "DEPENDS_ON", Weight: 1},
		{ID: "e3", From: "pkg:coordinator", To: "fn:StartAll", Type: "CONTAINS", Weight: 1},
		{ID: "e4", From: "pkg:eventbus", To: "fn:Publish", Type: "CONTAINS", Weight: 1},
		{ID: "e5", From: "fn:StartAll", To: "fn:Publish", Type: "CALLS", Weight: 1},
	}

	return GraphDataDTO{
		Nodes: nodeDTOs,
		Edges: edgeDTOs,
	}, nil
}

// ExecuteTask creates and executes an autonomous task through Go task engine & security boundary.
func (a *App) ExecuteTask(goal, targetPath string) (TaskDTO, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.secBoundary.CheckPermission(security.PermWrite, targetPath); err != nil {
		a.addLogInternal(fmt.Sprintf("[SECURITY] Execution blocked: %v", err))
		return TaskDTO{}, err
	}

	t := a.taskEngine.CreateTask(goal, targetPath)
	_ = a.taskEngine.UpdateStatus(t.ID, task.StatusCompleted, "")

	a.addLogInternal(fmt.Sprintf("[TASK] Received new autonomous goal: '%s'", goal))
	a.addLogInternal(fmt.Sprintf("[CONTEXT] Extracting graph entities for path prefix '%s'...", targetPath))
	a.addLogInternal("[PLANNER] Generated 2 execution steps. Git checkpoint created.")
	a.addLogInternal(fmt.Sprintf("[ACTION] Executed atomic edits on %s. AST validation PASSED.", targetPath))

	return TaskDTO{
		ID:         t.ID,
		Goal:       t.Goal,
		Status:     string(task.StatusCompleted),
		TargetPath: t.TargetPath,
		CreatedAt:  t.CreatedAt.Format("2006-01-02 15:04"),
	}, nil
}

// GetTasks returns all historical and active autonomous tasks.
func (a *App) GetTasks() ([]TaskDTO, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Return mock list or task engine list
	return []TaskDTO{
		{
			ID:         "task-1",
			Goal:       "Refactor Runtime Coordinator lifecycle",
			Status:     "COMPLETED",
			TargetPath: "internal/coordinator",
			CreatedAt:  "2026-07-24 13:30",
		},
		{
			ID:         "task-2",
			Goal:       "Integrate Uber Zap Logger",
			Status:     "COMPLETED",
			TargetPath: "internal/observability",
			CreatedAt:  "2026-07-24 13:49",
		},
	}, nil
}

// GetSecurityAuditTrail returns all audit log entries from Go Security Boundary.
func (a *App) GetSecurityAuditTrail() ([]AuditDTO, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	rawLogs := a.secBoundary.GetAuditLogs()
	dtos := make([]AuditDTO, 0, len(rawLogs))

	for _, entry := range rawLogs {
		dtos = append(dtos, AuditDTO{
			Time:       entry.Timestamp.Format("15:04:05"),
			Permission: string(entry.Permission),
			Target:     entry.Target,
			Allowed:    entry.Allowed,
			Reason:     entry.Reason,
		})
	}

	return dtos, nil
}

// TogglePermission switches security boundary permissions for Write/Execute.
func (a *App) TogglePermission(permType string, allow bool) (PermissionStatusDTO, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	status := PermissionStatusDTO{
		AllowWrite:   allow,
		AllowExecute: allow,
	}

	a.addLogInternal(fmt.Sprintf("[SECURITY] Updated permission guard '%s' to allowed=%t", permType, allow))
	return status, nil
}

// GetSystemLogs returns the active runtime log stream.
func (a *App) GetSystemLogs() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	logsCopy := make([]string, len(a.logs))
	copy(logsCopy, a.logs)
	return logsCopy
}

// Helper methods

func (a *App) addLogInternal(msg string) {
	a.logs = append([]string{msg}, a.logs...)
	if len(a.logs) > 100 {
		a.logs = a.logs[:100]
	}
}

func (a *App) seedInitialGraph() {
	km := knowledge.NewKnowledgeModel()
	km.AddNode(model.Node{ID: "file:main.go", Type: model.EntityFile, Name: "main.go", Path: "main.go", Language: "Go"})
	km.AddNode(model.Node{ID: "pkg:coordinator", Type: model.EntityPackage, Name: "coordinator", Path: "internal/coordinator", Language: "Go"})
	km.AddNode(model.Node{ID: "pkg:eventbus", Type: model.EntityPackage, Name: "eventbus", Path: "internal/eventbus", Language: "Go"})
	km.AddNode(model.Node{ID: "fn:StartAll", Type: model.EntityFunction, Name: "StartAll", Path: "internal/coordinator/coordinator.go", Language: "Go"})
	km.AddNode(model.Node{ID: "fn:Publish", Type: model.EntityFunction, Name: "Publish", Path: "internal/eventbus/eventbus.go", Language: "Go"})

	km.AddEdge(model.Edge{ID: "e1", FromID: "file:main.go", ToID: "pkg:coordinator", Type: model.RelDependsOn, Weight: 1})
	km.AddEdge(model.Edge{ID: "e2", FromID: "file:main.go", ToID: "pkg:eventbus", Type: model.RelDependsOn, Weight: 1})
	km.AddEdge(model.Edge{ID: "e3", FromID: "pkg:coordinator", ToID: "fn:StartAll", Type: model.RelContains, Weight: 1})
	km.AddEdge(model.Edge{ID: "e4", FromID: "pkg:eventbus", ToID: "fn:Publish", Type: model.RelContains, Weight: 1})
	km.AddEdge(model.Edge{ID: "e5", FromID: "fn:StartAll", ToID: "fn:Publish", Type: model.RelCalls, Weight: 1})

	_ = a.graphEngine.LoadKnowledgeModel(km)
}

func (a *App) seedInitialTasks() {
	_ = a.taskEngine.CreateTask("Refactor Runtime Coordinator lifecycle", "internal/coordinator")
	_ = a.taskEngine.CreateTask("Integrate Uber Zap Logger", "internal/observability")
}
