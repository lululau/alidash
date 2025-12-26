package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"aliyun-tui-viewer/internal/i18n"
)

// HeaderModel represents the header bar at the top
type HeaderModel struct {
	title   string
	profile string
	region  string
	width   int
	styles  HeaderStyles
}

// HeaderStyles defines styles for the header
type HeaderStyles struct {
	Background lipgloss.Style
	Title      lipgloss.Style
	Profile    lipgloss.Style
	Region     lipgloss.Style
	Separator  lipgloss.Style
}

// DefaultHeaderStyles returns default header styles
func DefaultHeaderStyles() HeaderStyles {
	return HeaderStyles{
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")).
			Foreground(lipgloss.Color("#E5E7EB")),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")),
		Profile: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06B6D4")),
		Region: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")),
	}
}

// NewHeaderModel creates a new header model
func NewHeaderModel(title, profile, region string) HeaderModel {
	return HeaderModel{
		title:   title,
		profile: profile,
		region:  region,
		styles:  DefaultHeaderStyles(),
	}
}

// SetTitle sets the header title
func (m HeaderModel) SetTitle(title string) HeaderModel {
	m.title = title
	return m
}

// SetProfile sets the current profile
func (m HeaderModel) SetProfile(profile string) HeaderModel {
	m.profile = profile
	return m
}

// SetRegion sets the current region
func (m HeaderModel) SetRegion(region string) HeaderModel {
	m.region = region
	return m
}

// SetWidth sets the header width
func (m HeaderModel) SetWidth(width int) HeaderModel {
	m.width = width
	return m
}

// View renders the header
func (m HeaderModel) View() string {
	// Title
	titlePart := m.styles.Title.Render(" " + m.title + " ")

	// Profile and Region
	sep := m.styles.Separator.Render(" | ")
	profilePart := m.styles.Profile.Render(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyHeaderProfile), m.profile))
	regionPart := m.styles.Region.Render(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyHeaderRegion), m.region))

	content := titlePart + sep + profilePart + sep + regionPart

	// Pad to full width
	contentLen := lipgloss.Width(content)
	if contentLen < m.width {
		content = content + lipgloss.NewStyle().
			Width(m.width-contentLen).
			Render("")
	}

	// Add empty line after header for spacing
	return m.styles.Background.
		Width(m.width).
		Render(content) + "\n"
}

