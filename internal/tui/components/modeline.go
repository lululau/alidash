package components

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"aliyun-tui-viewer/internal/tui/types"
)

// ModeLineModel represents the status bar at the bottom
type ModeLineModel struct {
	profile  string
	region   string
	page     types.PageType
	pageInfo string // Optional additional info (e.g., page number)
	width    int
	styles   ModeLineStyles
}

// ModeLineStyles defines styles for the mode line
type ModeLineStyles struct {
	Background lipgloss.Style
	Key        lipgloss.Style
	Help       lipgloss.Style
	Separator  lipgloss.Style
}

// DefaultModeLineStyles returns default mode line styles
func DefaultModeLineStyles() ModeLineStyles {
	return ModeLineStyles{
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")).
			Foreground(lipgloss.Color("#E5E7EB")),
		Key: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Bold(true),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")),
	}
}

// NewModeLineModel creates a new mode line model
func NewModeLineModel(profile string, region string, page types.PageType) ModeLineModel {
	return ModeLineModel{
		profile: profile,
		region:  region,
		page:    page,
		styles:  DefaultModeLineStyles(),
	}
}

// SetProfile sets the current profile
func (m ModeLineModel) SetProfile(profile string) ModeLineModel {
	m.profile = profile
	return m
}

// SetRegion sets the current region
func (m ModeLineModel) SetRegion(region string) ModeLineModel {
	m.region = region
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

// GetProfile returns the current profile
func (m ModeLineModel) GetProfile() string {
	return m.profile
}

// GetRegion returns the current region
func (m ModeLineModel) GetRegion() string {
	return m.region
}

// View renders the mode line (shortcuts only, left-aligned)
func (m ModeLineModel) View() string {
	// Get shortcuts and format with highlighted keys
	shortcuts := m.getShortcuts()
	formattedShortcuts := m.formatShortcuts(shortcuts)

	// Add page info if present
	content := " " + formattedShortcuts
	if m.pageInfo != "" {
		content = " " + m.styles.Help.Render(m.pageInfo) + m.styles.Separator.Render(" | ") + formattedShortcuts
	}

	return m.styles.Background.
		Width(m.width).
		Render(content)
}

// formatShortcuts formats shortcuts with highlighted keys
// Input format: "Enter: Select | j/k: Navigate | Q: Quit"
// Output: keys are highlighted, descriptions are normal
func (m ModeLineModel) formatShortcuts(shortcuts string) string {
	// Split by " | "
	parts := strings.Split(shortcuts, " | ")
	var formattedParts []string

	// Regex to match "key: description" or "key/key: description" patterns
	// Also handles special keys like "n/N", "yy", etc.
	keyPattern := regexp.MustCompile(`^([^:]+):\s*(.*)$`)

	for _, part := range parts {
		match := keyPattern.FindStringSubmatch(part)
		if len(match) == 3 {
			key := match[1]
			desc := match[2]
			formatted := m.styles.Key.Render(key) + m.styles.Help.Render(": "+desc)
			formattedParts = append(formattedParts, formatted)
		} else {
			// No colon found, treat as plain text
			formattedParts = append(formattedParts, m.styles.Help.Render(part))
		}
	}

	return strings.Join(formattedParts, m.styles.Separator.Render(" | "))
}

// getShortcuts returns the context-sensitive shortcuts for the current page
func (m ModeLineModel) getShortcuts() string {
	switch m.page {
	case types.PageMenu:
		return "Enter: Select | j/k: Navigate | Q: Quit | P: Profile | R: Region"

	case types.PageECSList:
		return "j/k: Navigate | Enter: Details | v: JSON | s: Disks | g: Security Groups | /: Search | yy: Copy | q: Back"

	case types.PageECSDetail:
		return "j/k: Row | Tab/S-Tab: Section | yy: Copy | q/Esc: Back"

	case types.PageECSJSONDetail:
		return "q/Esc: Back | yy: Copy | e: Edit | v: Pager | /: Search | n/N: Next/Prev"

	case types.PageECSDisks:
		return "j/k: Navigate | Enter: Details | /: Search | yy: Copy | q: Back"

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
		return "q/Esc: Back | Q: Quit | P: Profile | R: Region"
	}
}
