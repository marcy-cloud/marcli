package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"marcli/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// commandItem represents a command in the menu
type commandItem struct {
	name     string // canonical CLI name
	title    string
	desc     string
	run      func(context.Context) (string, error)
	selected bool
}

func (i commandItem) FilterValue() string {
	return i.title
}

func (i commandItem) IsSelected() bool {
	return i.selected
}

func (i *commandItem) SetSelected(selected bool) {
	i.selected = selected
}

func (i commandItem) DisplayText() string {
	return fmt.Sprintf("%s - %s", i.title, i.desc)
}

// tuiModel manages the main TUI menu
type tuiModel struct {
	listModel       *ui.Model
	items           []*commandItem
	selectedCommand *commandItem
	quitting        bool
	cancelled       bool // True if user pressed Ctrl+C
}

func initialTuiModel() tuiModel {
	osFlavor := "Linux"
	if runtime.GOOS == "windows" {
		osFlavor = "Windows"
	}

	// Create command items
	commandItems := []*commandItem{
		{
			name:  "go-echo",
			title: "Golang echo",
			desc:  `Echo "Golang echo" using native Go code`,
			run:   RunGoEcho,
		},
		{
			name:  "ps-echo",
			title: "PowerShell echo",
			desc:  `Echo "Powershell echo" by launching PowerShell`,
			run:   RunPSEcho,
		},
		{
			name:  "bash-echo",
			title: "Bash echo",
			desc:  `Echo "Bash echo" via bash (or sh)`,
			run:   RunBashEcho,
		},
		{
			name:  "build",
			title: "Build",
			desc:  `Run go build`,
			run:   RunBuild,
		},
		{
			name:  "version",
			title: "Version",
			desc:  `Show version and build number`,
			run:   RunVersion,
		},
		{
			name:  "mega-combine",
			title: "Mega Combine",
			desc:  `Select and combine video files from current directory`,
			run:   RunMegaCombine,
		},
	}

	// Convert to SelectableItem interface
	selectableItems := make([]ui.SelectableItem, len(commandItems))
	for i := range commandItems {
		selectableItems[i] = commandItems[i]
	}

	// Create selectable list using the UI component
	listModel := ui.New(ui.Config{
		Title:    fmt.Sprintf("marcli - Command Launcher [%s]", osFlavor),
		Items:    selectableItems,
		Width:    80,
		Height:   ui.DefaultListHeight,
		HelpText: "Space/Enter: run command, Ctrl+C: quit",
	})

	return tuiModel{
		listModel: listModel,
		items:     commandItems,
	}
}

func (m *tuiModel) Init() tea.Cmd {
	return m.listModel.Init()
}

func (m *tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle Ctrl+C before it reaches the list model
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "ctrl+c" {
		m.cancelled = true
		m.quitting = true
		return m, tea.Quit
	}

	// Handle spacebar before it reaches the list model - treat it like Enter
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == " " {
		// Get the currently highlighted command and select it
		if currentItem := m.listModel.GetCurrentItem(); currentItem != nil {
			if cmdItem, ok := currentItem.(*commandItem); ok {
				m.selectedCommand = cmdItem
				m.quitting = true
				return m, tea.Quit
			}
		}
		// If we can't get the item, fall through to normal handling
	}

	// Update the list model
	updatedModel, cmd := m.listModel.Update(msg)
	m.listModel = updatedModel.(*ui.Model)

	// Check if list was cancelled
	if m.listModel.IsCancelled() {
		m.cancelled = true
		m.quitting = true
		return m, tea.Quit
	}

	// If user confirmed (Enter), get the selected command
	if m.listModel.IsQuitting() {
		selectedItems := m.listModel.GetSelectedItems()
		if len(selectedItems) > 0 {
			// Use explicitly selected item
			if cmdItem, ok := selectedItems[0].(*commandItem); ok {
				m.selectedCommand = cmdItem
			}
		} else {
			// Use currently highlighted item if nothing was explicitly selected
			if currentItem := m.listModel.GetCurrentItem(); currentItem != nil {
				if cmdItem, ok := currentItem.(*commandItem); ok {
					m.selectedCommand = cmdItem
				}
			}
		}
		m.quitting = true
		return m, tea.Quit
	}

	return m, cmd
}

func (m *tuiModel) View() string {
	return m.listModel.View()
}

// GetSelectedCommand returns the selected command, if any
func (m *tuiModel) GetSelectedCommand() *commandItem {
	return m.selectedCommand
}

// waitForKeypress waits for a keypress without using alt screen
func waitForKeypress() error {
	// Small delay to ensure TUI has fully exited and terminal is ready
	time.Sleep(50 * time.Millisecond)

	fmt.Print("\n\nPress Enter (maybe twice honey) to continue...")

	// Create a fresh reader to avoid any buffered input issues
	reader := bufio.NewReader(os.Stdin)

	// Try to peek and discard any immediate newline that might be leftover
	if reader.Buffered() > 0 {
		peeked, _ := reader.Peek(1)
		if len(peeked) > 0 && peeked[0] == '\n' {
			reader.ReadByte() // Discard the leftover newline
		}
	}

	// Read until newline (Enter key on Windows)
	// This will wait for a fresh Enter press
	_, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Clear the continue message and add a newline
	fmt.Print("\r" + "                                         " + "\r\n")
	return nil
}

// RunCutiepieTUI starts the interactive cutiepie TUI - so cute and interactive! 🎀
// stayAliveOverride can be used to override the config setting (nil means use config)
func RunCutiepieTUI(stayAliveOverride *bool) error {
	// Load config to check StayAlive setting
	config, err := LoadConfig()
	if err != nil {
		// If config doesn't exist or can't be loaded, default to false (exit after command)
		config = &Config{
			StayAlive: false,
		}
	}

	// Default to false if not set (exit after command)
	stayAlive := false
	if config != nil {
		stayAlive = config.StayAlive
	}

	// Flag override wins if set
	if stayAliveOverride != nil {
		stayAlive = *stayAliveOverride
	}

	// Loop if StayAlive is true
	for {
		model := initialTuiModel()
		p := tea.NewProgram(&model, tea.WithAltScreen())
		finalModel, err := p.Run()
		if err != nil {
			return err
		}

		// Run the selected command only if not cancelled
		if tuiModel, ok := finalModel.(*tuiModel); ok {
			// Don't run command if user pressed Ctrl+C
			if tuiModel.cancelled {
				return nil
			}
			cmd := tuiModel.GetSelectedCommand()
			if cmd != nil {
				ctx := context.Background()
				out, err := cmd.run(ctx)
				if err != nil {
					return err
				}
				if out != "" {
					fmt.Print(out)
				}

				// If StayAlive is true, wait for keypress and loop
				if stayAlive {
					if err := waitForKeypress(); err != nil {
						return err
					}
					// Continue loop to show menu again
					continue
				}
			}
		}

		// If StayAlive is false, exit after running command
		if !stayAlive {
			return nil
		}
	}
}

// RunCutiepieTUICommand is a wrapper that matches the command registry signature - so organized! ✨
func RunCutiepieTUICommand(ctx context.Context) (string, error) {
	// Check for --stay-alive flag override
	var stayAliveOverride *bool
	if ctx.Value("stayAlive") != nil {
		if val, ok := ctx.Value("stayAlive").(bool); ok {
			stayAliveOverride = &val
		}
	}

	err := RunCutiepieTUI(stayAliveOverride)
	if err != nil {
		return "", err
	}
	return "", nil
}
