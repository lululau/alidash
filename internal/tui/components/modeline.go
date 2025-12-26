package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"aliyun-tui-viewer/internal/tui/types"
)

// ModeLineModel represents the status bar at the bottom
type ModeLineModel struct {
	profile   string
	page      types.PageType
	pageInfo  string // Optional additional info (e.g., page number)
	width     int
	styles    ModeLineStyles
}

// ModeLineStyles defines styles for the mode line
type ModeLineStyles struct {
	Background lipgloss.Style
	Profile    lipgloss.Style
	PageInfo   lipgloss.Style
	Help       lipgloss.Style
	Separator  lipgloss.Style
}

// DefaultModeLineStyles returns default mode line styles
func DefaultModeLineStyles() ModeLineStyles {
	return ModeLineStyles{
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")).
			Foreground(lipgloss.Color("#E5E7EB")),
		Profile: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true),
		PageInfo: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06B6D4")),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")),
	}
}

// NewModeLineModel creates a new mode line model
func NewModeLineModel(profile string, page types.PageType) ModeLineModel {
	return ModeLineModel{
		profile: profile,
		page:    page,
		styles:  DefaultModeLineStyles(),
	}
}

// SetProfile sets the current profile
func (m ModeLineModel) SetProfile(profile string) ModeLineModel {
	m.profile = profile
	return m
}

// SetPage sets the current page
func (m ModeLineModel) SetPage(page types.PageType) ModeLineModel {
	m.page = page
	m.pageInfo = "" // Clear page info when changing pages
	return m
}

// SetPageInfo sets additional page info (e.g., "Page 1/5")
func (m ModeLineModel) SetPageInfo(info string) ModeLineModel {
	m.pageInfo = info
	return m
}

// SetWidth sets the mode line width
func (m ModeLineModel) SetWidth(width int) ModeLineModel {
	m.width = width
	return m
}

// View renders the mode line
func (m ModeLineModel) View() string {
	// Left side: Profile
	leftPart := m.styles.Profile.Render(fmt.Sprintf(" Profile: %s ", m.profile))

	// Middle: Page info (if any)
	middlePart := ""
	if m.pageInfo != "" {
		middlePart = m.styles.PageInfo.Render(m.pageInfo)
	}

	// Right side: Shortcuts help
	shortcuts := m.getShortcuts()
	rightPart := m.styles.Help.Render(shortcuts + " ")

	// Calculate spacing
	leftLen := lipgloss.Width(leftPart)
	middleLen := lipgloss.Width(middlePart)
	rightLen := lipgloss.Width(rightPart)

	totalFixedLen := leftLen + middleLen + rightLen
	availableSpace := m.width - totalFixedLen

	if availableSpace < 0 {
		availableSpace = 0
	}

	// Distribute spacing
	var leftSpace, rightSpace int
	if middlePart != "" {
		leftSpace = availableSpace / 2
		rightSpace = availableSpace - leftSpace
	} else {
		leftSpace = availableSpace
		rightSpace = 0
	}

	var content strings.Builder
	content.WriteString(leftPart)
	content.WriteString(strings.Repeat(" ", leftSpace))
	if middlePart != "" {
		content.WriteString(middlePart)
		content.WriteString(strings.Repeat(" ", rightSpace))
	}
	content.WriteString(rightPart)

	return m.styles.Background.
		Width(m.width).
		Render(content.String())
}

// getShortcuts returns the context-sensitive shortcuts for the current page
func (m ModeLineModel) getShortcuts() string {
	switch m.page {
	case types.PageMenu:
		return "Enter: Select | j/k: Navigate | Q: Quit | O: Profile"

	case types.PageECSList:
		return "j/k: Navigate | Enter: Details | g: Security Groups | /: Search | yy: Copy | q: Back"

	case types.PageECSDetail:
		return "q/Esc: Back | yy: Copy | e: Edit | v: Pager | /: Search | n/N: Next/Prev"

	case types.PageSecurityGroups:
		return "j/k: Navigate | Enter: Rules | s: Instances | /: Search | yy: Copy | q: Back"

	case types.PageSecurityGroupRules:
		return "j/k: Navigate | /: Search | yy: Copy | q: Back"

	case types.PageSecurityGroupInstances, types.PageInstanceSecurityGroups:
		return "j/k: Navigate | Enter: Details | /: Search | yy: Copy | q: Back"

	case types.PageDNSDomains:
		return "j/k: Navigate | Enter: Records | /: Search | yy: Copy | q: Back"

	case types.PageDNSRecords:
		return "j/k: Navigate | /: Search | yy: Copy | q: Back"

	case types.PageSLBList:
		return "j/k: Navigate | Enter: Details | l: Listeners | v: VServer Groups | /: Search | q: Back"

	case types.PageSLBDetail:
		return "q/Esc: Back | yy: Copy | e: Edit | v: Pager | /: Search | n/N: Next/Prev"

	case types.PageSLBListeners:
		return "j/k: Navigate | /: Search | yy: Copy | q: Back"

	case types.PageSLBVServerGroups:
		return "j/k: Navigate | Enter: Backend Servers | /: Search | yy: Copy | q: Back"

	case types.PageSLBBackendServers:
		return "j/k: Navigate | /: Search | yy: Copy | q: Back"

	case types.PageOSSBuckets:
		return "j/k: Navigate | Enter: Objects | /: Search | q: Back"

	case types.PageOSSObjects:
		return "j/k: Navigate | Enter: Details | [/]: Prev/Next Page | 0: First | /: Search | q: Back"

	case types.PageOSSObjectDetail:
		return "q/Esc: Back | yy: Copy | e: Edit | v: Pager | /: Search"

	case types.PageRDSList:
		return "j/k: Navigate | Enter: Details | D: Databases | A: Accounts | /: Search | q: Back"

	case types.PageRDSDetail:
		return "q/Esc: Back | yy: Copy | e: Edit | v: Pager | /: Search"

	case types.PageRDSDatabases, types.PageRDSAccounts:
		return "j/k: Navigate | Enter: Details | /: Search | yy: Copy | q: Back"

	case types.PageRedisList:
		return "j/k: Navigate | Enter: Details | A: Accounts | /: Search | yy: Copy | q: Back"

	case types.PageRedisDetail:
		return "q/Esc: Back | yy: Copy | e: Edit | v: Pager | /: Search"

	case types.PageRedisAccounts:
		return "j/k: Navigate | Enter: Details | /: Search | yy: Copy | q: Back"

	case types.PageRocketMQList:
		return "j/k: Navigate | Enter: Details | T: Topics | G: Groups | /: Search | q: Back"

	case types.PageRocketMQDetail:
		return "q/Esc: Back | yy: Copy | e: Edit | v: Pager | /: Search"

	case types.PageRocketMQTopics, types.PageRocketMQGroups:
		return "j/k: Navigate | Enter: Details | /: Search | yy: Copy | q: Back"

	default:
		return "q/Esc: Back | Q: Quit | O: Profile"
	}
}
