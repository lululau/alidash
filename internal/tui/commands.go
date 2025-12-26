package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"aliyun-tui-viewer/internal/service"
)

// Services holds all service instances for data fetching
type Services struct {
	ECS      *service.ECSService
	DNS      *service.DNSService
	SLB      *service.SLBService
	RDS      *service.RDSService
	OSS      *service.OSSService
	Redis    *service.RedisService
	RocketMQ *service.RocketMQService
}

// --- ECS Commands ---

// LoadECSInstances creates a command to load ECS instances
func LoadECSInstances(svc *service.ECSService) tea.Cmd {
	return func() tea.Msg {
		instances, err := svc.FetchInstances()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ECSInstancesLoadedMsg{Instances: instances}
	}
}

// LoadSecurityGroups creates a command to load security groups
func LoadSecurityGroups(svc *service.ECSService) tea.Cmd {
	return func() tea.Msg {
		groups, err := svc.FetchSecurityGroups()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SecurityGroupsLoadedMsg{SecurityGroups: groups}
	}
}

// LoadSecurityGroupRules creates a command to load security group rules
func LoadSecurityGroupRules(svc *service.ECSService, securityGroupId string) tea.Cmd {
	return func() tea.Msg {
		response, err := svc.FetchSecurityGroupRules(securityGroupId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SecurityGroupRulesLoadedMsg{Response: response}
	}
}

// LoadSecurityGroupInstances creates a command to load instances for a security group
func LoadSecurityGroupInstances(svc *service.ECSService, securityGroupId string) tea.Cmd {
	return func() tea.Msg {
		instances, err := svc.FetchInstancesBySecurityGroup(securityGroupId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SecurityGroupInstancesLoadedMsg{
			Instances:       instances,
			SecurityGroupId: securityGroupId,
		}
	}
}

// LoadInstanceSecurityGroups creates a command to load security groups for an instance
func LoadInstanceSecurityGroups(svc *service.ECSService, instanceId string) tea.Cmd {
	return func() tea.Msg {
		groups, err := svc.FetchSecurityGroupsByInstance(instanceId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return InstanceSecurityGroupsLoadedMsg{
			SecurityGroups: groups,
			InstanceId:     instanceId,
		}
	}
}

// LoadECSDisks creates a command to load disks for an instance
func LoadECSDisks(svc *service.ECSService, instanceId string) tea.Cmd {
	return func() tea.Msg {
		disks, err := svc.FetchDisks(instanceId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ECSDisksLoadedMsg{
			Disks:      disks,
			InstanceId: instanceId,
		}
	}
}

// LoadECSNetworkInterfaces creates a command to load network interfaces for an instance
func LoadECSNetworkInterfaces(svc *service.ECSService, instanceId string) tea.Cmd {
	return func() tea.Msg {
		enis, err := svc.FetchNetworkInterfaces(instanceId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ECSNetworkInterfacesLoadedMsg{
			NetworkInterfaces: enis,
			InstanceId:        instanceId,
		}
	}
}

// --- DNS Commands ---

// LoadDNSDomains creates a command to load DNS domains
func LoadDNSDomains(svc *service.DNSService) tea.Cmd {
	return func() tea.Msg {
		domains, err := svc.FetchDomains()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return DNSDomainsLoadedMsg{Domains: domains}
	}
}

// LoadDNSRecords creates a command to load DNS records for a domain
func LoadDNSRecords(svc *service.DNSService, domainName string) tea.Cmd {
	return func() tea.Msg {
		records, err := svc.FetchDomainRecords(domainName)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return DNSRecordsLoadedMsg{
			Records:    records,
			DomainName: domainName,
		}
	}
}

// --- SLB Commands ---

// LoadSLBInstances creates a command to load SLB instances
func LoadSLBInstances(svc *service.SLBService) tea.Cmd {
	return func() tea.Msg {
		lbs, err := svc.FetchInstances()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SLBInstancesLoadedMsg{LoadBalancers: lbs}
	}
}

// LoadSLBListeners creates a command to load SLB listeners
func LoadSLBListeners(svc *service.SLBService, loadBalancerId string) tea.Cmd {
	return func() tea.Msg {
		listeners, err := svc.FetchDetailedListeners(loadBalancerId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SLBListenersLoadedMsg{
			Listeners:      listeners,
			LoadBalancerId: loadBalancerId,
		}
	}
}

// LoadSLBVServerGroups creates a command to load SLB VServer groups
func LoadSLBVServerGroups(svc *service.SLBService, loadBalancerId string) tea.Cmd {
	return func() tea.Msg {
		groups, err := svc.FetchDetailedVServerGroups(loadBalancerId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SLBVServerGroupsLoadedMsg{
			VServerGroups:  groups,
			LoadBalancerId: loadBalancerId,
		}
	}
}

// LoadSLBBackendServers creates a command to load backend servers
func LoadSLBBackendServers(svc *service.SLBService, vServerGroupId string, ecsClient *ecs.Client) tea.Cmd {
	return func() tea.Msg {
		servers, err := svc.FetchDetailedBackendServers(vServerGroupId, ecsClient)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SLBBackendServersLoadedMsg{
			BackendServers: servers,
			VServerGroupId: vServerGroupId,
		}
	}
}

// LoadSLBForwardingRules creates a command to load forwarding rules for a listener
func LoadSLBForwardingRules(svc *service.SLBService, loadBalancerId string, listenerPort int, listenerProtocol string) tea.Cmd {
	return func() tea.Msg {
		rules, err := svc.FetchForwardingRules(loadBalancerId, listenerPort, listenerProtocol)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SLBForwardingRulesLoadedMsg{
			Rules:            rules,
			LoadBalancerId:   loadBalancerId,
			ListenerPort:     listenerPort,
			ListenerProtocol: listenerProtocol,
		}
	}
}

// LoadSLBDefaultServers creates a command to load default backend servers
func LoadSLBDefaultServers(svc *service.SLBService, loadBalancerId string, ecsClient *ecs.Client) tea.Cmd {
	return func() tea.Msg {
		servers, err := svc.FetchDefaultBackendServers(loadBalancerId, ecsClient)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return SLBDefaultServersLoadedMsg{
			Servers:        servers,
			LoadBalancerId: loadBalancerId,
		}
	}
}

// --- OSS Commands ---

// LoadOSSBuckets creates a command to load OSS buckets
func LoadOSSBuckets(svc *service.OSSService) tea.Cmd {
	return func() tea.Msg {
		buckets, err := svc.FetchBuckets()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return OSSBucketsLoadedMsg{Buckets: buckets}
	}
}

// LoadOSSObjects creates a command to load OSS objects with pagination
func LoadOSSObjects(svc *service.OSSService, bucketName, marker string, pageSize, page int) tea.Cmd {
	return func() tea.Msg {
		result, err := svc.FetchObjects(bucketName, marker, pageSize)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return OSSObjectsLoadedMsg{
			Result:     result,
			BucketName: bucketName,
			Page:       page,
		}
	}
}

// --- RDS Commands ---

// LoadRDSInstances creates a command to load RDS instances
func LoadRDSInstances(svc *service.RDSService) tea.Cmd {
	return func() tea.Msg {
		instances, err := svc.FetchInstances()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RDSInstancesLoadedMsg{Instances: instances}
	}
}

// LoadRDSDetailedInstances creates a command to load RDS instances with network info
func LoadRDSDetailedInstances(svc *service.RDSService) tea.Cmd {
	return func() tea.Msg {
		instances, err := svc.FetchDetailedInstances()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RDSDetailedInstancesLoadedMsg{Instances: instances}
	}
}

// LoadRDSDatabases creates a command to load RDS databases
func LoadRDSDatabases(svc *service.RDSService, instanceId string) tea.Cmd {
	return func() tea.Msg {
		databases, err := svc.FetchDatabases(instanceId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RDSDatabasesLoadedMsg{
			Databases:  databases,
			InstanceId: instanceId,
		}
	}
}

// LoadRDSAccounts creates a command to load RDS accounts
func LoadRDSAccounts(svc *service.RDSService, instanceId string) tea.Cmd {
	return func() tea.Msg {
		accounts, err := svc.FetchAccounts(instanceId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RDSAccountsLoadedMsg{
			Accounts:   accounts,
			InstanceId: instanceId,
		}
	}
}

// --- Redis Commands ---

// LoadRedisInstances creates a command to load Redis instances
func LoadRedisInstances(svc *service.RedisService) tea.Cmd {
	return func() tea.Msg {
		instances, err := svc.FetchInstances()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RedisInstancesLoadedMsg{Instances: instances}
	}
}

// LoadRedisAccounts creates a command to load Redis accounts
func LoadRedisAccounts(svc *service.RedisService, instanceId string) tea.Cmd {
	return func() tea.Msg {
		accounts, err := svc.FetchAccounts(instanceId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RedisAccountsLoadedMsg{
			Accounts:   accounts,
			InstanceId: instanceId,
		}
	}
}

// --- RocketMQ Commands ---

// LoadRocketMQInstances creates a command to load RocketMQ instances
func LoadRocketMQInstances(svc *service.RocketMQService) tea.Cmd {
	return func() tea.Msg {
		instances, err := svc.FetchInstances()
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RocketMQInstancesLoadedMsg{Instances: instances}
	}
}

// LoadRocketMQTopics creates a command to load RocketMQ topics
func LoadRocketMQTopics(svc *service.RocketMQService, instanceId string) tea.Cmd {
	return func() tea.Msg {
		topics, err := svc.FetchTopics(instanceId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RocketMQTopicsLoadedMsg{
			Topics:     topics,
			InstanceId: instanceId,
		}
	}
}

// LoadRocketMQGroups creates a command to load RocketMQ groups
func LoadRocketMQGroups(svc *service.RocketMQService, instanceId string) tea.Cmd {
	return func() tea.Msg {
		groups, err := svc.FetchGroups(instanceId)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return RocketMQGroupsLoadedMsg{
			Groups:     groups,
			InstanceId: instanceId,
		}
	}
}

// --- Resource Finder Commands ---

// FindResources creates a command to find resources by IP or domain
func FindResources(svc *service.FinderService, query string) tea.Cmd {
	return func() tea.Msg {
		// Resolve the query to IPs
		ips, domain, err := svc.ResolveToIPs(query)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		// Find matching resources
		result, err := svc.FindResources(ips, domain)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return FindResourceResultMsg{Result: result}
	}
}

