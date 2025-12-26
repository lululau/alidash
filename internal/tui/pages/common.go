package pages

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"aliyun-tui-viewer/internal/tui/components"
)

// DetailModel is a generic detail view model for displaying JSON data
type DetailModel struct {
	viewport components.ViewportModel
	title    string
	data     interface{}
	width    int
	height   int
}

// NewDetailModel creates a new detail model
func NewDetailModel(title string, data interface{}) DetailModel {
	return DetailModel{
		viewport: components.NewViewportModel(title, data),
		title:    title,
		data:     data,
	}
}

// SetData sets the detail data
func (m DetailModel) SetData(data interface{}) DetailModel {
	m.data = data
	m.viewport = m.viewport.SetData(data)
	return m
}

// SetTitle sets the title
func (m DetailModel) SetTitle(title string) DetailModel {
	m.title = title
	m.viewport = m.viewport.SetTitle(title)
	return m
}

// SetSize sets the size
func (m DetailModel) SetSize(width, height int) DetailModel {
	m.width = width
	m.height = height
	m.viewport = m.viewport.SetSize(width, height)
	return m
}

// GetData returns the data
func (m DetailModel) GetData() interface{} {
	return m.data
}

// Init implements tea.Model
func (m DetailModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m DetailModel) View() string {
	return m.viewport.View()
}

// Search searches in the detail view
func (m DetailModel) Search(query string) DetailModel {
	m.viewport = m.viewport.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m DetailModel) NextSearchMatch() DetailModel {
	m.viewport = m.viewport.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m DetailModel) PrevSearchMatch() DetailModel {
	m.viewport = m.viewport.PrevSearchMatch()
	return m
}

// BaseListModel is a base model for list pages
type BaseListModel struct {
	table    components.TableModel
	title    string
	width    int
	height   int
	columns  []table.Column
	rowData  []interface{}
}

// NewBaseListModel creates a new base list model
func NewBaseListModel(title string, columns []table.Column) BaseListModel {
	return BaseListModel{
		table:   components.NewTableModel(columns, title),
		title:   title,
		columns: columns,
	}
}

// SetSize sets the size
func (m BaseListModel) SetSize(width, height int) BaseListModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SetTitle sets the title
func (m BaseListModel) SetTitle(title string) BaseListModel {
	m.title = title
	m.table = m.table.SetTitle(title)
	return m
}

// SetRows sets the table rows
func (m BaseListModel) SetRows(rows []table.Row) BaseListModel {
	m.table = m.table.SetRows(rows)
	return m
}

// SetRowData sets the row data for copying
func (m BaseListModel) SetRowData(data []interface{}) BaseListModel {
	m.rowData = data
	m.table = m.table.SetRowData(data)
	return m
}

// SelectedRow returns the selected row index
func (m BaseListModel) SelectedRow() int {
	return m.table.SelectedRow()
}

// SelectedRowData returns the selected row data
func (m BaseListModel) SelectedRowData() interface{} {
	return m.table.SelectedRowData()
}

// Init implements tea.Model
func (m BaseListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m BaseListModel) Update(msg tea.Msg) (BaseListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m BaseListModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m BaseListModel) Search(query string) BaseListModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m BaseListModel) NextSearchMatch() BaseListModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m BaseListModel) PrevSearchMatch() BaseListModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

