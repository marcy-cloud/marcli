package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/coder/websocket"
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
	http.HandleFunc("/ws", handleWebSocket)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Accept the WebSocket connection
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("Failed to accept WebSocket connection: %v", err)
		return
	}
	defer conn.CloseNow()

	log.Printf("WebSocket connection established")

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
	log.Printf("PTY started successfully")

	// Set initial terminal size - use a reasonable default that matches common browser sizes
	// We'll resize it properly once the client sends the actual terminal size
	// For now, use a larger default to prevent overflow
	if err := currentPTY.Resize(120, 30); err != nil {
		log.Printf("Failed to resize PTY: %v", err)
	}
	ptyMutex.Unlock()

	// Clean up on exit
	defer func() {
		log.Printf("WebSocket connection closing, cleaning up PTY")
		ptyMutex.Lock()
		if currentPTY != nil {
			currentPTY.Close()
			currentPTY = nil
		}
		ptyMutex.Unlock()
	}()

	// Channel to signal when copying is done
	done := make(chan bool, 2)
	ctx := context.Background()

	// Copy from PTY to WebSocket (terminal output -> browser)
	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := currentPTY.Read(buf)
			if n > 0 {
				if writeErr := conn.Write(ctx, websocket.MessageText, buf[:n]); writeErr != nil {
					log.Printf("Error writing to WebSocket: %v", writeErr)
					done <- true
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from PTY: %v", err)
				} else {
					log.Printf("PTY closed (EOF)")
				}
				done <- true
				return
			}
		}
	}()

	// Copy from WebSocket to PTY (browser input -> terminal)
	go func() {
		for {
			typ, buf, err := conn.Read(ctx)
			if len(buf) > 0 {
				if _, writeErr := currentPTY.Write(buf); writeErr != nil {
					log.Printf("Error writing to PTY: %v", writeErr)
					done <- true
					return
				}
			}
			if err != nil {
				// Check if it's a close error
				if err == io.EOF || websocket.CloseStatus(err) != -1 {
					log.Printf("WebSocket closed")
				} else {
					log.Printf("Error reading from WebSocket: %v", err)
				}
				done <- true
				return
			}
			_ = typ // Message type (binary/text), we handle both
		}
	}()

	// Wait for either copy to finish
	<-done
	log.Printf("One of the copy operations finished, closing WebSocket")
}
