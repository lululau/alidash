package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines all key bindings for the application
type KeyMap struct {
	// Navigation
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding

	// Vim-style navigation
	VimUp   key.Binding
	VimDown key.Binding

	// Search
	Search     key.Binding
	SearchNext key.Binding
	SearchPrev key.Binding

	// Actions
	Yank        key.Binding
	Edit        key.Binding
	ViewPager   key.Binding
	Refresh     key.Binding
	Profile     key.Binding
	Region      key.Binding
	Help        key.Binding

	// Pagination (for OSS)
	NextPage  key.Binding
	PrevPage  key.Binding
	FirstPage key.Binding

	// Service-specific shortcuts
	SecurityGroups    key.Binding // g - view security groups for ECS instance
	ViewInstances     key.Binding // s - view instances using security group
	ViewListeners     key.Binding // l - view SLB listeners
	ViewVServerGroups key.Binding // v - view SLB VServer groups
	ViewDatabases     key.Binding // D - view RDS databases
	ViewAccounts      key.Binding // A - view accounts (RDS/Redis)
	ViewTopics        key.Binding // T - view RocketMQ topics
	ViewGroups        key.Binding // G - view RocketMQ groups
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("Q", "ctrl+c"),
			key.WithHelp("Q", "quit"),
		),

		// Vim-style navigation
		VimUp: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "up"),
		),
		VimDown: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "down"),
		),

		// Search
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		SearchNext: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next match"),
		),
		SearchPrev: key.NewBinding(
			key.WithKeys("N"),
			key.WithHelp("N", "prev match"),
		),

		// Actions
		Yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("yy", "copy"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		ViewPager: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "view in pager"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "ctrl+r"),
			key.WithHelp("r", "refresh"),
		),
		Profile: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "switch profile"),
		),
		Region: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "switch region"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),

		// Pagination
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

		// Service-specific shortcuts
		SecurityGroups: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "security groups"),
		),
		ViewInstances: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "instances"),
		),
		ViewListeners: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "listeners"),
		),
		ViewVServerGroups: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "VServer groups"),
		),
		ViewDatabases: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "databases"),
		),
		ViewAccounts: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "accounts"),
		),
		ViewTopics: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "topics"),
		),
		ViewGroups: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "groups"),
		),
	}
}

// GlobalKeyMap is the default key map instance
var GlobalKeyMap = DefaultKeyMap()

// ShortHelp returns key bindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.VimDown, k.VimUp, k.Enter, k.Back, k.Search, k.Quit}
}

// FullHelp returns key bindings for the full help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.VimUp, k.VimDown, k.Enter, k.Back},        // Navigation
		{k.Search, k.SearchNext, k.SearchPrev},      // Search
		{k.Yank, k.Edit, k.ViewPager, k.Profile},    // Actions
		{k.Quit, k.Help, k.Refresh},                 // General
	}
}

// MenuKeyMap returns key bindings specific to menu pages
type MenuKeyMap struct {
	KeyMap
	ECS           key.Binding
	SecurityGroup key.Binding
	DNS           key.Binding
	SLB           key.Binding
	OSS           key.Binding
	RDS           key.Binding
	Redis         key.Binding
	RocketMQ      key.Binding
}

// DefaultMenuKeyMap returns menu-specific key bindings
func DefaultMenuKeyMap() MenuKeyMap {
	return MenuKeyMap{
		KeyMap: DefaultKeyMap(),
		ECS: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "ECS"),
		),
		SecurityGroup: key.NewBinding(
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
	}
}

// ShortHelp returns key bindings for menu short help
func (k MenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.VimDown, k.VimUp, k.Enter, k.Profile, k.Quit,
	}
}

// FullHelp returns key bindings for menu full help
func (k MenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ECS, k.SecurityGroup, k.DNS, k.SLB},
		{k.OSS, k.RDS, k.Redis, k.RocketMQ},
		{k.Enter, k.Profile, k.Quit},
	}
}

// TableKeyMap returns key bindings specific to table/list views
type TableKeyMap struct {
	KeyMap
}

// DefaultTableKeyMap returns table-specific key bindings
func DefaultTableKeyMap() TableKeyMap {
	return TableKeyMap{
		KeyMap: DefaultKeyMap(),
	}
}

// DetailKeyMap returns key bindings specific to detail views
type DetailKeyMap struct {
	KeyMap
}

// DefaultDetailKeyMap returns detail-specific key bindings
func DefaultDetailKeyMap() DetailKeyMap {
	return DetailKeyMap{
		KeyMap: DefaultKeyMap(),
	}
}

// ShortHelp returns key bindings for detail view short help
func (k DetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back, k.Yank, k.Edit, k.ViewPager, k.Search,
	}
}

// OSSKeyMap returns key bindings specific to OSS views with pagination
type OSSKeyMap struct {
	TableKeyMap
}

// DefaultOSSKeyMap returns OSS-specific key bindings
func DefaultOSSKeyMap() OSSKeyMap {
	return OSSKeyMap{
		TableKeyMap: DefaultTableKeyMap(),
	}
}

// ShortHelp returns key bindings for OSS short help
func (k OSSKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.VimDown, k.VimUp, k.Enter, k.PrevPage, k.NextPage, k.Back,
	}
}

