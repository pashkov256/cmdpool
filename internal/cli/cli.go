package cli

import (
	"fmt"
	"time"

	"github.com/pashkov256/cmdpool/internal/config"
	"github.com/pashkov256/cmdpool/internal/executor"
	"github.com/spf13/cobra"
)

var (
	configFile string
	commandSet string
	commands   []string
)

// Run initializes and runs the CLI
func Run() error {
	var rootCmd = &cobra.Command{
		Use:   "cmdpool [commands...]",
		Short: "Run multiple commands simultaneously with real-time monitoring",
		Long: `cmdpool is a powerful CLI/TUI utility that allows you to run multiple 
commands simultaneously while displaying their real-time output in separate terminal panels.

Examples:
  cmdpool "ping google.com" "ping github.com"
  cmdpool -config .cmdpool.yml
  cmdpool -set backend`,
		RunE: runCommands,
	}

	// Add flags
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	rootCmd.Flags().StringVarP(&commandSet, "set", "s", "", "Command set name from config")
	rootCmd.Flags().StringArrayVarP(&commands, "command", "e", []string{}, "Commands to execute")

	return rootCmd.Execute()
}

func runCommands(cmd *cobra.Command, args []string) error {
	var cmds []string

	// If commands provided via flags, use them
	if len(commands) > 0 {
		cmds = commands
	} else if len(args) > 0 {
		// If commands provided as arguments, use them
		cmds = args
	} else if configFile != "" {
		// Load from config file
		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if commandSet != "" {
			// Run specific command set
			set, exists := cfg.CommandSets[commandSet]
			if !exists {
				return fmt.Errorf("command set '%s' not found", commandSet)
			}
			cmds = set.Commands
		} else {
			// Run all commands from config
			for _, set := range cfg.CommandSets {
				cmds = append(cmds, set.Commands...)
			}
		}
	} else {
		return fmt.Errorf("no commands specified. Use -e flag, provide arguments, or use -config")
	}

	if len(cmds) == 0 {
		return fmt.Errorf("no commands to execute")
	}

	fmt.Printf("Starting %d commands...\n", len(cmds))
	for i, cmdStr := range cmds {
		fmt.Printf("[%d] %s\n", i+1, cmdStr)
	}
	fmt.Println()

	// Execute commands
	exec := executor.NewExecutor()

	// Start commands
	for i, cmdStr := range cmds {
		go exec.RunCommand(fmt.Sprintf("cmd_%d", i), cmdStr, ".", false)
	}

	// Monitor and display output
	return monitorCommands(exec, cmds)
}

// monitorCommands monitors running commands and displays their output
func monitorCommands(exec *executor.Executor, commands []string) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Keep track of completed commands
	completed := make(map[string]bool)

	for {
		select {
		case <-ticker.C:
			cmds := exec.GetCommands()

			// Check if all commands are completed
			allCompleted := true
			for _, cmd := range cmds {
				if cmd.Status != executor.StatusDone && cmd.Status != executor.StatusFailed && cmd.Status != executor.StatusStopped {
					allCompleted = false
					break
				}
			}

			if allCompleted && len(cmds) > 0 {
				// Show final results
				fmt.Println("\n=== Final Results ===")
				for _, cmd := range cmds {
					status := "✅"
					if cmd.Status == executor.StatusFailed {
						status = "🔴"
					} else if cmd.Status == executor.StatusStopped {
						status = "⏹️"
					}
					fmt.Printf("%s %s: %s\n", status, cmd.ID, cmd.Status)
					if cmd.Error != nil {
						fmt.Printf("   Error: %v\n", cmd.Error)
					}
				}
				return nil
			}

			// Display current output for each command
			for _, cmd := range cmds {
				if !completed[cmd.ID] {
					output := cmd.GetOutput()
					if len(output) > 0 {
						fmt.Printf("\n[%s] %s:\n", cmd.ID, cmd.Status)
						// Show last few lines of output
						start := 0
						if len(output) > 5 {
							start = len(output) - 5
						}
						for _, line := range output[start:] {
							fmt.Printf("  %s\n", line)
						}
					}

					// Mark as completed if done
					if cmd.Status == executor.StatusDone || cmd.Status == executor.StatusFailed || cmd.Status == executor.StatusStopped {
						completed[cmd.ID] = true
					}
				}
			}
		}
	}
}
