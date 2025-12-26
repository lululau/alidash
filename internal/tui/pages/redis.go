package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// RedisListModel represents the Redis instances list page
type RedisListModel struct {
	table     components.TableModel
	instances []r_kvstore.KVStoreInstance
	width     int
	height    int
	keys      RedisListKeyMap
}

// RedisListKeyMap defines key bindings
type RedisListKeyMap struct {
	Enter    key.Binding
	Accounts key.Binding
}

// DefaultRedisListKeyMap returns default key bindings
func DefaultRedisListKeyMap() RedisListKeyMap {
	return RedisListKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		Accounts: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "accounts"),
		),
	}
}

// NewRedisListModel creates a new Redis list model
func NewRedisListModel() RedisListModel {
	columns := []table.Column{
		{Title: "Instance ID", Width: 25},
		{Title: "Name", Width: 25},
		{Title: "Type", Width: 12},
		{Title: "Capacity (MB)", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "Connection", Width: 35},
	}

	return RedisListModel{
		table: components.NewTableModel(columns, "Redis Instances"),
		keys:  DefaultRedisListKeyMap(),
	}
}

// SetData sets the Redis instances data
func (m RedisListModel) SetData(instances []r_kvstore.KVStoreInstance) RedisListModel {
	m.instances = instances

	rows := make([]table.Row, len(instances))
	rowData := make([]interface{}, len(instances))

	for i, inst := range instances {
		connection := "N/A"
		if inst.ConnectionDomain != "" {
			connection = inst.ConnectionDomain
		}

		rows[i] = table.Row{
			inst.InstanceId,
			inst.InstanceName,
			inst.InstanceType,
			fmt.Sprintf("%d", inst.Capacity),
			inst.InstanceStatus,
			connection,
		}
		rowData[i] = inst
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m RedisListModel) SetSize(width, height int) RedisListModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SelectedInstance returns the selected instance
func (m RedisListModel) SelectedInstance() *r_kvstore.KVStoreInstance {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.instances) {
		return &m.instances[idx]
	}
	return nil
}

// Init implements tea.Model
func (m RedisListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m RedisListModel) Update(msg tea.Msg) (RedisListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageRedisDetail,
						Data: *inst,
					}
				}
			}

		case key.Matches(msg, m.keys.Accounts):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageRedisAccounts,
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
func (m RedisListModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m RedisListModel) Search(query string) RedisListModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m RedisListModel) NextSearchMatch() RedisListModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m RedisListModel) PrevSearchMatch() RedisListModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// RedisAccountsModel represents the Redis accounts page
type RedisAccountsModel struct {
	table      components.TableModel
	accounts   []r_kvstore.Account
	instanceId string
	width      int
	height     int
}

// NewRedisAccountsModel creates a new Redis accounts model
func NewRedisAccountsModel() RedisAccountsModel {
	columns := []table.Column{
		{Title: "Account Name", Width: 25},
		{Title: "Type", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "Privileges", Width: 40},
		{Title: "Description", Width: 30},
	}

	return RedisAccountsModel{
		table: components.NewTableModel(columns, "Redis Accounts"),
	}
}

// SetData sets the accounts data
func (m RedisAccountsModel) SetData(accounts []r_kvstore.Account, instanceId string) RedisAccountsModel {
	m.accounts = accounts
	m.instanceId = instanceId

	rows := make([]table.Row, len(accounts))
	rowData := make([]interface{}, len(accounts))

	for i, account := range accounts {
		rows[i] = table.Row{
			account.AccountName,
			account.AccountType,
			account.AccountStatus,
			"--", // Privileges not directly available
			account.AccountDescription,
		}
		rowData[i] = account
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Accounts for Redis: %s", instanceId))
	return m
}

// SetSize sets the size
func (m RedisAccountsModel) SetSize(width, height int) RedisAccountsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m RedisAccountsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m RedisAccountsModel) Update(msg tea.Msg) (RedisAccountsModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m RedisAccountsModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m RedisAccountsModel) Search(query string) RedisAccountsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m RedisAccountsModel) NextSearchMatch() RedisAccountsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m RedisAccountsModel) PrevSearchMatch() RedisAccountsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

