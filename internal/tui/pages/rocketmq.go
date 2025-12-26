package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"aliyun-tui-viewer/internal/service"
	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// RocketMQListModel represents the RocketMQ instances list page
type RocketMQListModel struct {
	table     components.TableModel
	instances []service.RocketMQInstance
	width     int
	height    int
	keys      RocketMQListKeyMap
}

// RocketMQListKeyMap defines key bindings
type RocketMQListKeyMap struct {
	Enter  key.Binding
	Topics key.Binding
	Groups key.Binding
}

// DefaultRocketMQListKeyMap returns default key bindings
func DefaultRocketMQListKeyMap() RocketMQListKeyMap {
	return RocketMQListKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		Topics: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "topics"),
		),
		Groups: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "groups"),
		),
	}
}

// NewRocketMQListModel creates a new RocketMQ list model
func NewRocketMQListModel() RocketMQListModel {
	columns := []table.Column{
		{Title: "Instance ID", Width: 40},
		{Title: "Name", Width: 35},
		{Title: "Type", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "Region", Width: 18},
	}

	return RocketMQListModel{
		table: components.NewTableModel(columns, "RocketMQ Instances"),
		keys:  DefaultRocketMQListKeyMap(),
	}
}

// SetData sets the RocketMQ instances data
func (m RocketMQListModel) SetData(instances []service.RocketMQInstance) RocketMQListModel {
	m.instances = instances

	rows := make([]table.Row, len(instances))
	rowData := make([]interface{}, len(instances))

	for i, inst := range instances {
		// Convert int32 values to string
		instanceType := fmt.Sprintf("%d", inst.InstanceType)
		switch inst.InstanceType {
		case 1:
			instanceType = "Standard"
		case 2:
			instanceType = "Professional"
		}

		instanceStatus := fmt.Sprintf("%d", inst.InstanceStatus)
		switch inst.InstanceStatus {
		case 0:
			instanceStatus = "Creating"
		case 5:
			instanceStatus = "Running"
		case 6:
			instanceStatus = "Expired"
		case 7:
			instanceStatus = "Releasing"
		}

		rows[i] = table.Row{
			inst.InstanceId,
			inst.InstanceName,
			instanceType,
			instanceStatus,
			inst.RegionId,
		}
		rowData[i] = inst
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m RocketMQListModel) SetSize(width, height int) RocketMQListModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SelectedInstance returns the selected instance
func (m RocketMQListModel) SelectedInstance() *service.RocketMQInstance {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.instances) {
		return &m.instances[idx]
	}
	return nil
}

// Init implements tea.Model
func (m RocketMQListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m RocketMQListModel) Update(msg tea.Msg) (RocketMQListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageRocketMQDetail,
						Data: *inst,
					}
				}
			}

		case key.Matches(msg, m.keys.Topics):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageRocketMQTopics,
						Data: inst.InstanceId,
					}
				}
			}

		case key.Matches(msg, m.keys.Groups):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageRocketMQGroups,
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
func (m RocketMQListModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m RocketMQListModel) Search(query string) RocketMQListModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m RocketMQListModel) NextSearchMatch() RocketMQListModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m RocketMQListModel) PrevSearchMatch() RocketMQListModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// RocketMQTopicsModel represents the RocketMQ topics page
type RocketMQTopicsModel struct {
	table      components.TableModel
	topics     []service.RocketMQTopic
	instanceId string
	width      int
	height     int
}

// NewRocketMQTopicsModel creates a new RocketMQ topics model
func NewRocketMQTopicsModel() RocketMQTopicsModel {
	columns := []table.Column{
		{Title: "Topic Name", Width: 40},
		{Title: "Message Type", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "Remark", Width: 40},
	}

	return RocketMQTopicsModel{
		table: components.NewTableModel(columns, "RocketMQ Topics"),
	}
}

// SetData sets the topics data
func (m RocketMQTopicsModel) SetData(topics []service.RocketMQTopic, instanceId string) RocketMQTopicsModel {
	m.topics = topics
	m.instanceId = instanceId

	rows := make([]table.Row, len(topics))
	rowData := make([]interface{}, len(topics))

	for i, topic := range topics {
		status := "Active"
		if topic.Status != 0 {
			status = fmt.Sprintf("%d", topic.Status)
		}

		messageType := fmt.Sprintf("%d", topic.MessageType)
		switch topic.MessageType {
		case 0:
			messageType = "Normal"
		case 1:
			messageType = "Partition"
		case 2:
			messageType = "Transaction"
		case 4:
			messageType = "Delay"
		case 5:
			messageType = "Ordered"
		}

		rows[i] = table.Row{
			topic.Topic,
			messageType,
			status,
			topic.Remark,
		}
		rowData[i] = topic
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Topics for RocketMQ: %s", instanceId))
	return m
}

// SetSize sets the size
func (m RocketMQTopicsModel) SetSize(width, height int) RocketMQTopicsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m RocketMQTopicsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m RocketMQTopicsModel) Update(msg tea.Msg) (RocketMQTopicsModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m RocketMQTopicsModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m RocketMQTopicsModel) Search(query string) RocketMQTopicsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m RocketMQTopicsModel) NextSearchMatch() RocketMQTopicsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m RocketMQTopicsModel) PrevSearchMatch() RocketMQTopicsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// RocketMQGroupsModel represents the RocketMQ groups page
type RocketMQGroupsModel struct {
	table      components.TableModel
	groups     []service.RocketMQGroup
	instanceId string
	width      int
	height     int
}

// NewRocketMQGroupsModel creates a new RocketMQ groups model
func NewRocketMQGroupsModel() RocketMQGroupsModel {
	columns := []table.Column{
		{Title: "Group ID", Width: 40},
		{Title: "Type", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "Remark", Width: 40},
	}

	return RocketMQGroupsModel{
		table: components.NewTableModel(columns, "RocketMQ Groups"),
	}
}

// SetData sets the groups data
func (m RocketMQGroupsModel) SetData(groups []service.RocketMQGroup, instanceId string) RocketMQGroupsModel {
	m.groups = groups
	m.instanceId = instanceId

	rows := make([]table.Row, len(groups))
	rowData := make([]interface{}, len(groups))

	for i, group := range groups {
		rows[i] = table.Row{
			group.GroupId,
			group.GroupType,
			"Active",
			group.Remark,
		}
		rowData[i] = group
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Groups for RocketMQ: %s", instanceId))
	return m
}

// SetSize sets the size
func (m RocketMQGroupsModel) SetSize(width, height int) RocketMQGroupsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m RocketMQGroupsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m RocketMQGroupsModel) Update(msg tea.Msg) (RocketMQGroupsModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m RocketMQGroupsModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m RocketMQGroupsModel) Search(query string) RocketMQGroupsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m RocketMQGroupsModel) NextSearchMatch() RocketMQGroupsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m RocketMQGroupsModel) PrevSearchMatch() RocketMQGroupsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

