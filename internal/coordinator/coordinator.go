package coordinator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"aether/internal/config"
	"aether/internal/eventbus"
	"go.uber.org/zap"
)

// Coordinator orchestrates all runtime services, managing their lifecycle and state.
type Coordinator struct {
	mu          sync.RWMutex
	cfg         config.Config
	logger      *zap.Logger
	eventBus    *eventbus.EventBus
	services    []Service
	serviceMap  map[string]Service
	statusMap   map[string]ServiceStatus
	isBooted    bool
}

// NewCoordinator creates a new Runtime Coordinator instance.
func NewCoordinator(cfg config.Config, logger *zap.Logger, eb *eventbus.EventBus) *Coordinator {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Coordinator{
		cfg:        cfg,
		logger:     logger,
		eventBus:   eb,
		serviceMap: make(map[string]Service),
		statusMap:  make(map[string]ServiceStatus),
	}
}

// RegisterService registers a service into the coordinator.
func (c *Coordinator) RegisterService(svc Service) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	name := svc.Name()
	if _, exists := c.serviceMap[name]; exists {
		return fmt.Errorf("service %s is already registered", name)
	}

	c.services = append(c.services, svc)
	c.serviceMap[name] = svc
	c.statusMap[name] = StatusStopped

	c.logger.Info("Registered service", zap.String("service", name))
	return nil
}

// StartAll initializes and boots all registered services in order.
func (c *Coordinator) StartAll(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isBooted {
		return fmt.Errorf("coordinator is already booted")
	}

	c.logger.Info("Starting runtime services...", zap.Int("count", len(c.services)))

	for _, svc := range c.services {
		name := svc.Name()
		c.logger.Info("Starting service", zap.String("service", name))

		if err := svc.Start(ctx); err != nil {
			c.statusMap[name] = StatusFailed
			c.logger.Error("Failed to start service", zap.Error(err), zap.String("service", name))
			return fmt.Errorf("service %s failed to start: %w", name, err)
		}

		c.statusMap[name] = StatusRunning
		c.eventBus.Publish(eventbus.NewEvent("coordinator.service_started", "coordinator", map[string]interface{}{
			"service": name,
		}))
	}

	c.isBooted = true
	c.logger.Info("All runtime services started successfully")
	return nil
}

// StopAll gracefully shuts down all services in reverse order.
func (c *Coordinator) StopAll(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isBooted {
		return nil
	}

	c.logger.Info("Stopping runtime services...")
	var stopErr error

	// Stop in reverse order
	for i := len(c.services) - 1; i >= 0; i-- {
		svc := c.services[i]
		name := svc.Name()

		c.logger.Info("Stopping service", zap.String("service", name))
		if err := svc.Stop(ctx); err != nil {
			c.statusMap[name] = StatusFailed
			c.logger.Error("Failed to stop service cleanly", zap.Error(err), zap.String("service", name))
			if stopErr == nil {
				stopErr = err
			}
		} else {
			c.statusMap[name] = StatusStopped
			c.eventBus.Publish(eventbus.NewEvent("coordinator.service_stopped", "coordinator", map[string]interface{}{
				"service": name,
			}))
		}
	}

	c.isBooted = false
	c.logger.Info("All runtime services stopped")
	return stopErr
}

// HealthCheck returns a map of all services and their current health status.
func (c *Coordinator) HealthCheck() map[string]ServiceStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	healthMap := make(map[string]ServiceStatus)
	for name, svc := range c.serviceMap {
		healthMap[name] = svc.Health()
	}
	return healthMap
}

// IsBooted returns true if the coordinator is booted.
func (c *Coordinator) IsBooted() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isBooted
}

// GetEventBus returns the event bus instance.
func (c *Coordinator) GetEventBus() *eventbus.EventBus {
	return c.eventBus
}

// GetConfig returns the runtime config.
func (c *Coordinator) GetConfig() config.Config {
	return c.cfg
}

// WaitTimeout waits for context done or returns after timeout.
func (c *Coordinator) WaitTimeout(timeout time.Duration) {
	time.Sleep(timeout)
}
