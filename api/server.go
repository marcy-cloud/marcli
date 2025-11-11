package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	currentPTY *PTYManager
	ptyMutex   sync.Mutex
)

// StartServer starts the HTTP server on the specified port
func StartServer(port int) error {
	// Create static file server
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	})

	// WebSocket endpoint for terminal I/O
	http.Handle("/ws", websocket.Handler(handleWebSocket))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func handleWebSocket(ws *websocket.Conn) {
	defer ws.Close()

	// Close existing PTY and create new one
	ptyMutex.Lock()
	if currentPTY != nil {
		currentPTY.Close()
	}
	currentPTY = NewPTYManager()
	if err := currentPTY.Start(); err != nil {
		ptyMutex.Unlock()
		log.Printf("Failed to start PTY: %v", err)
		return
	}

	// Set initial terminal size (80x24 is a common default)
	if err := currentPTY.Resize(80, 24); err != nil {
		log.Printf("Failed to resize PTY: %v", err)
	}
	ptyMutex.Unlock()

	// Clean up on exit
	defer func() {
		ptyMutex.Lock()
		if currentPTY != nil {
			currentPTY.Close()
			currentPTY = nil
		}
		ptyMutex.Unlock()
	}()

	// Channel to signal when copying is done
	done := make(chan bool)

	// Copy from PTY to WebSocket (terminal output -> browser)
	go func() {
		io.Copy(ws, currentPTY)
		done <- true
	}()

	// Copy from WebSocket to PTY (browser input -> terminal)
	go func() {
		io.Copy(currentPTY, ws)
		done <- true
	}()

	// Wait for either copy to finish
	<-done
}
