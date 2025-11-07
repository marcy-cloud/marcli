package cmd

import (
	"context"
	"fmt"
)

// RunVersion displays the current version and build number - so cute! ✨
func RunVersion(ctx context.Context) (string, error) {
	version, err := GetVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}
	return fmt.Sprintf("marcli %s\n", version), nil
}

