package coordinator_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"aether/internal/config"
	"aether/internal/coordinator"
	"aether/internal/eventbus"
	"aether/internal/observability"
)

type MockService struct {
	name      string
	status    coordinator.ServiceStatus
	failStart bool
}

func (m *MockService) Name() string {
	return m.name
}

func (m *MockService) Start(ctx context.Context) error {
	if m.failStart {
		m.status = coordinator.StatusFailed
		return fmt.Errorf("simulated start failure for %s", m.name)
	}
	m.status = coordinator.StatusRunning
	return nil
}

func (m *MockService) Stop(ctx context.Context) error {
	m.status = coordinator.StatusStopped
	return nil
}

func (m *MockService) Health() coordinator.ServiceStatus {
	return m.status
}

func TestCoordinator_Lifecycle(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := observability.InitLogger("development", "debug")
	eb := eventbus.NewEventBus(100)
	defer eb.Close()

	coord := coordinator.NewCoordinator(cfg, logger, eb)

	svc1 := &MockService{name: "ServiceA"}
	svc2 := &MockService{name: "ServiceB"}
	svc3 := &MockService{name: "ServiceC"}

	if err := coord.RegisterService(svc1); err != nil {
		t.Fatalf("failed to register svc1: %v", err)
	}
	if err := coord.RegisterService(svc2); err != nil {
		t.Fatalf("failed to register svc2: %v", err)
	}
	if err := coord.RegisterService(svc3); err != nil {
		t.Fatalf("failed to register svc3: %v", err)
	}

	// Verify duplicate registration fails
	if err := coord.RegisterService(svc1); err == nil {
		t.Fatal("expected error on duplicate service registration, got nil")
	}

	ctx := context.Background()

	// Track start events
	var eventMu sync.Mutex
	startedEvents := make([]string, 0)
	eb.Subscribe("coordinator.service_started", func(e eventbus.Event) {
		eventMu.Lock()
		startedEvents = append(startedEvents, e.Payload["service"].(string))
		eventMu.Unlock()
	})

	// Start all
	if err := coord.StartAll(ctx); err != nil {
		t.Fatalf("failed to start all services: %v", err)
	}

	if !coord.IsBooted() {
		t.Error("expected coordinator to be booted")
	}

	// Verify HealthCheck
	healthMap := coord.HealthCheck()
	if healthMap["ServiceA"] != coordinator.StatusRunning ||
		healthMap["ServiceB"] != coordinator.StatusRunning ||
		healthMap["ServiceC"] != coordinator.StatusRunning {
		t.Errorf("expected all services running, got %v", healthMap)
	}

	time.Sleep(50 * time.Millisecond)

	eventMu.Lock()
	if len(startedEvents) != 3 {
		t.Errorf("expected 3 start events, got %d", len(startedEvents))
	}
	eventMu.Unlock()

	// Stop all
	if err := coord.StopAll(ctx); err != nil {
		t.Fatalf("failed to stop all services: %v", err)
	}

	if coord.IsBooted() {
		t.Error("expected coordinator to be stopped")
	}

	healthMapAfter := coord.HealthCheck()
	if healthMapAfter["ServiceA"] != coordinator.StatusStopped {
		t.Errorf("expected ServiceA stopped, got %s", healthMapAfter["ServiceA"])
	}
}

func TestCoordinator_StartFailure(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := observability.InitLogger("development", "debug")
	eb := eventbus.NewEventBus(100)
	defer eb.Close()

	coord := coordinator.NewCoordinator(cfg, logger, eb)

	svcOK := &MockService{name: "GoodService"}
	svcBad := &MockService{name: "BadService", failStart: true}

	_ = coord.RegisterService(svcOK)
	_ = coord.RegisterService(svcBad)

	ctx := context.Background()
	err := coord.StartAll(ctx)
	if err == nil {
		t.Fatal("expected error on service start failure, got nil")
	}

	healthMap := coord.HealthCheck()
	if healthMap["BadService"] != coordinator.StatusFailed {
		t.Errorf("expected BadService status FAILED, got %s", healthMap["BadService"])
	}
}
