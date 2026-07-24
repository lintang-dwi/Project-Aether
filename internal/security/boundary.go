package security

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type PermissionLevel string

const (
	PermRead    PermissionLevel = "READ"
	PermWrite   PermissionLevel = "WRITE"
	PermExecute PermissionLevel = "EXECUTE"
)

// AuditLogEntry records a security audit log.
type AuditLogEntry struct {
	Timestamp  time.Time       `json:"timestamp"`
	Permission PermissionLevel `json:"permission"`
	Target     string          `json:"target"`
	Allowed    bool            `json:"allowed"`
	Reason     string          `json:"reason,omitempty"`
}

// Boundary enforces security controls and permission limits over system operations.
type Boundary struct {
	mu             sync.RWMutex
	allowWrite     bool
	allowExecute   bool
	logger         *zap.Logger
	auditLogs      []AuditLogEntry
}

// NewBoundary creates a Security Boundary instance.
func NewBoundary(allowWrite, allowExecute bool, logger *zap.Logger) *Boundary {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Boundary{
		allowWrite:   allowWrite,
		allowExecute: allowExecute,
		logger:       logger,
		auditLogs:    make([]AuditLogEntry, 0),
	}
}

// CheckPermission evaluates if a permission request is allowed.
func (b *Boundary) CheckPermission(perm PermissionLevel, target string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	allowed := false
	reason := ""

	switch perm {
	case PermRead:
		allowed = true
	case PermWrite:
		allowed = b.allowWrite
		if !allowed {
			reason = "Write operations disabled by security boundary"
		}
	case PermExecute:
		allowed = b.allowExecute
		if !allowed {
			reason = "Command execution disabled by security boundary"
		}
	}

	entry := AuditLogEntry{
		Timestamp:  time.Now().UTC(),
		Permission: perm,
		Target:     target,
		Allowed:    allowed,
		Reason:     reason,
	}
	b.auditLogs = append(b.auditLogs, entry)

	b.logger.Info("Security permission check",
		zap.String("permission", string(perm)),
		zap.String("target", target),
		zap.Bool("allowed", allowed),
	)

	if !allowed {
		return fmt.Errorf("security violation: %s on target '%s'", reason, target)
	}

	return nil
}

// GetAuditLogs returns a copy of all audit logs.
func (b *Boundary) GetAuditLogs() []AuditLogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	logs := make([]AuditLogEntry, len(b.auditLogs))
	copy(logs, b.auditLogs)
	return logs
}
