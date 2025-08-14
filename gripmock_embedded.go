package gripmock

import (
	"context"
	"fmt"
	"time"
)

// EmbeddedGripmockHelper provides a convenient way to use embedded gripmock in tests
type EmbeddedGripmockHelper struct {
	mocker *EmbeddedMocker
	server *Server
}

// NewEmbeddedGripmockHelper creates a new helper with the given port and proto files
func NewEmbeddedGripmockHelper(port int, protoFiles []string) (*EmbeddedGripmockHelper, error) {
	server, err := NewServer(port, protoFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	mocker := NewEmbeddedMocker(server)
	
	return &EmbeddedGripmockHelper{
		mocker: mocker,
		server: server,
	}, nil
}

// Start starts the gripmock server
func (h *EmbeddedGripmockHelper) Start(ctx context.Context) error {
	if err := h.server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	// Wait for server to be ready
	if err := h.server.WaitForReady(5 * time.Second); err != nil {
		h.server.Stop()
		return fmt.Errorf("server not ready: %w", err)
	}
	
	return nil
}

// Stop stops the gripmock server
func (h *EmbeddedGripmockHelper) Stop() {
	h.server.Stop()
}

// AddStub adds a stub using the simplified interface (compatible with your existing code)
func (h *EmbeddedGripmockHelper) AddStub(service, method string, input, output interface{}) error {
	return h.mocker.AddStub(service, method, input, output)
}

// Clear removes all stubs
func (h *EmbeddedGripmockHelper) Clear() {
	h.mocker.Clear()
}

// GetPort returns the server port
func (h *EmbeddedGripmockHelper) GetPort() int {
	return h.server.GetPort()
}

// IsRunning returns true if server is running
func (h *EmbeddedGripmockHelper) IsRunning() bool {
	return h.server.IsRunning()
}