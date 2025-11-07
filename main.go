package main

import (
	"context"
	"fmt"
	"os"

	"marcli/cmd"

	logger "github.com/charmbracelet/log"
)

// commandRegistry maps CLI names to command functions
var commandRegistry = make(map[string]func(context.Context) (string, error))

// initCommands populates the command registry
func initCommands() {
	commandRegistry["go-echo"] = cmd.RunGoEcho
	commandRegistry["ps-echo"] = cmd.RunPSEcho
	commandRegistry["bash-echo"] = cmd.RunBashEcho
	commandRegistry["build"] = cmd.RunBuild
}

func main() {
	// Initialize command registry
	initCommands()

	args := os.Args[1:]

	// CLI mode: if args provided, run command directly
	if len(args) > 0 {
		cmdName := args[0]
		cmd, exists := commandRegistry[cmdName]
		if !exists {
			logger.Fatal("unknown command", "command", cmdName)
		}

		ctx := context.Background()
		out, err := cmd(ctx)
		if err != nil {
			logger.Fatal("command failed", "err", err)
		}
		fmt.Print(out)
		return
	}

	// TUI mode: no args, show interactive menu (default)
	if err := cmd.RunTUI(); err != nil {
		logger.Fatal("error", "err", err)
	}
}
