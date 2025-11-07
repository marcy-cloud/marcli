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
	commandRegistry["mega-combine"] = cmd.RunMegaCombine
	commandRegistry["cutiepie"] = cmd.RunCutiepieTUICommand
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

		// Handle flags for specific commands
		if cmdName == "mega-combine" {
			for i := 1; i < len(args); i++ {
				if args[i] == "--test" {
					ctx = context.WithValue(ctx, "megaCombineTestMode", true)
				} else if args[i] == "--out" && i+1 < len(args) {
					ctx = context.WithValue(ctx, "megaCombineOutput", args[i+1])
					i++ // Skip the next argument since we consumed it
				} else if args[i] == "--waytoobig" {
					ctx = context.WithValue(ctx, "megaCombineWayTooBig", true)
				} else if args[i] == "--slowbutsmall" {
					ctx = context.WithValue(ctx, "megaCombineSlowButSmall", true)
				}
			}
		}
		if cmdName == "build" && len(args) > 1 && args[1] == "--fast" {
			ctx = context.WithValue(ctx, "buildFastMode", true)
		}

		out, err := cmd(ctx)
		if err != nil {
			logger.Fatal("command failed", "err", err)
		}
		fmt.Print(out)
		return
	}

	// TUI mode: no args, show the cutiepie interactive menu (default) 🎀
	if err := cmd.RunCutiepieTUI(); err != nil {
		logger.Fatal("error", "err", err)
	}
}
