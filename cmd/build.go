package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// RunBuild runs go build for macOS, Linux, and Windows - building everything with love! 💖
func RunBuild(ctx context.Context) (string, error) {
	// Check if fast mode is enabled
	fastMode := ctx.Value("buildFastMode") == true

	// Increment build number - we're so organized! 🎀
	if err := IncrementBuild(); err != nil {
		return "", fmt.Errorf("failed to increment build number: %w", err)
	}

	// Load and display version info - keeping track of our progress! ✨
	config, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	var results []string
	results = append(results, fmt.Sprintf("Building version %s (build %d)", config.Version, config.Build))
	results = append(results, "")

	var allErrors []string

	// Build ldflags to embed version and build - so embedded! ✨
	ldflags := fmt.Sprintf("-X marcli/cmd.Version=%s -X marcli/cmd.Build=%d", config.Version, config.Build)

	// Skip cross-platform builds in fast mode
	if !fastMode {
		// Build targets: [GOOS, GOARCH, output suffix] - building for everyone! 🌈
		targets := [][]string{
			{"darwin", "amd64", "darwin-amd64"},
			{"darwin", "arm64", "darwin-arm64"},
			{"linux", "amd64", "linux-amd64"},
			{"linux", "arm64", "linux-arm64"},
			{"windows", "amd64", "windows-amd64.exe"},
			{"windows", "arm64", "windows-arm64.exe"},
		}

		// Create releases directory if it doesn't exist - so organized! 💅
		releasesDir := "releases"
		if err := os.MkdirAll(releasesDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create releases directory: %w", err)
		}

		for _, target := range targets {
			goos, goarch, suffix := target[0], target[1], target[2]
			outputName := filepath.Join(releasesDir, fmt.Sprintf("marcli-%s", suffix))

			var out, errBuf bytes.Buffer
			buildCmd := exec.CommandContext(ctx, "go", "build", "-ldflags", ldflags, "-o", outputName)
			buildCmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", goos), fmt.Sprintf("GOARCH=%s", goarch))
			buildCmd.Stdout = &out
			buildCmd.Stderr = &errBuf
			err := buildCmd.Run()

			if err != nil {
				errorMsg := fmt.Sprintf("%s/%s: FAILED - %s", goos, goarch, strings.TrimSpace(errBuf.String()))
				allErrors = append(allErrors, errorMsg)
				results = append(results, errorMsg)
			} else {
				results = append(results, fmt.Sprintf("%s/%s: OK -> %s", goos, goarch, outputName))
			}
		}
	}

	// Build for current platform (no cross-compilation) - building locally! 💖
	var finalName string
	if runtime.GOOS == "windows" {
		finalName = "marcli.exe"
	} else {
		finalName = "marcli"
	}

	var out, errBuf bytes.Buffer
	buildCmd := exec.CommandContext(ctx, "go", "build", "-ldflags", ldflags, "-o", finalName)
	buildCmd.Stdout = &out
	buildCmd.Stderr = &errBuf
	err = buildCmd.Run()

	if err != nil {
		errorMsg := fmt.Sprintf("current platform (%s/%s): FAILED - %s", runtime.GOOS, runtime.GOARCH, strings.TrimSpace(errBuf.String()))
		allErrors = append(allErrors, errorMsg)
		results = append(results, errorMsg)
	} else {
		results = append(results, fmt.Sprintf("current platform (%s/%s): OK -> %s", runtime.GOOS, runtime.GOARCH, finalName))

		// Install to user's PATH
		installPath, err := getInstallPath()
		if err != nil {
			results = append(results, fmt.Sprintf("Warning: Failed to determine install path: %v", err))
		} else {
			installDir := filepath.Dir(installPath)
			if err := os.MkdirAll(installDir, 0755); err != nil {
				results = append(results, fmt.Sprintf("Warning: Failed to create install directory %s: %v", installDir, err))
			} else {
				// Copy binary to install location
				src, err := os.Open(finalName)
				if err != nil {
					results = append(results, fmt.Sprintf("Warning: Failed to open %s: %v", finalName, err))
				} else {
					defer src.Close()
					dst, err := os.Create(installPath)
					if err != nil {
						results = append(results, fmt.Sprintf("Warning: Failed to create %s: %v", installPath, err))
					} else {
						defer dst.Close()
						_, err = io.Copy(dst, src)
						if err != nil {
							results = append(results, fmt.Sprintf("Warning: Failed to copy to %s: %v", installPath, err))
						} else {
							// Make executable on Unix-like systems
							if runtime.GOOS != "windows" {
								os.Chmod(installPath, 0755)
							}
							results = append(results, fmt.Sprintf("Installed -> %s", installPath))

							// Check and add to PATH if needed
							if err := ensureInPath(installDir); err != nil {
								results = append(results, fmt.Sprintf("Note: %s may not be in PATH. Add it manually or restart your terminal.", installDir))
							} else {
								if runtime.GOOS == "windows" {
									results = append(results, fmt.Sprintf("Added %s to PATH", installDir))
									results = append(results, "Note: Restart terminal or run: $env:Path = [System.Environment]::GetEnvironmentVariable(\"Path\",\"Machine\") + \";\" + [System.Environment]::GetEnvironmentVariable(\"Path\",\"User\")")
								} else {
									results = append(results, fmt.Sprintf("Added %s to PATH (restart terminal or run: source ~/.bashrc)", installDir))
								}
							}
						}
					}
				}
			}
		}
	}

	output := strings.Join(results, "\n")
	if len(allErrors) > 0 {
		return output, fmt.Errorf("some builds failed:\n%s", strings.Join(allErrors, "\n"))
	}
	return output, nil
}

// getInstallPath returns the path where the binary should be installed.
func getInstallPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	var installPath string
	switch runtime.GOOS {
	case "windows":
		// Windows: %USERPROFILE%\bin\marcli.exe
		installPath = filepath.Join(homeDir, "bin", "marcli.exe")
	default:
		// Linux/macOS: ~/.local/bin/marcli (XDG standard)
		installPath = filepath.Join(homeDir, ".local", "bin", "marcli")
	}

	return installPath, nil
}

// ensureInPath checks if the directory is in PATH and adds it if not.
func ensureInPath(dir string) error {
	pathEnv := os.Getenv("PATH")
	pathDirs := strings.Split(pathEnv, string(os.PathListSeparator))

	// Normalize the directory path for comparison
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if already in PATH
	for _, pathDir := range pathDirs {
		pathDirAbs, err := filepath.Abs(pathDir)
		if err == nil && strings.EqualFold(pathDirAbs, dirAbs) {
			return nil // Already in PATH
		}
	}

	// Not in PATH, try to add it
	if runtime.GOOS == "windows" {
		return addToWindowsPath(dirAbs)
	}
	return addToUnixPath(dirAbs)
}

// addToWindowsPath adds a directory to the Windows user PATH via registry.
func addToWindowsPath(dir string) error {
	// Use PowerShell to modify the user PATH in the registry
	// This is more reliable than setx and properly handles existing PATH
	psScript := fmt.Sprintf(`
		$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
		$dir = "%s"
		if ($userPath -notlike "*$dir*") {
			$newPath = $userPath + ";$dir"
			[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
			# Refresh PATH in current session
			$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
		}
	`, dir)

	psCmd := exec.Command("powershell", "-Command", psScript)
	if err := psCmd.Run(); err != nil {
		// Fallback to setx if PowerShell fails
		pathEnv := os.Getenv("PATH")
		newPath := pathEnv + string(os.PathListSeparator) + dir
		setxCmd := exec.Command("setx", "PATH", newPath)
		if err := setxCmd.Run(); err != nil {
			return fmt.Errorf("failed to add to PATH (restart terminal or add %s manually): %w", dir, err)
		}
		// Refresh current session PATH
		os.Setenv("PATH", newPath)
	} else {
		// Refresh current session PATH after PowerShell update
		userPath := os.Getenv("PATH")
		if !strings.Contains(strings.ToLower(userPath), strings.ToLower(dir)) {
			os.Setenv("PATH", userPath+string(os.PathListSeparator)+dir)
		}
	}
	return nil
}

// addToUnixPath adds a directory to Unix PATH (typically in shell profile).
func addToUnixPath(dir string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Try common shell profiles
	profiles := []string{".bashrc", ".zshrc", ".profile", ".bash_profile"}
	for _, profile := range profiles {
		profilePath := filepath.Join(homeDir, profile)
		if _, err := os.Stat(profilePath); err == nil {
			// Check if already added
			content, err := os.ReadFile(profilePath)
			if err != nil {
				continue
			}
			if strings.Contains(string(content), dir) {
				return nil // Already added
			}

			// Add export PATH line
			exportLine := fmt.Sprintf("\nexport PATH=\"$PATH:%s\"\n", dir)
			if err := os.WriteFile(profilePath, append(content, []byte(exportLine)...), 0644); err != nil {
				continue
			}
			return nil
		}
	}

	return fmt.Errorf("could not find shell profile to modify")
}
