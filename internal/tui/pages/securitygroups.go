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

// SecurityGroupsModel represents the security groups list page
type SecurityGroupsModel struct {
	table          components.TableModel
	securityGroups []ecs.SecurityGroup
	title          string
	width          int
	height         int
	keys           SecurityGroupsKeyMap
}

// SecurityGroupsKeyMap defines key bindings
type SecurityGroupsKeyMap struct {
	Enter     key.Binding
	Instances key.Binding
}

// DefaultSecurityGroupsKeyMap returns default key bindings
func DefaultSecurityGroupsKeyMap() SecurityGroupsKeyMap {
	return SecurityGroupsKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "rules"),
		),
		Instances: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "instances"),
		),
	}
}

// NewSecurityGroupsModel creates a new security groups model
func NewSecurityGroupsModel() SecurityGroupsModel {
	columns := []table.Column{
		{Title: "Security Group ID", Width: 25},
		{Title: "Name", Width: 25},
		{Title: "Description", Width: 30},
		{Title: "VPC ID", Width: 25},
		{Title: "Type", Width: 12},
		{Title: "Created", Width: 20},
	}

	return SecurityGroupsModel{
		table: components.NewTableModel(columns, "Security Groups"),
		title: "Security Groups",
		keys:  DefaultSecurityGroupsKeyMap(),
	}
}

// SetData sets the security groups data
func (m SecurityGroupsModel) SetData(groups []ecs.SecurityGroup) SecurityGroupsModel {
	m.securityGroups = groups

	rows := make([]table.Row, len(groups))
	rowData := make([]interface{}, len(groups))

	for i, sg := range groups {
		rows[i] = table.Row{
			sg.SecurityGroupId,
			sg.SecurityGroupName,
			sg.Description,
			sg.VpcId,
			sg.SecurityGroupType,
			sg.CreationTime,
		}
		rowData[i] = sg
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m SecurityGroupsModel) SetSize(width, height int) SecurityGroupsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SetTitle sets the title
func (m SecurityGroupsModel) SetTitle(title string) SecurityGroupsModel {
	m.title = title
	m.table = m.table.SetTitle(title)
	return m
}

// SelectedSecurityGroup returns the selected security group
func (m SecurityGroupsModel) SelectedSecurityGroup() *ecs.SecurityGroup {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.securityGroups) {
		return &m.securityGroups[idx]
	}
	return nil
}

// Init implements tea.Model
func (m SecurityGroupsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SecurityGroupsModel) Update(msg tea.Msg) (SecurityGroupsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if sg := m.SelectedSecurityGroup(); sg != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageSecurityGroupRules,
						Data: sg.SecurityGroupId,
					}
				}
			}

		case key.Matches(msg, m.keys.Instances):
			if sg := m.SelectedSecurityGroup(); sg != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageSecurityGroupInstances,
						Data: sg.SecurityGroupId,
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
func (m SecurityGroupsModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m SecurityGroupsModel) Search(query string) SecurityGroupsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m SecurityGroupsModel) NextSearchMatch() SecurityGroupsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m SecurityGroupsModel) PrevSearchMatch() SecurityGroupsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// SecurityGroupRulesModel represents the security group rules page
type SecurityGroupRulesModel struct {
	table           components.TableModel
	securityGroupId string
	width           int
	height          int
}

// NewSecurityGroupRulesModel creates a new security group rules model
func NewSecurityGroupRulesModel(securityGroupId string) SecurityGroupRulesModel {
	columns := []table.Column{
		{Title: "Direction", Width: 10},
		{Title: "Protocol", Width: 10},
		{Title: "Port Range", Width: 12},
		{Title: "Source/Dest", Width: 25},
		{Title: "Policy", Width: 10},
		{Title: "Priority", Width: 10},
		{Title: "Description", Width: 30},
	}

	title := fmt.Sprintf("Security Group Rules: %s", securityGroupId)

	return SecurityGroupRulesModel{
		table:           components.NewTableModel(columns, title),
		securityGroupId: securityGroupId,
	}
}

// SetData sets the security group rules data
func (m SecurityGroupRulesModel) SetData(response *ecs.DescribeSecurityGroupAttributeResponse) SecurityGroupRulesModel {
	m.securityGroupId = response.SecurityGroupId

	var rows []table.Row
	var rowData []interface{}

	for _, rule := range response.Permissions.Permission {
		// Determine source/dest
		sourceDest := rule.SourceCidrIp
		if sourceDest == "" {
			sourceDest = rule.SourceGroupId
		}

		direction := "Ingress"
		if rule.Direction == "egress" {
			direction = "Egress"
			if rule.DestCidrIp != "" {
				sourceDest = rule.DestCidrIp
			} else if rule.DestGroupId != "" {
				sourceDest = rule.DestGroupId
			}
		}

		rows = append(rows, table.Row{
			direction,
			rule.IpProtocol,
			rule.PortRange,
			sourceDest,
			rule.Policy,
			rule.Priority,
			rule.Description,
		})
		rowData = append(rowData, rule)
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Security Group Rules: %s", m.securityGroupId))
	return m
}

// SetSize sets the size
func (m SecurityGroupRulesModel) SetSize(width, height int) SecurityGroupRulesModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m SecurityGroupRulesModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SecurityGroupRulesModel) Update(msg tea.Msg) (SecurityGroupRulesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m SecurityGroupRulesModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m SecurityGroupRulesModel) Search(query string) SecurityGroupRulesModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m SecurityGroupRulesModel) NextSearchMatch() SecurityGroupRulesModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m SecurityGroupRulesModel) PrevSearchMatch() SecurityGroupRulesModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

