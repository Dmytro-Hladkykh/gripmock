package gripmock

import (
	"context"
	"fmt"
	"sync"

	"github.com/gripmock/stuber"
)

var (
	// Global manager instance
	globalManager *MultiServerManager
	initOnce      sync.Once
)

// InitEmbeddedGripmock initializes embedded gripmock servers
// Call this once in TestMain or test setup
func InitEmbeddedGripmock(protoDir string, ports []int) error {
	var err error
	initOnce.Do(func() {
		globalManager = NewMultiServerManager()

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

// AddStub adds a stub to all gripmock servers
func AddStub(service, method string, input, output interface{}) error {
	if globalManager == nil {
		return fmt.Errorf("embedded gripmock not initialized - call InitEmbeddedGripmock first")
	}
	return globalManager.AddStub(service, method, input, output)
}

// Clear removes all stubs from all servers
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

// EmbeddedMocker provides a convenient interface for working with embedded gripmock
type EmbeddedMocker struct {
	server *Server
}

// NewEmbeddedMocker creates a new embedded mocker with the given server
func NewEmbeddedMocker(server *Server) *EmbeddedMocker {
	return &EmbeddedMocker{
		server: server,
	}
}

// AddStub adds a stub for the given service and method with input/output matching
func (m *EmbeddedMocker) AddStub(service, method string, input, output interface{}) error {
	stub := &stuber.Stub{
		Service: service,
		Method:  method,
		Input:   createInputData(input),
		Output:  createOutput(output),
	}

	return m.server.AddStub(stub)
}

// Clear removes all stubs from the server
func (m *EmbeddedMocker) Clear() {
	m.server.ClearStubs()
}

// GetServer returns the underlying server instance
func (m *EmbeddedMocker) GetServer() *Server {
	return m.server
}

// createInputData creates stuber.InputData from interface{}
func createInputData(input interface{}) stuber.InputData {
	if input == nil {
		return stuber.InputData{
			Matches: map[string]interface{}{},
		}
	}

	if inputMap, ok := input.(map[string]interface{}); ok {
		if matches, hasMatches := inputMap["matches"]; hasMatches {
			return stuber.InputData{
				Matches: matches.(map[string]interface{}),
			}
		}
		if equals, hasEquals := inputMap["equals"]; hasEquals {
			return stuber.InputData{
				Equals: equals.(map[string]interface{}),
			}
		}
	}

	return stuber.InputData{
		Matches: input.(map[string]interface{}),
	}
}

// createOutput creates stuber.Output from interface{}
func createOutput(output interface{}) stuber.Output {
	if output == nil {
		return stuber.Output{
			Data: map[string]interface{}{},
		}
	}

	if outputMap, ok := output.(map[string]interface{}); ok {
		if data, hasData := outputMap["data"]; hasData {
			return stuber.Output{
				Data: data.(map[string]interface{}),
			}
		}
		if errorMsg, hasError := outputMap["error"]; hasError {
			return stuber.Output{
				Error: errorMsg.(string),
			}
		}
	}

	return stuber.Output{
		Data: output.(map[string]interface{}),
	}
}
