package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ModalType represents different types of modals
type ModalType int

const (
	ModalTypeInfo ModalType = iota
	ModalTypeError
	ModalTypeSuccess
	ModalTypeConfirm
	ModalTypeProfileSelect
)

// ModalModel represents a modal dialog
type ModalModel struct {
	Visible       bool
	modalType     ModalType
	title         string
	message       string
	width         int
	height        int
	styles        ModalStyles

	// For profile selection
	profileList    list.Model
	profiles       []string
	currentProfile string
	selectedIndex  int
}

// ModalStyles defines styles for the modal
type ModalStyles struct {
	Overlay    lipgloss.Style
	Container  lipgloss.Style
	Title      lipgloss.Style
	Message    lipgloss.Style
	Button     lipgloss.Style
	InfoColor  lipgloss.Style
	ErrorColor lipgloss.Style
	SuccessColor lipgloss.Style
}

// DefaultModalStyles returns default modal styles
func DefaultModalStyles() ModalStyles {
	return ModalStyles{
		Overlay: lipgloss.NewStyle().
			Background(lipgloss.Color("#000000")),
		Container: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2).
			Background(lipgloss.Color("#1F2937")),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1),
		Message: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")),
		Button: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Background(lipgloss.Color("#374151")).
			Padding(0, 2),
		InfoColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6")),
		ErrorColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")),
		SuccessColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")),
	}
}

// NewModalModel creates a new modal model
func NewModalModel() ModalModel {
	return ModalModel{
		Visible: false,
		styles:  DefaultModalStyles(),
		width:   60,
		height:  10,
	}
}

// NewInfoModal creates an info modal
func NewInfoModal(message string) ModalModel {
	return ModalModel{
		Visible:   true,
		modalType: ModalTypeInfo,
		title:     "Info",
		message:   message,
		styles:    DefaultModalStyles(),
		width:     60,
		height:    10,
	}
}

// NewErrorModal creates an error modal
func NewErrorModal(message string) ModalModel {
	return ModalModel{
		Visible:   true,
		modalType: ModalTypeError,
		title:     "Error",
		message:   message,
		styles:    DefaultModalStyles(),
		width:     60,
		height:    10,
	}
}

// NewSuccessModal creates a success modal
func NewSuccessModal(message string) ModalModel {
	return ModalModel{
		Visible:   true,
		modalType: ModalTypeSuccess,
		title:     "Success",
		message:   message,
		styles:    DefaultModalStyles(),
		width:     60,
		height:    10,
	}
}

// NewProfileSelectionModal creates a profile selection modal
func NewProfileSelectionModal(profiles []string, currentProfile string) ModalModel {
	// Create list items
	items := make([]list.Item, len(profiles))
	selectedIdx := 0
	for i, p := range profiles {
		displayName := p
		if p == currentProfile {
			displayName = p + " (current)"
			selectedIdx = i
		}
		items[i] = profileItem{name: p, display: displayName}
	}

	// Create list
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#F59E0B")).
		Bold(true)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#E5E7EB"))

	l := list.New(items, delegate, 50, 15)
	l.Title = "Select Profile"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED"))

	// Select current profile
	l.Select(selectedIdx)

	return ModalModel{
		Visible:        true,
		modalType:      ModalTypeProfileSelect,
		title:          "Select Profile",
		profiles:       profiles,
		currentProfile: currentProfile,
		profileList:    l,
		selectedIndex:  selectedIdx,
		styles:         DefaultModalStyles(),
		width:          50,
		height:         min(len(profiles)+6, 20),
	}
}

// profileItem implements list.Item for profile selection
type profileItem struct {
	name    string
	display string
}

func (i profileItem) Title() string       { return i.display }
func (i profileItem) Description() string { return "" }
func (i profileItem) FilterValue() string { return i.name }

// Show shows the modal
func (m ModalModel) Show() ModalModel {
	m.Visible = true
	return m
}

// Hide hides the modal
func (m ModalModel) Hide() ModalModel {
	m.Visible = false
	return m
}

// SetSize sets the modal size
func (m ModalModel) SetSize(width, height int) ModalModel {
	m.width = width
	m.height = height
	return m
}

// Init implements tea.Model
func (m ModalModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ModalModel) Update(msg tea.Msg) (ModalModel, tea.Cmd) {
	if !m.Visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.modalType {
		case ModalTypeProfileSelect:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				// Select profile
				if item, ok := m.profileList.SelectedItem().(profileItem); ok {
					m.Visible = false
					return m, func() tea.Msg {
						return ProfileSelectedMsg{Profile: item.name}
					}
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))):
				m.Visible = false
				return m, func() tea.Msg {
					return ModalDismissedMsg{}
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("j", "down"))):
				var cmd tea.Cmd
				m.profileList, cmd = m.profileList.Update(msg)
				return m, cmd

			case key.Matches(msg, key.NewBinding(key.WithKeys("k", "up"))):
				var cmd tea.Cmd
				m.profileList, cmd = m.profileList.Update(msg)
				return m, cmd
			}

		default:
			// Info/Error/Success modals - dismiss on any key
			switch msg.Type {
			case tea.KeyEnter, tea.KeyEsc, tea.KeySpace:
				m.Visible = false
				return m, func() tea.Msg {
					return ModalDismissedMsg{}
				}
			}
		}
	}

	// Update profile list if applicable
	if m.modalType == ModalTypeProfileSelect {
		var cmd tea.Cmd
		m.profileList, cmd = m.profileList.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model
func (m ModalModel) View() string {
	if !m.Visible {
		return ""
	}

	var content strings.Builder

	switch m.modalType {
	case ModalTypeProfileSelect:
		return m.styles.Container.
			Width(m.width).
			Render(m.profileList.View())

	case ModalTypeInfo:
		title := m.styles.InfoColor.Render(m.title)
		content.WriteString(m.styles.Title.Render(title))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Message.Render(m.message))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Button.Render(" OK (Enter) "))

	case ModalTypeError:
		title := m.styles.ErrorColor.Render("⚠ " + m.title)
		content.WriteString(m.styles.Title.Render(title))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Message.Render(m.message))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Button.Render(" OK (Enter) "))

	case ModalTypeSuccess:
		title := m.styles.SuccessColor.Render("✓ " + m.title)
		content.WriteString(m.styles.Title.Render(title))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Message.Render(m.message))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Button.Render(" OK (Enter) "))
	}

	return m.styles.Container.
		Width(m.width).
		Render(content.String())
}

// ProfileSelectedMsg is sent when a profile is selected
type ProfileSelectedMsg struct {
	Profile string
}

// ModalDismissedMsg is sent when the modal is dismissed
type ModalDismissedMsg struct{}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

