package tui

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"aliyun-tui-viewer/internal/service"
	"aliyun-tui-viewer/internal/tui/types"
)

// Re-export PageType from types
type PageType = types.PageType

// Re-export page constants
const (
	PageMenu                   = types.PageMenu
	PageECSList                = types.PageECSList
	PageECSDetail              = types.PageECSDetail
	PageECSJSONDetail          = types.PageECSJSONDetail
	PageSecurityGroups         = types.PageSecurityGroups
	PageSecurityGroupRules     = types.PageSecurityGroupRules
	PageSecurityGroupInstances = types.PageSecurityGroupInstances
	PageInstanceSecurityGroups = types.PageInstanceSecurityGroups
	PageDNSDomains             = types.PageDNSDomains
	PageDNSRecords             = types.PageDNSRecords
	PageSLBList                = types.PageSLBList
	PageSLBDetail              = types.PageSLBDetail
	PageSLBListeners           = types.PageSLBListeners
	PageSLBVServerGroups       = types.PageSLBVServerGroups
	PageSLBBackendServers      = types.PageSLBBackendServers
	PageOSSBuckets             = types.PageOSSBuckets
	PageOSSObjects             = types.PageOSSObjects
	PageOSSObjectDetail        = types.PageOSSObjectDetail
	PageRDSList                = types.PageRDSList
	PageRDSDetail              = types.PageRDSDetail
	PageRDSDatabases           = types.PageRDSDatabases
	PageRDSAccounts            = types.PageRDSAccounts
	PageRedisList              = types.PageRedisList
	PageRedisDetail            = types.PageRedisDetail
	PageRedisAccounts          = types.PageRedisAccounts
	PageRocketMQList           = types.PageRocketMQList
	PageRocketMQDetail         = types.PageRocketMQDetail
	PageRocketMQTopics         = types.PageRocketMQTopics
	PageRocketMQGroups         = types.PageRocketMQGroups
)

// NavigateMsg requests navigation to a specific page
type NavigateMsg = types.NavigateMsg

// GoBackMsg requests navigation to the previous page  
type GoBackMsg = types.GoBackMsg

// --- Data Loading Messages ---

// LoadingMsg indicates data is being loaded
type LoadingMsg struct {
	Page PageType
}

// ErrorMsg indicates an error occurred
type ErrorMsg struct {
	Err error
}

// --- ECS Messages ---

// ECSInstancesLoadedMsg contains loaded ECS instances
type ECSInstancesLoadedMsg struct {
	Instances []ecs.Instance
}

// ECSInstanceSelectedMsg indicates an ECS instance was selected
type ECSInstanceSelectedMsg struct {
	Instance ecs.Instance
}

// --- Security Groups Messages ---

// SecurityGroupsLoadedMsg contains loaded security groups
type SecurityGroupsLoadedMsg struct {
	SecurityGroups []ecs.SecurityGroup
}

// SecurityGroupRulesLoadedMsg contains loaded security group rules
type SecurityGroupRulesLoadedMsg struct {
	Response *ecs.DescribeSecurityGroupAttributeResponse
}

// SecurityGroupInstancesLoadedMsg contains instances for a security group
type SecurityGroupInstancesLoadedMsg struct {
	Instances       []ecs.Instance
	SecurityGroupId string
}

// InstanceSecurityGroupsLoadedMsg contains security groups for an instance
type InstanceSecurityGroupsLoadedMsg struct {
	SecurityGroups []ecs.SecurityGroup
	InstanceId     string
}

// --- DNS Messages ---

// DNSDomainsLoadedMsg contains loaded DNS domains
type DNSDomainsLoadedMsg struct {
	Domains []alidns.DomainInDescribeDomains
}

// DNSRecordsLoadedMsg contains loaded DNS records
type DNSRecordsLoadedMsg struct {
	Records    []alidns.Record
	DomainName string
}

// --- SLB Messages ---

// SLBInstancesLoadedMsg contains loaded SLB instances
type SLBInstancesLoadedMsg struct {
	LoadBalancers []slb.LoadBalancer
}

// SLBListenersLoadedMsg contains loaded SLB listeners
type SLBListenersLoadedMsg struct {
	Listeners      []service.ListenerDetail
	LoadBalancerId string
}

// SLBVServerGroupsLoadedMsg contains loaded SLB VServer groups
type SLBVServerGroupsLoadedMsg struct {
	VServerGroups  []service.VServerGroupDetail
	LoadBalancerId string
}

// SLBBackendServersLoadedMsg contains loaded backend servers
type SLBBackendServersLoadedMsg struct {
	BackendServers []service.BackendServerDetail
	VServerGroupId string
}

// --- OSS Messages ---

// OSSBucketsLoadedMsg contains loaded OSS buckets
type OSSBucketsLoadedMsg struct {
	Buckets []oss.BucketProperties
}

// OSSObjectsLoadedMsg contains loaded OSS objects with pagination
type OSSObjectsLoadedMsg struct {
	Result     *service.ObjectListResult
	BucketName string
	Page       int
}

// OSSObjectSelectedMsg indicates an OSS object was selected
type OSSObjectSelectedMsg struct {
	Object oss.ObjectProperties
}

// --- RDS Messages ---

// RDSInstancesLoadedMsg contains loaded RDS instances
type RDSInstancesLoadedMsg struct {
	Instances []rds.DBInstance
}

// RDSDatabasesLoadedMsg contains loaded RDS databases
type RDSDatabasesLoadedMsg struct {
	Databases  []rds.Database
	InstanceId string
}

// RDSAccountsLoadedMsg contains loaded RDS accounts
type RDSAccountsLoadedMsg struct {
	Accounts   []rds.DBInstanceAccount
	InstanceId string
}

// --- Redis Messages ---

// RedisInstancesLoadedMsg contains loaded Redis instances
type RedisInstancesLoadedMsg struct {
	Instances []r_kvstore.KVStoreInstance
}

// RedisAccountsLoadedMsg contains loaded Redis accounts
type RedisAccountsLoadedMsg struct {
	Accounts   []r_kvstore.Account
	InstanceId string
}

// --- RocketMQ Messages ---

// RocketMQInstancesLoadedMsg contains loaded RocketMQ instances
type RocketMQInstancesLoadedMsg struct {
	Instances []service.RocketMQInstance
}

// RocketMQTopicsLoadedMsg contains loaded RocketMQ topics
type RocketMQTopicsLoadedMsg struct {
	Topics     []service.RocketMQTopic
	InstanceId string
}

// RocketMQGroupsLoadedMsg contains loaded RocketMQ groups
type RocketMQGroupsLoadedMsg struct {
	Groups     []service.RocketMQGroup
	InstanceId string
}

// --- Search Messages ---

// SearchStartMsg indicates search mode should start
type SearchStartMsg struct{}

// SearchQueryMsg contains a search query
type SearchQueryMsg struct {
	Query string
}

// SearchExitMsg indicates search mode should end
type SearchExitMsg struct{}

// SearchNextMsg requests navigation to next search result
type SearchNextMsg struct{}

// SearchPrevMsg requests navigation to previous search result
type SearchPrevMsg struct{}

// --- Action Messages ---

// CopyToClipboardMsg requests copying data to clipboard
type CopyToClipboardMsg struct {
	Data interface{}
}

// CopiedMsg indicates data was copied successfully
type CopiedMsg struct{}

// OpenEditorMsg requests opening data in external editor
type OpenEditorMsg struct {
	Data interface{}
}

// OpenPagerMsg requests opening data in external pager
type OpenPagerMsg struct {
	Data interface{}
}

// EditorClosedMsg indicates the external editor was closed
type EditorClosedMsg struct{}

// --- Profile Messages ---

// ProfileSwitchMsg requests switching to a different profile
type ProfileSwitchMsg struct {
	ProfileName string
}

// ProfileSwitchedMsg indicates profile was switched successfully
type ProfileSwitchedMsg struct {
	ProfileName string
}

// ProfileListMsg requests showing the profile selection dialog
type ProfileListMsg struct{}

// ProfileListLoadedMsg contains the list of available profiles
type ProfileListLoadedMsg struct {
	Profiles       []string
	CurrentProfile string
}

// --- Region Messages ---

// RegionsLoadedMsg contains the list of available regions with resources
type RegionsLoadedMsg struct {
	Regions       []string
	CurrentRegion string
}

// RegionSwitchedMsg indicates region was switched successfully
type RegionSwitchedMsg struct {
	RegionID string
}

// --- Modal Messages ---

// ShowModalMsg requests showing a modal dialog
type ShowModalMsg struct {
	Title   string
	Message string
	Type    ModalType
}

// ModalDismissedMsg indicates the modal was dismissed
type ModalDismissedMsg struct{}

// ModalType represents different types of modals
type ModalType int

const (
	ModalInfo ModalType = iota
	ModalError
	ModalSuccess
	ModalConfirm
)

// --- Window Messages ---

// WindowSizeMsg contains the terminal window size
type WindowSizeMsg struct {
	Width  int
	Height int
}

// --- Refresh Messages ---

// RefreshDataMsg requests refreshing the current page's data
type RefreshDataMsg struct{}

// --- Yank (Copy) Messages ---

// YankKeyMsg indicates 'y' key was pressed
type YankKeyMsg struct{}

// YankResetMsg resets the yank tracker
type YankResetMsg struct{}
