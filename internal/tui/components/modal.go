package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"aliyun-tui-viewer/internal/i18n"
)

// ModalType represents different types of modals
type ModalType int

const (
	ModalTypeInfo ModalType = iota
	ModalTypeError
	ModalTypeSuccess
	ModalTypeConfirm
	ModalTypeProfileSelect
	ModalTypeRegionSelect
	ModalTypeInput // Input dialog for user text input
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

	// For region selection
	regionList     list.Model
	regions        []string
	currentRegion  string
	regionsLoading bool

	// For input dialog
	inputField   textinput.Model
	inputPrompt  string
	inputHistory []string // History items
	historyIndex int      // Current position in history (-1 means not browsing)
	currentInput string   // Saved current input when browsing history
}

// ModalStyles defines styles for the modal
type ModalStyles struct {
	Overlay      lipgloss.Style
	Container    lipgloss.Style
	Title        lipgloss.Style
	Message      lipgloss.Style
	Button       lipgloss.Style
	Help         lipgloss.Style
	InfoColor    lipgloss.Style
	ErrorColor   lipgloss.Style
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
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")),
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
		title:     i18n.T(i18n.KeyModalInfo),
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
		title:     i18n.T(i18n.KeyModalError),
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
		title:     i18n.T(i18n.KeyModalSuccess),
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
			displayName = p + " (" + i18n.T(i18n.KeyModalCurrent) + ")"
			selectedIdx = i
		}
		items[i] = profileItem{name: p, display: displayName}
	}

	// Create compact delegate with minimal spacing
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(1)            // Single line items
	delegate.SetSpacing(0)           // No spacing between items
	delegate.ShowDescription = false // Don't show description
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#F59E0B")).
		Bold(true).
		Background(lipgloss.Color("#374151"))
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#E5E7EB"))

	// Calculate appropriate height based on number of profiles
	// Add extra height for filter input (3 lines: prompt + input + spacing)
	listHeight := min(len(profiles)+6, 18)

	l := list.New(items, delegate, 50, listHeight)
	l.Title = i18n.T(i18n.KeyModalSelectProfile)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true) // Enable filtering for search
	l.SetShowHelp(true)         // Show help to indicate / for filter
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED"))
	l.Styles.FilterPrompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))
	l.Styles.FilterCursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))
	l.FilterInput.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))
	l.FilterInput.Cursor.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))

	// Select current profile
	l.Select(selectedIdx)

	return ModalModel{
		Visible:        true,
		modalType:      ModalTypeProfileSelect,
		title:          i18n.T(i18n.KeyModalSelectProfile),
		profiles:       profiles,
		currentProfile: currentProfile,
		profileList:    l,
		selectedIndex:  selectedIdx,
		styles:         DefaultModalStyles(),
		width:          55,
		height:         listHeight + 4,
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

// regionItem implements list.Item for region selection
type regionItem struct {
	id      string
	display string
}

func (i regionItem) Title() string       { return i.display }
func (i regionItem) Description() string { return "" }
func (i regionItem) FilterValue() string { return i.id }

// NewRegionSelectionModal creates a region selection modal in loading state
func NewRegionSelectionModal(currentRegion string) ModalModel {
	return ModalModel{
		Visible:        true,
		modalType:      ModalTypeRegionSelect,
		title:          i18n.T(i18n.KeyModalSelectRegion),
		currentRegion:  currentRegion,
		regionsLoading: true,
		styles:         DefaultModalStyles(),
		width:          60,
		height:         15,
	}
}

// NewInputModal creates an input dialog modal
func NewInputModal(title, prompt, placeholder string) ModalModel {
	return NewInputModalWithHistory(title, prompt, placeholder, nil)
}

// NewInputModalWithHistory creates an input dialog modal with history support
func NewInputModalWithHistory(title, prompt, placeholder string, history []string) ModalModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.Width = 50
	ti.Prompt = ""
	ti.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B")).
		Bold(true)
	ti.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E5E7EB"))
	ti.PlaceholderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280"))
	ti.Focus()

	return ModalModel{
		Visible:      true,
		modalType:    ModalTypeInput,
		title:        title,
		inputPrompt:  prompt,
		inputField:   ti,
		inputHistory: history,
		historyIndex: -1, // Not browsing history
		styles:       DefaultModalStyles(),
		width:        60,
		height:       12,
	}
}

// SetRegions updates the region list and exits loading state
func (m ModalModel) SetRegions(regions []string, currentRegion string) ModalModel {
	if m.modalType != ModalTypeRegionSelect {
		return m
	}

	// Create list items
	items := make([]list.Item, len(regions))
	selectedIdx := 0
	for i, r := range regions {
		displayName := r
		if r == currentRegion {
			displayName = r + " (" + i18n.T(i18n.KeyModalCurrent) + ")"
			selectedIdx = i
		}
		items[i] = regionItem{id: r, display: displayName}
	}

	// Create compact delegate with minimal spacing
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#F59E0B")).
		Bold(true).
		Background(lipgloss.Color("#374151"))
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#E5E7EB"))

	// Calculate appropriate height based on number of regions
	// Add extra height for filter input (3 lines: prompt + input + spacing)
	listHeight := min(len(regions)+6, 18)

	l := list.New(items, delegate, 55, listHeight)
	l.Title = i18n.T(i18n.KeyModalSelectRegion)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true) // Show help to indicate / for filter
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED"))
	l.Styles.FilterPrompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))
	l.Styles.FilterCursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))
	l.FilterInput.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))
	l.FilterInput.Cursor.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F59E0B"))

	// Select current region
	l.Select(selectedIdx)

	m.regionList = l
	m.regions = regions
	m.currentRegion = currentRegion
	m.regionsLoading = false
	m.height = listHeight + 4

	return m
}

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
			// If list is filtering, let it handle all key events
			if m.profileList.FilterState() == list.Filtering {
				var cmd tea.Cmd
				m.profileList, cmd = m.profileList.Update(msg)
				return m, cmd
			}

			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				// Select profile
				if item, ok := m.profileList.SelectedItem().(profileItem); ok {
					m.Visible = false
					return m, func() tea.Msg {
						return ProfileSelectedMsg{Profile: item.name}
					}
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				// If filtering, let the list handle esc to cancel filter
				// Otherwise dismiss modal
				m.Visible = false
				return m, func() tea.Msg {
					return ModalDismissedMsg{}
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("q"))):
				m.Visible = false
				return m, func() tea.Msg {
					return ModalDismissedMsg{}
				}
			}

			// Forward all other keys to list (for j/k navigation and / filtering)
			var cmd tea.Cmd
			m.profileList, cmd = m.profileList.Update(msg)
			return m, cmd

		case ModalTypeRegionSelect:
			// If still loading, only allow dismiss
			if m.regionsLoading {
				switch {
				case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))):
					m.Visible = false
					return m, func() tea.Msg {
						return ModalDismissedMsg{}
					}
				}
				return m, nil
			}

			// If list is filtering, let it handle all key events
			if m.regionList.FilterState() == list.Filtering {
				var cmd tea.Cmd
				m.regionList, cmd = m.regionList.Update(msg)
				return m, cmd
			}

			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				// Select region
				if item, ok := m.regionList.SelectedItem().(regionItem); ok {
					m.Visible = false
					return m, func() tea.Msg {
						return RegionSelectedMsg{Region: item.id}
					}
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				m.Visible = false
				return m, func() tea.Msg {
					return ModalDismissedMsg{}
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("q"))):
				m.Visible = false
				return m, func() tea.Msg {
					return ModalDismissedMsg{}
				}
			}

			// Forward all other keys to list (for j/k navigation and / filtering)
			var cmd tea.Cmd
			m.regionList, cmd = m.regionList.Update(msg)
			return m, cmd

		case ModalTypeInput:
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				// Submit input
				value := m.inputField.Value()
				if value != "" {
					m.Visible = false
					m.historyIndex = -1
					return m, func() tea.Msg {
						return InputSubmittedMsg{Value: value}
					}
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				m.Visible = false
				m.historyIndex = -1
				return m, func() tea.Msg {
					return ModalDismissedMsg{}
				}

			case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+p"))):
				// Previous history item
				if len(m.inputHistory) > 0 {
					if m.historyIndex == -1 {
						// Save current input before browsing history
						m.currentInput = m.inputField.Value()
						m.historyIndex = 0
					} else if m.historyIndex < len(m.inputHistory)-1 {
						m.historyIndex++
					}
					m.inputField.SetValue(m.inputHistory[m.historyIndex])
					m.inputField.CursorEnd()
				}
				return m, nil

			case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+n"))):
				// Next history item (more recent)
				if m.historyIndex > 0 {
					m.historyIndex--
					m.inputField.SetValue(m.inputHistory[m.historyIndex])
					m.inputField.CursorEnd()
				} else if m.historyIndex == 0 {
					// Return to current input
					m.historyIndex = -1
					m.inputField.SetValue(m.currentInput)
					m.inputField.CursorEnd()
				}
				return m, nil
			}

			// Forward all other keys to text input
			var cmd tea.Cmd
			m.inputField, cmd = m.inputField.Update(msg)
			return m, cmd

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

	// Update region list if applicable
	if m.modalType == ModalTypeRegionSelect && !m.regionsLoading {
		var cmd tea.Cmd
		m.regionList, cmd = m.regionList.Update(msg)
		return m, cmd
	}

	// Update input field if applicable
	if m.modalType == ModalTypeInput {
		var cmd tea.Cmd
		m.inputField, cmd = m.inputField.Update(msg)
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

	case ModalTypeRegionSelect:
		if m.regionsLoading {
			content.WriteString(m.styles.Title.Render(i18n.T(i18n.KeyModalSelectRegion)))
			content.WriteString("\n\n")
			content.WriteString(m.styles.Message.Render(i18n.T(i18n.KeyModalLoading)))
			content.WriteString("\n\n")
			content.WriteString(m.styles.Help.Render("Esc: " + i18n.T(i18n.KeyModalCancel)))
			return m.styles.Container.
				Width(m.width).
				Render(content.String())
		}
		return m.styles.Container.
			Width(m.width).
			Render(m.regionList.View())

	case ModalTypeInfo:
		title := m.styles.InfoColor.Render(m.title)
		content.WriteString(m.styles.Title.Render(title))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Message.Render(m.message))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Button.Render(" " + i18n.T(i18n.KeyModalOK) + " (Enter) "))

	case ModalTypeError:
		title := m.styles.ErrorColor.Render("⚠ " + m.title)
		content.WriteString(m.styles.Title.Render(title))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Message.Render(m.message))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Button.Render(" " + i18n.T(i18n.KeyModalOK) + " (Enter) "))

	case ModalTypeSuccess:
		title := m.styles.SuccessColor.Render("✓ " + m.title)
		content.WriteString(m.styles.Title.Render(title))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Message.Render(m.message))
		content.WriteString("\n\n")
		content.WriteString(m.styles.Button.Render(" " + i18n.T(i18n.KeyModalOK) + " (Enter) "))

	case ModalTypeInput:
		content.WriteString(m.styles.Title.Render(m.title))
		content.WriteString("\n\n")
		if m.inputPrompt != "" {
			content.WriteString(m.styles.Message.Render(m.inputPrompt))
			content.WriteString("\n\n")
		}
		content.WriteString(m.inputField.View())
		content.WriteString("\n\n")
		helpText := "Enter: " + i18n.T(i18n.KeyModalConfirm) + " | Esc: " + i18n.T(i18n.KeyModalCancel)
		if len(m.inputHistory) > 0 {
			helpText += " | C-p/C-n: " + i18n.T(i18n.KeyModalHistory)
			if m.historyIndex >= 0 {
				helpText += fmt.Sprintf(" [%d/%d]", m.historyIndex+1, len(m.inputHistory))
			}
		}
		content.WriteString(m.styles.Help.Render(helpText))
	}

	return m.styles.Container.
		Width(m.width).
		Render(content.String())
}

// ProfileSelectedMsg is sent when a profile is selected
type ProfileSelectedMsg struct {
	Profile string
}

// RegionSelectedMsg is sent when a region is selected
type RegionSelectedMsg struct {
	Region string
}

// InputSubmittedMsg is sent when input is submitted
type InputSubmittedMsg struct {
	Value string
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

