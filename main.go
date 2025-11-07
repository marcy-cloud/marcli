package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"marcli/cmd"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	logger "github.com/charmbracelet/log"
)

/* ---------- styling ---------- */

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")) // cyan-ish
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	borderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	headerStyle = lipgloss.NewStyle().Padding(0, 1)
	footerStyle = lipgloss.NewStyle().Faint(true)
	itemTitle   = lipgloss.NewStyle().Bold(true)
	itemDesc    = lipgloss.NewStyle().Faint(true)
	appName     = "Charm Bubble Tea — Command Launcher"
)

/* ---------- menu items ---------- */

type menuItem struct {
	name  string // canonical CLI name
	title string
	desc  string
	run   func(context.Context) (string, error)
}

func (i menuItem) Title() string       { return itemTitle.Render(i.title) }
func (i menuItem) Description() string { return itemDesc.Render(i.desc) }
func (i menuItem) FilterValue() string { return i.title }

// commandRegistry maps CLI names to command functions
var commandRegistry = make(map[string]func(context.Context) (string, error))

// initCommands populates the command registry
func initCommands() {
	items := []menuItem{
		{
			name:  "go-echo",
			title: "Golang echo",
			desc:  `Echo "Golang echo" using native Go code`,
			run:   cmd.RunGoEcho,
		},
		{
			name:  "ps-echo",
			title: "PowerShell echo",
			desc:  `Echo "Powershell echo" by launching PowerShell`,
			run:   cmd.RunPSEcho,
		},
		{
			name:  "bash-echo",
			title: "Bash echo",
			desc:  `Echo "Bash echo" via bash (or sh)`,
			run:   cmd.RunBashEcho,
		},
		{
			name:  "build",
			title: "Build",
			desc:  `Run go build`,
			run:   cmd.RunBuild,
		},
	}

	for _, item := range items {
		commandRegistry[item.name] = item.run
	}
}

/* ---------- model ---------- */

type model struct {
	osFlavor string

	menu    list.Model
	view    viewport.Model
	spin    spinner.Model
	running bool
	cancel  context.CancelFunc
}

type outputMsg struct {
	out string
	err error
}

func initialModel() model {
	osFlavor := "Linux"
	if runtime.GOOS == "windows" {
		osFlavor = "Windows"
	}

	items := []list.Item{
		menuItem{
			name:  "go-echo",
			title: "Golang echo",
			desc:  `Echo "Golang echo" using native Go code`,
			run:   cmd.RunGoEcho,
		},
		menuItem{
			name:  "ps-echo",
			title: "PowerShell echo",
			desc:  `Echo "Powershell echo" by launching PowerShell`,
			run:   cmd.RunPSEcho,
		},
		menuItem{
			name:  "bash-echo",
			title: "Bash echo",
			desc:  `Echo "Bash echo" via bash (or sh)`,
			run:   cmd.RunBashEcho,
		},
		menuItem{
			name:  "build",
			title: "Build",
			desc:  `Run go build`,
			run:   cmd.RunBuild,
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 36, 10)
	l.Title = "Pick a command and press Enter"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	vp := viewport.New(0, 0)
	vp.SetContent("Ready.\n")

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return model{
		osFlavor: osFlavor,
		menu:     l,
		view:     vp,
		spin:     sp,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spin.Tick, nil)
}

/* ---------- tea update/view ---------- */

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.running {
				return m, nil
			}
			if it, ok := m.menu.SelectedItem().(menuItem); ok {
				// start run
				if m.cancel != nil {
					m.cancel()
				}
				ctx, cancel := context.WithCancel(context.Background())
				m.cancel = cancel
				m.running = true
				m.view.SetContent(m.view.View() + headerStyle.Render("> "+it.title) + "\n")
				return m, tea.Batch(m.spin.Tick, func() tea.Msg {
					start := time.Now()
					out, err := it.run(ctx)
					elapsed := time.Since(start)
					if out == "" {
						out = "(no output)"
					}
					out = fmt.Sprintf("%s\n\n%s (%.2fs)\n",
						out, okStyle.Render("Done"), elapsed.Seconds())
					return outputMsg{out: out, err: err}
				})
			}
		case "ctrl+c":
			if m.running && m.cancel != nil {
				m.cancel()
				m.running = false
				m.view.SetContent(m.view.View() + errorStyle.Render("Cancelled.") + "\n")
				return m, nil
			}
			return m, tea.Quit
		}

	case outputMsg:
		m.running = false
		if msg.err != nil {
			m.view.SetContent(m.view.View() + errorStyle.Render(msg.out+"\n"+msg.err.Error()) + "\n")
		} else {
			m.view.SetContent(m.view.View() + msg.out + "\n")
		}
		return m, nil

	case tea.WindowSizeMsg:
		// layout
		m.menu.SetWidth(msg.Width / 2)
		m.view.Width = msg.Width - m.menu.Width() - 4
		m.view.Height = msg.Height - 5
		return m, nil

	case spinner.TickMsg:
		if m.running {
			var cmd tea.Cmd
			m.spin, cmd = m.spin.Update(msg)
			return m, cmd
		}
	}
	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m model) View() string {
	header := titleStyle.Render(appName) +
		"  " + subtleStyle.Render(fmt.Sprintf("[%s]", m.osFlavor)) + "\n"
	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		borderStyle.Width(m.menu.Width()).Height(m.view.Height+2).Render(m.menu.View()),
		borderStyle.Width(m.view.Width).Height(m.view.Height+2).Render(m.view.View()),
	)
	foot := footerStyle.Render("\nEnter: run • Ctrl+C: cancel/quit\n")
	return header + body + foot
}

func main() {
	// Initialize command registry
	initCommands()

	args := os.Args[1:]

	// CLI mode: if args provided, run command directly
	if len(args) > 0 {
		cmdName := args[0]
		cmd, exists := commandRegistry[cmdName]
		if !exists {
			logger.Fatal("unknown command", "command", cmdName)
		}

		ctx := context.Background()
		out, err := cmd(ctx)
		if err != nil {
			logger.Fatal("command failed", "err", err)
		}
		fmt.Print(out)
		return
	}

	// TUI mode: no args, show interactive menu
	if err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start(); err != nil {
		logger.Fatal("error", "err", err)
	}
}
