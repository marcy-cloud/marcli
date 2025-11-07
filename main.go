package main

import (
	"context"
	"fmt"
	"os"

	"marcli/cmd"

	logger "github.com/charmbracelet/log"
)

// commandRegistry maps CLI names to command functions - so organized! 💕
var commandRegistry = make(map[string]func(context.Context) (string, error))

// initCommands populates the command registry with all our cute commands! ✨
func initCommands() {
	commandRegistry["go-echo"] = cmd.RunGoEcho
	commandRegistry["ps-echo"] = cmd.RunPSEcho
	commandRegistry["bash-echo"] = cmd.RunBashEcho
	commandRegistry["build"] = cmd.RunBuild
	commandRegistry["version"] = cmd.RunVersion
}

func main() {
	// Initialize our cute command registry! 💖
	initCommands()

	args := os.Args[1:]

	// CLI mode: if args provided, run command directly (so efficient!) 💅
	if len(args) > 0 {
		cmdName := args[0]

		// Handle flag aliases - we're so flexible! ✨
		if cmdName == "-v" || cmdName == "--version" {
			cmdName = "version"
		}

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

	// TUI mode: no args, show the cute interactive menu (default) 🎀
	if err := cmd.RunTUI(); err != nil {
		logger.Fatal("error", "err", err)
	}
}
