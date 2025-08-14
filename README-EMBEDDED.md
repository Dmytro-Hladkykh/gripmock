# Gripmock Embedded

**Fork of [tokopedia/gripmock](https://github.com/tokopedia/gripmock)** - Embedded gripmock library for Go that provides in-process gRPC mocking without external dependencies.

> This is a fork of the original gripmock project, modified to support embedded/in-process gRPC mocking for better integration testing experience.

## Quick Start

### 1. Initialize embedded gripmock in your `TestMain`:

```go
func TestMain(m *testing.M) {
    err := gripmock.InitEmbeddedGripmock("../../../protos", []int{4771, 4772, 4773, 4774, 4775})
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize embedded gripmock: %v", err))
    }
    defer gripmock.StopEmbeddedGripmock()

    // Run all tests
    os.Exit(m.Run())
}
```

### 2. Use in tests:

```go
func TestExample(t *testing.T) {
    // Clear any existing stubs
    err := gripmock.Clear()
    require.NoError(t, err)

    // Add stubs
    err = gripmock.AddStub("proto_internal.ProtoInternalService", "transfer_tokens", nil, nil)
    require.NoError(t, err)

    // Your logic tests go here...
}
```

## File Structure

- **`gripmock.go`** - Core server implementation
- **`embedded.go`** - Main API and drop-in replacement functions
- **`manager.go`** - Multi-server management  
- **`mocker.go`** - gRPC request handling and protobuf conversion

## API Reference

### Global Functions

- `InitEmbeddedGripmock(protoDir, ports)` - Initialize servers
- `StopEmbeddedGripmock()` - Stop all servers
- `AddStub(service, method, input, output)` - Add mock stub
- `Clear()` - Remove all stubs
- `GetActivePorts()` - Get running server ports
- `IsRunning()` - Check if servers are running

### Advanced Usage

For more control, you can use the underlying types:

```go
// Create individual servers
server, err := gripmock.NewServer(9001, []string{"path/to/service.proto"})
if err != nil {
    log.Fatal(err)
}

// Start server
ctx := context.Background()
if err := server.Start(ctx); err != nil {
    log.Fatal(err)
}
defer server.Stop()

// Create mocker for the server
mocker := gripmock.NewEmbeddedMocker(server)
err = mocker.AddStub("MyService", "MyMethod", inputData, outputData)
```

## Requirements

- Go 1.24+