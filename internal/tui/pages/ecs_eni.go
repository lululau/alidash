package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// ECSENIModel represents the ECS network interface page
type ECSENIModel struct {
	table      components.TableModel
	enis       []ecs.NetworkInterfaceSet
	instanceId string
	width      int
	height     int
	keys       ECSENIKeyMap
}

// ECSENIKeyMap defines key bindings for ECS ENI list
type ECSENIKeyMap struct {
	Enter key.Binding
}

// DefaultECSENIKeyMap returns default key bindings
func DefaultECSENIKeyMap() ECSENIKeyMap {
	return ECSENIKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
	}
}

// NewECSENIModel creates a new ECS ENI model
func NewECSENIModel(instanceId string) ECSENIModel {
	columns := []table.Column{
		{Title: "弹性网卡 ID", Width: 24},
		{Title: "名称", Width: 16},
		{Title: "网卡类型", Width: 10},
		{Title: "状态", Width: 8},
		{Title: "IP 地址", Width: 16},
		{Title: "专有网络", Width: 26},
		{Title: "可用区", Width: 16},
		{Title: "MAC 地址", Width: 18},
		{Title: "创建时间", Width: 22},
	}

	return ECSENIModel{
		table:      components.NewTableModel(columns, "弹性网卡列表"),
		instanceId: instanceId,
		keys:       DefaultECSENIKeyMap(),
	}
}

// SetData sets the ENI data
func (m ECSENIModel) SetData(enis []ecs.NetworkInterfaceSet) ECSENIModel {
	m.enis = enis

	rows := make([]table.Row, len(enis))
	rowData := make([]interface{}, len(enis))

	for i, eni := range enis {
		// ENI name
		eniName := eni.NetworkInterfaceName
		if eniName == "" {
			eniName = "-"
		}

		// NIC type
		nicType := m.formatNICType(eni.Type)

		// Status
		status := m.formatStatus(eni.Status)

		// Private IP
		privateIP := eni.PrivateIpAddress
		if privateIP == "" && len(eni.PrivateIpSets.PrivateIpSet) > 0 {
			privateIP = eni.PrivateIpSets.PrivateIpSet[0].PrivateIpAddress
		}
		if privateIP == "" {
			privateIP = "-"
		}

		// VPC / VSwitch
		vpcInfo := eni.VpcId
		if eni.VSwitchId != "" {
			vpcInfo = eni.VSwitchId
		}
		if vpcInfo == "" {
			vpcInfo = "-"
		}

		// Zone
		zone := eni.ZoneId
		if zone == "" {
			zone = "-"
		}

		// MAC Address
		macAddr := eni.MacAddress
		if macAddr == "" {
			macAddr = "-"
		}

		// Creation time
		creationTime := eni.CreationTime
		if creationTime == "" {
			creationTime = "-"
		}

		rows[i] = table.Row{
			eni.NetworkInterfaceId,
			eniName,
			nicType,
			status,
			privateIP,
			vpcInfo,
			zone,
			macAddr,
			creationTime,
		}
		rowData[i] = eni
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m ECSENIModel) SetSize(width, height int) ECSENIModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SetTitle sets the title
func (m ECSENIModel) SetTitle(title string) ECSENIModel {
	m.table = m.table.SetTitle(title)
	return m
}

// SelectedENI returns the selected ENI
func (m ECSENIModel) SelectedENI() *ecs.NetworkInterfaceSet {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.enis) {
		return &m.enis[idx]
	}
	return nil
}

// Init implements tea.Model
func (m ECSENIModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ECSENIModel) Update(msg tea.Msg) (ECSENIModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if eni := m.SelectedENI(); eni != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageECSJSONDetail,
						Data: *eni,
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
func (m ECSENIModel) View() string {
	return m.table.View()
}

// Helper functions
func (m ECSENIModel) formatNICType(nicType string) string {
	switch strings.ToLower(nicType) {
	case "primary":
		return "主网卡"
	case "secondary":
		return "辅助网卡"
	default:
		return nicType
	}
}

func (m ECSENIModel) formatStatus(status string) string {
	switch strings.ToLower(status) {
	case "available":
		return "可用"
	case "inuse":
		return "已绑定"
	case "attaching":
		return "绑定中"
	case "detaching":
		return "解绑中"
	case "creating":
		return "创建中"
	case "deleting":
		return "删除中"
	default:
		return status
	}
}

// GetSecurityGroups returns formatted security groups
func (m ECSENIModel) GetSecurityGroups(eni ecs.NetworkInterfaceSet) string {
	if len(eni.SecurityGroupIds.SecurityGroupId) == 0 {
		return "-"
	}
	return fmt.Sprintf("%d 个安全组", len(eni.SecurityGroupIds.SecurityGroupId))
}

// Search searches in the list
func (m ECSENIModel) Search(query string) ECSENIModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m ECSENIModel) NextSearchMatch() ECSENIModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m ECSENIModel) PrevSearchMatch() ECSENIModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

