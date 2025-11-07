package cmd

import (
	"context"

	logger "github.com/charmbracelet/log"
)

// RunGoEcho runs a pure Go echo command without external processes - so clean! 💕
func RunGoEcho(ctx context.Context) (string, error) {
	logger.Info("Running Go echo")
	return "Golang echo", nil
}

