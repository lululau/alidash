package pages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"aliyun-tui-viewer/internal/i18n"
	"aliyun-tui-viewer/internal/service"
	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// FinderSection represents a section in the finder results
type FinderSection struct {
	Title    string
	Columns  []string
	ColWidths []int // Fixed column widths
	Rows     [][]string
	Data     []interface{} // Original data for navigation
	PageType types.PageType
}

// FinderModel represents the resource finder results page
type FinderModel struct {
	result         *service.FindResult
	sections       []FinderSection
	viewport       viewport.Model // Scrollable viewport
	currentSection int
	currentRow     int
	width          int
	height         int
	keys           FinderKeyMap
	yankLastTime   time.Time
	yankCount      int
	styles         FinderStyles
}

// FinderStyles defines styles for the finder
type FinderStyles struct {
	Header       lipgloss.Style
	Cell         lipgloss.Style
	Selected     lipgloss.Style
	Border       lipgloss.Style
	FocusedBorder lipgloss.Style
	Title        lipgloss.Style
	SectionTitle lipgloss.Style
	Separator    lipgloss.Style
	Empty        lipgloss.Style
}

// DefaultFinderStyles returns default finder styles matching the table component
func DefaultFinderStyles() FinderStyles {
	return FinderStyles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F59E0B")).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Padding(0, 1),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Bold(true).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")),
		FocusedBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")),
		SectionTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#374151")),
		Empty: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true).
			Padding(0, 1),
	}
}

// FinderKeyMap defines key bindings for the finder
type FinderKeyMap struct {
	Up          key.Binding
	Down        key.Binding
	NextSection key.Binding
	PrevSection key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	Top         key.Binding
	Bottom      key.Binding
	Enter       key.Binding
	Yank        key.Binding
}

// DefaultFinderKeyMap returns default key bindings
func DefaultFinderKeyMap() FinderKeyMap {
	return FinderKeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		NextSection: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next section"),
		),
		PrevSection: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("S-tab", "prev section"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("ctrl+b", "pgup"),
			key.WithHelp("ctrl+b", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("ctrl+f", "pgdown"),
			key.WithHelp("ctrl+f", "page down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "go to top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "go to bottom"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		Yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("yy", "copy"),
		),
	}
}

// NewFinderModel creates a new finder model
func NewFinderModel(result *service.FindResult) FinderModel {
	m := FinderModel{
		result:   result,
		keys:     DefaultFinderKeyMap(),
		styles:   DefaultFinderStyles(),
		viewport: viewport.New(80, 20), // Initial size, will be updated by SetSize
	}
	m.buildSections()
	m.updateViewportContent()
	return m
}

// buildSections builds the display sections from the result - always shows all resource types
func (m *FinderModel) buildSections() {
	m.sections = nil

	if m.result == nil {
		return
	}

	// ECS Instances Section - always show
	ecsSection := FinderSection{
		Title:     fmt.Sprintf("%s (%d)", i18n.T(i18n.KeyFinderECS), len(m.result.ECSInstances)),
		Columns:   []string{i18n.T(i18n.KeyColInstanceID), i18n.T(i18n.KeyColName), i18n.T(i18n.KeyColPublicIP), i18n.T(i18n.KeyColPrivateIP), i18n.T(i18n.KeyColStatus)},
		ColWidths: []int{24, 20, 16, 16, 10},
		PageType:  types.PageECSDetail,
	}
	for _, inst := range m.result.ECSInstances {
		publicIP := "-"
		if len(inst.PublicIpAddress.IpAddress) > 0 {
			publicIP = inst.PublicIpAddress.IpAddress[0]
		} else if inst.EipAddress.IpAddress != "" {
			publicIP = inst.EipAddress.IpAddress
		}
		privateIP := "-"
		if len(inst.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
			privateIP = inst.VpcAttributes.PrivateIpAddress.IpAddress[0]
		}
		ecsSection.Rows = append(ecsSection.Rows, []string{
			inst.InstanceId,
			inst.InstanceName,
			publicIP,
			privateIP,
			inst.Status,
		})
		ecsSection.Data = append(ecsSection.Data, inst)
	}
	m.sections = append(m.sections, ecsSection)

	// ENI Section - always show
	eniSection := FinderSection{
		Title:     fmt.Sprintf("%s (%d)", i18n.T(i18n.KeyFinderENI), len(m.result.ENIs)),
		Columns:   []string{i18n.T(i18n.KeyColENIID), i18n.T(i18n.KeyColPrivateIP), i18n.T(i18n.KeyColType), i18n.T(i18n.KeyColStatus), i18n.T(i18n.KeyColAttachedInst)},
		ColWidths: []int{24, 16, 10, 10, 24},
		PageType:  types.PageECSJSONDetail,
	}
	for _, eni := range m.result.ENIs {
		eniType := i18n.T(i18n.KeyENISecondary)
		if eni.Type == "Primary" {
			eniType = i18n.T(i18n.KeyENIPrimary)
		}
		eniSection.Rows = append(eniSection.Rows, []string{
			eni.NetworkInterfaceId,
			eni.PrivateIpAddress,
			eniType,
			eni.Status,
			eni.InstanceId,
		})
		eniSection.Data = append(eniSection.Data, eni)
	}
	m.sections = append(m.sections, eniSection)

	// SLB Section - always show
	slbSection := FinderSection{
		Title:     fmt.Sprintf("%s (%d)", i18n.T(i18n.KeyFinderSLB), len(m.result.SLBInstances)),
		Columns:   []string{i18n.T(i18n.KeyColSLBID), i18n.T(i18n.KeyColName), i18n.T(i18n.KeyColAddress), i18n.T(i18n.KeyColType), i18n.T(i18n.KeyColStatus)},
		ColWidths: []int{24, 20, 16, 16, 10},
		PageType:  types.PageSLBDetail,
	}
	for _, lb := range m.result.SLBInstances {
		slbSection.Rows = append(slbSection.Rows, []string{
			lb.LoadBalancerId,
			lb.LoadBalancerName,
			lb.Address,
			lb.LoadBalancerSpec,
			lb.LoadBalancerStatus,
		})
		slbSection.Data = append(slbSection.Data, lb)
	}
	m.sections = append(m.sections, slbSection)

	// DNS Records Section - always show
	dnsSection := FinderSection{
		Title:     fmt.Sprintf("%s (%d)", i18n.T(i18n.KeyFinderDNS), len(m.result.DNSRecords)),
		Columns:   []string{i18n.T(i18n.KeyColDomain), i18n.T(i18n.KeyColRR), i18n.T(i18n.KeyColType), i18n.T(i18n.KeyColRecordValue), i18n.T(i18n.KeyColTTL)},
		ColWidths: []int{24, 20, 8, 20, 8},
		PageType:  types.PageECSJSONDetail,
	}
	for _, match := range m.result.DNSRecords {
		dnsSection.Rows = append(dnsSection.Rows, []string{
			match.DomainName,
			match.Record.RR,
			match.Record.Type,
			match.Record.Value,
			fmt.Sprintf("%d", match.Record.TTL),
		})
		dnsSection.Data = append(dnsSection.Data, match.Record)
	}
	m.sections = append(m.sections, dnsSection)

	// RDS Section - always show
	rdsSection := FinderSection{
		Title:     fmt.Sprintf("%s (%d)", i18n.T(i18n.KeyFinderRDS), len(m.result.RDSInstances)),
		Columns:   []string{i18n.T(i18n.KeyColInstanceID), i18n.T(i18n.KeyColDescription), i18n.T(i18n.KeyColEngine), i18n.T(i18n.KeyColConnString), i18n.T(i18n.KeyColConnString), i18n.T(i18n.KeyColStatus)},
		ColWidths: []int{24, 16, 12, 30, 30, 10},
		PageType:  types.PageRDSDetail,
	}
	for _, detail := range m.result.RDSInstances {
		internalAddr := detail.InternalConnectionStr
		if internalAddr == "" {
			internalAddr = "-"
		}
		publicAddr := detail.PublicConnectionStr
		if publicAddr == "" {
			publicAddr = "-"
		}
		rdsSection.Rows = append(rdsSection.Rows, []string{
			detail.Instance.DBInstanceId,
			detail.Instance.DBInstanceDescription,
			fmt.Sprintf("%s %s", detail.Instance.Engine, detail.Instance.EngineVersion),
			internalAddr,
			publicAddr,
			detail.Instance.DBInstanceStatus,
		})
		rdsSection.Data = append(rdsSection.Data, detail.Instance)
	}
	m.sections = append(m.sections, rdsSection)

	// Redis Section - always show
	redisSection := FinderSection{
		Title:     fmt.Sprintf("%s (%d)", i18n.T(i18n.KeyFinderRedis), len(m.result.RedisInstances)),
		Columns:   []string{i18n.T(i18n.KeyColInstanceID), i18n.T(i18n.KeyColName), i18n.T(i18n.KeyColConnDomain), i18n.T(i18n.KeyColPrivateIP), i18n.T(i18n.KeyColStatus)},
		ColWidths: []int{24, 20, 30, 16, 10},
		PageType:  types.PageRedisDetail,
	}
	for _, inst := range m.result.RedisInstances {
		redisSection.Rows = append(redisSection.Rows, []string{
			inst.InstanceId,
			inst.InstanceName,
			inst.ConnectionDomain,
			inst.PrivateIp,
			inst.InstanceStatus,
		})
		redisSection.Data = append(redisSection.Data, inst)
	}
	m.sections = append(m.sections, redisSection)

	// RocketMQ Section - always show
	rocketmqSection := FinderSection{
		Title:     fmt.Sprintf("%s (%d)", i18n.T(i18n.KeyFinderRocketMQ), len(m.result.RocketMQInstances)),
		Columns:   []string{i18n.T(i18n.KeyColInstanceID), i18n.T(i18n.KeyColName), i18n.T(i18n.KeyColStatus)},
		ColWidths: []int{30, 30, 12},
		PageType:  types.PageRocketMQDetail,
	}
	for _, inst := range m.result.RocketMQInstances {
		status := i18n.T(i18n.KeyStatusUnknown)
		switch inst.InstanceStatus {
		case 0:
			status = i18n.T(i18n.KeyStatusCreating)
		case 2:
			status = i18n.T(i18n.KeyStatusRunning)
		case 5:
			status = i18n.T(i18n.KeyStatusReleased)
		}
		rocketmqSection.Rows = append(rocketmqSection.Rows, []string{
			inst.InstanceId,
			inst.InstanceName,
			status,
		})
		rocketmqSection.Data = append(rocketmqSection.Data, inst)
	}
	m.sections = append(m.sections, rocketmqSection)
}

// SetSize sets the size of the finder view
func (m FinderModel) SetSize(width, height int) FinderModel {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height
	m.updateViewportContent()
	return m
}

// updateViewportContent renders all content and sets it to viewport
func (m *FinderModel) updateViewportContent() {
	if m.result == nil {
		return
	}

	var b strings.Builder

	// Header
	query := m.result.Query
	if len(m.result.ResolvedIPs) > 0 {
		query = fmt.Sprintf("%s → %s", m.result.Query, strings.Join(m.result.ResolvedIPs, ", "))
	}
	b.WriteString(m.styles.Title.Render(fmt.Sprintf("%s: %s", i18n.T(i18n.KeyFinderResult), query)))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")).Bold(true).
		Render(fmt.Sprintf(i18n.T(i18n.KeyFinderTotalMatches), m.result.TotalCount())))
	b.WriteString("\n\n")

	// Calculate section width
	sectionWidth := m.width - 4
	if sectionWidth < 80 {
		sectionWidth = 80
	}

	// Render each section
	for i, section := range m.sections {
		isFocused := i == m.currentSection

		// Render section content
		sectionContent := m.renderSection(section, isFocused, sectionWidth-4)

		// Apply border style based on focus
		var bordered string
		if isFocused {
			bordered = m.styles.FocusedBorder.Width(sectionWidth).Render(sectionContent)
		} else {
			bordered = m.styles.Border.Width(sectionWidth).Render(sectionContent)
		}

		b.WriteString(bordered)
		b.WriteString("\n")
	}

	m.viewport.SetContent(b.String())
}

// ensureSelectedVisible adjusts viewport scroll to keep selected row visible
func (m *FinderModel) ensureSelectedVisible() {
	// Calculate approximate line position of current selection
	// Header: 3 lines (title + count + empty line)
	linePos := 3

	for i := 0; i < m.currentSection; i++ {
		section := m.sections[i]
		// Each section has: border top (1) + title (1) + header (1) + separator (1) + rows + border bottom (1) + margin (1)
		rowCount := len(section.Rows)
		if rowCount == 0 {
			rowCount = 1 // Empty message
		}
		linePos += 1 + 1 + 1 + 1 + rowCount + 1 + 1
	}

	// Add current section's header lines (border + title + header + separator)
	linePos += 1 + 1 + 1 + 1
	// Add current row position
	linePos += m.currentRow

	// Get viewport visible range
	viewportTop := m.viewport.YOffset
	viewportBottom := viewportTop + m.viewport.Height - 1

	// Scroll if needed
	if linePos < viewportTop {
		m.viewport.SetYOffset(linePos)
	} else if linePos > viewportBottom {
		m.viewport.SetYOffset(linePos - m.viewport.Height + 1)
	}
}

// Init implements tea.Model
func (m FinderModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m FinderModel) Update(msg tea.Msg) (FinderModel, tea.Cmd) {
	var cmd tea.Cmd
	needsUpdate := false

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Down):
			m = m.moveDown()
			needsUpdate = true
		case key.Matches(msg, m.keys.Up):
			m = m.moveUp()
			needsUpdate = true
		case key.Matches(msg, m.keys.NextSection):
			m = m.nextSection()
			needsUpdate = true
		case key.Matches(msg, m.keys.PrevSection):
			m = m.prevSection()
			needsUpdate = true
		case key.Matches(msg, m.keys.PageUp):
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case key.Matches(msg, m.keys.PageDown):
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		case key.Matches(msg, m.keys.Top):
			m.currentSection = 0
			m.currentRow = 0
			m.viewport.GotoTop()
			needsUpdate = true
		case key.Matches(msg, m.keys.Bottom):
			// Go to last section with rows
			for i := len(m.sections) - 1; i >= 0; i-- {
				if len(m.sections[i].Rows) > 0 {
					m.currentSection = i
					m.currentRow = len(m.sections[i].Rows) - 1
					break
				}
			}
			m.viewport.GotoBottom()
			needsUpdate = true
		case key.Matches(msg, m.keys.Enter):
			return m, m.handleEnter()
		case key.Matches(msg, m.keys.Yank):
			// Handle double-y for yank
			now := time.Now()
			if now.Sub(m.yankLastTime) < 500*time.Millisecond {
				m.yankCount++
			} else {
				m.yankCount = 1
			}
			m.yankLastTime = now

			if m.yankCount >= 2 {
				m.yankCount = 0
				// Copy current row to clipboard
				if m.currentSection < len(m.sections) {
					section := m.sections[m.currentSection]
					if m.currentRow < len(section.Rows) {
						value := strings.Join(section.Rows[m.currentRow], " | ")
						return m, func() tea.Msg {
							return components.CopyDataMsg{Data: value}
						}
					}
				}
			}
		default:
			// Delegate other keys to viewport for scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	default:
		// Delegate other messages to viewport
		m.viewport, cmd = m.viewport.Update(msg)
	}

	if needsUpdate {
		m.updateViewportContent()
		m.ensureSelectedVisible()
	}

	return m, cmd
}

// handleEnter handles the enter key to navigate to detail view
func (m FinderModel) handleEnter() tea.Cmd {
	if m.currentSection >= len(m.sections) {
		return nil
	}
	section := m.sections[m.currentSection]
	if len(section.Rows) == 0 || m.currentRow >= len(section.Data) {
		return nil
	}

	data := section.Data[m.currentRow]
	return func() tea.Msg {
		return types.NavigateMsg{
			Page: section.PageType,
			Data: data,
		}
	}
}

// moveDown moves cursor down
func (m FinderModel) moveDown() FinderModel {
	if len(m.sections) == 0 {
		return m
	}

	// Find current section with rows
	section := m.sections[m.currentSection]
	
	if len(section.Rows) > 0 && m.currentRow < len(section.Rows)-1 {
		m.currentRow++
	} else {
		// Try to move to next section with rows
		for i := m.currentSection + 1; i < len(m.sections); i++ {
			if len(m.sections[i].Rows) > 0 {
				m.currentSection = i
				m.currentRow = 0
				return m
			}
		}
	}

	return m
}

// moveUp moves cursor up
func (m FinderModel) moveUp() FinderModel {
	if len(m.sections) == 0 {
		return m
	}

	if m.currentRow > 0 {
		m.currentRow--
	} else {
		// Try to move to previous section with rows
		for i := m.currentSection - 1; i >= 0; i-- {
			if len(m.sections[i].Rows) > 0 {
				m.currentSection = i
				m.currentRow = len(m.sections[i].Rows) - 1
				return m
			}
		}
	}

	return m
}

// nextSection moves to the next section
func (m FinderModel) nextSection() FinderModel {
	if m.currentSection < len(m.sections)-1 {
		m.currentSection++
		m.currentRow = 0
	}
	return m
}

// prevSection moves to the previous section
func (m FinderModel) prevSection() FinderModel {
	if m.currentSection > 0 {
		m.currentSection--
		m.currentRow = 0
	}
	return m
}

// View implements tea.Model
func (m FinderModel) View() string {
	if m.result == nil {
		return i18n.T(i18n.KeyFinderNoMatch)
	}
	return m.viewport.View()
}

// renderSection renders a single section's table content
func (m FinderModel) renderSection(section FinderSection, isFocused bool, width int) string {
	var b strings.Builder

	// Section title
	if isFocused {
		b.WriteString(m.styles.SectionTitle.Render(section.Title))
	} else {
		b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9CA3AF")).Render(section.Title))
	}
	b.WriteString("\n")

	// Use predefined column widths
	colWidths := section.ColWidths

	// Render header row
	headerCells := make([]string, len(section.Columns))
	for i, col := range section.Columns {
		cell := truncateStr(col, colWidths[i])
		cell = padStr(cell, colWidths[i])
		headerCells[i] = m.styles.Header.Render(cell)
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))
	b.WriteString("\n")

	// Header separator
	totalWidth := 0
	for _, w := range colWidths {
		totalWidth += w + 2 // +2 for padding
	}
	b.WriteString(m.styles.Separator.Render(strings.Repeat("─", totalWidth)))
	b.WriteString("\n")

	// Render data rows or empty message
	if len(section.Rows) == 0 {
		emptyMsg := padStr(i18n.T(i18n.KeyFinderNoMatch), totalWidth-2)
		b.WriteString(m.styles.Empty.Render(emptyMsg))
	} else {
		for j, row := range section.Rows {
			rowCells := make([]string, len(colWidths))
			for k, w := range colWidths {
				cellContent := ""
				if k < len(row) {
					cellContent = row[k]
				}
				
				// Truncate and pad
				displayContent := truncateStr(cellContent, w)
				displayContent = padStr(displayContent, w)

				// Apply style based on selection
				if isFocused && j == m.currentRow {
					rowCells[k] = m.styles.Selected.Render(displayContent)
				} else {
					rowCells[k] = m.styles.Cell.Render(displayContent)
				}
			}
			b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, rowCells...))
			if j < len(section.Rows)-1 {
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

// truncateStr truncates a string to fit within the specified display width
func truncateStr(s string, width int) string {
	displayWidth := runewidth.StringWidth(s)
	if displayWidth <= width {
		return s
	}
	if width <= 3 {
		return runewidth.Truncate(s, width, "")
	}
	return runewidth.Truncate(s, width-2, "") + ".."
}

// padStr pads a string to the specified display width
func padStr(s string, width int) string {
	displayWidth := runewidth.StringWidth(s)
	if displayWidth >= width {
		return s
	}
	return s + strings.Repeat(" ", width-displayWidth)
}

// Search is not applicable to finder
func (m FinderModel) Search(query string) FinderModel {
	return m
}

// NextSearchMatch is not applicable to finder
func (m FinderModel) NextSearchMatch() FinderModel {
	return m
}

// PrevSearchMatch is not applicable to finder
func (m FinderModel) PrevSearchMatch() FinderModel {
	return m
}
