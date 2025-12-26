package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"aliyun-tui-viewer/internal/i18n"
	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// Colors for the disk page
var (
	diskPrimaryColor    = lipgloss.Color("#7C3AED") // Purple
	diskSecondaryColor  = lipgloss.Color("#06B6D4") // Cyan
	diskSuccessColor    = lipgloss.Color("#10B981") // Green
	diskWarningColor    = lipgloss.Color("#F59E0B") // Amber
	diskTextColor       = lipgloss.Color("#E5E7EB") // Light gray
	diskSubtleTextColor = lipgloss.Color("#9CA3AF") // Gray
	diskBorderColor     = lipgloss.Color("#374151") // Dark border
)

// ECSDiskModel represents the ECS disk/storage page
type ECSDiskModel struct {
	table      components.TableModel
	disks      []ecs.Disk
	instanceId string
	width      int
	height     int
	keys       ECSDiskKeyMap
}

// ECSDiskKeyMap defines key bindings for ECS disk list
type ECSDiskKeyMap struct {
	Enter key.Binding
}

// DefaultECSDiskKeyMap returns default key bindings
func DefaultECSDiskKeyMap() ECSDiskKeyMap {
	return ECSDiskKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
	}
}

// NewECSDiskModel creates a new ECS disk model
func NewECSDiskModel(instanceId string) ECSDiskModel {
	columns := []table.Column{
		{Title: i18n.T(i18n.KeyColDiskID), Width: 24},
		{Title: i18n.T(i18n.KeyColName), Width: 20},
		{Title: i18n.T(i18n.KeyColDiskAttribute), Width: 8},
		{Title: i18n.T(i18n.KeyColStatus), Width: 8},
		{Title: i18n.T(i18n.KeyColType), Width: 22},
		{Title: i18n.T(i18n.KeyColSize), Width: 12},
		{Title: i18n.T(i18n.KeyColDiskIOPS), Width: 10},
		{Title: i18n.T(i18n.KeyColDiskDeleteBehavior), Width: 14},
		{Title: i18n.T(i18n.KeyLabelChargeType), Width: 12},
		{Title: i18n.T(i18n.KeyColDiskPortable), Width: 8},
	}

	return ECSDiskModel{
		table:      components.NewTableModel(columns, i18n.T(i18n.KeyPageECSDisks)),
		instanceId: instanceId,
		keys:       DefaultECSDiskKeyMap(),
	}
}

// SetData sets the disk data
func (m ECSDiskModel) SetData(disks []ecs.Disk) ECSDiskModel {
	m.disks = disks

	rows := make([]table.Row, len(disks))
	rowData := make([]interface{}, len(disks))

	for i, disk := range disks {
		// Format disk type
		diskType := m.formatDiskCategory(disk.Category)
		if disk.PerformanceLevel != "" {
			diskType += " " + disk.PerformanceLevel
		}

		// Format size with IOPS
		sizeStr := fmt.Sprintf("%d GiB", disk.Size)

		// Format IOPS
		iopsStr := "-"
		if disk.IOPS > 0 {
			iopsStr = fmt.Sprintf("%d", disk.IOPS)
		}

		// Disk name
		diskName := disk.DiskName
		if diskName == "" {
			diskName = "-"
		}

		// Disk property (system/data disk)
		diskProp := m.formatDiskType(disk.Type)

		// Delete with instance
		deleteWithInst := i18n.T(i18n.KeyDiskDeleteWithInst)
		if !disk.DeleteWithInstance {
			deleteWithInst = i18n.T(i18n.KeyDiskKeepAfterInst)
		}

		// Charge type
		chargeType := m.formatChargeType(disk.DiskChargeType)

		// Portable
		portable := i18n.T(i18n.KeyDiskPortable)
		if !disk.Portable {
			portable = i18n.T(i18n.KeyDiskNotPortable)
		}

		rows[i] = table.Row{
			disk.DiskId,
			diskName,
			diskProp,
			disk.Status,
			diskType,
			sizeStr,
			iopsStr,
			deleteWithInst,
			chargeType,
			portable,
		}
		rowData[i] = disk
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m ECSDiskModel) SetSize(width, height int) ECSDiskModel {
	m.width = width
	m.height = height
	// Reserve space for overview section (about 8 lines)
	tableHeight := height - 10
	if tableHeight < 5 {
		tableHeight = 5
	}
	m.table = m.table.SetSize(width, tableHeight)
	return m
}

// SetTitle sets the title
func (m ECSDiskModel) SetTitle(title string) ECSDiskModel {
	m.table = m.table.SetTitle(title)
	return m
}

// SelectedDisk returns the selected disk
func (m ECSDiskModel) SelectedDisk() *ecs.Disk {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.disks) {
		return &m.disks[idx]
	}
	return nil
}

// Init implements tea.Model
func (m ECSDiskModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ECSDiskModel) Update(msg tea.Msg) (ECSDiskModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if disk := m.SelectedDisk(); disk != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageECSJSONDetail,
						Data: *disk,
					}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m ECSDiskModel) View() string {
	// Build overview section
	overview := m.renderOverviewSection()

	// Build table
	tableView := m.table.View()

	// Combine overview and table
	return lipgloss.JoinVertical(lipgloss.Left, overview, "", tableView)
}

// renderOverviewSection renders the usage overview section
func (m ECSDiskModel) renderOverviewSection() string {
	// Title style
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(diskPrimaryColor).
		MarginBottom(1)

	// Label style
	labelStyle := lipgloss.NewStyle().
		Foreground(diskSubtleTextColor)

	// Value style
	valueStyle := lipgloss.NewStyle().
		Foreground(diskSecondaryColor).
		Bold(true)

	// Calculate total and system/data disk counts
	totalDisks := len(m.disks)
	var systemDisks, dataDisks int
	var totalSize int
	for _, disk := range m.disks {
		totalSize += disk.Size
		if disk.Type == "system" {
			systemDisks++
		} else {
			dataDisks++
		}
	}

	// Build left section - Usage Overview
	leftTitle := titleStyle.Render(i18n.T(i18n.KeySectionUsageOverview))
	leftContent := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(i18n.T(i18n.KeyLabelTotalDisks)),
		valueStyle.Render(fmt.Sprintf(i18n.T(i18n.KeyCountDisks), totalDisks)),
	)

	// Build right section - Storage Overview
	rightTitle := titleStyle.Render(i18n.T(i18n.KeySectionStorageOverview))
	rightContent := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(i18n.T(i18n.KeyLabelTotalCapacity)),
		valueStyle.Render(fmt.Sprintf("%d GiB", totalSize)),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render(i18n.T(i18n.KeyDiskSystem)+": "),
			valueStyle.Render(fmt.Sprintf("%d", systemDisks)),
			labelStyle.Render("  "+i18n.T(i18n.KeyDiskData)+": "),
			valueStyle.Render(fmt.Sprintf("%d", dataDisks)),
		),
	)

	// Section border style
	innerWidth := (m.width - 12) / 2
	if innerWidth < 30 {
		innerWidth = 30
	}

	sectionStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(diskBorderColor).
		Padding(1, 2).
		Width(innerWidth)

	leftSection := lipgloss.JoinVertical(lipgloss.Left, leftTitle, sectionStyle.Render(leftContent))
	rightSection := lipgloss.JoinVertical(lipgloss.Left, rightTitle, sectionStyle.Render(rightContent))

	// Join sections horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, leftSection, "  ", rightSection)
}

// Helper functions
func (m ECSDiskModel) formatDiskCategory(category string) string {
	switch category {
	case "cloud":
		return i18n.T(i18n.KeyDiskCloud)
	case "cloud_efficiency":
		return i18n.T(i18n.KeyDiskCloudEfficiency)
	case "cloud_ssd":
		return i18n.T(i18n.KeyDiskCloudSSD)
	case "cloud_essd":
		return i18n.T(i18n.KeyDiskCloudEssd)
	case "cloud_auto":
		return i18n.T(i18n.KeyDiskCloudEssdAuto)
	case "cloud_essd_entry":
		return i18n.T(i18n.KeyDiskCloudEssdEntry)
	case "ephemeral_ssd":
		return "Local SSD"
	default:
		return category
	}
}

func (m ECSDiskModel) formatDiskType(diskType string) string {
	switch diskType {
	case "system":
		return i18n.T(i18n.KeyDiskSystem)
	case "data":
		return i18n.T(i18n.KeyDiskData)
	default:
		return diskType
	}
}

func (m ECSDiskModel) formatChargeType(chargeType string) string {
	switch chargeType {
	case "PrePaid":
		return i18n.T(i18n.KeyChargePrePaid)
	case "PostPaid":
		return i18n.T(i18n.KeyChargePostPaid)
	default:
		return chargeType
	}
}

// Search searches in the list
func (m ECSDiskModel) Search(query string) ECSDiskModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m ECSDiskModel) NextSearchMatch() ECSDiskModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m ECSDiskModel) PrevSearchMatch() ECSDiskModel {
	m.table = m.table.PrevSearchMatch()
	return m
}
