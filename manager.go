package gripmock

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// MultiServerManager manages multiple embedded gripmock servers
type MultiServerManager struct {
	servers map[int]*EmbeddedMocker // map[port]*EmbeddedMocker
	mu      sync.RWMutex
}

// ServerConfig represents configuration for a single gripmock server
type ServerConfig struct {
	Port       int
	ProtoDir   string
	Identifier string // optional identifier for logging
}

// NewMultiServerManager creates a new manager for multiple gripmock servers
func NewMultiServerManager() *MultiServerManager {
	return &MultiServerManager{
		servers: make(map[int]*EmbeddedMocker),
	}
}

// StartServers starts multiple gripmock servers with the given configurations
func (m *MultiServerManager) StartServers(ctx context.Context, configs []ServerConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, config := range configs {
		// Discover proto files from directory
		protoFiles, err := discoverProtoFiles(config.ProtoDir)
		if err != nil {
			return fmt.Errorf("failed to discover proto files in %s: %w", config.ProtoDir, err)
		}

		// Create server
		server, err := NewServer(config.Port, protoFiles)
		if err != nil {
			return fmt.Errorf("failed to create server on port %d: %w", config.Port, err)
		}

		// Start server
		if err := server.Start(ctx); err != nil {
			return fmt.Errorf("failed to start server on port %d: %w", config.Port, err)
		}

		// Wait for server to be ready
		if err := server.WaitForReady(5 * time.Second); err != nil {
			server.Stop()
			return fmt.Errorf("server on port %d not ready: %w", config.Port, err)
		}

		// Create embedded mocker
		mocker := NewEmbeddedMocker(server)
		m.servers[config.Port] = mocker

		fmt.Printf("Started gripmock server on port %d with %d proto files\n", config.Port, len(protoFiles))
	}

	return nil
}

// StopAll stops all running servers
func (m *MultiServerManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for port, mocker := range m.servers {
		mocker.GetServer().Stop()
		fmt.Printf("Stopped gripmock server on port %d\n", port)
	}
	m.servers = make(map[int]*EmbeddedMocker)
}

// AddStub adds a stub to all running servers
func (m *MultiServerManager) AddStub(service, method string, input, output interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.servers) == 0 {
		return fmt.Errorf("no servers running")
	}

	var lastErr error
	for port, mocker := range m.servers {
		if err := mocker.AddStub(service, method, input, output); err != nil {
			lastErr = fmt.Errorf("failed to add stub to server on port %d: %w", port, err)
		}
	}

	return lastErr
}

// Clear removes all stubs from all servers
func (m *MultiServerManager) Clear() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, mocker := range m.servers {
		mocker.Clear()
	}
}

// GetServerPorts returns all active server ports
func (m *MultiServerManager) GetServerPorts() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ports := make([]int, 0, len(m.servers))
	for port := range m.servers {
		ports = append(ports, port)
	}
	return ports
}

// GetServer returns the embedded mocker for a specific port
func (m *MultiServerManager) GetServer(port int) (*EmbeddedMocker, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mocker, exists := m.servers[port]
	return mocker, exists
}

// IsRunning returns true if all servers are running
func (m *MultiServerManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, mocker := range m.servers {
		if !mocker.GetServer().IsRunning() {
			return false
		}
	}
	return len(m.servers) > 0
}

// discoverProtoFiles recursively finds all .proto files in the given directory
func discoverProtoFiles(protoDir string) ([]string, error) {
	if protoDir == "" {
		return nil, fmt.Errorf("proto directory not specified")
	}

	if _, err := os.Stat(protoDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("proto directory does not exist: %s", protoDir)
	}

	var protoFiles []string

	err := filepath.Walk(protoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".proto") {
			protoFiles = append(protoFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking proto directory: %w", err)
	}

	if len(protoFiles) == 0 {
		return nil, fmt.Errorf("no .proto files found in directory: %s", protoDir)
	}

	return protoFiles, nil
}
