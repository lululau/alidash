package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"aliyun-tui-viewer/internal/client"
	"aliyun-tui-viewer/internal/config"
	"aliyun-tui-viewer/internal/service"
	"aliyun-tui-viewer/internal/tui/components"
	"aliyun-tui-viewer/internal/tui/pages"
)

// Model is the root application model
type Model struct {
	// Current state
	currentPage   PageType
	previousPages []PageType // Navigation stack for back navigation

	// Profile
	profile  string
	profiles []string

	// Services and clients
	services *Services
	clients  *client.AliyunClients

	// Page models
	menuPage           pages.MenuModel
	ecsListPage        pages.ECSListModel
	ecsDetailPage      pages.DetailModel
	sgListPage         pages.SecurityGroupsModel
	sgRulesPage        pages.SecurityGroupRulesModel
	sgInstancesPage    pages.ECSListModel
	instSGPage         pages.SecurityGroupsModel
	dnsDomainsPage     pages.DNSDomainsModel
	dnsRecordsPage     pages.DNSRecordsModel
	slbListPage        pages.SLBListModel
	slbDetailPage      pages.DetailModel
	slbListenersPage   pages.SLBListenersModel
	slbVServerPage     pages.SLBVServerGroupsModel
	slbBackendPage     pages.SLBBackendServersModel
	ossBucketsPage     pages.OSSBucketsModel
	ossObjectsPage     pages.OSSObjectsModel
	ossDetailPage      pages.DetailModel
	rdsListPage        pages.RDSListModel
	rdsDetailPage      pages.DetailModel
	rdsDatabasesPage   pages.RDSDatabasesModel
	rdsAccountsPage    pages.RDSAccountsModel
	redisListPage      pages.RedisListModel
	redisDetailPage    pages.DetailModel
	redisAccountsPage  pages.RedisAccountsModel
	rocketmqListPage   pages.RocketMQListModel
	rocketmqDetailPage pages.DetailModel
	rocketmqTopicsPage pages.RocketMQTopicsModel
	rocketmqGroupsPage pages.RocketMQGroupsModel

	// Shared components
	modeLine components.ModeLineModel
	search   components.SearchModel
	modal    components.ModalModel

	// UI state
	width, height int
	loading       bool
	err           error

	// Yank tracker for double-y
	yankLastTime time.Time
	yankCount    int

	// Styles
	styles *Styles
	keys   KeyMap
}

// New creates a new application model
func New() (*Model, error) {
	// Load configuration
	cfg, err := config.LoadAliyunConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// Get current profile name
	currentProfile, err := config.GetCurrentProfileName()
	if err != nil {
		return nil, fmt.Errorf("getting current profile: %w", err)
	}

	// Get all profiles
	profiles, err := config.ListAllProfiles()
	if err != nil {
		profiles = []string{currentProfile}
	}

	// Create clients
	clientConfig := &client.Config{
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		RegionID:        cfg.RegionID,
		OssEndpoint:     cfg.OssEndpoint,
	}

	clients, err := client.NewAliyunClients(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("creating clients: %w", err)
	}

	// Create services
	services := &Services{
		ECS:      service.NewECSService(clients.ECS),
		DNS:      service.NewDNSService(clients.DNS),
		SLB:      service.NewSLBService(clients.SLB),
		RDS:      service.NewRDSService(clients.RDS),
		OSS:      service.NewOSSServiceWithCredentials(clients.OSS, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.OssEndpoint),
		Redis:    service.NewRedisService(clients.Redis),
		RocketMQ: service.NewRocketMQService(clients.RocketMQ),
	}

	m := &Model{
		currentPage:   PageMenu,
		previousPages: []PageType{},
		profile:       currentProfile,
		profiles:      profiles,
		services:      services,
		clients:       clients,
		styles:        GlobalStyles,
		keys:          GlobalKeyMap,
	}

	// Initialize page models
	m.menuPage = pages.NewMenuModel()
	m.modeLine = components.NewModeLineModel(currentProfile, PageMenu)
	m.search = components.NewSearchModel()
	m.modal = components.NewModalModel()

	return m, nil
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.menuPage.Init(),
	)
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle modal first if active
		if m.modal.Visible {
			var cmd tea.Cmd
			m.modal, cmd = m.modal.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		// Handle search input if active
		if m.search.Active {
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		// Global key handling
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Profile):
			m.modal = components.NewProfileSelectionModal(m.profiles, m.profile)
			return m, nil

		case key.Matches(msg, m.keys.Back):
			if m.currentPage != PageMenu {
				return m.navigateBack()
			}

		case key.Matches(msg, m.keys.Search):
			m.search = m.search.Activate()
			return m, m.search.Focus()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update all components with new size
		headerHeight := 1  // Mode line
		contentHeight := m.height - headerHeight

		m.menuPage = m.menuPage.SetSize(m.width, contentHeight)
		m.modeLine = m.modeLine.SetWidth(m.width)

		// Update current page size
		m = m.updateCurrentPageSize(contentHeight)

	case ErrorMsg:
		m.err = msg.Err
		m.loading = false
		m.modal = components.NewErrorModal(msg.Err.Error())

	case ModalDismissedMsg:
		m.modal = m.modal.Hide()

	case ProfileSwitchedMsg:
		m.profile = msg.ProfileName
		m.modeLine = m.modeLine.SetProfile(msg.ProfileName)
		// Clear all cached data by resetting page models
		m = m.clearCachedData()
		m.currentPage = PageMenu
		m.previousPages = []PageType{}

	case NavigateMsg:
		return m.navigateTo(msg.Page, msg.Data)

	case GoBackMsg:
		return m.navigateBack()

	// Handle search messages
	case SearchQueryMsg:
		return m.handleSearchQuery(msg.Query)

	case SearchExitMsg:
		m.search = m.search.Deactivate()

	case SearchNextMsg:
		return m.handleSearchNext()

	case SearchPrevMsg:
		return m.handleSearchPrev()

	// Handle data loaded messages
	case ECSInstancesLoadedMsg:
		m.loading = false
		m.ecsListPage = m.ecsListPage.SetData(msg.Instances)
		m.ecsListPage = m.ecsListPage.SetSize(m.width, m.height-1)

	case SecurityGroupsLoadedMsg:
		m.loading = false
		if m.currentPage == PageSecurityGroups {
			m.sgListPage = m.sgListPage.SetData(msg.SecurityGroups)
			m.sgListPage = m.sgListPage.SetSize(m.width, m.height-1)
		} else if m.currentPage == PageInstanceSecurityGroups {
			m.instSGPage = m.instSGPage.SetData(msg.SecurityGroups)
			m.instSGPage = m.instSGPage.SetSize(m.width, m.height-1)
		}

	case SecurityGroupRulesLoadedMsg:
		m.loading = false
		m.sgRulesPage = m.sgRulesPage.SetData(msg.Response)
		m.sgRulesPage = m.sgRulesPage.SetSize(m.width, m.height-1)

	case SecurityGroupInstancesLoadedMsg:
		m.loading = false
		m.sgInstancesPage = m.sgInstancesPage.SetData(msg.Instances)
		m.sgInstancesPage = m.sgInstancesPage.SetTitle(fmt.Sprintf("Instances using Security Group: %s", msg.SecurityGroupId))
		m.sgInstancesPage = m.sgInstancesPage.SetSize(m.width, m.height-1)

	case InstanceSecurityGroupsLoadedMsg:
		m.loading = false
		m.instSGPage = m.instSGPage.SetData(msg.SecurityGroups)
		m.instSGPage = m.instSGPage.SetTitle(fmt.Sprintf("Security Groups for Instance: %s", msg.InstanceId))
		m.instSGPage = m.instSGPage.SetSize(m.width, m.height-1)

	case DNSDomainsLoadedMsg:
		m.loading = false
		m.dnsDomainsPage = m.dnsDomainsPage.SetData(msg.Domains)
		m.dnsDomainsPage = m.dnsDomainsPage.SetSize(m.width, m.height-1)

	case DNSRecordsLoadedMsg:
		m.loading = false
		m.dnsRecordsPage = m.dnsRecordsPage.SetData(msg.Records, msg.DomainName)
		m.dnsRecordsPage = m.dnsRecordsPage.SetSize(m.width, m.height-1)

	case SLBInstancesLoadedMsg:
		m.loading = false
		m.slbListPage = m.slbListPage.SetData(msg.LoadBalancers)
		m.slbListPage = m.slbListPage.SetSize(m.width, m.height-1)

	case SLBListenersLoadedMsg:
		m.loading = false
		m.slbListenersPage = m.slbListenersPage.SetData(msg.Listeners, msg.LoadBalancerId)
		m.slbListenersPage = m.slbListenersPage.SetSize(m.width, m.height-1)

	case SLBVServerGroupsLoadedMsg:
		m.loading = false
		m.slbVServerPage = m.slbVServerPage.SetData(msg.VServerGroups, msg.LoadBalancerId)
		m.slbVServerPage = m.slbVServerPage.SetSize(m.width, m.height-1)

	case SLBBackendServersLoadedMsg:
		m.loading = false
		m.slbBackendPage = m.slbBackendPage.SetData(msg.BackendServers, msg.VServerGroupId)
		m.slbBackendPage = m.slbBackendPage.SetSize(m.width, m.height-1)

	case OSSBucketsLoadedMsg:
		m.loading = false
		m.ossBucketsPage = m.ossBucketsPage.SetData(msg.Buckets)
		m.ossBucketsPage = m.ossBucketsPage.SetSize(m.width, m.height-1)

	case OSSObjectsLoadedMsg:
		m.loading = false
		m.ossObjectsPage = m.ossObjectsPage.SetData(msg.Result, msg.BucketName, msg.Page)
		m.ossObjectsPage = m.ossObjectsPage.SetSize(m.width, m.height-1)

	case RDSInstancesLoadedMsg:
		m.loading = false
		m.rdsListPage = m.rdsListPage.SetData(msg.Instances)
		m.rdsListPage = m.rdsListPage.SetSize(m.width, m.height-1)

	case RDSDatabasesLoadedMsg:
		m.loading = false
		m.rdsDatabasesPage = m.rdsDatabasesPage.SetData(msg.Databases, msg.InstanceId)
		m.rdsDatabasesPage = m.rdsDatabasesPage.SetSize(m.width, m.height-1)

	case RDSAccountsLoadedMsg:
		m.loading = false
		m.rdsAccountsPage = m.rdsAccountsPage.SetData(msg.Accounts, msg.InstanceId)
		m.rdsAccountsPage = m.rdsAccountsPage.SetSize(m.width, m.height-1)

	case RedisInstancesLoadedMsg:
		m.loading = false
		m.redisListPage = m.redisListPage.SetData(msg.Instances)
		m.redisListPage = m.redisListPage.SetSize(m.width, m.height-1)

	case RedisAccountsLoadedMsg:
		m.loading = false
		m.redisAccountsPage = m.redisAccountsPage.SetData(msg.Accounts, msg.InstanceId)
		m.redisAccountsPage = m.redisAccountsPage.SetSize(m.width, m.height-1)

	case RocketMQInstancesLoadedMsg:
		m.loading = false
		m.rocketmqListPage = m.rocketmqListPage.SetData(msg.Instances)
		m.rocketmqListPage = m.rocketmqListPage.SetSize(m.width, m.height-1)

	case RocketMQTopicsLoadedMsg:
		m.loading = false
		m.rocketmqTopicsPage = m.rocketmqTopicsPage.SetData(msg.Topics, msg.InstanceId)
		m.rocketmqTopicsPage = m.rocketmqTopicsPage.SetSize(m.width, m.height-1)

	case RocketMQGroupsLoadedMsg:
		m.loading = false
		m.rocketmqGroupsPage = m.rocketmqGroupsPage.SetData(msg.Groups, msg.InstanceId)
		m.rocketmqGroupsPage = m.rocketmqGroupsPage.SetSize(m.width, m.height-1)

	// Handle copy messages
	case CopiedMsg:
		m.modal = components.NewInfoModal("Copied to clipboard!")
	}

	// Update current page
	cmd := m.updateCurrentPage(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var content string

	// Render current page
	switch m.currentPage {
	case PageMenu:
		content = m.menuPage.View()
	case PageECSList:
		content = m.ecsListPage.View()
	case PageECSDetail:
		content = m.ecsDetailPage.View()
	case PageSecurityGroups:
		content = m.sgListPage.View()
	case PageSecurityGroupRules:
		content = m.sgRulesPage.View()
	case PageSecurityGroupInstances:
		content = m.sgInstancesPage.View()
	case PageInstanceSecurityGroups:
		content = m.instSGPage.View()
	case PageDNSDomains:
		content = m.dnsDomainsPage.View()
	case PageDNSRecords:
		content = m.dnsRecordsPage.View()
	case PageSLBList:
		content = m.slbListPage.View()
	case PageSLBDetail:
		content = m.slbDetailPage.View()
	case PageSLBListeners:
		content = m.slbListenersPage.View()
	case PageSLBVServerGroups:
		content = m.slbVServerPage.View()
	case PageSLBBackendServers:
		content = m.slbBackendPage.View()
	case PageOSSBuckets:
		content = m.ossBucketsPage.View()
	case PageOSSObjects:
		content = m.ossObjectsPage.View()
	case PageOSSObjectDetail:
		content = m.ossDetailPage.View()
	case PageRDSList:
		content = m.rdsListPage.View()
	case PageRDSDetail:
		content = m.rdsDetailPage.View()
	case PageRDSDatabases:
		content = m.rdsDatabasesPage.View()
	case PageRDSAccounts:
		content = m.rdsAccountsPage.View()
	case PageRedisList:
		content = m.redisListPage.View()
	case PageRedisDetail:
		content = m.redisDetailPage.View()
	case PageRedisAccounts:
		content = m.redisAccountsPage.View()
	case PageRocketMQList:
		content = m.rocketmqListPage.View()
	case PageRocketMQDetail:
		content = m.rocketmqDetailPage.View()
	case PageRocketMQTopics:
		content = m.rocketmqTopicsPage.View()
	case PageRocketMQGroups:
		content = m.rocketmqGroupsPage.View()
	default:
		content = "Unknown page"
	}

	// Show loading spinner
	if m.loading {
		content = Center("Loading...", m.width, m.height-1)
	}

	// Build the full view
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		m.modeLine.View(),
	)

	// Overlay search bar if active
	if m.search.Active {
		view = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			m.search.View(),
			m.modeLine.View(),
		)
	}

	// Overlay modal if visible
	if m.modal.Visible {
		modalView := m.modal.View()
		view = Center(modalView, m.width, m.height)
	}

	return view
}

// navigateTo handles navigation to a specific page
func (m Model) navigateTo(page PageType, data interface{}) (Model, tea.Cmd) {
	// Push current page to stack
	m.previousPages = append(m.previousPages, m.currentPage)
	m.currentPage = page
	m.loading = true

	// Update mode line
	m.modeLine = m.modeLine.SetPage(page)

	var cmd tea.Cmd

	switch page {
	case PageECSList:
		m.ecsListPage = pages.NewECSListModel()
		cmd = LoadECSInstances(m.services.ECS)

	case PageECSDetail:
		if inst, ok := data.(interface{}); ok {
			m.ecsDetailPage = pages.NewDetailModel("ECS Detail", inst)
			m.loading = false
		}

	case PageSecurityGroups:
		m.sgListPage = pages.NewSecurityGroupsModel()
		cmd = LoadSecurityGroups(m.services.ECS)

	case PageSecurityGroupRules:
		if sgId, ok := data.(string); ok {
			m.sgRulesPage = pages.NewSecurityGroupRulesModel(sgId)
			cmd = LoadSecurityGroupRules(m.services.ECS, sgId)
		}

	case PageSecurityGroupInstances:
		if sgId, ok := data.(string); ok {
			m.sgInstancesPage = pages.NewECSListModel()
			cmd = LoadSecurityGroupInstances(m.services.ECS, sgId)
		}

	case PageInstanceSecurityGroups:
		if instId, ok := data.(string); ok {
			m.instSGPage = pages.NewSecurityGroupsModel()
			cmd = LoadInstanceSecurityGroups(m.services.ECS, instId)
		}

	case PageDNSDomains:
		m.dnsDomainsPage = pages.NewDNSDomainsModel()
		cmd = LoadDNSDomains(m.services.DNS)

	case PageDNSRecords:
		if domain, ok := data.(string); ok {
			m.dnsRecordsPage = pages.NewDNSRecordsModel()
			cmd = LoadDNSRecords(m.services.DNS, domain)
		}

	case PageSLBList:
		m.slbListPage = pages.NewSLBListModel()
		cmd = LoadSLBInstances(m.services.SLB)

	case PageSLBDetail:
		if lb, ok := data.(interface{}); ok {
			m.slbDetailPage = pages.NewDetailModel("SLB Detail", lb)
			m.loading = false
		}

	case PageSLBListeners:
		if lbId, ok := data.(string); ok {
			m.slbListenersPage = pages.NewSLBListenersModel()
			cmd = LoadSLBListeners(m.services.SLB, lbId)
		}

	case PageSLBVServerGroups:
		if lbId, ok := data.(string); ok {
			m.slbVServerPage = pages.NewSLBVServerGroupsModel()
			cmd = LoadSLBVServerGroups(m.services.SLB, lbId)
		}

	case PageSLBBackendServers:
		if vsgId, ok := data.(string); ok {
			m.slbBackendPage = pages.NewSLBBackendServersModel()
			cmd = LoadSLBBackendServers(m.services.SLB, vsgId, m.clients.ECS)
		}

	case PageOSSBuckets:
		m.ossBucketsPage = pages.NewOSSBucketsModel()
		cmd = LoadOSSBuckets(m.services.OSS)

	case PageOSSObjects:
		if bucket, ok := data.(string); ok {
			m.ossObjectsPage = pages.NewOSSObjectsModel(m.services.OSS, bucket)
			cmd = LoadOSSObjects(m.services.OSS, bucket, "", 20, 1)
		}

	case PageOSSObjectDetail:
		if obj, ok := data.(interface{}); ok {
			m.ossDetailPage = pages.NewDetailModel("OSS Object Detail", obj)
			m.loading = false
		}

	case PageRDSList:
		m.rdsListPage = pages.NewRDSListModel()
		cmd = LoadRDSInstances(m.services.RDS)

	case PageRDSDetail:
		if inst, ok := data.(interface{}); ok {
			m.rdsDetailPage = pages.NewDetailModel("RDS Detail", inst)
			m.loading = false
		}

	case PageRDSDatabases:
		if instId, ok := data.(string); ok {
			m.rdsDatabasesPage = pages.NewRDSDatabasesModel()
			cmd = LoadRDSDatabases(m.services.RDS, instId)
		}

	case PageRDSAccounts:
		if instId, ok := data.(string); ok {
			m.rdsAccountsPage = pages.NewRDSAccountsModel()
			cmd = LoadRDSAccounts(m.services.RDS, instId)
		}

	case PageRedisList:
		m.redisListPage = pages.NewRedisListModel()
		cmd = LoadRedisInstances(m.services.Redis)

	case PageRedisDetail:
		if inst, ok := data.(interface{}); ok {
			m.redisDetailPage = pages.NewDetailModel("Redis Detail", inst)
			m.loading = false
		}

	case PageRedisAccounts:
		if instId, ok := data.(string); ok {
			m.redisAccountsPage = pages.NewRedisAccountsModel()
			cmd = LoadRedisAccounts(m.services.Redis, instId)
		}

	case PageRocketMQList:
		m.rocketmqListPage = pages.NewRocketMQListModel()
		cmd = LoadRocketMQInstances(m.services.RocketMQ)

	case PageRocketMQDetail:
		if inst, ok := data.(interface{}); ok {
			m.rocketmqDetailPage = pages.NewDetailModel("RocketMQ Detail", inst)
			m.loading = false
		}

	case PageRocketMQTopics:
		if instId, ok := data.(string); ok {
			m.rocketmqTopicsPage = pages.NewRocketMQTopicsModel()
			cmd = LoadRocketMQTopics(m.services.RocketMQ, instId)
		}

	case PageRocketMQGroups:
		if instId, ok := data.(string); ok {
			m.rocketmqGroupsPage = pages.NewRocketMQGroupsModel()
			cmd = LoadRocketMQGroups(m.services.RocketMQ, instId)
		}

	default:
		m.loading = false
	}

	return m, cmd
}

// navigateBack handles back navigation
func (m Model) navigateBack() (Model, tea.Cmd) {
	if len(m.previousPages) == 0 {
		return m, nil
	}

	// Pop from stack
	lastIdx := len(m.previousPages) - 1
	prevPage := m.previousPages[lastIdx]
	m.previousPages = m.previousPages[:lastIdx]
	m.currentPage = prevPage

	// Update mode line
	m.modeLine = m.modeLine.SetPage(prevPage)

	return m, nil
}

// updateCurrentPage delegates update to the current page
func (m Model) updateCurrentPage(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch m.currentPage {
	case PageMenu:
		var newModel pages.MenuModel
		newModel, cmd = m.menuPage.Update(msg)
		m.menuPage = newModel

	case PageECSList:
		var newModel pages.ECSListModel
		newModel, cmd = m.ecsListPage.Update(msg)
		m.ecsListPage = newModel

	case PageECSDetail:
		var newModel pages.DetailModel
		newModel, cmd = m.ecsDetailPage.Update(msg)
		m.ecsDetailPage = newModel

	case PageSecurityGroups:
		var newModel pages.SecurityGroupsModel
		newModel, cmd = m.sgListPage.Update(msg)
		m.sgListPage = newModel

	case PageSecurityGroupRules:
		var newModel pages.SecurityGroupRulesModel
		newModel, cmd = m.sgRulesPage.Update(msg)
		m.sgRulesPage = newModel

	case PageDNSDomains:
		var newModel pages.DNSDomainsModel
		newModel, cmd = m.dnsDomainsPage.Update(msg)
		m.dnsDomainsPage = newModel

	case PageDNSRecords:
		var newModel pages.DNSRecordsModel
		newModel, cmd = m.dnsRecordsPage.Update(msg)
		m.dnsRecordsPage = newModel

	case PageSLBList:
		var newModel pages.SLBListModel
		newModel, cmd = m.slbListPage.Update(msg)
		m.slbListPage = newModel

	case PageSLBListeners:
		var newModel pages.SLBListenersModel
		newModel, cmd = m.slbListenersPage.Update(msg)
		m.slbListenersPage = newModel

	case PageSLBVServerGroups:
		var newModel pages.SLBVServerGroupsModel
		newModel, cmd = m.slbVServerPage.Update(msg)
		m.slbVServerPage = newModel

	case PageSLBBackendServers:
		var newModel pages.SLBBackendServersModel
		newModel, cmd = m.slbBackendPage.Update(msg)
		m.slbBackendPage = newModel

	case PageOSSBuckets:
		var newModel pages.OSSBucketsModel
		newModel, cmd = m.ossBucketsPage.Update(msg)
		m.ossBucketsPage = newModel

	case PageOSSObjects:
		var newModel pages.OSSObjectsModel
		newModel, cmd = m.ossObjectsPage.Update(msg)
		m.ossObjectsPage = newModel

	case PageRDSList:
		var newModel pages.RDSListModel
		newModel, cmd = m.rdsListPage.Update(msg)
		m.rdsListPage = newModel

	case PageRDSDatabases:
		var newModel pages.RDSDatabasesModel
		newModel, cmd = m.rdsDatabasesPage.Update(msg)
		m.rdsDatabasesPage = newModel

	case PageRDSAccounts:
		var newModel pages.RDSAccountsModel
		newModel, cmd = m.rdsAccountsPage.Update(msg)
		m.rdsAccountsPage = newModel

	case PageRedisList:
		var newModel pages.RedisListModel
		newModel, cmd = m.redisListPage.Update(msg)
		m.redisListPage = newModel

	case PageRedisAccounts:
		var newModel pages.RedisAccountsModel
		newModel, cmd = m.redisAccountsPage.Update(msg)
		m.redisAccountsPage = newModel

	case PageRocketMQList:
		var newModel pages.RocketMQListModel
		newModel, cmd = m.rocketmqListPage.Update(msg)
		m.rocketmqListPage = newModel

	case PageRocketMQTopics:
		var newModel pages.RocketMQTopicsModel
		newModel, cmd = m.rocketmqTopicsPage.Update(msg)
		m.rocketmqTopicsPage = newModel

	case PageRocketMQGroups:
		var newModel pages.RocketMQGroupsModel
		newModel, cmd = m.rocketmqGroupsPage.Update(msg)
		m.rocketmqGroupsPage = newModel

	// Detail pages
	case PageSLBDetail, PageOSSObjectDetail, PageRDSDetail, PageRedisDetail, PageRocketMQDetail:
		// Detail pages handle their own updates
	}

	return cmd
}

// updateCurrentPageSize updates the current page's size
func (m Model) updateCurrentPageSize(height int) Model {
	switch m.currentPage {
	case PageMenu:
		m.menuPage = m.menuPage.SetSize(m.width, height)
	case PageECSList:
		m.ecsListPage = m.ecsListPage.SetSize(m.width, height)
	case PageSecurityGroups:
		m.sgListPage = m.sgListPage.SetSize(m.width, height)
	// Add other pages as needed
	}
	return m
}

// clearCachedData resets all page models
func (m Model) clearCachedData() Model {
	m.ecsListPage = pages.NewECSListModel()
	m.sgListPage = pages.NewSecurityGroupsModel()
	m.dnsDomainsPage = pages.NewDNSDomainsModel()
	m.slbListPage = pages.NewSLBListModel()
	m.ossBucketsPage = pages.NewOSSBucketsModel()
	m.rdsListPage = pages.NewRDSListModel()
	m.redisListPage = pages.NewRedisListModel()
	m.rocketmqListPage = pages.NewRocketMQListModel()
	return m
}

// handleSearchQuery handles search query
func (m Model) handleSearchQuery(query string) (Model, tea.Cmd) {
	// Delegate to current page
	return m, nil
}

// handleSearchNext handles next search result
func (m Model) handleSearchNext() (Model, tea.Cmd) {
	return m, nil
}

// handleSearchPrev handles previous search result
func (m Model) handleSearchPrev() (Model, tea.Cmd) {
	return m, nil
}

// GetServices returns the services
func (m *Model) GetServices() *Services {
	return m.services
}

// GetClients returns the clients
func (m *Model) GetClients() *client.AliyunClients {
	return m.clients
}

