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

	// Set initial terminal size (80x24 is a common default)
	if err := currentPTY.Resize(80, 24); err != nil {
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

	// Copy from PTY to WebSocket (terminal output -> browser)
	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, err := currentPTY.Read(buf)
			if n > 0 {
				if _, writeErr := ws.Write(buf[:n]); writeErr != nil {
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
		buf := make([]byte, 32*1024)
		for {
			n, err := ws.Read(buf)
			if n > 0 {
				if _, writeErr := currentPTY.Write(buf[:n]); writeErr != nil {
					log.Printf("Error writing to PTY: %v", writeErr)
					done <- true
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from WebSocket: %v", err)
				} else {
					log.Printf("WebSocket closed (EOF)")
				}
				done <- true
				return
			}
		}
	}()

	// Wait for either copy to finish
	<-done
	log.Printf("One of the copy operations finished, closing WebSocket")
}
