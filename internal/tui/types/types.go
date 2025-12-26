package types

// PageType represents different pages in the application
type PageType int

const (
	PageMenu PageType = iota
	PageECSList
	PageECSDetail     // New formatted detail view
	PageECSJSONDetail // JSON detail view (previously PageECSDetail)
	PageECSDisks      // ECS Disk/Storage page
	PageSecurityGroups
	PageSecurityGroupRules
	PageSecurityGroupInstances
	PageInstanceSecurityGroups
	PageDNSDomains
	PageDNSRecords
	PageSLBList
	PageSLBDetail
	PageSLBListeners
	PageSLBVServerGroups
	PageSLBBackendServers
	PageOSSBuckets
	PageOSSObjects
	PageOSSObjectDetail
	PageRDSList
	PageRDSDetail
	PageRDSDatabases
	PageRDSAccounts
	PageRedisList
	PageRedisDetail
	PageRedisAccounts
	PageRocketMQList
	PageRocketMQDetail
	PageRocketMQTopics
	PageRocketMQGroups
)

// String returns the string representation of PageType
func (p PageType) String() string {
	switch p {
	case PageMenu:
		return "Main Menu"
	case PageECSList:
		return "ECS Instances"
	case PageECSDetail:
		return "ECS Detail"
	case PageECSJSONDetail:
		return "ECS JSON Detail"
	case PageECSDisks:
		return "ECS Disks"
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
		return "SLB VServer Groups"
	case PageSLBBackendServers:
		return "SLB Backend Servers"
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
		return "Unknown"
	}
}

// --- Navigation Messages ---

// NavigateMsg requests navigation to a specific page
type NavigateMsg struct {
	Page PageType
	Data interface{} // Optional data to pass to the page
}

// GoBackMsg requests navigation to the previous page
type GoBackMsg struct{}

