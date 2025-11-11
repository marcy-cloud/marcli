package api

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

// PTYManager manages a single PTY instance
type PTYManager struct {
	ptmx   *os.File
	cmd    *exec.Cmd
	mu     sync.Mutex
	closed bool
}

// NewPTYManager creates a new PTY manager
func NewPTYManager() *PTYManager {
	return &PTYManager{}
}

// Start starts a new PTY with marcli running cutiepie-tui
func (p *PTYManager) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Close existing PTY if any
	if p.ptmx != nil {
		p.closeLocked()
	}

	// Find the marcli executable
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Create command to run marcli (which will show cutiepie-tui by default)
	// We need to set ExitAfterCommand=false, but since we're spawning the process,
	// we'll need to modify config or pass an env var. For now, let's just run it.
	// The TUI will use the config file's setting.
	p.cmd = exec.Command(execPath)

	// Start the command with a PTY
	ptmx, err := pty.Start(p.cmd)
	if err != nil {
		return err
	}

	p.ptmx = ptmx
	p.closed = false

	return nil
}

// Write writes data to the PTY stdin
func (p *PTYManager) Write(data []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed || p.ptmx == nil {
		return 0, io.EOF
	}

	return p.ptmx.Write(data)
}

// Read reads data from the PTY stdout
func (p *PTYManager) Read(data []byte) (int, error) {
	p.mu.Lock()
	ptmx := p.ptmx
	closed := p.closed
	p.mu.Unlock()

	if closed || ptmx == nil {
		return 0, io.EOF
	}

	return ptmx.Read(data)
}

// Resize resizes the PTY
func (p *PTYManager) Resize(cols, rows uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed || p.ptmx == nil {
		return io.EOF
	}

	size := &pty.Winsize{
		Rows: rows,
		Cols: cols,
	}

	return pty.Setsize(p.ptmx, size)
}

// Close closes the PTY and kills the command
func (p *PTYManager) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closeLocked()
}

func (p *PTYManager) closeLocked() error {
	if p.closed {
		return nil
	}

	p.closed = true

	if p.ptmx != nil {
		p.ptmx.Close()
		p.ptmx = nil
	}

	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Wait()
		p.cmd = nil
	}

	return nil
}

// IsClosed returns whether the PTY is closed
func (p *PTYManager) IsClosed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closed
}
