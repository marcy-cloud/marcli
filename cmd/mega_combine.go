package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	logger "github.com/charmbracelet/log"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// videoFileItem represents a video file in the list
type videoFileItem struct {
	title    string
	filePath string
	modTime  time.Time
	selected bool
}

func (i videoFileItem) Title() string {
	checkmark := " "
	if i.selected {
		checkmark = "✓"
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Render(checkmark + " " + i.title)
}

func (i videoFileItem) Description() string {
	return lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("241")).
		Render(i.modTime.Format("2006-01-02 15:04:05"))
}

func (i videoFileItem) FilterValue() string {
	return i.title
}

// megaCombineModel manages the state of the mega-combine TUI
type megaCombineModel struct {
	list         list.Model
	items        []videoFileItem
	selected     map[int]struct{}
	quitting     bool
	selectedFiles []string // Store selected file paths for return
}

func initialMegaCombineModel() (megaCombineModel, error) {
	// Get video files from current directory
	items, err := getVideoFiles(".")
	if err != nil {
		return megaCombineModel{}, err
	}

	if len(items) == 0 {
		return megaCombineModel{}, fmt.Errorf("no video files found in current directory")
	}

	// Convert to list items
	listItems := make([]list.Item, len(items))
	for i := range items {
		listItems[i] = items[i]
	}

	// Create list with simple delegate - use default size that will be updated by WindowSizeMsg
	l := list.New(listItems, list.NewDefaultDelegate(), 80, 20)
	l.Title = "Select Video Files (Space: toggle, Enter: confirm, Ctrl+C: quit)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(lipgloss.Color("212")).
		Bold(true)

	selected := make(map[int]struct{})

	return megaCombineModel{
		list:     l,
		items:    items,
		selected: selected,
	}, nil
}

func (m *megaCombineModel) Init() tea.Cmd {
	return nil
}

func (m *megaCombineModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle key messages before passing to list
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case " ":
			// Toggle selection - handle BEFORE list gets it
			idx := m.list.Index()
			if _, ok := m.selected[idx]; ok {
				delete(m.selected, idx)
				m.items[idx].selected = false
			} else {
				m.selected[idx] = struct{}{}
				m.items[idx].selected = true
			}
			// Update the list items to reflect selection changes
			listItems := make([]list.Item, len(m.items))
			for i := range m.items {
				listItems[i] = m.items[i]
			}
			m.list.SetItems(listItems)
			// Don't pass spacebar to list, we handled it
			return m, nil

		case "enter":
			// Log selected files and quit
			m.logSelectedFiles()
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Handle window size
	if winSizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.list.SetWidth(winSizeMsg.Width)
		m.list.SetHeight(winSizeMsg.Height - 2)
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	// Pass all other messages to the list
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *megaCombineModel) View() string {
	if m.quitting {
		return ""
	}
	return "\n" + m.list.View()
}

func (m *megaCombineModel) logSelectedFiles() {
	if len(m.selected) == 0 {
		logger.Info("No files selected")
		m.selectedFiles = []string{}
		return
	}

	// Sort selected indices for consistent output
	indices := make([]int, 0, len(m.selected))
	for idx := range m.selected {
		indices = append(indices, idx)
	}
	sort.Ints(indices)

	// Collect selected files
	m.selectedFiles = make([]string, 0, len(indices))
	logger.Info("Selected video files:")
	for _, idx := range indices {
		filePath := m.items[idx].filePath
		m.selectedFiles = append(m.selectedFiles, filePath)
		logger.Info("  " + filePath)
	}
}

// getVideoFiles scans the current directory for video files and sorts by modification time
func getVideoFiles(dir string) ([]videoFileItem, error) {
	videoExtensions := map[string]bool{
		".mp4":  true,
		".avi":  true,
		".mov":  true,
		".mkv":  true,
		".webm": true,
		".flv":  true,
		".wmv":  true,
		".m4v":  true,
		".mpg":  true,
		".mpeg": true,
		".3gp":  true,
		".ogv":  true,
	}

	var items []videoFileItem

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if !videoExtensions[ext] {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())
		items = append(items, videoFileItem{
			title:    entry.Name(),
			filePath: fullPath,
			modTime:  info.ModTime(),
			selected: false,
		})
	}

	// Sort by modification time (oldest first)
	sort.Slice(items, func(i, j int) bool {
		return items[i].modTime.Before(items[j].modTime)
	})

	return items, nil
}

// RunMegaCombine runs the mega-combine TUI command
func RunMegaCombine(ctx context.Context) (string, error) {
	model, err := initialMegaCombineModel()
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(&model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	// Get the final model and extract selected files
	if m, ok := finalModel.(*megaCombineModel); ok {
		if len(m.selectedFiles) == 0 {
			return "No files selected.", nil
		}
		
		var output strings.Builder
		output.WriteString("Selected video files:\n")
		for _, file := range m.selectedFiles {
			output.WriteString("  " + file + "\n")
		}
		return output.String(), nil
	}

	return "Video file selection completed. Check logs for selected files.", nil
}

