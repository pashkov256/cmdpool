package executor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Command represents a running command
type Command struct {
	ID          string
	Name        string
	Command     string
	Dir         string
	Status      CommandStatus
	Output      []string
	Error       error
	StartTime   time.Time
	EndTime     time.Time
	Process     *os.Process
	AutoRestart bool
	mu          sync.RWMutex
}

// CommandStatus represents the status of a command
type CommandStatus string

const (
	StatusPending CommandStatus = "pending"
	StatusRunning CommandStatus = "running"
	StatusDone    CommandStatus = "done"
	StatusFailed  CommandStatus = "failed"
	StatusStopped CommandStatus = "stopped"
)

// Executor manages multiple command executions
type Executor struct {
	commands map[string]*Command
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewExecutor creates a new command executor
func NewExecutor() *Executor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Executor{
		commands: make(map[string]*Command),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// RunCommands executes multiple commands simultaneously
func (e *Executor) RunCommands(commands []string) error {
	var wg sync.WaitGroup

	for i, cmdStr := range commands {
		wg.Add(1)
		go func(id int, command string) {
			defer wg.Done()
			e.runCommand(fmt.Sprintf("cmd_%d", id), command, ".", false)
		}(i, cmdStr)
	}

	wg.Wait()
	return nil
}

// RunCommand executes a single command (public method)
func (e *Executor) RunCommand(id, command, dir string, autoRestart bool) {
	e.runCommand(id, command, dir, autoRestart)
}

// runCommand executes a single command (private implementation)
func (e *Executor) runCommand(id, command, dir string, autoRestart bool) {
	cmd := &Command{
		ID:          id,
		Name:        command,
		Command:     command,
		Dir:         dir,
		Status:      StatusPending,
		Output:      make([]string, 0),
		AutoRestart: autoRestart,
		StartTime:   time.Now(),
	}

	e.mu.Lock()
	e.commands[id] = cmd
	e.mu.Unlock()

	// Execute command
	e.executeCommand(cmd)
}

// executeCommand runs the actual command
func (e *Executor) executeCommand(cmd *Command) {
	// Parse command and arguments
	args := parseCommand(cmd.Command)
	if len(args) == 0 {
		cmd.setError(fmt.Errorf("empty command"))
		return
	}

	// Create exec.Cmd
	execCmd := exec.CommandContext(e.ctx, args[0], args[1:]...)
	execCmd.Dir = cmd.Dir

	// Set up pipes for stdout and stderr
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		cmd.setError(fmt.Errorf("failed to create stdout pipe: %w", err))
		return
	}

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		cmd.setError(fmt.Errorf("failed to create stderr pipe: %w", err))
		return
	}

	// Start command
	if err := execCmd.Start(); err != nil {
		cmd.setError(fmt.Errorf("failed to start command: %w", err))
		return
	}

	cmd.Process = execCmd.Process
	cmd.Status = StatusRunning

	// Read output in separate goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			cmd.addOutput(scanner.Text())
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			cmd.addErrorOutput(scanner.Text())
		}
	}()

	// Wait for command to complete
	err = execCmd.Wait()
	wg.Wait()

	cmd.EndTime = time.Now()

	if err != nil {
		cmd.setError(err)
	} else {
		cmd.Status = StatusDone
	}
}

// parseCommand splits a command string into command and arguments
func parseCommand(cmdStr string) []string {
	// Simple parsing - split by spaces
	// In a real implementation, you might want more sophisticated parsing
	var args []string
	var current string
	var inQuotes bool

	for _, char := range cmdStr {
		if char == '"' {
			inQuotes = !inQuotes
		} else if char == ' ' && !inQuotes {
			if current != "" {
				args = append(args, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		args = append(args, current)
	}

	return args
}

// setError sets the error status and message
func (c *Command) setError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Error = err
	c.Status = StatusFailed
}

// addOutput adds a line to the output
func (c *Command) addOutput(line string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Output = append(c.Output, line)

	// Keep only last 1000 lines
	if len(c.Output) > 1000 {
		c.Output = c.Output[len(c.Output)-1000:]
	}
}

// addErrorOutput adds a line to the output (treating stderr as output)
func (c *Command) addErrorOutput(line string) {
	c.addOutput("[STDERR] " + line)
}

// GetOutput returns a copy of the command output
func (c *Command) GetOutput() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]string, len(c.Output))
	copy(result, c.Output)
	return result
}

// GetCommands returns all commands
func (e *Executor) GetCommands() map[string]*Command {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make(map[string]*Command)
	for k, v := range e.commands {
		result[k] = v
	}
	return result
}

// StopCommand stops a running command
func (e *Executor) StopCommand(id string) error {
	e.mu.RLock()
	cmd, exists := e.commands[id]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("command %s not found", id)
	}

	if cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
		cmd.Status = StatusStopped
		cmd.EndTime = time.Now()
	}

	return nil
}

// RestartCommand restarts a command
func (e *Executor) RestartCommand(id string) error {
	e.mu.RLock()
	cmd, exists := e.commands[id]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("command %s not found", id)
	}

	// Stop if running
	if cmd.Status == StatusRunning {
		if err := e.StopCommand(id); err != nil {
			return err
		}
	}

	// Reset command state
	cmd.mu.Lock()
	cmd.Status = StatusPending
	cmd.Output = make([]string, 0)
	cmd.Error = nil
	cmd.StartTime = time.Now()
	cmd.EndTime = time.Time{}
	cmd.Process = nil
	cmd.mu.Unlock()

	// Restart
	go e.executeCommand(cmd)
	return nil
}

// Stop stops all running commands
func (e *Executor) Stop() {
	e.cancel()

	e.mu.RLock()
	for _, cmd := range e.commands {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}
	e.mu.RUnlock()
}
