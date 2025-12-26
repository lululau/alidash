package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

	"aliyun-tui-viewer/internal/service"
	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// RDSListModel represents the RDS instances list page
type RDSListModel struct {
	table            components.TableModel
	instances        []rds.DBInstance
	detailedInstances []service.RDSInstanceDetail
	width            int
	height           int
	keys             RDSListKeyMap
}

// RDSListKeyMap defines key bindings
type RDSListKeyMap struct {
	Enter     key.Binding
	Databases key.Binding
	Accounts  key.Binding
}

// DefaultRDSListKeyMap returns default key bindings
func DefaultRDSListKeyMap() RDSListKeyMap {
	return RDSListKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		Databases: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "databases"),
		),
		Accounts: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "accounts"),
		),
	}
}

// NewRDSListModel creates a new RDS list model
func NewRDSListModel() RDSListModel {
	columns := []table.Column{
		{Title: "Instance ID", Width: 25},
		{Title: "Engine", Width: 12},
		{Title: "Version", Width: 8},
		{Title: "Class", Width: 18},
		{Title: "Internal Addr", Width: 35},
		{Title: "Public Addr", Width: 35},
		{Title: "Status", Width: 10},
		{Title: "Description", Width: 20},
	}

	return RDSListModel{
		table: components.NewTableModel(columns, "RDS Instances"),
		keys:  DefaultRDSListKeyMap(),
	}
}

// SetData sets the RDS instances data (basic, without network info)
func (m RDSListModel) SetData(instances []rds.DBInstance) RDSListModel {
	m.instances = instances
	m.detailedInstances = nil

	rows := make([]table.Row, len(instances))
	rowData := make([]interface{}, len(instances))

	for i, inst := range instances {
		rows[i] = table.Row{
			inst.DBInstanceId,
			inst.Engine,
			inst.EngineVersion,
			inst.DBInstanceClass,
			inst.ConnectionString, // Internal address from basic info
			"-",                   // Public address not available without detailed fetch
			inst.DBInstanceStatus,
			inst.DBInstanceDescription,
		}
		rowData[i] = inst
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetDetailedData sets the RDS instances data with network info
func (m RDSListModel) SetDetailedData(detailedInstances []service.RDSInstanceDetail) RDSListModel {
	m.detailedInstances = detailedInstances
	m.instances = make([]rds.DBInstance, len(detailedInstances))

	rows := make([]table.Row, len(detailedInstances))
	rowData := make([]interface{}, len(detailedInstances))

	for i, detail := range detailedInstances {
		m.instances[i] = detail.Instance

		internalAddr := detail.InternalConnectionStr
		if internalAddr == "" {
			internalAddr = "-"
		}
		publicAddr := detail.PublicConnectionStr
		if publicAddr == "" {
			publicAddr = "-"
		}

		rows[i] = table.Row{
			detail.Instance.DBInstanceId,
			detail.Instance.Engine,
			detail.Instance.EngineVersion,
			detail.Instance.DBInstanceClass,
			internalAddr,
			publicAddr,
			detail.Instance.DBInstanceStatus,
			detail.Instance.DBInstanceDescription,
		}
		rowData[i] = detail
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m RDSListModel) SetSize(width, height int) RDSListModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SelectedInstance returns the selected instance
func (m RDSListModel) SelectedInstance() *rds.DBInstance {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.instances) {
		return &m.instances[idx]
	}
	return nil
}

// Init implements tea.Model
func (m RDSListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m RDSListModel) Update(msg tea.Msg) (RDSListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageRDSDetail,
						Data: *inst,
					}
				}
			}

		case key.Matches(msg, m.keys.Databases):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageRDSDatabases,
						Data: inst.DBInstanceId,
					}
				}
			}

		case key.Matches(msg, m.keys.Accounts):
			if inst := m.SelectedInstance(); inst != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageRDSAccounts,
						Data: inst.DBInstanceId,
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
func (m RDSListModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m RDSListModel) Search(query string) RDSListModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m RDSListModel) NextSearchMatch() RDSListModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m RDSListModel) PrevSearchMatch() RDSListModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// RDSDatabasesModel represents the RDS databases page
type RDSDatabasesModel struct {
	table      components.TableModel
	databases  []rds.Database
	instanceId string
	width      int
	height     int
}

// NewRDSDatabasesModel creates a new RDS databases model
func NewRDSDatabasesModel() RDSDatabasesModel {
	columns := []table.Column{
		{Title: "Database Name", Width: 25},
		{Title: "Status", Width: 12},
		{Title: "Character Set", Width: 15},
		{Title: "Bound Accounts", Width: 30},
		{Title: "Description", Width: 30},
	}

	return RDSDatabasesModel{
		table: components.NewTableModel(columns, "RDS Databases"),
	}
}

// SetData sets the databases data
func (m RDSDatabasesModel) SetData(databases []rds.Database, instanceId string) RDSDatabasesModel {
	m.databases = databases
	m.instanceId = instanceId

	rows := make([]table.Row, len(databases))
	rowData := make([]interface{}, len(databases))

	for i, db := range databases {
		// Format bound accounts
		boundAccounts := "--"
		if len(db.Accounts.AccountPrivilegeInfo) > 0 {
			boundAccounts = ""
			for j, account := range db.Accounts.AccountPrivilegeInfo {
				if j > 0 {
					boundAccounts += ", "
				}
				boundAccounts += account.Account
			}
		}

		rows[i] = table.Row{
			db.DBName,
			db.DBStatus,
			db.CharacterSetName,
			boundAccounts,
			db.DBDescription,
		}
		rowData[i] = db
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Databases for RDS: %s", instanceId))
	return m
}

// SetSize sets the size
func (m RDSDatabasesModel) SetSize(width, height int) RDSDatabasesModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m RDSDatabasesModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m RDSDatabasesModel) Update(msg tea.Msg) (RDSDatabasesModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m RDSDatabasesModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m RDSDatabasesModel) Search(query string) RDSDatabasesModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m RDSDatabasesModel) NextSearchMatch() RDSDatabasesModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m RDSDatabasesModel) PrevSearchMatch() RDSDatabasesModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// RDSAccountsModel represents the RDS accounts page
type RDSAccountsModel struct {
	table      components.TableModel
	accounts   []rds.DBInstanceAccount
	instanceId string
	width      int
	height     int
}

// NewRDSAccountsModel creates a new RDS accounts model
func NewRDSAccountsModel() RDSAccountsModel {
	columns := []table.Column{
		{Title: "Account Name", Width: 25},
		{Title: "Type", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "Bound Databases", Width: 40},
		{Title: "Description", Width: 25},
	}

	return RDSAccountsModel{
		table: components.NewTableModel(columns, "RDS Accounts"),
	}
}

// SetData sets the accounts data
func (m RDSAccountsModel) SetData(accounts []rds.DBInstanceAccount, instanceId string) RDSAccountsModel {
	m.accounts = accounts
	m.instanceId = instanceId

	rows := make([]table.Row, len(accounts))
	rowData := make([]interface{}, len(accounts))

	for i, account := range accounts {
		// Format bound databases
		boundDatabases := "--"
		if len(account.DatabasePrivileges.DatabasePrivilege) > 0 {
			boundDatabases = ""
			for j, dbPriv := range account.DatabasePrivileges.DatabasePrivilege {
				if j > 0 {
					boundDatabases += ", "
				}
				boundDatabases += fmt.Sprintf("%s(%s)", dbPriv.DBName, dbPriv.AccountPrivilege)
			}
		}

		rows[i] = table.Row{
			account.AccountName,
			account.AccountType,
			account.AccountStatus,
			boundDatabases,
			account.AccountDescription,
		}
		rowData[i] = account
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Accounts for RDS: %s", instanceId))
	return m
}

// SetSize sets the size
func (m RDSAccountsModel) SetSize(width, height int) RDSAccountsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m RDSAccountsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m RDSAccountsModel) Update(msg tea.Msg) (RDSAccountsModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m RDSAccountsModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m RDSAccountsModel) Search(query string) RDSAccountsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m RDSAccountsModel) NextSearchMatch() RDSAccountsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m RDSAccountsModel) PrevSearchMatch() RDSAccountsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

