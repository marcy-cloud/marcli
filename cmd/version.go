package cmd

import (
	"context"
	"fmt"
)

// RunVersion displays the current version and build number - so cute! ✨
func RunVersion(ctx context.Context) (string, error) {
	return fmt.Sprintf("marcli %s (build %s)\n", Version, Build), nil
}

