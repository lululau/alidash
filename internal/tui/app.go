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

	// Region
	region        string
	regions       []string
	regionService *service.RegionService

	// Services and clients
	services *Services
	clients  *client.AliyunClients

	// Page models
	menuPage           pages.MenuModel
	ecsListPage        pages.ECSListModel
	ecsDetailPage      pages.ECSDetailModel // Formatted detail view
	ecsJSONDetailPage  pages.DetailModel    // JSON detail view
	ecsDiskPage        pages.ECSDiskModel   // Disk/storage page
	ecsENIPage         pages.ECSENIModel    // Network interfaces page
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
	header   components.HeaderModel
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

	// Create region service
	regionService := service.NewRegionService(cfg.AccessKeyID, cfg.AccessKeySecret, currentProfile)

	m := &Model{
		currentPage:   PageMenu,
		previousPages: []PageType{},
		profile:       currentProfile,
		profiles:      profiles,
		region:        cfg.RegionID,
		regionService: regionService,
		services:      services,
		clients:       clients,
		styles:        GlobalStyles,
		keys:          GlobalKeyMap,
	}

	// Initialize page models
	m.menuPage = pages.NewMenuModel()
	m.header = components.NewHeaderModel("Aliyun TUI Dashboard", currentProfile, cfg.RegionID)
	m.modeLine = components.NewModeLineModel(currentProfile, cfg.RegionID, PageMenu)
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
			// Q (capital) always quits
			return m, tea.Quit

		case key.Matches(msg, m.keys.Profile):
			m.modal = components.NewProfileSelectionModal(m.profiles, m.profile)
			return m, nil

		case key.Matches(msg, m.keys.Region):
			// Show loading modal and start async region loading
			m.modal = components.NewRegionSelectionModal(m.region)
			return m, m.loadRegions()

		case key.Matches(msg, m.keys.Back):
			// q/esc goes back, but not on menu page (menu uses Q to quit)
			if m.currentPage != PageMenu {
				return m.navigateBack()
			}
			// On menu page, esc does nothing, q is handled by menu shortcuts

		case key.Matches(msg, m.keys.Search):
			// Don't activate search on menu page
			if m.currentPage != PageMenu {
				m.search = m.search.Activate()
				return m, m.search.Focus()
			}

		case key.Matches(msg, m.keys.SearchNext):
			// n for next search match
			if m.search.Query() != "" {
				return m.handleSearchNext()
			}

		case key.Matches(msg, m.keys.SearchPrev):
			// N for previous search match
			if m.search.Query() != "" {
				return m.handleSearchPrev()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update all components with new size
		// Layout: header (1) + empty line (1) + content + modeline (1)
		chromeHeight := 3 // header + empty line + modeline
		contentHeight := m.height - chromeHeight

		m.header = m.header.SetWidth(m.width)
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

	case components.ProfileSelectedMsg:
		// Actual profile switching - called when user selects a profile from modal
		if err := config.SwitchProfile(msg.Profile); err != nil {
			m.modal = components.NewErrorModal(fmt.Sprintf("Failed to switch profile: %v", err))
			return m, nil
		}

		// Reload configuration with new profile
		cfg, err := config.LoadAliyunConfig()
		if err != nil {
			m.modal = components.NewErrorModal(fmt.Sprintf("Failed to reload config: %v", err))
			return m, nil
		}

		// Recreate clients with new credentials
		clientConfig := &client.Config{
			AccessKeyID:     cfg.AccessKeyID,
			AccessKeySecret: cfg.AccessKeySecret,
			RegionID:        cfg.RegionID,
			OssEndpoint:     cfg.OssEndpoint,
		}

		newClients, err := client.NewAliyunClients(clientConfig)
		if err != nil {
			m.modal = components.NewErrorModal(fmt.Sprintf("Failed to create clients: %v", err))
			return m, nil
		}

		// Update clients and recreate services
		m.clients = newClients
		m.services = &Services{
			ECS:      service.NewECSService(newClients.ECS),
			DNS:      service.NewDNSService(newClients.DNS),
			SLB:      service.NewSLBService(newClients.SLB),
			RDS:      service.NewRDSService(newClients.RDS),
			OSS:      service.NewOSSServiceWithCredentials(newClients.OSS, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.OssEndpoint),
			Redis:    service.NewRedisService(newClients.Redis),
			RocketMQ: service.NewRocketMQService(newClients.RocketMQ),
		}

		// Clear cached data first
		m = m.clearCachedData()

		// Update profile, region and mode line AFTER clearing cache
		m.profile = msg.Profile
		m.region = cfg.RegionID // Reset to profile's default region
		m.header = m.header.SetProfile(msg.Profile).SetRegion(cfg.RegionID).SetTitle("Aliyun TUI Dashboard")
		m.modeLine = m.modeLine.SetProfile(msg.Profile).SetRegion(cfg.RegionID)

		// Update region service for new profile (cache is per-profile)
		m.regionService = service.NewRegionService(cfg.AccessKeyID, cfg.AccessKeySecret, msg.Profile)

		// Set page state
		m.currentPage = PageMenu
		m.previousPages = []PageType{}

		// Show success message
		m.modal = components.NewSuccessModal(fmt.Sprintf("Switched to profile: %s (region: %s)", msg.Profile, cfg.RegionID))
		return m, nil

	case ProfileSwitchedMsg:
		// Clear all cached data by resetting page models
		m = m.clearCachedData()
		// Update profile, header and mode line AFTER clearing cache
		m.profile = msg.ProfileName
		m.header = m.header.SetProfile(msg.ProfileName).SetTitle("Aliyun TUI Dashboard")
		m.modeLine = m.modeLine.SetProfile(msg.ProfileName)
		m.currentPage = PageMenu
		m.previousPages = []PageType{}

	case RegionsLoadedMsg:
		// Update modal with loaded regions
		m.regions = msg.Regions
		m.modal = m.modal.SetRegions(msg.Regions, m.region)
		return m, nil

	case components.RegionSelectedMsg:
		// Region switching - called when user selects a region from modal
		if msg.Region == m.region {
			// Same region, just dismiss modal
			m.modal = m.modal.Hide()
			return m, nil
		}

		// Recreate clients with new region using UpdateRegion
		newClients, err := m.clients.UpdateRegion(msg.Region)
		if err != nil {
			m.modal = components.NewErrorModal(fmt.Sprintf("Failed to create clients: %v", err))
			return m, nil
		}

		// Get config for OSS service
		clientCfg := newClients.GetConfig()

		// Clear cached data first
		m = m.clearCachedData()

		// Update region AFTER clearing cache
		m.region = msg.Region
		m.header = m.header.SetRegion(msg.Region).SetTitle("Aliyun TUI Dashboard")
		m.modeLine = m.modeLine.SetRegion(msg.Region)

		// Update clients and recreate services
		m.clients = newClients
		m.services = &Services{
			ECS:      service.NewECSService(newClients.ECS),
			DNS:      service.NewDNSService(newClients.DNS),
			SLB:      service.NewSLBService(newClients.SLB),
			RDS:      service.NewRDSService(newClients.RDS),
			OSS:      service.NewOSSServiceWithCredentials(newClients.OSS, clientCfg.AccessKeyID, clientCfg.AccessKeySecret, clientCfg.OssEndpoint),
			Redis:    service.NewRedisService(newClients.Redis),
			RocketMQ: service.NewRocketMQService(newClients.RocketMQ),
		}

		// Set page state
		m.currentPage = PageMenu
		m.previousPages = []PageType{}

		// Show success message
		m.modal = components.NewSuccessModal(fmt.Sprintf("Switched to region: %s", msg.Region))
		return m, nil

	case NavigateMsg:
		return m.navigateTo(msg.Page, msg.Data)

	case GoBackMsg:
		return m.navigateBack()

	// Handle search messages
	case components.SearchExecuteMsg:
		m.search = m.search.Deactivate()
		return m.handleSearchQuery(msg.Query)

	case components.SearchCancelMsg:
		m.search = m.search.Deactivate()
		return m, nil

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

	case ECSDisksLoadedMsg:
		m.loading = false
		m.ecsDiskPage = m.ecsDiskPage.SetData(msg.Disks)
		m.ecsDiskPage = m.ecsDiskPage.SetTitle(fmt.Sprintf("云盘 - 实例: %s", msg.InstanceId))
		m.ecsDiskPage = m.ecsDiskPage.SetSize(m.width, m.height-1)

	case ECSNetworkInterfacesLoadedMsg:
		m.loading = false
		m.ecsENIPage = m.ecsENIPage.SetData(msg.NetworkInterfaces)
		m.ecsENIPage = m.ecsENIPage.SetTitle(fmt.Sprintf("弹性网卡 - 实例: %s", msg.InstanceId))
		m.ecsENIPage = m.ecsENIPage.SetSize(m.width, m.height-1)

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

	// Handle pagination messages from pages package
	case pages.OSSObjectsLoadedMsg:
		m.loading = false
		m.ossObjectsPage = m.ossObjectsPage.SetData(msg.Result, msg.BucketName, msg.Page)
		m.ossObjectsPage = m.ossObjectsPage.SetSize(m.width, m.height-1)

	case pages.OSSErrorMsg:
		m.loading = false
		m.modal = components.NewErrorModal(msg.Err.Error())

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

	// Handle component messages (from table and viewport)
	case components.CopyDataMsg:
		return m, CopyToClipboard(msg.Data)

	case components.OpenEditorMsg:
		return m, OpenInEditor(msg.Data)

	case components.OpenPagerMsg:
		return m, OpenInPager(msg.Data)

	// Handle copy messages
	case CopiedMsg:
		m.modal = components.NewInfoModal("Copied to clipboard!")
	}

	// Forward non-key messages to modal if visible (for list filtering to work)
	if m.modal.Visible {
		var cmd tea.Cmd
		m.modal, cmd = m.modal.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update current page
	var cmd tea.Cmd
	m, cmd = m.updateCurrentPage(msg)
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
	case PageECSJSONDetail:
		content = m.ecsJSONDetailPage.View()
	case PageECSDisks:
		content = m.ecsDiskPage.View()
	case PageECSNetworkInterfaces:
		content = m.ecsENIPage.View()
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
		content = Center("Loading...", m.width, m.height-2)
	}

	// Build the full view: header + content + modeline
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		m.header.View(),
		content,
		m.modeLine.View(),
	)

	// Overlay search bar if active
	if m.search.Active {
		view = lipgloss.JoinVertical(
			lipgloss.Left,
			m.header.View(),
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

	// Update mode line and header
	m.modeLine = m.modeLine.SetPage(page)
	m.header = m.header.SetTitle(m.getPageTitle(page))

	var cmd tea.Cmd

	switch page {
	case PageECSList:
		m.ecsListPage = pages.NewECSListModel()
		cmd = LoadECSInstances(m.services.ECS)

	case PageECSDetail:
		// Try to get ecs.Instance for formatted detail view using the pages package function
		// This ensures type assertion happens in the same package where ecs.Instance is defined
		if detailModel, ok := pages.NewECSDetailModelFromInterface(data); ok {
			m.ecsDetailPage = detailModel
			m.ecsDetailPage = m.ecsDetailPage.SetSize(m.width, m.height-1)
			m.loading = false
		} else {
			// Fallback: if type assertion fails, navigate to JSON detail instead
			m.ecsJSONDetailPage = pages.NewDetailModel("ECS JSON Detail", data)
			m.ecsJSONDetailPage = m.ecsJSONDetailPage.SetSize(m.width, m.height-1)
			m.currentPage = PageECSJSONDetail
			m.loading = false
		}

	case PageECSJSONDetail:
		m.ecsJSONDetailPage = pages.NewDetailModel("ECS JSON Detail", data)
		m.ecsJSONDetailPage = m.ecsJSONDetailPage.SetSize(m.width, m.height-1)
		m.loading = false

	case PageECSDisks:
		if instanceId, ok := data.(string); ok {
			m.ecsDiskPage = pages.NewECSDiskModel(instanceId)
			cmd = LoadECSDisks(m.services.ECS, instanceId)
		}

	case PageECSNetworkInterfaces:
		if instanceId, ok := data.(string); ok {
			m.ecsENIPage = pages.NewECSENIModel(instanceId)
			cmd = LoadECSNetworkInterfaces(m.services.ECS, instanceId)
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
			m.slbDetailPage = m.slbDetailPage.SetSize(m.width, m.height-1)
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
			m.ossDetailPage = m.ossDetailPage.SetSize(m.width, m.height-1)
			m.loading = false
		}

	case PageRDSList:
		m.rdsListPage = pages.NewRDSListModel()
		cmd = LoadRDSInstances(m.services.RDS)

	case PageRDSDetail:
		if inst, ok := data.(interface{}); ok {
			m.rdsDetailPage = pages.NewDetailModel("RDS Detail", inst)
			m.rdsDetailPage = m.rdsDetailPage.SetSize(m.width, m.height-1)
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
			m.redisDetailPage = m.redisDetailPage.SetSize(m.width, m.height-1)
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
			m.rocketmqDetailPage = m.rocketmqDetailPage.SetSize(m.width, m.height-1)
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

	// Update mode line and header
	m.modeLine = m.modeLine.SetPage(prevPage)
	m.header = m.header.SetTitle(m.getPageTitle(prevPage))

	return m, nil
}

// getPageTitle returns the title for a given page type
func (m Model) getPageTitle(page PageType) string {
	switch page {
	case PageMenu:
		return "Aliyun TUI Dashboard"
	case PageECSList:
		return "ECS Instances"
	case PageECSDetail:
		return "ECS Detail"
	case PageECSJSONDetail:
		return "ECS JSON Detail"
	case PageECSDisks:
		return "ECS Disks"
	case PageECSNetworkInterfaces:
		return "ECS Network Interfaces"
	case PageSecurityGroups:
		return "Security Groups"
	case PageSecurityGroupRules:
		return "Security Group Rules"
	case PageSecurityGroupInstances:
		return "Security Group Instances"
	case PageInstanceSecurityGroups:
		return "Instance Security Groups"
	case PageDNSDomains:
		return "DNS Domains"
	case PageDNSRecords:
		return "DNS Records"
	case PageSLBList:
		return "SLB Instances"
	case PageSLBDetail:
		return "SLB Detail"
	case PageSLBListeners:
		return "SLB Listeners"
	case PageSLBVServerGroups:
		return "VServer Groups"
	case PageSLBBackendServers:
		return "Backend Servers"
	case PageOSSBuckets:
		return "OSS Buckets"
	case PageOSSObjects:
		return "OSS Objects"
	case PageOSSObjectDetail:
		return "OSS Object Detail"
	case PageRDSList:
		return "RDS Instances"
	case PageRDSDetail:
		return "RDS Detail"
	case PageRDSDatabases:
		return "RDS Databases"
	case PageRDSAccounts:
		return "RDS Accounts"
	case PageRedisList:
		return "Redis Instances"
	case PageRedisDetail:
		return "Redis Detail"
	case PageRedisAccounts:
		return "Redis Accounts"
	case PageRocketMQList:
		return "RocketMQ Instances"
	case PageRocketMQDetail:
		return "RocketMQ Detail"
	case PageRocketMQTopics:
		return "RocketMQ Topics"
	case PageRocketMQGroups:
		return "RocketMQ Groups"
	default:
		return "Aliyun TUI"
	}
}

// updateCurrentPage delegates update to the current page
func (m Model) updateCurrentPage(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.currentPage {
	case PageMenu:
		m.menuPage, cmd = m.menuPage.Update(msg)

	case PageECSList:
		m.ecsListPage, cmd = m.ecsListPage.Update(msg)

	case PageECSDetail:
		m.ecsDetailPage, cmd = m.ecsDetailPage.Update(msg)

	case PageECSJSONDetail:
		m.ecsJSONDetailPage, cmd = m.ecsJSONDetailPage.Update(msg)

	case PageECSDisks:
		m.ecsDiskPage, cmd = m.ecsDiskPage.Update(msg)

	case PageECSNetworkInterfaces:
		m.ecsENIPage, cmd = m.ecsENIPage.Update(msg)

	case PageSecurityGroups:
		m.sgListPage, cmd = m.sgListPage.Update(msg)

	case PageSecurityGroupRules:
		m.sgRulesPage, cmd = m.sgRulesPage.Update(msg)

	case PageSecurityGroupInstances:
		m.sgInstancesPage, cmd = m.sgInstancesPage.Update(msg)

	case PageInstanceSecurityGroups:
		m.instSGPage, cmd = m.instSGPage.Update(msg)

	case PageDNSDomains:
		m.dnsDomainsPage, cmd = m.dnsDomainsPage.Update(msg)

	case PageDNSRecords:
		m.dnsRecordsPage, cmd = m.dnsRecordsPage.Update(msg)

	case PageSLBList:
		m.slbListPage, cmd = m.slbListPage.Update(msg)

	case PageSLBDetail:
		m.slbDetailPage, cmd = m.slbDetailPage.Update(msg)

	case PageSLBListeners:
		m.slbListenersPage, cmd = m.slbListenersPage.Update(msg)

	case PageSLBVServerGroups:
		m.slbVServerPage, cmd = m.slbVServerPage.Update(msg)

	case PageSLBBackendServers:
		m.slbBackendPage, cmd = m.slbBackendPage.Update(msg)

	case PageOSSBuckets:
		m.ossBucketsPage, cmd = m.ossBucketsPage.Update(msg)

	case PageOSSObjects:
		m.ossObjectsPage, cmd = m.ossObjectsPage.Update(msg)

	case PageOSSObjectDetail:
		m.ossDetailPage, cmd = m.ossDetailPage.Update(msg)

	case PageRDSList:
		m.rdsListPage, cmd = m.rdsListPage.Update(msg)

	case PageRDSDetail:
		m.rdsDetailPage, cmd = m.rdsDetailPage.Update(msg)

	case PageRDSDatabases:
		m.rdsDatabasesPage, cmd = m.rdsDatabasesPage.Update(msg)

	case PageRDSAccounts:
		m.rdsAccountsPage, cmd = m.rdsAccountsPage.Update(msg)

	case PageRedisList:
		m.redisListPage, cmd = m.redisListPage.Update(msg)

	case PageRedisDetail:
		m.redisDetailPage, cmd = m.redisDetailPage.Update(msg)

	case PageRedisAccounts:
		m.redisAccountsPage, cmd = m.redisAccountsPage.Update(msg)

	case PageRocketMQList:
		m.rocketmqListPage, cmd = m.rocketmqListPage.Update(msg)

	case PageRocketMQDetail:
		m.rocketmqDetailPage, cmd = m.rocketmqDetailPage.Update(msg)

	case PageRocketMQTopics:
		m.rocketmqTopicsPage, cmd = m.rocketmqTopicsPage.Update(msg)

	case PageRocketMQGroups:
		m.rocketmqGroupsPage, cmd = m.rocketmqGroupsPage.Update(msg)
	}

	return m, cmd
}

// updateCurrentPageSize updates the current page's size
func (m Model) updateCurrentPageSize(height int) Model {
	switch m.currentPage {
	case PageMenu:
		m.menuPage = m.menuPage.SetSize(m.width, height)
	case PageECSList:
		m.ecsListPage = m.ecsListPage.SetSize(m.width, height)
	case PageECSDetail:
		m.ecsDetailPage = m.ecsDetailPage.SetSize(m.width, height)
	case PageECSJSONDetail:
		m.ecsJSONDetailPage = m.ecsJSONDetailPage.SetSize(m.width, height)
	case PageECSDisks:
		m.ecsDiskPage = m.ecsDiskPage.SetSize(m.width, height)
	case PageECSNetworkInterfaces:
		m.ecsENIPage = m.ecsENIPage.SetSize(m.width, height)
	case PageSecurityGroups:
		m.sgListPage = m.sgListPage.SetSize(m.width, height)
	case PageSecurityGroupRules:
		m.sgRulesPage = m.sgRulesPage.SetSize(m.width, height)
	case PageSecurityGroupInstances:
		m.sgInstancesPage = m.sgInstancesPage.SetSize(m.width, height)
	case PageInstanceSecurityGroups:
		m.instSGPage = m.instSGPage.SetSize(m.width, height)
	case PageDNSDomains:
		m.dnsDomainsPage = m.dnsDomainsPage.SetSize(m.width, height)
	case PageDNSRecords:
		m.dnsRecordsPage = m.dnsRecordsPage.SetSize(m.width, height)
	case PageSLBList:
		m.slbListPage = m.slbListPage.SetSize(m.width, height)
	case PageSLBDetail:
		m.slbDetailPage = m.slbDetailPage.SetSize(m.width, height)
	case PageSLBListeners:
		m.slbListenersPage = m.slbListenersPage.SetSize(m.width, height)
	case PageSLBVServerGroups:
		m.slbVServerPage = m.slbVServerPage.SetSize(m.width, height)
	case PageSLBBackendServers:
		m.slbBackendPage = m.slbBackendPage.SetSize(m.width, height)
	case PageOSSBuckets:
		m.ossBucketsPage = m.ossBucketsPage.SetSize(m.width, height)
	case PageOSSObjects:
		m.ossObjectsPage = m.ossObjectsPage.SetSize(m.width, height)
	case PageOSSObjectDetail:
		m.ossDetailPage = m.ossDetailPage.SetSize(m.width, height)
	case PageRDSList:
		m.rdsListPage = m.rdsListPage.SetSize(m.width, height)
	case PageRDSDetail:
		m.rdsDetailPage = m.rdsDetailPage.SetSize(m.width, height)
	case PageRDSDatabases:
		m.rdsDatabasesPage = m.rdsDatabasesPage.SetSize(m.width, height)
	case PageRDSAccounts:
		m.rdsAccountsPage = m.rdsAccountsPage.SetSize(m.width, height)
	case PageRedisList:
		m.redisListPage = m.redisListPage.SetSize(m.width, height)
	case PageRedisDetail:
		m.redisDetailPage = m.redisDetailPage.SetSize(m.width, height)
	case PageRedisAccounts:
		m.redisAccountsPage = m.redisAccountsPage.SetSize(m.width, height)
	case PageRocketMQList:
		m.rocketmqListPage = m.rocketmqListPage.SetSize(m.width, height)
	case PageRocketMQDetail:
		m.rocketmqDetailPage = m.rocketmqDetailPage.SetSize(m.width, height)
	case PageRocketMQTopics:
		m.rocketmqTopicsPage = m.rocketmqTopicsPage.SetSize(m.width, height)
	case PageRocketMQGroups:
		m.rocketmqGroupsPage = m.rocketmqGroupsPage.SetSize(m.width, height)
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
	if query == "" {
		return m, nil
	}

	// Route search to the current page's table or viewport
	switch m.currentPage {
	case PageECSList:
		m.ecsListPage = m.ecsListPage.Search(query)
	case PageECSDetail:
		m.ecsDetailPage = m.ecsDetailPage.Search(query)
	case PageECSJSONDetail:
		m.ecsJSONDetailPage = m.ecsJSONDetailPage.Search(query)
	case PageECSDisks:
		m.ecsDiskPage = m.ecsDiskPage.Search(query)
	case PageECSNetworkInterfaces:
		m.ecsENIPage = m.ecsENIPage.Search(query)
	case PageSecurityGroups:
		m.sgListPage = m.sgListPage.Search(query)
	case PageSecurityGroupRules:
		m.sgRulesPage = m.sgRulesPage.Search(query)
	case PageSecurityGroupInstances:
		m.sgInstancesPage = m.sgInstancesPage.Search(query)
	case PageInstanceSecurityGroups:
		m.instSGPage = m.instSGPage.Search(query)
	case PageDNSDomains:
		m.dnsDomainsPage = m.dnsDomainsPage.Search(query)
	case PageDNSRecords:
		m.dnsRecordsPage = m.dnsRecordsPage.Search(query)
	case PageSLBList:
		m.slbListPage = m.slbListPage.Search(query)
	case PageSLBDetail:
		m.slbDetailPage = m.slbDetailPage.Search(query)
	case PageSLBListeners:
		m.slbListenersPage = m.slbListenersPage.Search(query)
	case PageSLBVServerGroups:
		m.slbVServerPage = m.slbVServerPage.Search(query)
	case PageSLBBackendServers:
		m.slbBackendPage = m.slbBackendPage.Search(query)
	case PageOSSBuckets:
		m.ossBucketsPage = m.ossBucketsPage.Search(query)
	case PageOSSObjects:
		m.ossObjectsPage = m.ossObjectsPage.Search(query)
	case PageOSSObjectDetail:
		m.ossDetailPage = m.ossDetailPage.Search(query)
	case PageRDSList:
		m.rdsListPage = m.rdsListPage.Search(query)
	case PageRDSDetail:
		m.rdsDetailPage = m.rdsDetailPage.Search(query)
	case PageRDSDatabases:
		m.rdsDatabasesPage = m.rdsDatabasesPage.Search(query)
	case PageRDSAccounts:
		m.rdsAccountsPage = m.rdsAccountsPage.Search(query)
	case PageRedisList:
		m.redisListPage = m.redisListPage.Search(query)
	case PageRedisDetail:
		m.redisDetailPage = m.redisDetailPage.Search(query)
	case PageRedisAccounts:
		m.redisAccountsPage = m.redisAccountsPage.Search(query)
	case PageRocketMQList:
		m.rocketmqListPage = m.rocketmqListPage.Search(query)
	case PageRocketMQDetail:
		m.rocketmqDetailPage = m.rocketmqDetailPage.Search(query)
	case PageRocketMQTopics:
		m.rocketmqTopicsPage = m.rocketmqTopicsPage.Search(query)
	case PageRocketMQGroups:
		m.rocketmqGroupsPage = m.rocketmqGroupsPage.Search(query)
	}

	return m, nil
}

// handleSearchNext handles next search result
func (m Model) handleSearchNext() (Model, tea.Cmd) {
	switch m.currentPage {
	case PageECSList:
		m.ecsListPage = m.ecsListPage.NextSearchMatch()
	case PageECSDetail:
		m.ecsDetailPage = m.ecsDetailPage.NextSearchMatch()
	case PageECSJSONDetail:
		m.ecsJSONDetailPage = m.ecsJSONDetailPage.NextSearchMatch()
	case PageECSDisks:
		m.ecsDiskPage = m.ecsDiskPage.NextSearchMatch()
	case PageECSNetworkInterfaces:
		m.ecsENIPage = m.ecsENIPage.NextSearchMatch()
	case PageSecurityGroups:
		m.sgListPage = m.sgListPage.NextSearchMatch()
	case PageSecurityGroupRules:
		m.sgRulesPage = m.sgRulesPage.NextSearchMatch()
	case PageSecurityGroupInstances:
		m.sgInstancesPage = m.sgInstancesPage.NextSearchMatch()
	case PageInstanceSecurityGroups:
		m.instSGPage = m.instSGPage.NextSearchMatch()
	case PageDNSDomains:
		m.dnsDomainsPage = m.dnsDomainsPage.NextSearchMatch()
	case PageDNSRecords:
		m.dnsRecordsPage = m.dnsRecordsPage.NextSearchMatch()
	case PageSLBList:
		m.slbListPage = m.slbListPage.NextSearchMatch()
	case PageSLBDetail:
		m.slbDetailPage = m.slbDetailPage.NextSearchMatch()
	case PageSLBListeners:
		m.slbListenersPage = m.slbListenersPage.NextSearchMatch()
	case PageSLBVServerGroups:
		m.slbVServerPage = m.slbVServerPage.NextSearchMatch()
	case PageSLBBackendServers:
		m.slbBackendPage = m.slbBackendPage.NextSearchMatch()
	case PageOSSBuckets:
		m.ossBucketsPage = m.ossBucketsPage.NextSearchMatch()
	case PageOSSObjects:
		m.ossObjectsPage = m.ossObjectsPage.NextSearchMatch()
	case PageOSSObjectDetail:
		m.ossDetailPage = m.ossDetailPage.NextSearchMatch()
	case PageRDSList:
		m.rdsListPage = m.rdsListPage.NextSearchMatch()
	case PageRDSDetail:
		m.rdsDetailPage = m.rdsDetailPage.NextSearchMatch()
	case PageRDSDatabases:
		m.rdsDatabasesPage = m.rdsDatabasesPage.NextSearchMatch()
	case PageRDSAccounts:
		m.rdsAccountsPage = m.rdsAccountsPage.NextSearchMatch()
	case PageRedisList:
		m.redisListPage = m.redisListPage.NextSearchMatch()
	case PageRedisDetail:
		m.redisDetailPage = m.redisDetailPage.NextSearchMatch()
	case PageRedisAccounts:
		m.redisAccountsPage = m.redisAccountsPage.NextSearchMatch()
	case PageRocketMQList:
		m.rocketmqListPage = m.rocketmqListPage.NextSearchMatch()
	case PageRocketMQDetail:
		m.rocketmqDetailPage = m.rocketmqDetailPage.NextSearchMatch()
	case PageRocketMQTopics:
		m.rocketmqTopicsPage = m.rocketmqTopicsPage.NextSearchMatch()
	case PageRocketMQGroups:
		m.rocketmqGroupsPage = m.rocketmqGroupsPage.NextSearchMatch()
	}

	return m, nil
}

// handleSearchPrev handles previous search result
func (m Model) handleSearchPrev() (Model, tea.Cmd) {
	switch m.currentPage {
	case PageECSList:
		m.ecsListPage = m.ecsListPage.PrevSearchMatch()
	case PageECSDetail:
		m.ecsDetailPage = m.ecsDetailPage.PrevSearchMatch()
	case PageECSJSONDetail:
		m.ecsJSONDetailPage = m.ecsJSONDetailPage.PrevSearchMatch()
	case PageECSDisks:
		m.ecsDiskPage = m.ecsDiskPage.PrevSearchMatch()
	case PageECSNetworkInterfaces:
		m.ecsENIPage = m.ecsENIPage.PrevSearchMatch()
	case PageSecurityGroups:
		m.sgListPage = m.sgListPage.PrevSearchMatch()
	case PageSecurityGroupRules:
		m.sgRulesPage = m.sgRulesPage.PrevSearchMatch()
	case PageSecurityGroupInstances:
		m.sgInstancesPage = m.sgInstancesPage.PrevSearchMatch()
	case PageInstanceSecurityGroups:
		m.instSGPage = m.instSGPage.PrevSearchMatch()
	case PageDNSDomains:
		m.dnsDomainsPage = m.dnsDomainsPage.PrevSearchMatch()
	case PageDNSRecords:
		m.dnsRecordsPage = m.dnsRecordsPage.PrevSearchMatch()
	case PageSLBList:
		m.slbListPage = m.slbListPage.PrevSearchMatch()
	case PageSLBDetail:
		m.slbDetailPage = m.slbDetailPage.PrevSearchMatch()
	case PageSLBListeners:
		m.slbListenersPage = m.slbListenersPage.PrevSearchMatch()
	case PageSLBVServerGroups:
		m.slbVServerPage = m.slbVServerPage.PrevSearchMatch()
	case PageSLBBackendServers:
		m.slbBackendPage = m.slbBackendPage.PrevSearchMatch()
	case PageOSSBuckets:
		m.ossBucketsPage = m.ossBucketsPage.PrevSearchMatch()
	case PageOSSObjects:
		m.ossObjectsPage = m.ossObjectsPage.PrevSearchMatch()
	case PageOSSObjectDetail:
		m.ossDetailPage = m.ossDetailPage.PrevSearchMatch()
	case PageRDSList:
		m.rdsListPage = m.rdsListPage.PrevSearchMatch()
	case PageRDSDetail:
		m.rdsDetailPage = m.rdsDetailPage.PrevSearchMatch()
	case PageRDSDatabases:
		m.rdsDatabasesPage = m.rdsDatabasesPage.PrevSearchMatch()
	case PageRDSAccounts:
		m.rdsAccountsPage = m.rdsAccountsPage.PrevSearchMatch()
	case PageRedisList:
		m.redisListPage = m.redisListPage.PrevSearchMatch()
	case PageRedisDetail:
		m.redisDetailPage = m.redisDetailPage.PrevSearchMatch()
	case PageRedisAccounts:
		m.redisAccountsPage = m.redisAccountsPage.PrevSearchMatch()
	case PageRocketMQList:
		m.rocketmqListPage = m.rocketmqListPage.PrevSearchMatch()
	case PageRocketMQDetail:
		m.rocketmqDetailPage = m.rocketmqDetailPage.PrevSearchMatch()
	case PageRocketMQTopics:
		m.rocketmqTopicsPage = m.rocketmqTopicsPage.PrevSearchMatch()
	case PageRocketMQGroups:
		m.rocketmqGroupsPage = m.rocketmqGroupsPage.PrevSearchMatch()
	}

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

// loadRegions returns a command that loads regions with resources asynchronously
func (m *Model) loadRegions() tea.Cmd {
	return func() tea.Msg {
		regions, err := m.regionService.GetRegionsWithResources()
		if err != nil {
			return ErrorMsg{Err: fmt.Errorf("failed to load regions: %w", err)}
		}
		return RegionsLoadedMsg{
			Regions:       regions,
			CurrentRegion: m.region,
		}
	}
}

