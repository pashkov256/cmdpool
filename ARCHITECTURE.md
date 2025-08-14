# cmdpool Architecture

## 🏗️ Project Structure

```
cmdpool/
├── cmd/
│   └── cmdpool/          # Main application entry point
│       └── main.go
├── internal/              # Internal packages (not importable from outside)
│   ├── app/              # TUI application logic
│   │   └── tui.go
│   ├── cli/              # CLI command handling
│   │   └── cli.go
│   ├── config/           # Configuration management
│   │   └── config.go
│   └── executor/         # Command execution engine
│       └── executor.go
├── .cmdpool.yml          # Example configuration
├── go.mod                # Go module definition
├── Makefile              # Build and development tasks
└── README.md             # Project documentation
```

## 🎯 Design Principles

### 1. **Separation of Concerns**
- **CLI**: Handles command-line argument parsing and execution
- **TUI**: Manages the interactive terminal interface
- **Executor**: Handles command execution and output capture
- **Config**: Manages configuration loading and validation

### 2. **Modular Architecture**
- Each package has a single responsibility
- Clear interfaces between components
- Easy to test individual components
- Simple to extend with new features

### 3. **Concurrent Execution**
- Commands run in separate goroutines
- Thread-safe output collection
- Non-blocking UI updates
- Graceful shutdown handling

## 🔧 Core Components

### Executor Package (`internal/executor/`)

The executor is the heart of the system, responsible for:

- **Command Management**: Creating, starting, stopping, and restarting commands
- **Output Capture**: Capturing stdout and stderr in real-time
- **Status Tracking**: Monitoring command execution status
- **Process Control**: Managing OS processes and cleanup

```go
type Executor struct {
    commands map[string]*Command
    mu       sync.RWMutex
    ctx      context.Context
    cancel   context.CancelFunc
}
```

**Key Features:**
- Concurrent command execution
- Real-time output streaming
- Automatic output buffering (last 1000 lines)
- Process lifecycle management
- Error handling and recovery

### TUI Package (`internal/app/`)

The TUI provides an interactive interface built with `tview`:

- **Panel Management**: Dynamic command panels with borders and titles
- **Navigation**: Arrow key navigation between panels
- **Real-time Updates**: Live status and output updates
- **Keyboard Shortcuts**: Quick actions (restart, stop, add, quit)

**UI Layout:**
```
┌─────────────────────────────────────────────────────────┐
│                    Command Panel 1                      │
│  [Status] 🟢 Running                                  │
│  Output line 1...                                      │
│  Output line 2...                                      │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    Command Panel 2                      │
│  [Status] ✅ Done                                      │
│  Output line 1...                                      │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│              Status: Running: 2 | Done: 1              │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│ ↑↓: Navigate | Enter: Expand | r: Restart | s: Stop   │
└─────────────────────────────────────────────────────────┘
```

### CLI Package (`internal/cli/`)

The CLI provides command-line functionality:

- **Argument Parsing**: Using Cobra for robust CLI handling
- **Configuration Loading**: Support for YAML config files
- **Command Sets**: Running predefined command groups
- **Flexible Input**: Commands via arguments, flags, or config files

**Usage Patterns:**
```bash
# Direct commands
cmdpool "ping google.com" "ping github.com"

# From config file
cmdpool -config .cmdpool.yml

# Specific command set
cmdpool -config .cmdpool.yml -set backend
```

### Config Package (`internal/config/`)

Configuration management with YAML support:

- **Command Sets**: Grouped commands with metadata
- **Global Settings**: Application-wide configuration
- **Environment Variables**: Per-command environment setup
- **Auto-restart**: Automatic restart on failure

## 🚀 Data Flow

### 1. **Command Execution Flow**
```
User Input → CLI/TUI → Executor → OS Process → Output Capture → UI Update
```

### 2. **Output Streaming**
```
Process stdout/stderr → Buffered Scanner → Command Output → Panel Display
```

### 3. **Status Updates**
```
Process State Change → Status Update → UI Refresh → Visual Feedback
```

## 🔒 Thread Safety

- **Executor**: Uses `sync.RWMutex` for command map access
- **Commands**: Individual commands have their own mutex for output
- **UI Updates**: All UI updates are queued through `QueueUpdateDraw`
- **Process Management**: Safe process creation and termination

## 📊 Performance Considerations

### **Memory Management**
- Output buffering (configurable limit)
- Automatic cleanup of completed commands
- Efficient string handling for large outputs

### **CPU Usage**
- Non-blocking UI updates (100ms refresh rate)
- Efficient output scanning with `bufio.Scanner`
- Minimal goroutine overhead

### **I/O Optimization**
- Direct process pipe reading
- Buffered output collection
- Minimal string copying

## 🧪 Testing Strategy

### **Unit Tests**
- Individual package testing
- Mock interfaces for external dependencies
- Isolated command execution testing

### **Integration Tests**
- End-to-end command execution
- Configuration file loading
- TUI interaction testing

### **Performance Tests**
- Concurrent command execution
- Memory usage under load
- Output buffering efficiency

## 🔮 Future Enhancements

### **Short Term**
- [ ] WebSocket support for remote monitoring
- [ ] Plugin system for output formatters
- [ ] Advanced filtering and search

### **Medium Term**
- [ ] Integration with CI/CD tools
- [ ] Metrics collection and export
- [ ] Custom panel layouts

### **Long Term**
- [ ] Web-based interface
- [ ] Distributed command execution
- [ ] Machine learning for command optimization

## 🛠️ Development Workflow

### **Local Development**
```bash
# Install dependencies
make deps

# Run in development mode
make dev

# Run with CLI arguments
make dev-cli

# Run tests
make test
```

### **Building**
```bash
# Build for current platform
make build

# Build for all platforms
make release

# Install to system
make install
```

### **Code Quality**
```bash
# Format code
make fmt

# Lint code
make lint

# Run tests with coverage
make test-coverage
```

## 📚 Dependencies

### **Core Dependencies**
- **tview**: Terminal UI framework
- **tcell**: Terminal cell manipulation
- **cobra**: CLI command framework
- **viper**: Configuration management
- **yaml.v3**: YAML parsing

### **Development Dependencies**
- **golangci-lint**: Code linting
- **godoc**: Documentation generation

## 🔧 Configuration Schema

The configuration file supports:

```yaml
commands:
  <set_name>:
    name: "Display Name"
    description: "Description"
    commands: ["command1", "command2"]
    dir: "./working/directory"
    auto_restart: true
    env: ["KEY=value"]

global:
  log_file: "filename.log"
  max_output_lines: 1000
  refresh_rate_ms: 100
```

This architecture provides a solid foundation for a robust, scalable command execution tool with both CLI and TUI interfaces. 