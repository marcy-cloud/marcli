package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// RunBuild runs go build for macOS, Linux, and Windows.
func RunBuild(ctx context.Context) (string, error) {
	var results []string
	var allErrors []string

	// Build targets: [GOOS, GOARCH, output suffix]
	targets := [][]string{
		{"darwin", "amd64", "darwin-amd64"},
		{"darwin", "arm64", "darwin-arm64"},
		{"linux", "amd64", "linux-amd64"},
		{"linux", "arm64", "linux-arm64"},
		{"windows", "amd64", "windows-amd64.exe"},
		{"windows", "arm64", "windows-arm64.exe"},
	}

	// Create releases directory if it doesn't exist
	releasesDir := "releases"
	if err := os.MkdirAll(releasesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create releases directory: %w", err)
	}

	for _, target := range targets {
		goos, goarch, suffix := target[0], target[1], target[2]
		outputName := filepath.Join(releasesDir, fmt.Sprintf("marcli-%s", suffix))

		var out, errBuf bytes.Buffer
		cmd := exec.CommandContext(ctx, "go", "build", "-o", outputName)
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", goos), fmt.Sprintf("GOARCH=%s", goarch))
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()

		if err != nil {
			errorMsg := fmt.Sprintf("%s/%s: FAILED - %s", goos, goarch, strings.TrimSpace(errBuf.String()))
			allErrors = append(allErrors, errorMsg)
			results = append(results, errorMsg)
		} else {
			results = append(results, fmt.Sprintf("%s/%s: OK -> %s", goos, goarch, outputName))
		}
	}

	// Build for current platform (no cross-compilation)
	var finalName string
	if runtime.GOOS == "windows" {
		finalName = "marcli.exe"
	} else {
		finalName = "marcli"
	}

	var out, errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, "go", "build", "-o", finalName)
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()

	if err != nil {
		errorMsg := fmt.Sprintf("current platform (%s/%s): FAILED - %s", runtime.GOOS, runtime.GOARCH, strings.TrimSpace(errBuf.String()))
		allErrors = append(allErrors, errorMsg)
		results = append(results, errorMsg)
	} else {
		results = append(results, fmt.Sprintf("current platform (%s/%s): OK -> %s", runtime.GOOS, runtime.GOARCH, finalName))
	}

	output := strings.Join(results, "\n")
	if len(allErrors) > 0 {
		return output, fmt.Errorf("some builds failed:\n%s", strings.Join(allErrors, "\n"))
	}
	return output, nil
}
