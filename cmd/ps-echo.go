package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

// RunPSEcho runs a PowerShell echo command - so powerful! 💪
func RunPSEcho(ctx context.Context) (string, error) {
	// Prefer PowerShell 7+ if available - we're so modern! ✨
	ps := "pwsh"
	if _, err := exec.LookPath(ps); err != nil {
		// Fallbacks - we're flexible like that! 💅
		if runtime.GOOS == "windows" {
			ps = "powershell.exe"
		} else {
			// Non-Windows without pwsh installed - we'll help them out! 💖
			return "", fmt.Errorf("PowerShell (pwsh) not found. Install from https://github.com/PowerShell/PowerShell")
		}
	}
	args := []string{"-NoLogo", "-NoProfile"}
	if runtime.GOOS == "windows" {
		args = append(args, "-ExecutionPolicy", "Bypass")
	}
	args = append(args, "-Command", "Write-Output 'Powershell echo'")

	var out, errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, ps, args...)
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if errBuf.Len() > 0 {
		out.WriteString("\n" + errBuf.String())
	}
	return out.String(), err
}

