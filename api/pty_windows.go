//go:build windows
// +build windows

package api

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/UserExistsError/conpty"
)

// PTYManager manages a single PTY instance on Windows using ConPTY
type PTYManager struct {
	cpty   *conpty.ConPty
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
	if p.cpty != nil {
		p.closeLocked()
	}

	// Find the marcli executable
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Create command line with --stay-alive flag
	commandLine := fmt.Sprintf(`"%s" --stay-alive`, execPath)

	// Start the command with ConPTY
	cpty, err := conpty.Start(commandLine)
	if err != nil {
		return fmt.Errorf("unsupported: %w", err)
	}

	p.cpty = cpty
	p.closed = false

	return nil
}

// Write writes data to the PTY stdin
func (p *PTYManager) Write(data []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed || p.cpty == nil {
		return 0, io.EOF
	}

	return p.cpty.Write(data)
}

// Read reads data from the PTY stdout
func (p *PTYManager) Read(data []byte) (int, error) {
	p.mu.Lock()
	cpty := p.cpty
	closed := p.closed
	p.mu.Unlock()

	if closed || cpty == nil {
		return 0, io.EOF
	}

	return cpty.Read(data)
}

// Resize resizes the PTY
func (p *PTYManager) Resize(cols, rows uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed || p.cpty == nil {
		return io.EOF
	}

	return p.cpty.Resize(int(cols), int(rows))
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

	if p.cpty != nil {
		p.cpty.Close()
		p.cpty = nil
	}

	return nil
}

// IsClosed returns whether the PTY is closed
func (p *PTYManager) IsClosed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closed
}

