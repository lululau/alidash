package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SearchModel represents a vim-style search bar
type SearchModel struct {
	input    textinput.Model
	Active   bool
	query    string
	width    int
	styles   SearchStyles
}

// SearchStyles defines styles for the search bar
type SearchStyles struct {
	Label      lipgloss.Style
	Input      lipgloss.Style
	Cursor     lipgloss.Style
	Background lipgloss.Style
}

// DefaultSearchStyles returns default search styles
func DefaultSearchStyles() SearchStyles {
	return SearchStyles{
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Bold(true),
		Input: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")),
		Cursor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")),
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")),
	}
}

// NewSearchModel creates a new search model
func NewSearchModel() SearchModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.CharLimit = 100
	ti.Width = 50
	ti.Prompt = "/"
	ti.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Bold(true)
	ti.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E5E7EB"))

	return SearchModel{
		input:  ti,
		Active: false,
		styles: DefaultSearchStyles(),
	}
}

// Activate activates the search bar
func (m SearchModel) Activate() SearchModel {
	m.Active = true
	m.input.SetValue("")
	m.input.Focus()
	return m
}

// Deactivate deactivates the search bar
func (m SearchModel) Deactivate() SearchModel {
	m.Active = false
	m.input.Blur()
	return m
}

// Focus focuses the search input
func (m SearchModel) Focus() tea.Cmd {
	return m.input.Focus()
}

// SetWidth sets the search bar width
func (m SearchModel) SetWidth(width int) SearchModel {
	m.width = width
	m.input.Width = width - 5 // Account for prompt and padding
	return m
}

// Query returns the current search query
func (m SearchModel) Query() string {
	return m.query
}

// SetQuery sets the search query
func (m SearchModel) SetQuery(query string) SearchModel {
	m.query = query
	m.input.SetValue(query)
	return m
}

// Init implements tea.Model
func (m SearchModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	if !m.Active {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Execute search
			m.query = m.input.Value()
			m.Active = false
			m.input.Blur()
			return m, func() tea.Msg {
				return SearchExecuteMsg{Query: m.query}
			}

		case tea.KeyEsc:
			// Cancel search
			m.Active = false
			m.input.Blur()
			return m, func() tea.Msg {
				return SearchCancelMsg{}
			}

		case tea.KeyCtrlC:
			// Cancel search
			m.Active = false
			m.input.Blur()
			return m, func() tea.Msg {
				return SearchCancelMsg{}
			}
		}
	}

	// Delegate to text input
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m SearchModel) View() string {
	if !m.Active {
		return ""
	}

	return m.styles.Background.
		Width(m.width).
		Render(m.input.View())
}

// SearchExecuteMsg is sent when search is executed
type SearchExecuteMsg struct {
	Query string
}

// SearchCancelMsg is sent when search is cancelled
type SearchCancelMsg struct{}

