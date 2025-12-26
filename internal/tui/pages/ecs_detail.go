package pages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"aliyun-tui-viewer/internal/tui/components"
)

// Colors matching the main interface's purple theme
var (
	primaryColor    = lipgloss.Color("#7C3AED") // Purple
	secondaryColor  = lipgloss.Color("#06B6D4") // Cyan
	accentColor     = lipgloss.Color("#F59E0B") // Amber
	successColor    = lipgloss.Color("#10B981") // Green
	errorColor      = lipgloss.Color("#EF4444") // Red
	warningColor    = lipgloss.Color("#F59E0B") // Amber
	textColor       = lipgloss.Color("#E5E7EB") // Light gray
	subtleTextColor = lipgloss.Color("#9CA3AF") // Gray
	mutedTextColor  = lipgloss.Color("#6B7280") // Dark gray
	borderColor     = lipgloss.Color("#374151") // Dark border
	selectedBg      = lipgloss.Color("#374151") // Selected background
)

// DetailRow represents a single row in a section
type DetailRow struct {
	Label string
	Value string
}

// DetailSection represents a section with multiple rows
type DetailSection struct {
	Title string
	Rows  []DetailRow
}

// ECSDetailModel represents the ECS instance detail page with formatted view
type ECSDetailModel struct {
	instance       ecs.Instance
	sections       []DetailSection
	width          int
	height         int
	keys           ECSDetailKeyMap
	currentSection int // Currently focused section
	currentRow     int // Currently focused row within section
	scrollOffset   int // Vertical scroll offset
	yankLastTime   time.Time
	yankCount      int
}

// ECSDetailKeyMap defines key bindings for ECS detail view
type ECSDetailKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	NextSection key.Binding
	PrevSection key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	Top        key.Binding
	Bottom     key.Binding
	Yank       key.Binding
}

// DefaultECSDetailKeyMap returns default key bindings
func DefaultECSDetailKeyMap() ECSDetailKeyMap {
	return ECSDetailKeyMap{
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
		Yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("yy", "copy value"),
		),
	}
}

// NewECSDetailModel creates a new ECS detail model
func NewECSDetailModel(instance ecs.Instance) ECSDetailModel {
	m := ECSDetailModel{
		instance: instance,
		keys:     DefaultECSDetailKeyMap(),
	}
	m.buildSections()
	return m
}

// NewECSDetailModelFromInterface creates a new ECS detail model from interface{}
func NewECSDetailModelFromInterface(data interface{}) (ECSDetailModel, bool) {
	if inst, ok := data.(ecs.Instance); ok {
		return NewECSDetailModel(inst), true
	}
	return ECSDetailModel{}, false
}

// buildSections builds the detail sections from the instance data
func (m *ECSDetailModel) buildSections() {
	inst := m.instance

	// Basic Info Section
	basicInfo := DetailSection{
		Title: "基本信息",
		Rows: []DetailRow{
			{Label: "实例 ID", Value: inst.InstanceId},
			{Label: "实例名称", Value: inst.InstanceName},
			{Label: "实例状态", Value: inst.Status},
			{Label: "所在可用区", Value: inst.ZoneId},
			{Label: "付费类型", Value: m.formatChargeType(inst.InstanceChargeType)},
			{Label: "到期时间", Value: m.formatValue(inst.ExpiredTime)},
		},
	}

	// Config Info Section
	configInfo := DetailSection{
		Title: "配置信息",
		Rows: []DetailRow{
			{Label: "实例规格", Value: inst.InstanceType},
			{Label: "CPU & 内存", Value: fmt.Sprintf("%d 核 (vCPU)  %d GiB", inst.Cpu, inst.Memory/1024)},
			{Label: "公网 IP", Value: m.getPublicIP()},
			{Label: "主私网 IP", Value: m.getPrivateIPs()},
			{Label: "镜像 ID", Value: inst.ImageId},
			{Label: "操作系统", Value: m.formatValue(inst.OSName)},
			{Label: "专有网络", Value: m.formatValue(inst.VpcAttributes.VpcId)},
			{Label: "虚拟交换机", Value: m.formatValue(inst.VpcAttributes.VSwitchId)},
			{Label: "网络类型", Value: m.formatNetworkType(inst.InstanceNetworkType)},
			{Label: "公网带宽", Value: fmt.Sprintf("入: %d Mbps / 出: %d Mbps", inst.InternetMaxBandwidthIn, inst.InternetMaxBandwidthOut)},
			{Label: "带宽计费方式", Value: m.formatValue(inst.InternetChargeType)},
		},
	}

	// Bound Resources Section
	boundResources := DetailSection{
		Title: "绑定资源",
		Rows: []DetailRow{
			{Label: "安全组", Value: fmt.Sprintf("%d 个安全组", len(inst.SecurityGroupIds.SecurityGroupId))},
			{Label: "弹性网卡", Value: fmt.Sprintf("%d 个", len(inst.NetworkInterfaces.NetworkInterface))},
			{Label: "弹性公网 IP ID", Value: m.formatValue(inst.EipAddress.AllocationId)},
			{Label: "辅助私网 IP", Value: m.getSecondaryIPs()},
		},
	}

	// Group Info Section
	groupInfo := DetailSection{
		Title: "分组信息",
		Rows: []DetailRow{
			{Label: "资源组", Value: m.formatValue(inst.ResourceGroupId)},
			{Label: "标签", Value: m.getTags()},
		},
	}

	// Other Info Section
	otherInfo := DetailSection{
		Title: "其他信息",
		Rows: []DetailRow{
			{Label: "主机名", Value: m.formatValue(inst.HostName)},
			{Label: "描述", Value: m.formatValue(inst.Description)},
			{Label: "创建时间", Value: m.formatValue(inst.CreationTime)},
			{Label: "密钥对", Value: m.formatValue(inst.KeyPairName)},
			{Label: "序列号", Value: m.formatValue(inst.SerialNumber)},
		},
	}

	m.sections = []DetailSection{basicInfo, configInfo, boundResources, groupInfo, otherInfo}
}

// SetSize sets the size of the detail view
func (m ECSDetailModel) SetSize(width, height int) ECSDetailModel {
	m.width = width
	m.height = height
	return m
}

// Init implements tea.Model
func (m ECSDetailModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ECSDetailModel) Update(msg tea.Msg) (ECSDetailModel, tea.Cmd) {
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
		case key.Matches(msg, m.keys.Top):
			m.currentSection = 0
			m.currentRow = 0
			m.scrollOffset = 0
		case key.Matches(msg, m.keys.Bottom):
			m.currentSection = len(m.sections) - 1
			m.currentRow = len(m.sections[m.currentSection].Rows) - 1
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
				// Copy current row value to clipboard
				if m.currentSection < len(m.sections) {
					section := m.sections[m.currentSection]
					if m.currentRow < len(section.Rows) {
						value := section.Rows[m.currentRow].Value
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

// moveDown moves cursor down within section or to next section
func (m ECSDetailModel) moveDown() ECSDetailModel {
	if m.currentSection >= len(m.sections) {
		return m
	}

	section := m.sections[m.currentSection]
	if m.currentRow < len(section.Rows)-1 {
		m.currentRow++
	} else if m.currentSection < len(m.sections)-1 {
		// Move to next section
		m.currentSection++
		m.currentRow = 0
	}

	return m
}

// moveUp moves cursor up within section or to previous section
func (m ECSDetailModel) moveUp() ECSDetailModel {
	if m.currentRow > 0 {
		m.currentRow--
	} else if m.currentSection > 0 {
		// Move to previous section
		m.currentSection--
		m.currentRow = len(m.sections[m.currentSection].Rows) - 1
	}

	return m
}

// nextSection moves to the next section
func (m ECSDetailModel) nextSection() ECSDetailModel {
	if m.currentSection < len(m.sections)-1 {
		m.currentSection++
		m.currentRow = 0
	}
	return m
}

// prevSection moves to the previous section
func (m ECSDetailModel) prevSection() ECSDetailModel {
	if m.currentSection > 0 {
		m.currentSection--
		m.currentRow = 0
	}
	return m
}

// View implements tea.Model
func (m ECSDetailModel) View() string {
	if len(m.sections) == 0 {
		return "Loading..."
	}

	var sections []string
	for i, section := range m.sections {
		isFocused := (i == m.currentSection)
		sections = append(sections, m.renderSection(section, i, isFocused))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Create a viewport-like container
	viewportStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height)

	return viewportStyle.Render(content)
}

// renderSection renders a single section
func (m ECSDetailModel) renderSection(section DetailSection, sectionIdx int, isFocused bool) string {
	// Section title style - purple when focused
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(subtleTextColor).
		MarginBottom(1)

	if isFocused {
		titleStyle = titleStyle.Foreground(primaryColor)
	}

	// Section border style - purple when focused
	borderFg := borderColor
	if isFocused {
		borderFg = primaryColor
	}

	// Calculate inner width
	innerWidth := m.width - 8
	if innerWidth < 40 {
		innerWidth = 40
	}

	sectionBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderFg).
		Padding(1, 2).
		MarginBottom(1).
		Width(innerWidth)

	// Render rows
	var rows []string
	for rowIdx, row := range section.Rows {
		isRowSelected := isFocused && (rowIdx == m.currentRow)
		rows = append(rows, m.renderRow(row, isRowSelected))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	titleRendered := titleStyle.Render(section.Title)
	boxContent := sectionBorderStyle.Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, titleRendered, boxContent, "")
}

// renderRow renders a single row
func (m ECSDetailModel) renderRow(row DetailRow, isSelected bool) string {
	// Calculate row width (account for borders and padding)
	rowWidth := m.width - 12
	if rowWidth < 40 {
		rowWidth = 40
	}

	if isSelected {
		// Selected row style: purple background, white text, bold
		selectedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Bold(true)

		labelStyle := selectedStyle.Width(18)
		valueStyle := selectedStyle

		// For status, keep the indicator but use white text
		value := row.Value
		if row.Label == "实例状态" {
			switch value {
			case "Running":
				value = "● " + value
			case "Stopped":
				value = "● " + value
			default:
				value = "● " + value
			}
		}

		rowContent := lipgloss.JoinHorizontal(
			lipgloss.Top,
			labelStyle.Render(row.Label),
			valueStyle.Render(value),
		)

		// Ensure the entire row has purple background
		rowStyle := lipgloss.NewStyle().
			Background(primaryColor).
			Width(rowWidth)

		return rowStyle.Render(rowContent)
	}

	// Normal row style
	labelStyle := lipgloss.NewStyle().
		Foreground(subtleTextColor).
		Width(18)

	valueStyle := lipgloss.NewStyle().
		Foreground(textColor)

	// Style for status values
	value := row.Value
	styledValue := valueStyle.Render(value)

	// Special styling for status
	if row.Label == "实例状态" {
		switch value {
		case "Running":
			styledValue = lipgloss.NewStyle().Foreground(successColor).Bold(true).Render("● " + value)
		case "Stopped":
			styledValue = lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render("● " + value)
		default:
			styledValue = lipgloss.NewStyle().Foreground(warningColor).Bold(true).Render("● " + value)
		}
	}

	rowContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render(row.Label),
		styledValue,
	)

	return rowContent
}

// Helper functions
func (m ECSDetailModel) formatValue(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func (m ECSDetailModel) formatChargeType(chargeType string) string {
	switch chargeType {
	case "PrePaid":
		return "包年包月"
	case "PostPaid":
		return "按量付费"
	default:
		return chargeType
	}
}

func (m ECSDetailModel) formatNetworkType(networkType string) string {
	switch networkType {
	case "vpc":
		return "专有网络"
	case "classic":
		return "经典网络"
	default:
		return networkType
	}
}

func (m ECSDetailModel) getPrivateIPs() string {
	inst := m.instance
	var ips []string

	if len(inst.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		ips = append(ips, inst.VpcAttributes.PrivateIpAddress.IpAddress...)
	}
	if len(inst.InnerIpAddress.IpAddress) > 0 {
		ips = append(ips, inst.InnerIpAddress.IpAddress...)
	}

	if len(ips) == 0 {
		return "-"
	}
	return strings.Join(ips, ", ")
}

func (m ECSDetailModel) getPublicIP() string {
	inst := m.instance

	if len(inst.PublicIpAddress.IpAddress) > 0 {
		return strings.Join(inst.PublicIpAddress.IpAddress, ", ")
	}
	if inst.EipAddress.IpAddress != "" {
		return inst.EipAddress.IpAddress + " (EIP)"
	}
	return "-"
}

func (m ECSDetailModel) getSecondaryIPs() string {
	inst := m.instance
	var ips []string

	for _, ni := range inst.NetworkInterfaces.NetworkInterface {
		for _, ip := range ni.PrivateIpSets.PrivateIpSet {
			if !ip.Primary {
				ips = append(ips, ip.PrivateIpAddress)
			}
		}
	}

	if len(ips) == 0 {
		return "-"
	}
	return strings.Join(ips, ", ")
}

func (m ECSDetailModel) getTags() string {
	inst := m.instance
	if len(inst.Tags.Tag) == 0 {
		return "-"
	}

	var tags []string
	for _, tag := range inst.Tags.Tag {
		tags = append(tags, fmt.Sprintf("%s: %s", tag.TagKey, tag.TagValue))
	}
	return strings.Join(tags, ", ")
}

// Search placeholder for interface compatibility
func (m ECSDetailModel) Search(query string) ECSDetailModel {
	return m
}

// NextSearchMatch placeholder
func (m ECSDetailModel) NextSearchMatch() ECSDetailModel {
	return m
}

// PrevSearchMatch placeholder
func (m ECSDetailModel) PrevSearchMatch() ECSDetailModel {
	return m
}
