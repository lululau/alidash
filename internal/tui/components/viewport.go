package components

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewportModel wraps bubbles/viewport for JSON detail views
type ViewportModel struct {
	viewport    viewport.Model
	title       string
	content     string
	rawContent  string // Original unhighlighted content
	data        interface{}
	width       int
	height      int
	focused     bool
	searchQuery string
	searchIndex int
	searchCount int
	keys        ViewportKeyMap

	// Yank tracking
	yankLastTime time.Time
	yankCount    int

	// Styles
	styles ViewportStyles
}

// ViewportKeyMap defines key bindings for the viewport
type ViewportKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Yank     key.Binding
	Edit     key.Binding
	Pager    key.Binding
}

// DefaultViewportKeyMap returns default key bindings
func DefaultViewportKeyMap() ViewportKeyMap {
	return ViewportKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("pgdn", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "top"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "bottom"),
		),
		Yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("yy", "copy"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Pager: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "pager"),
		),
	}
}

// ViewportStyles defines styles for the viewport
type ViewportStyles struct {
	Border      lipgloss.Style
	Title       lipgloss.Style
	JSONKey     lipgloss.Style
	JSONString  lipgloss.Style
	JSONNumber  lipgloss.Style
	JSONBoolean lipgloss.Style
	JSONNull    lipgloss.Style
	SearchMatch lipgloss.Style
}

// DefaultViewportStyles returns default viewport styles
func DefaultViewportStyles() ViewportStyles {
	return ViewportStyles{
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")),
		JSONKey: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06B6D4")),
		JSONString: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")),
		JSONNumber: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")),
		JSONBoolean: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")),
		JSONNull: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true),
		SearchMatch: lipgloss.NewStyle().
			Background(lipgloss.Color("#CA8A04")).
			Foreground(lipgloss.Color("#000000")),
	}
}

// NewViewportModel creates a new viewport model
func NewViewportModel(title string, data interface{}) ViewportModel {
	vp := viewport.New(80, 20)
	// No border style here - we add it in View() for full-width control

	m := ViewportModel{
		viewport: vp,
		title:    title,
		data:     data,
		keys:     DefaultViewportKeyMap(),
		styles:   DefaultViewportStyles(),
		focused:  true,
	}

	// Format and set content
	m = m.formatContent()

	return m
}

// formatContent formats the data as JSON with syntax highlighting
func (m ViewportModel) formatContent() ViewportModel {
	jsonData, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		m.rawContent = fmt.Sprintf("Error marshaling JSON: %v", err)
		m.content = m.rawContent
		m.viewport.SetContent(m.content)
		return m
	}

	m.rawContent = string(jsonData)
	m.content = m.highlightJSON(m.rawContent)
	m.viewport.SetContent(m.content)
	return m
}

// highlightJSON applies syntax highlighting to JSON
func (m ViewportModel) highlightJSON(jsonStr string) string {
	var result strings.Builder

	// Process line by line for better control
	lines := strings.Split(jsonStr, "\n")
	for i, line := range lines {
		highlighted := m.highlightLine(line)
		result.WriteString(highlighted)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// highlightLine highlights a single line of JSON
func (m ViewportModel) highlightLine(line string) string {
	// Match JSON key
	keyRegex := regexp.MustCompile(`"([^"]+)"(\s*):`)
	line = keyRegex.ReplaceAllStringFunc(line, func(match string) string {
		// Extract key name
		parts := keyRegex.FindStringSubmatch(match)
		if len(parts) >= 3 {
			key := m.styles.JSONKey.Render(`"` + parts[1] + `"`)
			return key + parts[2] + ":"
		}
		return match
	})

	// Match string values (not keys)
	stringRegex := regexp.MustCompile(`:\s*"([^"]*)"`)
	line = stringRegex.ReplaceAllStringFunc(line, func(match string) string {
		parts := stringRegex.FindStringSubmatch(match)
		if len(parts) >= 2 {
			value := m.styles.JSONString.Render(`"` + parts[1] + `"`)
			return ": " + value
		}
		return match
	})

	// Match numbers
	numberRegex := regexp.MustCompile(`:\s*(-?\d+\.?\d*)([,\s\n]|$)`)
	line = numberRegex.ReplaceAllStringFunc(line, func(match string) string {
		parts := numberRegex.FindStringSubmatch(match)
		if len(parts) >= 3 {
			value := m.styles.JSONNumber.Render(parts[1])
			return ": " + value + parts[2]
		}
		return match
	})

	// Match booleans
	boolRegex := regexp.MustCompile(`:\s*(true|false)([,\s\n]|$)`)
	line = boolRegex.ReplaceAllStringFunc(line, func(match string) string {
		parts := boolRegex.FindStringSubmatch(match)
		if len(parts) >= 3 {
			value := m.styles.JSONBoolean.Render(parts[1])
			return ": " + value + parts[2]
		}
		return match
	})

	// Match null
	nullRegex := regexp.MustCompile(`:\s*(null)([,\s\n]|$)`)
	line = nullRegex.ReplaceAllStringFunc(line, func(match string) string {
		parts := nullRegex.FindStringSubmatch(match)
		if len(parts) >= 3 {
			value := m.styles.JSONNull.Render(parts[1])
			return ": " + value + parts[2]
		}
		return match
	})

	return line
}

// SetData sets the data and reformats
func (m ViewportModel) SetData(data interface{}) ViewportModel {
	m.data = data
	return m.formatContent()
}

// SetTitle sets the title
func (m ViewportModel) SetTitle(title string) ViewportModel {
	m.title = title
	return m
}

// SetSize sets the viewport size
func (m ViewportModel) SetSize(width, height int) ViewportModel {
	m.width = width
	m.height = height
	// Account for: title (1) + help (1) + border (2) + search info (2)
	vpHeight := height - 6
	if vpHeight < 1 {
		vpHeight = 1
	}
	// Account for border width (2) and padding (2)
	vpWidth := width - 6
	if vpWidth < 10 {
		vpWidth = 10
	}
	m.viewport.Width = vpWidth
	m.viewport.Height = vpHeight
	return m
}

// SetFocused sets the focus state
func (m ViewportModel) SetFocused(focused bool) ViewportModel {
	m.focused = focused
	return m
}

// GetData returns the raw data
func (m ViewportModel) GetData() interface{} {
	return m.data
}

// Init implements tea.Model
func (m ViewportModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ViewportModel) Update(msg tea.Msg) (ViewportModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Yank):
			// Handle yank (double-y for copy)
			now := time.Now()
			if now.Sub(m.yankLastTime) < 500*time.Millisecond {
				m.yankCount++
			} else {
				m.yankCount = 1
			}
			m.yankLastTime = now

			if m.yankCount >= 2 {
				m.yankCount = 0
				return m, func() tea.Msg {
					return CopyDataMsg{Data: m.data}
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Edit):
			return m, func() tea.Msg {
				return OpenEditorMsg{Data: m.data}
			}

		case key.Matches(msg, m.keys.Pager):
			return m, func() tea.Msg {
				return OpenPagerMsg{Data: m.data}
			}
		}
	}

	// Delegate to underlying viewport
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m ViewportModel) View() string {
	var b strings.Builder

	// Title
	if m.title != "" {
		title := m.styles.Title.Render(m.title)
		b.WriteString(title)
		b.WriteString("\n")
	}

	// Help text
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Render("q/Esc: back | yy: copy | e: edit | v: pager | /: search | n/N: next/prev")
	b.WriteString(help)
	b.WriteString("\n")

	// Viewport with border that fills the width
	viewportContent := m.viewport.View()
	bordered := m.styles.Border.
		Width(m.width - 2).
		Render(viewportContent)
	b.WriteString(bordered)

	// Search info
	if m.searchQuery != "" {
		searchInfo := fmt.Sprintf(" Search: %s (%d/%d) ", m.searchQuery, m.searchIndex+1, m.searchCount)
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Render(searchInfo))
	}

	return b.String()
}

// Search searches for a query in the content
func (m ViewportModel) Search(query string) ViewportModel {
	if query == "" {
		m.searchQuery = ""
		m.searchIndex = -1
		m.searchCount = 0
		m.content = m.highlightJSON(m.rawContent)
		m.viewport.SetContent(m.content)
		return m
	}

	m.searchQuery = query
	lowerQuery := strings.ToLower(query)

	// Count matches
	m.searchCount = strings.Count(strings.ToLower(m.rawContent), lowerQuery)

	if m.searchCount > 0 {
		m.searchIndex = 0
		// Highlight matches
		m.content = m.highlightSearchMatches(query)
		m.viewport.SetContent(m.content)
		// Scroll to first match
		m.scrollToMatch(0)
	} else {
		m.searchIndex = -1
	}

	return m
}

// highlightSearchMatches highlights search matches in the content
func (m ViewportModel) highlightSearchMatches(query string) string {
	// First apply JSON highlighting
	highlighted := m.highlightJSON(m.rawContent)

	// Then highlight search matches (case insensitive)
	lowerHighlighted := strings.ToLower(highlighted)
	lowerQuery := strings.ToLower(query)

	var result strings.Builder
	lastIdx := 0

	for {
		idx := strings.Index(lowerHighlighted[lastIdx:], lowerQuery)
		if idx == -1 {
			break
		}

		actualIdx := lastIdx + idx
		result.WriteString(highlighted[lastIdx:actualIdx])

		// Get original case match
		match := highlighted[actualIdx : actualIdx+len(query)]
		result.WriteString(m.styles.SearchMatch.Render(match))

		lastIdx = actualIdx + len(query)
	}

	result.WriteString(highlighted[lastIdx:])
	return result.String()
}

// scrollToMatch scrolls to the nth match
func (m *ViewportModel) scrollToMatch(matchIndex int) {
	if matchIndex < 0 || matchIndex >= m.searchCount {
		return
	}

	lowerContent := strings.ToLower(m.rawContent)
	lowerQuery := strings.ToLower(m.searchQuery)

	// Find the position of the nth match
	pos := 0
	for i := 0; i <= matchIndex; i++ {
		idx := strings.Index(lowerContent[pos:], lowerQuery)
		if idx == -1 {
			return
		}
		pos += idx
		if i < matchIndex {
			pos += len(m.searchQuery)
		}
	}

	// Count newlines before this position to find line number
	lineNum := strings.Count(m.rawContent[:pos], "\n")

	// Scroll to that line
	m.viewport.SetYOffset(lineNum)
}

// NextSearchMatch moves to the next search match
func (m ViewportModel) NextSearchMatch() ViewportModel {
	if m.searchQuery == "" || m.searchCount == 0 {
		return m
	}

	m.searchIndex = (m.searchIndex + 1) % m.searchCount
	m.scrollToMatch(m.searchIndex)
	return m
}

// PrevSearchMatch moves to the previous search match
func (m ViewportModel) PrevSearchMatch() ViewportModel {
	if m.searchQuery == "" || m.searchCount == 0 {
		return m
	}

	m.searchIndex = (m.searchIndex - 1 + m.searchCount) % m.searchCount
	m.scrollToMatch(m.searchIndex)
	return m
}

// ClearSearch clears the search
func (m ViewportModel) ClearSearch() ViewportModel {
	m.searchQuery = ""
	m.searchIndex = -1
	m.searchCount = 0
	m.content = m.highlightJSON(m.rawContent)
	m.viewport.SetContent(m.content)
	return m
}

// OpenEditorMsg requests opening in external editor
type OpenEditorMsg struct {
	Data interface{}
}

// OpenPagerMsg requests opening in external pager
type OpenPagerMsg struct {
	Data interface{}
}

