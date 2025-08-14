package gripmock

import (
	"github.com/gripmock/stuber"
)

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
