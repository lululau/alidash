package pages

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"aliyun-tui-viewer/internal/tui/types"
)

// MenuItem represents a menu item
type MenuItem struct {
	title       string
	description string
	shortcut    rune
	page        types.PageType
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }
func (i MenuItem) FilterValue() string { return i.title }

// MenuModel represents the main menu page
type MenuModel struct {
	list   list.Model
	width  int
	height int
	keys   MenuKeyMap
}

// MenuKeyMap defines key bindings for the menu
type MenuKeyMap struct {
	Enter    key.Binding
	ECS      key.Binding
	SG       key.Binding
	DNS      key.Binding
	SLB      key.Binding
	OSS      key.Binding
	RDS      key.Binding
	Redis    key.Binding
	RocketMQ key.Binding
	Quit     key.Binding
}

// DefaultMenuKeyMap returns default menu key bindings
func DefaultMenuKeyMap() MenuKeyMap {
	return MenuKeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		ECS: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "ECS"),
		),
		SG: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "Security Groups"),
		),
		DNS: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "DNS"),
		),
		SLB: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "SLB"),
		),
		OSS: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "OSS"),
		),
		RDS: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "RDS"),
		),
		Redis: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "Redis"),
		),
		RocketMQ: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "RocketMQ"),
		),
		Quit: key.NewBinding(
			key.WithKeys("Q"),
			key.WithHelp("Q", "quit"),
		),
	}
}

// NewMenuModel creates a new menu model
func NewMenuModel() MenuModel {
	items := []list.Item{
		MenuItem{title: "(s) ECS Instances", description: "View ECS instances", shortcut: 's', page: types.PageECSList},
		MenuItem{title: "(g) Security Groups", description: "View ECS security groups", shortcut: 'g', page: types.PageSecurityGroups},
		MenuItem{title: "(d) DNS Management", description: "View AliDNS domains and records", shortcut: 'd', page: types.PageDNSDomains},
		MenuItem{title: "(b) SLB Instances", description: "View SLB instances", shortcut: 'b', page: types.PageSLBList},
		MenuItem{title: "(o) OSS Management", description: "Browse OSS buckets and objects", shortcut: 'o', page: types.PageOSSBuckets},
		MenuItem{title: "(r) RDS Instances", description: "View RDS instances", shortcut: 'r', page: types.PageRDSList},
		MenuItem{title: "(i) Redis Instances", description: "View Redis instances", shortcut: 'i', page: types.PageRedisList},
		MenuItem{title: "(m) RocketMQ Instances", description: "View RocketMQ instances", shortcut: 'm', page: types.PageRocketMQList},
	}

	// Create delegate
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#F59E0B")).
		Bold(true).
		BorderLeftForeground(lipgloss.Color("#7C3AED"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#9CA3AF"))
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#E5E7EB"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(lipgloss.Color("#6B7280"))

	l := list.New(items, delegate, 0, 0)
	l.Title = "Aliyun TUI Dashboard"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginLeft(2)

	return MenuModel{
		list: l,
		keys: DefaultMenuKeyMap(),
	}
}

// SetSize sets the menu size
func (m MenuModel) SetSize(width, height int) MenuModel {
	m.width = width
	m.height = height
	m.list.SetWidth(width)
	m.list.SetHeight(height)
	return m
}

// Init implements tea.Model
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if item, ok := m.list.SelectedItem().(MenuItem); ok {
				return m, func() tea.Msg {
					return types.NavigateMsg{Page: item.page}
				}
			}

		case key.Matches(msg, m.keys.ECS):
			return m, func() tea.Msg {
				return types.NavigateMsg{Page: types.PageECSList}
			}

		case key.Matches(msg, m.keys.SG):
			return m, func() tea.Msg {
				return types.NavigateMsg{Page: types.PageSecurityGroups}
			}

		case key.Matches(msg, m.keys.DNS):
			return m, func() tea.Msg {
				return types.NavigateMsg{Page: types.PageDNSDomains}
			}

		case key.Matches(msg, m.keys.SLB):
			return m, func() tea.Msg {
				return types.NavigateMsg{Page: types.PageSLBList}
			}

		case key.Matches(msg, m.keys.OSS):
			return m, func() tea.Msg {
				return types.NavigateMsg{Page: types.PageOSSBuckets}
			}

		case key.Matches(msg, m.keys.RDS):
			return m, func() tea.Msg {
				return types.NavigateMsg{Page: types.PageRDSList}
			}

		case key.Matches(msg, m.keys.Redis):
			return m, func() tea.Msg {
				return types.NavigateMsg{Page: types.PageRedisList}
			}

		case key.Matches(msg, m.keys.RocketMQ):
			return m, func() tea.Msg {
				return types.NavigateMsg{Page: types.PageRocketMQList}
			}

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View implements tea.Model
func (m MenuModel) View() string {
	return m.list.View()
}
