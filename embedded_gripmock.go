// Package gripmock provides embedded gripmock functionality as a drop-in replacement
// for the original HTTP-based gripmock client
package gripmock

import (
	"context"
	"fmt"
	"sync"
)

var (
	// Global manager instance - similar to your original global variables
	globalManager *MultiServerManager
	initOnce      sync.Once
)

// InitEmbeddedGripmock initializes embedded gripmock servers to replace your docker-compose setup
// Call this once in TestMain or test setup, similar to how you started docker containers
func InitEmbeddedGripmock(protoDir string, ports []int) error {
	var err error
	initOnce.Do(func() {
		globalManager = NewMultiServerManager()
		
		// Create server configurations matching your docker-compose setup
		configs := make([]ServerConfig, len(ports))
		for i, port := range ports {
			configs[i] = ServerConfig{
				Port:       port,
				ProtoDir:   protoDir,
				Identifier: fmt.Sprintf("gripmock-%d", i),
			}
		}

		// Start all servers
		ctx := context.Background()
		err = globalManager.StartServers(ctx, configs)
	})
	
	return err
}

// StopEmbeddedGripmock stops all embedded gripmock servers
// Call this in test teardown
func StopEmbeddedGripmock() {
	if globalManager != nil {
		globalManager.StopAll()
	}
}

// AddStub adds a stub to all gripmock servers - DROP-IN REPLACEMENT for your original function
// This maintains the exact same signature and behavior as your original AddStub function
func AddStub(service, method string, input, output interface{}) error {
	if globalManager == nil {
		return fmt.Errorf("embedded gripmock not initialized - call InitEmbeddedGripmock first")
	}
	return globalManager.AddStub(service, method, input, output)
}

// Clear removes all stubs from all servers - DROP-IN REPLACEMENT for your original function
// This maintains the exact same signature and behavior as your original Clear function
func Clear() error {
	if globalManager == nil {
		return fmt.Errorf("embedded gripmock not initialized - call InitEmbeddedGripmock first")
	}
	globalManager.Clear()
	return nil
}

// GetActivePorts returns the ports of all running gripmock servers
// Useful for debugging or integration with other services
func GetActivePorts() []int {
	if globalManager == nil {
		return nil
	}
	return globalManager.GetServerPorts()
}

// IsRunning returns true if all gripmock servers are running
func IsRunning() bool {
	if globalManager == nil {
		return false
	}
	return globalManager.IsRunning()
}

// Helper functions are defined in helper.go