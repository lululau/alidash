package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"aliyun-tui-viewer/internal/i18n"
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
		{Title: i18n.T(i18n.KeyColENIID), Width: 24},
		{Title: i18n.T(i18n.KeyColName), Width: 16},
		{Title: i18n.T(i18n.KeyColENIType), Width: 10},
		{Title: i18n.T(i18n.KeyColStatus), Width: 8},
		{Title: i18n.T(i18n.KeyColPrivateIP), Width: 16},
		{Title: i18n.T(i18n.KeyLabelVPC), Width: 26},
		{Title: i18n.T(i18n.KeyColZone), Width: 16},
		{Title: i18n.T(i18n.KeyColMAC), Width: 18},
		{Title: i18n.T(i18n.KeyColCreatedAt), Width: 22},
	}

	return ECSENIModel{
		table:      components.NewTableModel(columns, i18n.T(i18n.KeyPageECSENIs)),
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
		return i18n.T(i18n.KeyENIPrimary)
	case "secondary":
		return i18n.T(i18n.KeyENISecondary)
	default:
		return nicType
	}
}

func (m ECSENIModel) formatStatus(status string) string {
	switch strings.ToLower(status) {
	case "available":
		return i18n.T(i18n.KeyENIAvailable)
	case "inuse":
		return i18n.T(i18n.KeyENIInUse)
	case "attaching":
		return i18n.T(i18n.KeyENIAttaching)
	case "detaching":
		return i18n.T(i18n.KeyENIDetaching)
	case "creating":
		return i18n.T(i18n.KeyStatusCreating)
	case "deleting":
		return i18n.T(i18n.KeyENIDeleting)
	default:
		return status
	}
}

// GetSecurityGroups returns formatted security groups
func (m ECSENIModel) GetSecurityGroups(eni ecs.NetworkInterfaceSet) string {
	if len(eni.SecurityGroupIds.SecurityGroupId) == 0 {
		return "-"
	}
	return fmt.Sprintf(i18n.T(i18n.KeyCountSG), len(eni.SecurityGroupIds.SecurityGroupId))
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

