package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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
	// Check if test mode is enabled
	testMode := ctx.Value("megaCombineTestMode") == true

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

		// In test mode, show the ffmpeg command that would be run
		if testMode {
			cmd, err := generateFFmpegCommand(m.selectedFiles)
			if err != nil {
				return "", err
			}
			return cmd, nil
		}

		// Main mode - actually run the ffmpeg command
		return runFFmpegCommand(m.selectedFiles)
	}

	return "Video file selection completed. Check logs for selected files.", nil
}

// generateFFmpegCommand creates an ffmpeg command using the selected files without a filelist
func generateFFmpegCommand(selectedFiles []string) (string, error) {
	if len(selectedFiles) == 0 {
		return "", fmt.Errorf("no files selected")
	}

	var cmd strings.Builder

	// Add all input files with -i flag
	for _, file := range selectedFiles {
		// Get absolute path for each file to ensure ffmpeg can find them
		absFilePath, err := filepath.Abs(file)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path for %s: %w", file, err)
		}
		cmd.WriteString(fmt.Sprintf(" -i \"%s\"", absFilePath))
	}

	// Build concat filter complex with timestamp normalization (robust version)
	// Normalizes timestamps to handle VFR/mismatched starts safely
	// Format: [0:v]setpts=PTS-STARTPTS[v0];[0:a]asetpts=PTS-STARTPTS[a0];...concat=n=N:v=1:a=1[outv][outa]
	numFiles := len(selectedFiles)
	var filterComplex strings.Builder

	// Add timestamp normalization for each input
	for i := 0; i < numFiles; i++ {
		if i > 0 {
			filterComplex.WriteString(";")
		}
		filterComplex.WriteString(fmt.Sprintf("[%d:v]setpts=PTS-STARTPTS[v%d];[%d:a]asetpts=PTS-STARTPTS[a%d]", i, i, i, i))
	}

	// Add concat with normalized streams
	filterComplex.WriteString(";")
	for i := 0; i < numFiles; i++ {
		filterComplex.WriteString(fmt.Sprintf("[v%d][a%d]", i, i))
	}
	filterComplex.WriteString(fmt.Sprintf("concat=n=%d:v=1:a=1[outv][outa]", numFiles))

	// Generate the complete ffmpeg command
	cmd.WriteString(" \\\n")
	cmd.WriteString("  -filter_complex \"")
	cmd.WriteString(filterComplex.String())
	cmd.WriteString("\" \\\n")
	cmd.WriteString("  -map \"[outv]\" -map \"[outa]\" \\\n")
	cmd.WriteString("  -c:v prores_ks -profile:v 1 -pix_fmt yuv422p10le -threads 0 \\\n")
	cmd.WriteString("  -c:a pcm_s16le -ar 48000 -ac 2 \\\n")
	cmd.WriteString("  \"out.mov\"")

	// Prepend "ffmpeg" to the command
	fullCmd := "ffmpeg" + cmd.String()
	return fullCmd, nil
}

// runFFmpegCommand executes the ffmpeg command with the selected files
func runFFmpegCommand(selectedFiles []string) (string, error) {
	if len(selectedFiles) == 0 {
		return "", fmt.Errorf("no files selected")
	}

	// Build the ffmpeg command arguments
	var args []string

	// Add all input files with -i flag
	for _, file := range selectedFiles {
		// Get absolute path for each file to ensure ffmpeg can find them
		absFilePath, err := filepath.Abs(file)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path for %s: %w", file, err)
		}
		args = append(args, "-i", absFilePath)
	}

	// Build concat filter complex with timestamp normalization (robust version)
	numFiles := len(selectedFiles)
	var filterComplex strings.Builder

	// Add timestamp normalization for each input
	for i := 0; i < numFiles; i++ {
		if i > 0 {
			filterComplex.WriteString(";")
		}
		filterComplex.WriteString(fmt.Sprintf("[%d:v]setpts=PTS-STARTPTS[v%d];[%d:a]asetpts=PTS-STARTPTS[a%d]", i, i, i, i))
	}

	// Add concat with normalized streams
	filterComplex.WriteString(";")
	for i := 0; i < numFiles; i++ {
		filterComplex.WriteString(fmt.Sprintf("[v%d][a%d]", i, i))
	}
	filterComplex.WriteString(fmt.Sprintf("concat=n=%d:v=1:a=1[outv][outa]", numFiles))

	// Add filter_complex and other arguments
	args = append(args, "-filter_complex", filterComplex.String())
	args = append(args, "-map", "[outv]")
	args = append(args, "-map", "[outa]")
	args = append(args, "-c:v", "prores_ks")
	args = append(args, "-profile:v", "1")
	args = append(args, "-pix_fmt", "yuv422p10le")
	args = append(args, "-threads", "0")
	args = append(args, "-c:a", "pcm_s16le")
	args = append(args, "-ar", "48000")
	args = append(args, "-ac", "2")
	args = append(args, "out.mov")

	// Execute ffmpeg command
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info("Running ffmpeg command to combine videos...")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg command failed: %w", err)
	}

	return "Video files successfully combined into out.mov", nil
}
