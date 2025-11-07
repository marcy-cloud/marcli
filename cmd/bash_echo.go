package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

// RunBashEcho runs a bash echo command.
func RunBashEcho(ctx context.Context) (string, error) {
	// Try bash; fallback to sh if present (Linux/macOS). On Windows, suggest Git Bash/WSL.
	if _, err := exec.LookPath("bash"); err == nil {
		return runShell(ctx, "bash", []string{"-lc", "echo 'Bash echo'"})
	}
	if _, err := exec.LookPath("sh"); err == nil {
		return runShell(ctx, "sh", []string{"-lc", "echo 'Bash echo'"})
	}
	if runtime.GOOS == "windows" {
		return "", fmt.Errorf("bash not found. Install Git Bash (Git for Windows) or enable WSL.")
	}
	return "", fmt.Errorf("neither bash nor sh found in PATH")
}

func runShell(ctx context.Context, bin string, args []string) (string, error) {
	var out, errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if errBuf.Len() > 0 {
		out.WriteString("\n" + errBuf.String())
	}
	return out.String(), err
}
