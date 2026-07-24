package coordinator

import "context"

// ServiceStatus represents the current state of a registered service.
type ServiceStatus string

const (
	StatusStopped ServiceStatus = "STOPPED"
	StatusRunning ServiceStatus = "RUNNING"
	StatusFailed  ServiceStatus = "FAILED"
)

// Service is the unified interface that all runtime subsystems must implement.
type Service interface {
	// Name returns the unique identifier for the service.
	Name() string
	// Start initializes and boots the service.
	Start(ctx context.Context) error
	// Stop gracefully shuts down the service.
	Stop(ctx context.Context) error
	// Health returns the current health status of the service.
	Health() ServiceStatus
}
