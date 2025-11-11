# Commands 💕

Hey there! ✨ This directory contains all the cute commands for **marcli**! 🎀

We keep this list updated as we add new commands - so organized! 💅

## Command List (Newest First) 🎀

### cutiepie-tty 🌐
**File:** `cutiepie-tty.go`  
**Description:** Serves a web-based terminal interface for remote access to cutiepie-tui - so accessible! 🌐  
**Usage:** `marcli cutiepie-tty [--port 8080]`  
**Details:** Starts an HTTP server that serves a web terminal using HTMx, Alpine.js, and xterm.js. The terminal connects to a PTY running cutiepie-tui with `--stay-alive` enabled, allowing remote access via browser. Uses WebSocket for real-time bidirectional communication. Default port is 8080. The web interface features a beautiful terminal emulator with proper overflow handling and responsive design.

### cutiepie 🎀
**File:** `cutiepie-tui.go`  
**Description:** The main interactive TUI menu with a cute purple border - so adorable! 💜  
**Usage:** `marcli` or `marcli cutiepie [--stay-alive]`  
**Details:** Launches the interactive terminal UI with a beautiful purple rounded border. Navigate with arrow keys, select with Enter/Space, quit with Ctrl+C or 'q'. The `--stay-alive` flag keeps the TUI open after running commands, returning to the menu instead of exiting. Can also be configured via `stayAlive` in `config.yml`.

### version ✨
**File:** `version.go`  
**Description:** Shows the current version and build number - so organized! 💖  
**Usage:** `marcli version` or `marcli -v` or `marcli --version`

### build 💪
**File:** `build.go`  
**Description:** Builds for all platforms (macOS, Linux, Windows) and installs to PATH - building everything with love! 💖  
**Usage:** `marcli build`

### bash-echo 🎀
**File:** `bash-echo.go`  
**Description:** Echo using bash/sh - classic and cute! 🎀  
**Usage:** `marcli bash-echo`

### ps-echo 💪
**File:** `ps-echo.go`  
**Description:** Echo using PowerShell - so powerful! 💪  
**Usage:** `marcli ps-echo`

### go-echo 💕
**File:** `go-echo.go`  
**Description:** Echo using pure Go (no external processes) - so clean! 💕  
**Usage:** `marcli go-echo`

---

**Remember** ✨: When adding a new command, update this README with the newest command at the top! We're so organized! 💖

