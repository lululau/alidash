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

// Search searches in the list
func (m SLBListModel) Search(query string) SLBListModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m SLBListModel) NextSearchMatch() SLBListModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m SLBListModel) PrevSearchMatch() SLBListModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// SLBListenersModel represents the SLB listeners page
type SLBListenersModel struct {
	table          components.TableModel
	listeners      []service.ListenerDetail
	loadBalancerId string
	width          int
	height         int
	keys           SLBListenersKeyMap
}

// SLBListenersKeyMap defines key bindings for SLB listeners
type SLBListenersKeyMap struct {
	Enter key.Binding
}

// DefaultSLBListenersKeyMap returns default key bindings
func DefaultSLBListenersKeyMap() SLBListenersKeyMap {
	return SLBListenersKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "forwarding rules"),
		),
	}
}

// ListenerNavData contains the data needed to navigate to forwarding rules
type ListenerNavData struct {
	LoadBalancerId string
	ListenerPort   int
	ListenerProtocol string
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
		keys:  DefaultSLBListenersKeyMap(),
	}
}

// SelectedListener returns the selected listener
func (m SLBListenersModel) SelectedListener() *service.ListenerDetail {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.listeners) {
		return &m.listeners[idx]
	}
	return nil
}

// GetLoadBalancerId returns the load balancer ID
func (m SLBListenersModel) GetLoadBalancerId() string {
	return m.loadBalancerId
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if listener := m.SelectedListener(); listener != nil {
				// Only HTTP/HTTPS listeners have forwarding rules
				if listener.Protocol == "HTTP" || listener.Protocol == "HTTPS" {
					return m, func() tea.Msg {
						return types.NavigateMsg{
							Page: types.PageSLBForwardingRules,
							Data: ListenerNavData{
								LoadBalancerId:   m.loadBalancerId,
								ListenerPort:     listener.Port,
								ListenerProtocol: listener.Protocol,
							},
						}
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
func (m SLBListenersModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m SLBListenersModel) Search(query string) SLBListenersModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m SLBListenersModel) NextSearchMatch() SLBListenersModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m SLBListenersModel) PrevSearchMatch() SLBListenersModel {
	m.table = m.table.PrevSearchMatch()
	return m
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
		{Title: "服务器组ID/名称", Width: 45},
		{Title: "关联监听", Width: 15},
		{Title: "关联转发策略", Width: 40},
		{Title: "后端服务器数量", Width: 14},
	}

	return SLBVServerGroupsModel{
		table: components.NewTableModel(columns, "虚拟服务器组"),
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
		// ID / Name combined (use space separator instead of newline)
		idName := vsg.VServerGroupId
		if vsg.VServerGroupName != "" {
			idName = fmt.Sprintf("%s / %s", vsg.VServerGroupId, vsg.VServerGroupName)
		}

		// Associated listeners
		listeners := "-"
		if len(vsg.AssociatedListeners) > 0 {
			listeners = ""
			for j, l := range vsg.AssociatedListeners {
				if j > 0 {
					listeners += ", "
				}
				listeners += l
			}
		}

		// Associated forwarding rules (use comma separator instead of newline)
		rules := "-"
		if len(vsg.AssociatedRules) > 0 {
			rules = ""
			for j, r := range vsg.AssociatedRules {
				if j > 0 {
					rules += ", "
				}
				rules += r
			}
		}

		rows[i] = table.Row{
			idName,
			listeners,
			rules,
			fmt.Sprintf("%d", vsg.BackendServerCount),
		}
		rowData[i] = vsg
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("虚拟服务器组 - SLB: %s", loadBalancerId))
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

// Search searches in the list
func (m SLBVServerGroupsModel) Search(query string) SLBVServerGroupsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m SLBVServerGroupsModel) NextSearchMatch() SLBVServerGroupsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m SLBVServerGroupsModel) PrevSearchMatch() SLBVServerGroupsModel {
	m.table = m.table.PrevSearchMatch()
	return m
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

// Search searches in the list
func (m SLBBackendServersModel) Search(query string) SLBBackendServersModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m SLBBackendServersModel) NextSearchMatch() SLBBackendServersModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m SLBBackendServersModel) PrevSearchMatch() SLBBackendServersModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// SLBForwardingRulesModel represents the SLB forwarding rules page
type SLBForwardingRulesModel struct {
	table            components.TableModel
	rules            []service.ForwardingRuleDetail
	loadBalancerId   string
	listenerPort     int
	listenerProtocol string
	width            int
	height           int
	keys             SLBForwardingRulesKeyMap
}

// SLBForwardingRulesKeyMap defines key bindings for forwarding rules
type SLBForwardingRulesKeyMap struct {
	Enter key.Binding
}

// DefaultSLBForwardingRulesKeyMap returns default key bindings
func DefaultSLBForwardingRulesKeyMap() SLBForwardingRulesKeyMap {
	return SLBForwardingRulesKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
	}
}

// NewSLBForwardingRulesModel creates a new SLB forwarding rules model
func NewSLBForwardingRulesModel() SLBForwardingRulesModel {
	columns := []table.Column{
		{Title: "域名", Width: 35},
		{Title: "URL", Width: 15},
		{Title: "虚拟服务器组", Width: 30},
		{Title: "备注", Width: 20},
	}

	return SLBForwardingRulesModel{
		table: components.NewTableModel(columns, "转发策略列表"),
		keys:  DefaultSLBForwardingRulesKeyMap(),
	}
}

// SetData sets the forwarding rules data
func (m SLBForwardingRulesModel) SetData(rules []service.ForwardingRuleDetail, loadBalancerId string, listenerPort int, listenerProtocol string) SLBForwardingRulesModel {
	m.rules = rules
	m.loadBalancerId = loadBalancerId
	m.listenerPort = listenerPort
	m.listenerProtocol = listenerProtocol

	rows := make([]table.Row, len(rules))
	rowData := make([]interface{}, len(rules))

	for i, rule := range rules {
		domain := rule.Domain
		if domain == "" {
			domain = "-"
		}

		url := rule.Url
		if url == "" {
			url = "/"
		}

		vServerGroup := rule.VServerGroupName
		if vServerGroup == "" {
			vServerGroup = rule.VServerGroupId
		}
		if vServerGroup == "" {
			vServerGroup = "-"
		}

		remark := rule.RuleName
		if remark == "" {
			remark = "-"
		}

		rows[i] = table.Row{
			domain,
			url,
			vServerGroup,
			remark,
		}
		rowData[i] = rule
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("转发策略列表 - %s:%d", loadBalancerId, listenerPort))
	return m
}

// SetSize sets the size
func (m SLBForwardingRulesModel) SetSize(width, height int) SLBForwardingRulesModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SelectedRule returns the selected forwarding rule
func (m SLBForwardingRulesModel) SelectedRule() *service.ForwardingRuleDetail {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.rules) {
		return &m.rules[idx]
	}
	return nil
}

// Init implements tea.Model
func (m SLBForwardingRulesModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SLBForwardingRulesModel) Update(msg tea.Msg) (SLBForwardingRulesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if rule := m.SelectedRule(); rule != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageSLBDetail,
						Data: *rule,
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
func (m SLBForwardingRulesModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m SLBForwardingRulesModel) Search(query string) SLBForwardingRulesModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m SLBForwardingRulesModel) NextSearchMatch() SLBForwardingRulesModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m SLBForwardingRulesModel) PrevSearchMatch() SLBForwardingRulesModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

