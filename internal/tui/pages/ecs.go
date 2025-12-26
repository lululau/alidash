package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// ECSListModel represents the ECS instances list page
type ECSListModel struct {
	table     components.TableModel
	instances []ecs.Instance
	width     int
	height    int
	keys      ECSListKeyMap
}

// ECSListKeyMap defines key bindings for ECS list
type ECSListKeyMap struct {
	Enter          key.Binding
	SecurityGroups key.Binding
}

// DefaultECSListKeyMap returns default key bindings
func DefaultECSListKeyMap() ECSListKeyMap {
	return ECSListKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		SecurityGroups: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "security groups"),
		),
	}
}

// NewECSListModel creates a new ECS list model
func NewECSListModel() ECSListModel {
	columns := []table.Column{
		{Title: "Instance ID", Width: 26},
		{Title: "Status", Width: 10},
		{Title: "Zone", Width: 18},
		{Title: "CPU/RAM", Width: 10},
		{Title: "Private IP", Width: 16},
		{Title: "Public IP", Width: 16},
		{Title: "Name", Width: 30},
		{Title: "Expired", Width: 22},
	}

	return ECSListModel{
		table: components.NewTableModel(columns, "ECS Instances"),
		keys:  DefaultECSListKeyMap(),
	}
}

// SetData sets the ECS instances data
func (m ECSListModel) SetData(instances []ecs.Instance) ECSListModel {
	m.instances = instances

	rows := make([]table.Row, len(instances))
	rowData := make([]interface{}, len(instances))

	for i, inst := range instances {
		// Private IP
		privateIP := "N/A"
		if len(inst.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
			privateIP = inst.VpcAttributes.PrivateIpAddress.IpAddress[0]
		} else if len(inst.InnerIpAddress.IpAddress) > 0 {
			privateIP = inst.InnerIpAddress.IpAddress[0]
		}

		// Public IP
		publicIP := "N/A"
		if len(inst.PublicIpAddress.IpAddress) > 0 {
			publicIP = inst.PublicIpAddress.IpAddress[0]
		} else if inst.EipAddress.IpAddress != "" {
			publicIP = inst.EipAddress.IpAddress
		}

		// CPU/RAM
		cpuRam := fmt.Sprintf("%dC/%dG", inst.Cpu, inst.Memory/1024)

		// Expired Time
		expiredTime := "N/A"
		if inst.ExpiredTime != "" {
			expiredTime = inst.ExpiredTime
		}

		rows[i] = table.Row{
			inst.InstanceId,
			inst.Status,
			inst.ZoneId,
			cpuRam,
			privateIP,
			publicIP,
			inst.InstanceName,
			expiredTime,
		}
		rowData[i] = inst
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the list size
func (m ECSListModel) SetSize(width, height int) ECSListModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SetTitle sets the title
func (m ECSListModel) SetTitle(title string) ECSListModel {
	m.table = m.table.SetTitle(title)
	return m
}

// SelectedInstance returns the selected ECS instance
func (m ECSListModel) SelectedInstance() *ecs.Instance {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.instances) {
		return &m.instances[idx]
	}
	return nil
}

// Init implements tea.Model
func (m ECSListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ECSListModel) Update(msg tea.Msg) (ECSListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageECSDetail,
						Data: *inst,
					}
				}
			}

		case key.Matches(msg, m.keys.SecurityGroups):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageInstanceSecurityGroups,
						Data: inst.InstanceId,
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
func (m ECSListModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m ECSListModel) Search(query string) ECSListModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m ECSListModel) NextSearchMatch() ECSListModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m ECSListModel) PrevSearchMatch() ECSListModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

