package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"

	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// DNSDomainsModel represents the DNS domains list page
type DNSDomainsModel struct {
	table   components.TableModel
	domains []alidns.DomainInDescribeDomains
	width   int
	height  int
	keys    DNSDomainsKeyMap
}

// DNSDomainsKeyMap defines key bindings
type DNSDomainsKeyMap struct {
	Enter key.Binding
}

// DefaultDNSDomainsKeyMap returns default key bindings
func DefaultDNSDomainsKeyMap() DNSDomainsKeyMap {
	return DNSDomainsKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "records"),
		),
	}
}

// NewDNSDomainsModel creates a new DNS domains model
func NewDNSDomainsModel() DNSDomainsModel {
	columns := []table.Column{
		{Title: "Domain Name", Width: 40},
		{Title: "Record Count", Width: 15},
		{Title: "Version", Width: 15},
	}

	return DNSDomainsModel{
		table: components.NewTableModel(columns, "DNS Domains"),
		keys:  DefaultDNSDomainsKeyMap(),
	}
}

// SetData sets the domains data
func (m DNSDomainsModel) SetData(domains []alidns.DomainInDescribeDomains) DNSDomainsModel {
	m.domains = domains

	rows := make([]table.Row, len(domains))
	rowData := make([]interface{}, len(domains))

	for i, domain := range domains {
		rows[i] = table.Row{
			domain.DomainName,
			fmt.Sprintf("%d", domain.RecordCount),
			domain.VersionCode,
		}
		rowData[i] = domain
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m DNSDomainsModel) SetSize(width, height int) DNSDomainsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SelectedDomain returns the selected domain
func (m DNSDomainsModel) SelectedDomain() *alidns.DomainInDescribeDomains {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.domains) {
		return &m.domains[idx]
	}
	return nil
}

// Init implements tea.Model
func (m DNSDomainsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m DNSDomainsModel) Update(msg tea.Msg) (DNSDomainsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if domain := m.SelectedDomain(); domain != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageDNSRecords,
						Data: domain.DomainName,
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
func (m DNSDomainsModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m DNSDomainsModel) Search(query string) DNSDomainsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m DNSDomainsModel) NextSearchMatch() DNSDomainsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m DNSDomainsModel) PrevSearchMatch() DNSDomainsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// DNSRecordsModel represents the DNS records list page
type DNSRecordsModel struct {
	table      components.TableModel
	records    []alidns.Record
	domainName string
	width      int
	height     int
}

// NewDNSRecordsModel creates a new DNS records model
func NewDNSRecordsModel() DNSRecordsModel {
	columns := []table.Column{
		{Title: "Record ID", Width: 25},
		{Title: "RR", Width: 20},
		{Title: "Type", Width: 10},
		{Title: "Value", Width: 40},
		{Title: "TTL", Width: 10},
		{Title: "Status", Width: 10},
	}

	return DNSRecordsModel{
		table: components.NewTableModel(columns, "DNS Records"),
	}
}

// SetData sets the records data
func (m DNSRecordsModel) SetData(records []alidns.Record, domainName string) DNSRecordsModel {
	m.records = records
	m.domainName = domainName

	rows := make([]table.Row, len(records))
	rowData := make([]interface{}, len(records))

	for i, record := range records {
		rows[i] = table.Row{
			record.RecordId,
			record.RR,
			record.Type,
			record.Value,
			fmt.Sprintf("%d", record.TTL),
			record.Status,
		}
		rowData[i] = record
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("DNS Records for %s", domainName))
	return m
}

// SetSize sets the size
func (m DNSRecordsModel) SetSize(width, height int) DNSRecordsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// Init implements tea.Model
func (m DNSRecordsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m DNSRecordsModel) Update(msg tea.Msg) (DNSRecordsModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m DNSRecordsModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m DNSRecordsModel) Search(query string) DNSRecordsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m DNSRecordsModel) NextSearchMatch() DNSRecordsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m DNSRecordsModel) PrevSearchMatch() DNSRecordsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

