package pages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

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
	currentSection int
	currentRow     int
	scrollOffset   int
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
		result: result,
		keys:   DefaultFinderKeyMap(),
		styles: DefaultFinderStyles(),
	}
	m.buildSections()
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
		Title:     fmt.Sprintf("ECS 实例 (%d)", len(m.result.ECSInstances)),
		Columns:   []string{"实例 ID", "名称", "公网 IP", "私网 IP", "状态"},
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
		Title:     fmt.Sprintf("弹性网卡 (%d)", len(m.result.ENIs)),
		Columns:   []string{"网卡 ID", "私网 IP", "类型", "状态", "绑定实例"},
		ColWidths: []int{24, 16, 10, 10, 24},
		PageType:  types.PageECSJSONDetail,
	}
	for _, eni := range m.result.ENIs {
		eniType := "辅助网卡"
		if eni.Type == "Primary" {
			eniType = "主网卡"
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
		Title:     fmt.Sprintf("负载均衡 (%d)", len(m.result.SLBInstances)),
		Columns:   []string{"SLB ID", "名称", "IP 地址", "类型", "状态"},
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
		Title:     fmt.Sprintf("DNS 记录 (%d)", len(m.result.DNSRecords)),
		Columns:   []string{"域名", "主机记录", "类型", "记录值", "TTL"},
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
		Title:     fmt.Sprintf("RDS 实例 (%d)", len(m.result.RDSInstances)),
		Columns:   []string{"实例 ID", "描述", "引擎", "内网地址", "外网地址", "状态"},
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
		Title:     fmt.Sprintf("Redis 实例 (%d)", len(m.result.RedisInstances)),
		Columns:   []string{"实例 ID", "名称", "连接地址", "私网 IP", "状态"},
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
		Title:     fmt.Sprintf("RocketMQ 实例 (%d)", len(m.result.RocketMQInstances)),
		Columns:   []string{"实例 ID", "名称", "状态"},
		ColWidths: []int{30, 30, 12},
		PageType:  types.PageRocketMQDetail,
	}
	for _, inst := range m.result.RocketMQInstances {
		status := "未知"
		switch inst.InstanceStatus {
		case 0:
			status = "创建中"
		case 2:
			status = "运行中"
		case 5:
			status = "已释放"
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
	return m
}

// Init implements tea.Model
func (m FinderModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m FinderModel) Update(msg tea.Msg) (FinderModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Down):
			m = m.moveDown()
		case key.Matches(msg, m.keys.Up):
			m = m.moveUp()
		case key.Matches(msg, m.keys.NextSection):
			m = m.nextSection()
		case key.Matches(msg, m.keys.PrevSection):
			m = m.prevSection()
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
		}
	}

	return m, nil
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
		return "No results"
	}

	var b strings.Builder

	// Header
	query := m.result.Query
	if len(m.result.ResolvedIPs) > 0 {
		query = fmt.Sprintf("%s → %s", m.result.Query, strings.Join(m.result.ResolvedIPs, ", "))
	}
	b.WriteString(m.styles.Title.Render(fmt.Sprintf("资源查找结果: %s", query)))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")).Bold(true).
		Render(fmt.Sprintf("共找到 %d 个匹配资源", m.result.TotalCount())))
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

	return b.String()
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
		emptyMsg := padStr("暂无匹配数据", totalWidth-2)
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
