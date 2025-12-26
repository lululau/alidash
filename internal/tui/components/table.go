package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableModel wraps bubbles/table with additional features
type TableModel struct {
	table       table.Model
	columns     []table.Column
	rows        []table.Row
	title       string
	width       int
	height      int
	focused     bool
	searchQuery string
	searchIndex int
	searchCount int
	matchRows   []int // Row indices that match search
	keys        TableKeyMap

	// Cursor and scroll for custom rendering
	cursor     int
	scrollOffset int

	// Yank tracking
	yankLastTime time.Time
	yankCount    int

	// Row data for copying
	rowData []interface{}

	// Styles
	styles TableStyles
}

// TableKeyMap defines key bindings for the table
type TableKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Enter    key.Binding
	Yank     key.Binding
}

// DefaultTableKeyMap returns default key bindings
func DefaultTableKeyMap() TableKeyMap {
	return TableKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("pgdn", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "first"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "last"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("yy", "copy"),
		),
	}
}

// TableStyles defines styles for the table
type TableStyles struct {
	Header      lipgloss.Style
	Cell        lipgloss.Style
	Selected    lipgloss.Style
	Border      lipgloss.Style
	Title       lipgloss.Style
	SearchMatch lipgloss.Style
}

// DefaultTableStyles returns default table styles
func DefaultTableStyles() TableStyles {
	return TableStyles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F59E0B")).
			Padding(0, 1),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Padding(0, 1),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Bold(true).
			Padding(0, 1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#374151")),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")),
		SearchMatch: lipgloss.NewStyle().
			Background(lipgloss.Color("#CA8A04")).
			Foreground(lipgloss.Color("#000000")),
	}
}

// NewTableModel creates a new table model
func NewTableModel(columns []table.Column, title string) TableModel {
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Set table styles - ensure Selected style highlights entire row
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#374151")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("#F59E0B"))
	// Selected style for the entire row - background color will extend to cell width
	s.Selected = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7C3AED")).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(lipgloss.Color("#E5E7EB"))

	t.SetStyles(s)
	t.Focus() // Ensure table starts focused

	return TableModel{
		table:   t,
		columns: columns,
		title:   title,
		keys:    DefaultTableKeyMap(),
		styles:  DefaultTableStyles(),
		focused: true,
	}
}

// SetRows sets the table rows
func (m TableModel) SetRows(rows []table.Row) TableModel {
	m.rows = rows
	m.cursor = 0
	m.scrollOffset = 0
	// Clear search when data changes
	m.searchQuery = ""
	m.searchIndex = -1
	m.searchCount = 0
	m.matchRows = nil
	return m
}

// SetRowData sets the underlying data for each row (for copying)
func (m TableModel) SetRowData(data []interface{}) TableModel {
	m.rowData = data
	return m
}

// SetTitle sets the table title
func (m TableModel) SetTitle(title string) TableModel {
	m.title = title
	return m
}

// SetSize sets the table size
func (m TableModel) SetSize(width, height int) TableModel {
	m.width = width
	m.height = height
	// Account for title and borders
	tableHeight := height - 4
	if tableHeight < 1 {
		tableHeight = 1
	}
	m.table.SetWidth(width - 2)
	m.table.SetHeight(tableHeight)
	return m
}

// SetFocused sets the focus state
func (m TableModel) SetFocused(focused bool) TableModel {
	m.focused = focused
	m.table.Focus()
	return m
}

// SetColumns sets the table columns
func (m TableModel) SetColumns(columns []table.Column) TableModel {
	m.columns = columns
	m.table.SetColumns(columns)
	return m
}

// SelectedRow returns the currently selected row index
func (m TableModel) SelectedRow() int {
	return m.cursor
}

// SelectedRowData returns the data for the selected row
func (m TableModel) SelectedRowData() interface{} {
	if m.cursor >= 0 && m.cursor < len(m.rowData) {
		return m.rowData[m.cursor]
	}
	return nil
}

// RowCount returns the number of rows
func (m TableModel) RowCount() int {
	return len(m.rows)
}

// Init implements tea.Model
func (m TableModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m TableModel) Update(msg tea.Msg) (TableModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Yank):
			// Handle yank (double-y for copy)
			now := time.Now()
			if now.Sub(m.yankLastTime) < 500*time.Millisecond {
				m.yankCount++
			} else {
				m.yankCount = 1
			}
			m.yankLastTime = now

			if m.yankCount >= 2 {
				// Double-y detected, return copy command
				m.yankCount = 0
				if data := m.SelectedRowData(); data != nil {
					return m, func() tea.Msg {
						return CopyDataMsg{Data: data}
					}
				}
			}
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			// Return selection message
			return m, func() tea.Msg {
				return TableSelectMsg{
					Index: m.cursor,
					Data:  m.SelectedRowData(),
				}
			}

		case key.Matches(msg, m.keys.Up):
			m.moveCursor(-1)
			return m, nil

		case key.Matches(msg, m.keys.Down):
			m.moveCursor(1)
			return m, nil

		case key.Matches(msg, m.keys.PageUp):
			m.moveCursor(-m.visibleRows())
			return m, nil

		case key.Matches(msg, m.keys.PageDown):
			m.moveCursor(m.visibleRows())
			return m, nil

		case key.Matches(msg, m.keys.Home):
			m.cursor = 0
			m.scrollOffset = 0
			return m, nil

		case key.Matches(msg, m.keys.End):
			if len(m.rows) > 0 {
				m.cursor = len(m.rows) - 1
				m.ensureCursorVisible()
			}
			return m, nil
		}
	}

	return m, nil
}

// moveCursor moves the cursor by delta and ensures it stays within bounds
func (m *TableModel) moveCursor(delta int) {
	if len(m.rows) == 0 {
		return
	}

	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.rows) {
		m.cursor = len(m.rows) - 1
	}

	m.ensureCursorVisible()
}

// ensureCursorVisible adjusts scroll offset to keep cursor in view
func (m *TableModel) ensureCursorVisible() {
	visible := m.visibleRows()
	if visible <= 0 {
		return
	}

	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+visible {
		m.scrollOffset = m.cursor - visible + 1
	}
}

// visibleRows returns the number of visible data rows
func (m TableModel) visibleRows() int {
	// Account for header row and borders
	rows := m.height - 6
	if rows < 1 {
		rows = 1
	}
	return rows
}

// View implements tea.Model
func (m TableModel) View() string {
	var b strings.Builder

	// Title
	if m.title != "" {
		title := m.styles.Title.Render(m.title)
		b.WriteString(title)
		b.WriteString("\n")
	}

	// Render custom table
	tableContent := m.renderTable()

	// Add border
	bordered := m.styles.Border.
		Width(m.width - 2).
		Render(tableContent)

	b.WriteString(bordered)

	// Search info
	if m.searchQuery != "" {
		searchInfo := fmt.Sprintf(" Search: %s (%d/%d) ", m.searchQuery, m.searchIndex+1, m.searchCount)
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Render(searchInfo))
	}

	return b.String()
}

// renderTable renders the table with custom highlighting
func (m TableModel) renderTable() string {
	var b strings.Builder

	// Render header
	headerCells := make([]string, len(m.columns))
	for i, col := range m.columns {
		cell := truncateString(col.Title, col.Width)
		cell = padString(cell, col.Width)
		headerCells[i] = m.styles.Header.Render(cell)
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))
	b.WriteString("\n")

	// Header separator
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#374151"))
	totalWidth := 0
	for _, col := range m.columns {
		totalWidth += col.Width + 2 // +2 for padding
	}
	b.WriteString(separatorStyle.Render(strings.Repeat("─", totalWidth)))
	b.WriteString("\n")

	// Render visible rows
	visible := m.visibleRows()
	endIdx := m.scrollOffset + visible
	if endIdx > len(m.rows) {
		endIdx = len(m.rows)
	}

	for rowIdx := m.scrollOffset; rowIdx < endIdx; rowIdx++ {
		row := m.rows[rowIdx]
		isSelected := rowIdx == m.cursor

		rowCells := make([]string, len(m.columns))
		for colIdx, col := range m.columns {
			cellContent := ""
			if colIdx < len(row) {
				cellContent = row[colIdx]
			}

			// Truncate and pad cell content
			displayContent := truncateString(cellContent, col.Width)
			displayContent = padString(displayContent, col.Width)

			// Apply search highlighting to the content
			if m.searchQuery != "" {
				displayContent = m.highlightSearchMatch(displayContent)
			}

			// Apply row style
			if isSelected {
				// For selected row, apply selected style
				rowCells[colIdx] = m.styles.Selected.Render(displayContent)
			} else {
				rowCells[colIdx] = m.styles.Cell.Render(displayContent)
			}
		}

		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, rowCells...))
		if rowIdx < endIdx-1 {
			b.WriteString("\n")
		}
	}

	// Fill empty rows if needed
	for i := endIdx - m.scrollOffset; i < visible; i++ {
		b.WriteString("\n")
		emptyCells := make([]string, len(m.columns))
		for colIdx, col := range m.columns {
			emptyCells[colIdx] = m.styles.Cell.Render(strings.Repeat(" ", col.Width))
		}
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, emptyCells...))
	}

	return b.String()
}

// highlightSearchMatch highlights search matches in a string
func (m TableModel) highlightSearchMatch(s string) string {
	if m.searchQuery == "" {
		return s
	}

	lowerS := strings.ToLower(s)
	lowerQuery := strings.ToLower(m.searchQuery)

	var result strings.Builder
	lastIdx := 0

	for {
		idx := strings.Index(lowerS[lastIdx:], lowerQuery)
		if idx == -1 {
			break
		}

		actualIdx := lastIdx + idx

		// Add text before match
		result.WriteString(s[lastIdx:actualIdx])

		// Add highlighted match (preserving original case)
		match := s[actualIdx : actualIdx+len(m.searchQuery)]
		result.WriteString(m.styles.SearchMatch.Render(match))

		lastIdx = actualIdx + len(m.searchQuery)
	}

	// Add remaining text
	result.WriteString(s[lastIdx:])

	return result.String()
}

// truncateString truncates a string to fit within width
func truncateString(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width <= 3 {
		return s[:width]
	}
	return s[:width-3] + "..."
}

// padString pads a string to the specified width
func padString(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// Search searches for a query in the table
func (m TableModel) Search(query string) TableModel {
	if query == "" {
		m.searchQuery = ""
		m.searchIndex = -1
		m.searchCount = 0
		m.matchRows = nil
		return m
	}

	m.searchQuery = query
	lowerQuery := strings.ToLower(query)

	// Find matching rows
	m.matchRows = nil
	for i, row := range m.rows {
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), lowerQuery) {
				m.matchRows = append(m.matchRows, i)
				break
			}
		}
	}

	m.searchCount = len(m.matchRows)
	if m.searchCount > 0 {
		m.searchIndex = 0
		m.cursor = m.matchRows[0]
		m.ensureCursorVisible()
	} else {
		m.searchIndex = -1
	}

	return m
}

// NextSearchMatch moves to the next search match
func (m TableModel) NextSearchMatch() TableModel {
	if m.searchQuery == "" || m.searchCount == 0 {
		return m
	}

	// Find current position in matchRows
	currentMatchIdx := -1
	for i, rowIdx := range m.matchRows {
		if rowIdx == m.cursor {
			currentMatchIdx = i
			break
		}
	}

	// Move to next match
	nextMatchIdx := (currentMatchIdx + 1) % m.searchCount
	m.searchIndex = nextMatchIdx
	m.cursor = m.matchRows[nextMatchIdx]
	m.ensureCursorVisible()

	return m
}

// PrevSearchMatch moves to the previous search match
func (m TableModel) PrevSearchMatch() TableModel {
	if m.searchQuery == "" || m.searchCount == 0 {
		return m
	}

	// Find current position in matchRows
	currentMatchIdx := -1
	for i, rowIdx := range m.matchRows {
		if rowIdx == m.cursor {
			currentMatchIdx = i
			break
		}
	}

	// Move to previous match
	prevMatchIdx := currentMatchIdx - 1
	if prevMatchIdx < 0 {
		prevMatchIdx = m.searchCount - 1
	}
	m.searchIndex = prevMatchIdx
	m.cursor = m.matchRows[prevMatchIdx]
	m.ensureCursorVisible()

	return m
}

// ClearSearch clears the search
func (m TableModel) ClearSearch() TableModel {
	m.searchQuery = ""
	m.searchIndex = -1
	m.searchCount = 0
	m.matchRows = nil
	return m
}

// TableSelectMsg is sent when a row is selected
type TableSelectMsg struct {
	Index int
	Data  interface{}
}

// CopyDataMsg is sent when data should be copied
type CopyDataMsg struct {
	Data interface{}
}

// Helper function to create columns from headers
func CreateColumns(headers []string, widths []int) []table.Column {
	cols := make([]table.Column, len(headers))
	for i, h := range headers {
		width := 20 // default width
		if i < len(widths) && widths[i] > 0 {
			width = widths[i]
		}
		cols[i] = table.Column{
			Title: h,
			Width: width,
		}
	}
	return cols
}

// Helper function to auto-calculate column widths
func AutoColumnWidths(headers []string, rows []table.Row, maxWidth int) []int {
	widths := make([]int, len(headers))

	// Start with header widths
	for i, h := range headers {
		widths[i] = len(h)
	}

	// Check row contents
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Apply max width constraint and add padding
	totalWidth := 0
	for i := range widths {
		widths[i] += 2 // padding
		if widths[i] > 50 {
			widths[i] = 50
		}
		totalWidth += widths[i]
	}

	// Scale down if too wide
	if totalWidth > maxWidth && maxWidth > 0 {
		scale := float64(maxWidth) / float64(totalWidth)
		for i := range widths {
			widths[i] = int(float64(widths[i]) * scale)
			if widths[i] < 5 {
				widths[i] = 5
			}
		}
	}

	return widths
}

