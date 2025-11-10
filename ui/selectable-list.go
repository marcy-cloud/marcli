package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const DefaultListHeight = 20

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

// SelectableItem is the interface that items must implement to be used in a selectable list
type SelectableItem interface {
	list.Item
	IsSelected() bool
	SetSelected(bool)
	DisplayText() string // Returns the text to display for this item
}

// itemDelegate handles rendering of selectable list items
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(SelectableItem)
	if !ok {
		return
	}

	// Format: checkbox/checkmark + display text
	// ⬜ for unselected, ✅ for selected
	checkbox := "⬜"
	if item.IsSelected() {
		checkbox = "✅"
	}
	
	str := fmt.Sprintf("%s %s", checkbox, item.DisplayText())

	fn := itemStyle.Render
	if index == m.Index() {
		// Heart (💖) is the cursor indicator for highlighted items
		fn = func(s ...string) string {
			return selectedItemStyle.Render("💖 " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// Model represents a selectable list model
type Model struct {
	list      list.Model
	items     []SelectableItem
	selected  map[int]struct{}
	quitting  bool
	cancelled bool // True if user pressed Ctrl+C to quit
}

// Config holds configuration for creating a selectable list
type Config struct {
	Title      string
	Items      []SelectableItem
	Width      int
	Height     int
	HelpText   string
}

// New creates a new selectable list model
func New(cfg Config) *Model {
	// Convert items to list.Item interface
	listItems := make([]list.Item, len(cfg.Items))
	for i := range cfg.Items {
		listItems[i] = cfg.Items[i]
	}

	width := cfg.Width
	if width == 0 {
		width = 80
	}

	height := cfg.Height
	if height == 0 {
		height = DefaultListHeight
	}

	// Create list with custom delegate
	l := list.New(listItems, itemDelegate{}, width, height)
	l.Title = cfg.Title
	if cfg.HelpText != "" {
		l.Title = cfg.Title + " (" + cfg.HelpText + ")"
	}
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	selected := make(map[int]struct{})

	return &Model{
		list:     l,
		items:    cfg.Items,
		selected: selected,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle key messages before passing to list
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c":
			m.quitting = true
			m.cancelled = true // Mark as cancelled so commands don't run
			return m, tea.Quit

		case " ":
			// Toggle selection - handle BEFORE list gets it
			idx := m.list.Index()
			if _, ok := m.selected[idx]; ok {
				delete(m.selected, idx)
				m.items[idx].SetSelected(false)
			} else {
				m.selected[idx] = struct{}{}
				m.items[idx].SetSelected(true)
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
			// Confirm selection and quit
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Handle window size
	if winSizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.list.SetWidth(winSizeMsg.Width)
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	// Pass all other messages to the list
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.quitting {
		return ""
	}
	return "\n" + m.list.View()
}

// GetSelectedIndices returns the indices of all selected items
func (m *Model) GetSelectedIndices() []int {
	indices := make([]int, 0, len(m.selected))
	for idx := range m.selected {
		indices = append(indices, idx)
	}
	return indices
}

// GetSelectedItems returns the selected items
func (m *Model) GetSelectedItems() []SelectableItem {
	indices := m.GetSelectedIndices()
	selected := make([]SelectableItem, 0, len(indices))
	for _, idx := range indices {
		if idx < len(m.items) {
			selected = append(selected, m.items[idx])
		}
	}
	return selected
}

// IsQuitting returns whether the user has quit the list
func (m *Model) IsQuitting() bool {
	return m.quitting
}

// IsCancelled returns whether the user cancelled with Ctrl+C
func (m *Model) IsCancelled() bool {
	return m.cancelled
}

// GetCurrentIndex returns the index of the currently highlighted item
func (m *Model) GetCurrentIndex() int {
	return m.list.Index()
}

// GetCurrentItem returns the currently highlighted item
func (m *Model) GetCurrentItem() SelectableItem {
	idx := m.list.Index()
	if idx >= 0 && idx < len(m.items) {
		return m.items[idx]
	}
	return nil
}

