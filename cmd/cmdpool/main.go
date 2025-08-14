package main

import (
	"fmt"
	"os"

	"github.com/pashkov256/cmdpool/internal/app"
	"github.com/pashkov256/cmdpool/internal/cli"
)

func main() {
	// Parse command line arguments
	if len(os.Args) > 1 {
		// CLI mode
		if err := cli.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// TUI mode (default)
	if err := app.RunTUI(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
} 