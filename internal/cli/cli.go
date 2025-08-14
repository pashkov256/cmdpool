package cli

import (
	"fmt"

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

	// Execute commands
	exec := executor.NewExecutor()
	return exec.RunCommands(cmds)
}
