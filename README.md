# marcli 💕

Hey there! ✨ Welcome to **marcli** - the cutest CLI tool for Marcy.cloud! 🎀

![mega-combine demo](./assets/demo.gif)

This is a super cute Terminal UI (a cutiepie TUI if you will *wink*) app that represents my personal CLI knowledge, built from the ground up with lots of love! 💖 It's got everything you need - a fancy TUI for when you want to explore, and direct CLI commands for when you're being that hacker girl from the 80's on the terminal! 💪

> **Note** 💅: This style of working (naming a CLI after myself as an internal marketing tool) was originally developed by me, Marcy, during my time at RealEyes Media while working for NBC Sports! ARE YOU IMPRESSED? While none of the actual code from the RealEyes Media MarCLI is present in this repo (we've come a long way, baby! ✨), the philosophy and approach live on here with lots of love! 💖  
>  
> And of course, **special thank you to Josh Sprow** for *all* the coding help, all the time — you make everything better! 💖

## What's Inside 🎀

- **Cutiepie TUI** 🎨 - An interactive menu with a cute purple border that's just adorable!
- **Web Terminal** 🌐 - Access the TUI remotely via browser with `cutiepie-tty`!
- **CLI Commands** 💅 - Run commands directly from the terminal
- **Version Tracking** ✨ - We keep track of builds because we're organized like that!
- **Cross-Platform Builds** 🌈 - Works everywhere because we're inclusive!

## Commands 💅

Here are all the cute commands available:

- `cutiepie` / (no args) - Launch the interactive TUI menu - so cute! 🎀
  - `--stay-alive` - Keep TUI open after running commands (returns to menu)
- `cutiepie-tty` 🌐 - Serve a web-based terminal interface for remote access
  - `--port <port>` - Specify port (default: 8080)
- `go-echo` - Echo using pure Go (no external processes) - so clean! 💕
- `ps-echo` - Echo using PowerShell - so powerful! 💪
- `bash-echo` - Echo using bash/sh - classic and cute! 🎀
- `build` - Build for all platforms (macOS, Linux, Windows) and install to PATH - building everything with love! 💖
  - `--fast` - Skip updating static JavaScript files and cross platform Go binaries for faster builds
- `version` - Show version and build number - so organized! ✨
- `-v` / `--version` - Quick version check (aliases for `version`) - we're so flexible! 💅
- `mega-combine` - Select and combine video files into ProRes for DaVinci Resolve on iPad - so efficient! 🎨 See [cmd/mega-combine-README.md](cmd/mega-combine-README.md) for details! 💕

## Quick Start 💖

Just run `marcli` with no args to see the cutie pie TUI, or use commands directly:

```bash
marcli                    # Launch the interactive TUI menu! 🎀
marcli --stay-alive       # Launch TUI that stays open after commands
marcli cutiepie-tty       # Start web terminal server on port 8080 🌐
marcli cutiepie-tty --port 3000  # Start on custom port
marcli mega-combine       # Combine videos for DaVinci Resolve! 🎨
marcli version            # See the version (so fancy!)
marcli build              # Build everything! 💪
marcli build --fast       # Fast build (skip JS updates)
marcli go-echo            # Try a command! 🎀
```

### Web Terminal 🌐

Access your TUI remotely via browser:

```bash
marcli cutiepie-tty --port 8080
# Then open http://localhost:8080 in your browser!
```

The web terminal uses HTMx, Alpine.js, and xterm.js for a full terminal experience in your browser. Perfect for remote access! ✨

Enjoy! 💕
