package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"marcli/ui"

	tea "github.com/charmbracelet/bubbletea"
	logger "github.com/charmbracelet/log"
)

// videoFileItem represents a video file in the list
type videoFileItem struct {
	title    string
	filePath string
	modTime  time.Time
	selected bool
}

func (i videoFileItem) FilterValue() string {
	return i.title
}

func (i videoFileItem) IsSelected() bool {
	return i.selected
}

func (i *videoFileItem) SetSelected(selected bool) {
	i.selected = selected
}

func (i videoFileItem) DisplayText() string {
	dateStr := i.modTime.Format("2006-01-02 15:04")
	return fmt.Sprintf("%s  %s", i.title, dateStr)
}

// megaCombineModel manages the state of the mega-combine TUI
type megaCombineModel struct {
	listModel     *ui.Model
	items         []*videoFileItem
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

	// Convert to pointers and SelectableItem interface
	itemPtrs := make([]*videoFileItem, len(items))
	selectableItems := make([]ui.SelectableItem, len(items))
	for i := range items {
		itemPtrs[i] = &items[i]
		selectableItems[i] = itemPtrs[i]
	}

	// Create selectable list using the UI component
	listModel := ui.New(ui.Config{
		Title:    "Select Video Files",
		Items:    selectableItems,
		Width:    80,
		Height:   ui.DefaultListHeight,
		HelpText: "Space: toggle, Enter: confirm, Ctrl+C: quit",
	})

	return megaCombineModel{
		listModel: listModel,
		items:     itemPtrs,
	}, nil
}

func (m *megaCombineModel) Init() tea.Cmd {
	return m.listModel.Init()
}

func (m *megaCombineModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update the list model
	updatedModel, cmd := m.listModel.Update(msg)
	m.listModel = updatedModel.(*ui.Model)

	// If user confirmed (Enter), log selected files
	if m.listModel.IsQuitting() {
		m.logSelectedFiles()
	}

	return m, cmd
}

func (m *megaCombineModel) View() string {
	return m.listModel.View()
}

func (m *megaCombineModel) logSelectedFiles() {
	selectedItems := m.listModel.GetSelectedItems()
	if len(selectedItems) == 0 {
		logger.Info("No files selected")
		m.selectedFiles = []string{}
		return
	}

	// Collect selected files
	m.selectedFiles = make([]string, 0, len(selectedItems))
	logger.Info("Selected video files:")
	for _, item := range selectedItems {
		if videoItem, ok := item.(*videoFileItem); ok {
			m.selectedFiles = append(m.selectedFiles, videoItem.filePath)
			logger.Info("  " + videoItem.filePath)
		}
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
