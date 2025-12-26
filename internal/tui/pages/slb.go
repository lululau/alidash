package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"

	"aliyun-tui-viewer/internal/service"
	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// SLBListModel represents the SLB instances list page
type SLBListModel struct {
	table         components.TableModel
	loadBalancers []slb.LoadBalancer
	width         int
	height        int
	keys          SLBListKeyMap
}

// SLBListKeyMap defines key bindings
type SLBListKeyMap struct {
	Enter        key.Binding
	Listeners    key.Binding
	VServerGroups key.Binding
}

// DefaultSLBListKeyMap returns default key bindings
func DefaultSLBListKeyMap() SLBListKeyMap {
	return SLBListKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		Listeners: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "listeners"),
		),
		VServerGroups: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "VServer groups"),
		),
	}
}

// NewSLBListModel creates a new SLB list model
func NewSLBListModel() SLBListModel {
	columns := []table.Column{
		{Title: "SLB ID", Width: 25},
		{Title: "Name", Width: 30},
		{Title: "IP Address", Width: 20},
		{Title: "Type", Width: 15},
		{Title: "Status", Width: 12},
	}

	return SLBListModel{
		table: components.NewTableModel(columns, "SLB Instances"),
		keys:  DefaultSLBListKeyMap(),
	}
}

// SetData sets the load balancers data
func (m SLBListModel) SetData(lbs []slb.LoadBalancer) SLBListModel {
	m.loadBalancers = lbs

	rows := make([]table.Row, len(lbs))
	rowData := make([]interface{}, len(lbs))

	for i, lb := range lbs {
		rows[i] = table.Row{
			lb.LoadBalancerId,
			lb.LoadBalancerName,
			lb.Address,
			lb.LoadBalancerSpec,
			lb.LoadBalancerStatus,
		}
		rowData[i] = lb
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m SLBListModel) SetSize(width, height int) SLBListModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SelectedLoadBalancer returns the selected load balancer
func (m SLBListModel) SelectedLoadBalancer() *slb.LoadBalancer {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.loadBalancers) {
		return &m.loadBalancers[idx]
	}
	return nil
}

// Init implements tea.Model
func (m SLBListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SLBListModel) Update(msg tea.Msg) (SLBListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if lb := m.SelectedLoadBalancer(); lb != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageSLBDetail,
						Data: *lb,
					}
				}
			}

		case key.Matches(msg, m.keys.Listeners):
			if lb := m.SelectedLoadBalancer(); lb != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageSLBListeners,
						Data: lb.LoadBalancerId,
					}
				}
			}

		case key.Matches(msg, m.keys.VServerGroups):
			if lb := m.SelectedLoadBalancer(); lb != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageSLBVServerGroups,
						Data: lb.LoadBalancerId,
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
func (m SLBListModel) View() string {
	return m.table.View()
}

// SLBListenersModel represents the SLB listeners page
type SLBListenersModel struct {
	table          components.TableModel
	listeners      []service.ListenerDetail
	loadBalancerId string
	width          int
	height         int
}

// NewSLBListenersModel creates a new SLB listeners model
func NewSLBListenersModel() SLBListenersModel {
	columns := []table.Column{
		{Title: "Protocol", Width: 10},
		{Title: "Port", Width: 10},
		{Title: "Backend Port", Width: 12},
		{Title: "Status", Width: 12},
		{Title: "Health Check", Width: 12},
		{Title: "Scheduler", Width: 12},
		{Title: "VServer Group", Width: 30},
	}

	return SLBListenersModel{
		table: components.NewTableModel(columns, "SLB Listeners"),
	}
}

// SetData sets the listeners data
func (m SLBListenersModel) SetData(listeners []service.ListenerDetail, loadBalancerId string) SLBListenersModel {
	m.listeners = listeners
	m.loadBalancerId = loadBalancerId

	rows := make([]table.Row, len(listeners))
	rowData := make([]interface{}, len(listeners))

	for i, listener := range listeners {
		backendPort := "--"
		if listener.BackendPort > 0 {
			backendPort = fmt.Sprintf("%d", listener.BackendPort)
		}

		vServerGroup := "--"
		if listener.VServerGroupName != "" {
			vServerGroup = listener.VServerGroupName
		} else if listener.VServerGroupId != "" {
			vServerGroup = listener.VServerGroupId
		}

		rows[i] = table.Row{
			listener.Protocol,
			fmt.Sprintf("%d", listener.Port),
			backendPort,
			listener.Status,
			listener.HealthCheck,
			listener.Scheduler,
			vServerGroup,
		}
		rowData[i] = listener
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Listeners for SLB: %s", loadBalancerId))
	return m
}

// SetSize sets the size
func (m SLBListenersModel) SetSize(width, height int) SLBListenersModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m SLBListenersModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SLBListenersModel) Update(msg tea.Msg) (SLBListenersModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m SLBListenersModel) View() string {
	return m.table.View()
}

// SLBVServerGroupsModel represents the SLB VServer groups page
type SLBVServerGroupsModel struct {
	table          components.TableModel
	vServerGroups  []service.VServerGroupDetail
	loadBalancerId string
	width          int
	height         int
	keys           SLBVServerGroupsKeyMap
}

// SLBVServerGroupsKeyMap defines key bindings
type SLBVServerGroupsKeyMap struct {
	Enter key.Binding
}

// DefaultSLBVServerGroupsKeyMap returns default key bindings
func DefaultSLBVServerGroupsKeyMap() SLBVServerGroupsKeyMap {
	return SLBVServerGroupsKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "backend servers"),
		),
	}
}

// NewSLBVServerGroupsModel creates a new SLB VServer groups model
func NewSLBVServerGroupsModel() SLBVServerGroupsModel {
	columns := []table.Column{
		{Title: "VServer Group ID", Width: 30},
		{Title: "Name", Width: 30},
		{Title: "Server Count", Width: 15},
		{Title: "Associated Listeners", Width: 40},
	}

	return SLBVServerGroupsModel{
		table: components.NewTableModel(columns, "SLB VServer Groups"),
		keys:  DefaultSLBVServerGroupsKeyMap(),
	}
}

// SetData sets the VServer groups data
func (m SLBVServerGroupsModel) SetData(groups []service.VServerGroupDetail, loadBalancerId string) SLBVServerGroupsModel {
	m.vServerGroups = groups
	m.loadBalancerId = loadBalancerId

	rows := make([]table.Row, len(groups))
	rowData := make([]interface{}, len(groups))

	for i, vsg := range groups {
		listeners := "--"
		if len(vsg.AssociatedListeners) > 0 {
			listeners = ""
			for j, l := range vsg.AssociatedListeners {
				if j > 0 {
					listeners += ", "
				}
				listeners += l
			}
		}

		rows[i] = table.Row{
			vsg.VServerGroupId,
			vsg.VServerGroupName,
			fmt.Sprintf("%d", vsg.BackendServerCount),
			listeners,
		}
		rowData[i] = vsg
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("VServer Groups for SLB: %s", loadBalancerId))
	return m
}

// SetSize sets the size
func (m SLBVServerGroupsModel) SetSize(width, height int) SLBVServerGroupsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SelectedVServerGroup returns the selected VServer group
func (m SLBVServerGroupsModel) SelectedVServerGroup() *service.VServerGroupDetail {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.vServerGroups) {
		return &m.vServerGroups[idx]
	}
	return nil
}

// Init implements tea.Model
func (m SLBVServerGroupsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SLBVServerGroupsModel) Update(msg tea.Msg) (SLBVServerGroupsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if vsg := m.SelectedVServerGroup(); vsg != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageSLBBackendServers,
						Data: vsg.VServerGroupId,
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
func (m SLBVServerGroupsModel) View() string {
	return m.table.View()
}

// SLBBackendServersModel represents the SLB backend servers page
type SLBBackendServersModel struct {
	table          components.TableModel
	backendServers []service.BackendServerDetail
	vServerGroupId string
	width          int
	height         int
}

// NewSLBBackendServersModel creates a new SLB backend servers model
func NewSLBBackendServersModel() SLBBackendServersModel {
	columns := []table.Column{
		{Title: "Server ID", Width: 25},
		{Title: "ECS Name", Width: 25},
		{Title: "Port", Width: 8},
		{Title: "Weight", Width: 8},
		{Title: "Type", Width: 10},
		{Title: "Private IP", Width: 15},
		{Title: "Public IP", Width: 15},
		{Title: "Description", Width: 20},
	}

	return SLBBackendServersModel{
		table: components.NewTableModel(columns, "Backend Servers"),
	}
}

// SetData sets the backend servers data
func (m SLBBackendServersModel) SetData(servers []service.BackendServerDetail, vServerGroupId string) SLBBackendServersModel {
	m.backendServers = servers
	m.vServerGroupId = vServerGroupId

	rows := make([]table.Row, len(servers))
	rowData := make([]interface{}, len(servers))

	for i, server := range servers {
		rows[i] = table.Row{
			server.ServerId,
			server.InstanceName,
			fmt.Sprintf("%d", server.Port),
			fmt.Sprintf("%d", server.Weight),
			server.Type,
			server.PrivateIpAddress,
			server.PublicIpAddress,
			server.Description,
		}
		rowData[i] = server
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Backend Servers for VServer Group: %s", vServerGroupId))
	return m
}

// SetSize sets the size
func (m SLBBackendServersModel) SetSize(width, height int) SLBBackendServersModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m SLBBackendServersModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SLBBackendServersModel) Update(msg tea.Msg) (SLBBackendServersModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m SLBBackendServersModel) View() string {
	return m.table.View()
}

