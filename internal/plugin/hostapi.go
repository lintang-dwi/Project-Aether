package plugin

import "context"

// HostAPI defines runtime functions exposed to WASM plugins.
type HostAPI interface {
	LogInfo(ctx context.Context, message string)
	GetVersion() string
}

type DefaultHostAPI struct{}

func (h *DefaultHostAPI) LogInfo(ctx context.Context, message string) {}
func (h *DefaultHostAPI) GetVersion() string                          { return "0.1.0" }
