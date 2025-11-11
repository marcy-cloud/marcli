package cmd

import (
	"context"
	"fmt"
	"marcli/api"
)

// RunCutiepieTTY starts the web-based terminal server
func RunCutiepieTTY(ctx context.Context) (string, error) {
	// Get port from context if available
	port := 8080
	if ctx.Value("port") != nil {
		if p, ok := ctx.Value("port").(int); ok {
			port = p
		}
	}

	// Start the server (this will block)
	err := api.StartServer(port)
	if err != nil {
		return "", fmt.Errorf("server error: %w", err)
	}

	return "", nil
}

