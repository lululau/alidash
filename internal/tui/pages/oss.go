package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"aliyun-tui-viewer/internal/service"
	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/types"
)

// OSSBucketsModel represents the OSS buckets list page
type OSSBucketsModel struct {
	table   components.TableModel
	buckets []oss.BucketProperties
	width   int
	height  int
	keys    OSSBucketsKeyMap
}

// OSSBucketsKeyMap defines key bindings
type OSSBucketsKeyMap struct {
	Enter key.Binding
}

// DefaultOSSBucketsKeyMap returns default key bindings
func DefaultOSSBucketsKeyMap() OSSBucketsKeyMap {
	return OSSBucketsKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "objects"),
		),
	}
}

// NewOSSBucketsModel creates a new OSS buckets model
func NewOSSBucketsModel() OSSBucketsModel {
	columns := []table.Column{
		{Title: "Bucket Name", Width: 40},
		{Title: "Location", Width: 25},
		{Title: "Created", Width: 25},
		{Title: "Storage Class", Width: 15},
	}

	return OSSBucketsModel{
		table: components.NewTableModel(columns, "OSS Buckets"),
		keys:  DefaultOSSBucketsKeyMap(),
	}
}

// SetData sets the buckets data
func (m OSSBucketsModel) SetData(buckets []oss.BucketProperties) OSSBucketsModel {
	m.buckets = buckets

	rows := make([]table.Row, len(buckets))
	rowData := make([]interface{}, len(buckets))

	for i, bucket := range buckets {
		rows[i] = table.Row{
			bucket.Name,
			bucket.Location,
			bucket.CreationDate.Format("2006-01-02 15:04:05"),
			bucket.StorageClass,
		}
		rowData[i] = bucket
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	return m
}

// SetSize sets the size
func (m OSSBucketsModel) SetSize(width, height int) OSSBucketsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height)
	return m
}

// SelectedBucket returns the selected bucket
func (m OSSBucketsModel) SelectedBucket() *oss.BucketProperties {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.buckets) {
		return &m.buckets[idx]
	}
	return nil
}

// Init implements tea.Model
func (m OSSBucketsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m OSSBucketsModel) Update(msg tea.Msg) (OSSBucketsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if bucket := m.SelectedBucket(); bucket != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageOSSObjects,
						Data: bucket.Name,
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
func (m OSSBucketsModel) View() string {
	return m.table.View()
}

// Search searches in the list
func (m OSSBucketsModel) Search(query string) OSSBucketsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m OSSBucketsModel) NextSearchMatch() OSSBucketsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m OSSBucketsModel) PrevSearchMatch() OSSBucketsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// OSSObjectsModel represents the OSS objects list page with pagination
type OSSObjectsModel struct {
	table      components.TableModel
	objects    []oss.ObjectProperties
	bucketName string
	width      int
	height     int
	keys       OSSObjectsKeyMap

	// Pagination state
	currentPage     int
	hasNextPage     bool
	hasPrevPage     bool
	currentMarker   string
	nextMarker      string
	previousMarkers []string
	pageSize        int
	ossSvc          *service.OSSService
}

// OSSObjectsKeyMap defines key bindings
type OSSObjectsKeyMap struct {
	Enter     key.Binding
	NextPage  key.Binding
	PrevPage  key.Binding
	FirstPage key.Binding
}

// DefaultOSSObjectsKeyMap returns default key bindings
func DefaultOSSObjectsKeyMap() OSSObjectsKeyMap {
	return OSSObjectsKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "details"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("]"),
			key.WithHelp("]", "next page"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("["),
			key.WithHelp("[", "prev page"),
		),
		FirstPage: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "first page"),
		),
	}
}

// NewOSSObjectsModel creates a new OSS objects model
func NewOSSObjectsModel(svc *service.OSSService, bucketName string) OSSObjectsModel {
	columns := []table.Column{
		{Title: "Object Key", Width: 80},
		{Title: "Size", Width: 12},
		{Title: "Last Modified", Width: 22},
		{Title: "Storage Class", Width: 14},
		{Title: "ETag", Width: 36},
	}

	return OSSObjectsModel{
		table:           components.NewTableModel(columns, fmt.Sprintf("Objects in %s", bucketName)),
		bucketName:      bucketName,
		keys:            DefaultOSSObjectsKeyMap(),
		pageSize:        20,
		currentPage:     1,
		previousMarkers: []string{},
		ossSvc:          svc,
	}
}

// SetData sets the objects data with pagination info
func (m OSSObjectsModel) SetData(result *service.ObjectListResult, bucketName string, page int) OSSObjectsModel {
	m.objects = result.Objects
	m.bucketName = bucketName
	m.currentPage = page
	m.hasNextPage = result.IsTruncated
	m.nextMarker = result.NextMarker
	m.hasPrevPage = len(m.previousMarkers) > 0

	rows := make([]table.Row, len(m.objects))
	rowData := make([]interface{}, len(m.objects))

	for i, obj := range m.objects {
		rows[i] = table.Row{
			obj.Key,
			formatSize(obj.Size),
			obj.LastModified.Format("2006-01-02 15:04:05"),
			obj.StorageClass,
			obj.ETag,
		}
		rowData[i] = obj
	}

	m.table = m.table.SetRows(rows)
	m.table = m.table.SetRowData(rowData)
	m.table = m.table.SetTitle(fmt.Sprintf("Objects in %s (Page %d)", bucketName, page))
	return m
}

// SetSize sets the size
func (m OSSObjectsModel) SetSize(width, height int) OSSObjectsModel {
	m.width = width
	m.height = height
	m.table = m.table.SetSize(width, height-2) // Account for pagination info
	return m
}

// SelectedObject returns the selected object
func (m OSSObjectsModel) SelectedObject() *oss.ObjectProperties {
	idx := m.table.SelectedRow()
	if idx >= 0 && idx < len(m.objects) {
		return &m.objects[idx]
	}
	return nil
}

// Init implements tea.Model
func (m OSSObjectsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m OSSObjectsModel) Update(msg tea.Msg) (OSSObjectsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if obj := m.SelectedObject(); obj != nil {
				return m, func() tea.Msg {
					return types.NavigateMsg{
						Page: types.PageOSSObjectDetail,
						Data: *obj,
					}
				}
			}

		case key.Matches(msg, m.keys.NextPage):
			if m.hasNextPage {
				m.previousMarkers = append(m.previousMarkers, m.currentMarker)
				m.currentMarker = m.nextMarker
				return m, loadOSSObjects(m.ossSvc, m.bucketName, m.nextMarker, m.pageSize, m.currentPage+1)
			}

		case key.Matches(msg, m.keys.PrevPage):
			if len(m.previousMarkers) > 0 {
				lastIdx := len(m.previousMarkers) - 1
				m.currentMarker = m.previousMarkers[lastIdx]
				m.previousMarkers = m.previousMarkers[:lastIdx]
				return m, loadOSSObjects(m.ossSvc, m.bucketName, m.currentMarker, m.pageSize, m.currentPage-1)
			}

		case key.Matches(msg, m.keys.FirstPage):
			m.currentMarker = ""
			m.previousMarkers = []string{}
			return m, loadOSSObjects(m.ossSvc, m.bucketName, "", m.pageSize, 1)
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m OSSObjectsModel) View() string {
	// Build pagination info
	pageInfo := fmt.Sprintf("Page %d", m.currentPage)
	if m.hasNextPage {
		pageInfo += "+"
	}

	navHelp := ""
	if m.hasPrevPage {
		navHelp += "[ Prev | "
	}
	if m.hasNextPage {
		navHelp += "] Next | "
	}
	navHelp += "0 First"

	paginationLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06B6D4")).
		Render(fmt.Sprintf(" %s | %s ", pageInfo, navHelp))

	return m.table.View() + "\n" + paginationLine
}

// Search searches in the list
func (m OSSObjectsModel) Search(query string) OSSObjectsModel {
	m.table = m.table.Search(query)
	return m
}

// NextSearchMatch moves to next search match
func (m OSSObjectsModel) NextSearchMatch() OSSObjectsModel {
	m.table = m.table.NextSearchMatch()
	return m
}

// PrevSearchMatch moves to previous search match
func (m OSSObjectsModel) PrevSearchMatch() OSSObjectsModel {
	m.table = m.table.PrevSearchMatch()
	return m
}

// Helper function to format file size
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// OSSObjectsLoadedMsg contains loaded OSS objects with pagination
type OSSObjectsLoadedMsg struct {
	Result     *service.ObjectListResult
	BucketName string
	Page       int
}

// OSSErrorMsg indicates an OSS error
type OSSErrorMsg struct {
	Err error
}

// loadOSSObjects creates a command to load OSS objects with pagination
func loadOSSObjects(svc *service.OSSService, bucketName, marker string, pageSize, page int) tea.Cmd {
	return func() tea.Msg {
		result, err := svc.FetchObjects(bucketName, marker, pageSize)
		if err != nil {
			return OSSErrorMsg{Err: err}
		}
		return OSSObjectsLoadedMsg{
			Result:     result,
			BucketName: bucketName,
			Page:       page,
		}
	}
}

