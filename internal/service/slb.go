package service

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

// SLBService handles SLB operations
type SLBService struct {
	client *slb.Client
}

// NewSLBService creates a new SLB service
func NewSLBService(client *slb.Client) *SLBService {
	return &SLBService{client: client}
}

// FetchInstances retrieves all SLB instances using pagination
func (s *SLBService) FetchInstances() ([]slb.LoadBalancer, error) {
	var allLoadBalancers []slb.LoadBalancer
	pageNumber := int64(1)
	pageSize := int64(100)

	for {
		request := slb.CreateDescribeLoadBalancersRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(int(pageNumber))
		request.PageSize = requests.NewInteger(int(pageSize))

		response, err := s.client.DescribeLoadBalancers(request)
		if err != nil {
			return nil, fmt.Errorf("describing SLB instances (page %d): %w", pageNumber, err)
		}

		allLoadBalancers = append(allLoadBalancers, response.LoadBalancers.LoadBalancer...)

		if pageNumber*pageSize >= int64(response.TotalCount) {
			break
		}

		if len(response.LoadBalancers.LoadBalancer) < int(pageSize) {
			break
		}

		pageNumber++
	}
	return allLoadBalancers, nil
}

// ListenerDetail contains detailed information about a listener
type ListenerDetail struct {
	Protocol         string
	Port             int
	BackendPort      int
	Status           string
	HealthCheck      string
	Scheduler        string
	VServerGroupId   string
	VServerGroupName string
}

// FetchListeners retrieves all listeners for a specific SLB instance
func (s *SLBService) FetchListeners(loadBalancerId string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.Scheme = "https"
	request.LoadBalancerId = loadBalancerId

	response, err := s.client.DescribeLoadBalancerAttribute(request)
	if err != nil {
		return nil, fmt.Errorf("describing listeners for SLB %s: %w", loadBalancerId, err)
	}

	return response, nil
}

// FetchDetailedListeners retrieves detailed information for all listeners of an SLB instance
func (s *SLBService) FetchDetailedListeners(loadBalancerId string) ([]ListenerDetail, error) {
	// First get the basic listener info
	basicResponse, err := s.FetchListeners(loadBalancerId)
	if err != nil {
		return nil, err
	}

	var detailedListeners []ListenerDetail

	// For each listener port, try to get detailed information
	for _, port := range basicResponse.ListenerPorts.ListenerPort {
		// Try HTTP listener first
		if httpDetail := s.fetchHTTPListenerDetail(loadBalancerId, port); httpDetail != nil {
			detailedListeners = append(detailedListeners, *httpDetail)
			continue
		}

		// Try HTTPS listener
		if httpsDetail := s.fetchHTTPSListenerDetail(loadBalancerId, port); httpsDetail != nil {
			detailedListeners = append(detailedListeners, *httpsDetail)
			continue
		}

		// Try TCP listener
		if tcpDetail := s.fetchTCPListenerDetail(loadBalancerId, port); tcpDetail != nil {
			detailedListeners = append(detailedListeners, *tcpDetail)
			continue
		}

		// Try UDP listener
		if udpDetail := s.fetchUDPListenerDetail(loadBalancerId, port); udpDetail != nil {
			detailedListeners = append(detailedListeners, *udpDetail)
			continue
		}

		// If no specific listener type found, create a basic entry
		detailedListeners = append(detailedListeners, ListenerDetail{
			Protocol:    "Unknown",
			Port:        port,
			BackendPort: 0,
			Status:      "Unknown",
			HealthCheck: "Unknown",
			Scheduler:   "Unknown",
		})
	}

	return detailedListeners, nil
}

// fetchHTTPListenerDetail tries to fetch HTTP listener details
func (s *SLBService) fetchHTTPListenerDetail(loadBalancerId string, port int) *ListenerDetail {
	request := slb.CreateDescribeLoadBalancerHTTPListenerAttributeRequest()
	request.Scheme = "https"
	request.LoadBalancerId = loadBalancerId
	request.ListenerPort = requests.NewInteger(port)

	response, err := s.client.DescribeLoadBalancerHTTPListenerAttribute(request)
	if err != nil {
		return nil
	}

	// Get VServer group name if available
	vsgName := ""
	if response.VServerGroupId != "" {
		if vsg, err := s.getVServerGroupName(response.VServerGroupId); err == nil {
			vsgName = vsg
		}
	}

	return &ListenerDetail{
		Protocol:         "HTTP",
		Port:             port,
		BackendPort:      response.BackendServerPort,
		Status:           response.Status,
		HealthCheck:      response.HealthCheck,
		Scheduler:        response.Scheduler,
		VServerGroupId:   response.VServerGroupId,
		VServerGroupName: vsgName,
	}
}

// fetchHTTPSListenerDetail tries to fetch HTTPS listener details
func (s *SLBService) fetchHTTPSListenerDetail(loadBalancerId string, port int) *ListenerDetail {
	request := slb.CreateDescribeLoadBalancerHTTPSListenerAttributeRequest()
	request.Scheme = "https"
	request.LoadBalancerId = loadBalancerId
	request.ListenerPort = requests.NewInteger(port)

	response, err := s.client.DescribeLoadBalancerHTTPSListenerAttribute(request)
	if err != nil {
		return nil
	}

	// Get VServer group name if available
	vsgName := ""
	if response.VServerGroupId != "" {
		if vsg, err := s.getVServerGroupName(response.VServerGroupId); err == nil {
			vsgName = vsg
		}
	}

	return &ListenerDetail{
		Protocol:         "HTTPS",
		Port:             port,
		BackendPort:      response.BackendServerPort,
		Status:           response.Status,
		HealthCheck:      response.HealthCheck,
		Scheduler:        response.Scheduler,
		VServerGroupId:   response.VServerGroupId,
		VServerGroupName: vsgName,
	}
}

// fetchTCPListenerDetail tries to fetch TCP listener details
func (s *SLBService) fetchTCPListenerDetail(loadBalancerId string, port int) *ListenerDetail {
	request := slb.CreateDescribeLoadBalancerTCPListenerAttributeRequest()
	request.Scheme = "https"
	request.LoadBalancerId = loadBalancerId
	request.ListenerPort = requests.NewInteger(port)

	response, err := s.client.DescribeLoadBalancerTCPListenerAttribute(request)
	if err != nil {
		return nil
	}

	// Get VServer group name if available
	vsgName := ""
	if response.VServerGroupId != "" {
		if vsg, err := s.getVServerGroupName(response.VServerGroupId); err == nil {
			vsgName = vsg
		}
	}

	return &ListenerDetail{
		Protocol:         "TCP",
		Port:             port,
		BackendPort:      response.BackendServerPort,
		Status:           response.Status,
		HealthCheck:      response.HealthCheck,
		Scheduler:        response.Scheduler,
		VServerGroupId:   response.VServerGroupId,
		VServerGroupName: vsgName,
	}
}

// fetchUDPListenerDetail tries to fetch UDP listener details
func (s *SLBService) fetchUDPListenerDetail(loadBalancerId string, port int) *ListenerDetail {
	request := slb.CreateDescribeLoadBalancerUDPListenerAttributeRequest()
	request.Scheme = "https"
	request.LoadBalancerId = loadBalancerId
	request.ListenerPort = requests.NewInteger(port)

	response, err := s.client.DescribeLoadBalancerUDPListenerAttribute(request)
	if err != nil {
		return nil
	}

	// Get VServer group name if available
	vsgName := ""
	if response.VServerGroupId != "" {
		if vsg, err := s.getVServerGroupName(response.VServerGroupId); err == nil {
			vsgName = vsg
		}
	}

	return &ListenerDetail{
		Protocol:         "UDP",
		Port:             port,
		BackendPort:      response.BackendServerPort,
		Status:           response.Status,
		HealthCheck:      response.HealthCheck,
		Scheduler:        response.Scheduler,
		VServerGroupId:   response.VServerGroupId,
		VServerGroupName: vsgName,
	}
}

// getVServerGroupName retrieves the name of a virtual server group
func (s *SLBService) getVServerGroupName(vServerGroupId string) (string, error) {
	request := slb.CreateDescribeVServerGroupAttributeRequest()
	request.Scheme = "https"
	request.VServerGroupId = vServerGroupId

	response, err := s.client.DescribeVServerGroupAttribute(request)
	if err != nil {
		return "", err
	}

	return response.VServerGroupName, nil
}

// VServerGroupDetail contains detailed information about a virtual server group
type VServerGroupDetail struct {
	VServerGroupId      string
	VServerGroupName    string
	BackendServerCount  int
	AssociatedListeners []string // List of listener ports that use this VServer group
	AssociatedRules     []string // List of forwarding rule IDs that use this VServer group
}

// FetchVServerGroups retrieves all virtual server groups for a specific SLB instance
func (s *SLBService) FetchVServerGroups(loadBalancerId string) ([]slb.VServerGroup, error) {
	request := slb.CreateDescribeVServerGroupsRequest()
	request.Scheme = "https"
	request.LoadBalancerId = loadBalancerId

	response, err := s.client.DescribeVServerGroups(request)
	if err != nil {
		return nil, fmt.Errorf("describing virtual server groups for SLB %s: %w", loadBalancerId, err)
	}

	return response.VServerGroups.VServerGroup, nil
}

// FetchDetailedVServerGroups retrieves detailed information for all virtual server groups
func (s *SLBService) FetchDetailedVServerGroups(loadBalancerId string) ([]VServerGroupDetail, error) {
	// Get basic VServer groups
	vServerGroups, err := s.FetchVServerGroups(loadBalancerId)
	if err != nil {
		return nil, err
	}

	// Get detailed listeners to find associations
	listeners, err := s.FetchDetailedListeners(loadBalancerId)
	if err != nil {
		return nil, err
	}

	// Build a map of VServerGroupId -> associated rule IDs
	vsgRulesMap := make(map[string][]string)

	// Fetch forwarding rules for each HTTP/HTTPS listener
	for _, listener := range listeners {
		if listener.Protocol == "HTTP" || listener.Protocol == "HTTPS" {
			rules, err := s.FetchForwardingRules(loadBalancerId, listener.Port, listener.Protocol)
			if err != nil {
				continue // Skip if we can't fetch rules
			}
			for _, rule := range rules {
				if rule.VServerGroupId != "" {
					vsgRulesMap[rule.VServerGroupId] = append(vsgRulesMap[rule.VServerGroupId], rule.RuleId)
				}
			}
		}
	}

	var detailedVServerGroups []VServerGroupDetail

	for _, vsg := range vServerGroups {
		// Get backend server count
		backendServers, err := s.FetchVServerGroupBackendServers(vsg.VServerGroupId)
		if err != nil {
			// If we can't get backend servers, set count to 0
			backendServers = []slb.BackendServerInDescribeVServerGroupAttribute{}
		}

		// Find associated listeners
		var associatedListeners []string
		for _, listener := range listeners {
			if listener.VServerGroupId == vsg.VServerGroupId {
				associatedListeners = append(associatedListeners, fmt.Sprintf("%s:%d", listener.Protocol, listener.Port))
			}
		}

		// Get associated rules
		associatedRules := vsgRulesMap[vsg.VServerGroupId]

		detailedVServerGroups = append(detailedVServerGroups, VServerGroupDetail{
			VServerGroupId:      vsg.VServerGroupId,
			VServerGroupName:    vsg.VServerGroupName,
			BackendServerCount:  len(backendServers),
			AssociatedListeners: associatedListeners,
			AssociatedRules:     associatedRules,
		})
	}

	return detailedVServerGroups, nil
}

// BackendServerDetail contains detailed information about a backend server
type BackendServerDetail struct {
	ServerId         string
	Port             int
	Weight           int
	Type             string
	Description      string
	InstanceName     string
	PrivateIpAddress string
	PublicIpAddress  string
}

// FetchVServerGroupBackendServers retrieves backend servers for a specific virtual server group
func (s *SLBService) FetchVServerGroupBackendServers(vServerGroupId string) ([]slb.BackendServerInDescribeVServerGroupAttribute, error) {
	request := slb.CreateDescribeVServerGroupAttributeRequest()
	request.Scheme = "https"
	request.VServerGroupId = vServerGroupId

	response, err := s.client.DescribeVServerGroupAttribute(request)
	if err != nil {
		return nil, fmt.Errorf("describing backend servers for virtual server group %s: %w", vServerGroupId, err)
	}

	return response.BackendServers.BackendServer, nil
}

// FetchDetailedBackendServers retrieves detailed information for backend servers including ECS details
func (s *SLBService) FetchDetailedBackendServers(vServerGroupId string, ecsClient *ecs.Client) ([]BackendServerDetail, error) {
	// Get basic backend servers
	backendServers, err := s.FetchVServerGroupBackendServers(vServerGroupId)
	if err != nil {
		return nil, err
	}

	var detailedServers []BackendServerDetail

	for _, server := range backendServers {
		detail := BackendServerDetail{
			ServerId:         server.ServerId,
			Port:             server.Port,
			Weight:           server.Weight,
			Type:             server.Type,
			Description:      server.Description,
			InstanceName:     "N/A",
			PrivateIpAddress: "N/A",
			PublicIpAddress:  "N/A",
		}

		// Try to get ECS instance details if ecsClient is provided
		if ecsClient != nil {
			if ecsInstanceDetail := s.getECSInstanceDetail(server.ServerId, ecsClient); ecsInstanceDetail != nil {
				detail.InstanceName = ecsInstanceDetail.InstanceName
				detail.PrivateIpAddress = ecsInstanceDetail.PrivateIpAddress
				detail.PublicIpAddress = ecsInstanceDetail.PublicIpAddress
			}
		}

		detailedServers = append(detailedServers, detail)
	}

	return detailedServers, nil
}

// ECSInstanceDetail contains ECS instance information needed for backend server details
type ECSInstanceDetail struct {
	InstanceName     string
	PrivateIpAddress string
	PublicIpAddress  string
}

// ForwardingRuleDetail contains detailed information about a forwarding rule
type ForwardingRuleDetail struct {
	RuleId           string
	RuleName         string
	Domain           string
	Url              string
	VServerGroupId   string
	VServerGroupName string
}

// FetchForwardingRules retrieves all forwarding rules for a specific listener
func (s *SLBService) FetchForwardingRules(loadBalancerId string, listenerPort int, listenerProtocol string) ([]ForwardingRuleDetail, error) {
	// Forwarding rules only apply to HTTP/HTTPS listeners
	if listenerProtocol != "HTTP" && listenerProtocol != "HTTPS" {
		return nil, nil
	}

	request := slb.CreateDescribeRulesRequest()
	request.Scheme = "https"
	request.LoadBalancerId = loadBalancerId
	request.ListenerPort = requests.NewInteger(listenerPort)
	request.ListenerProtocol = listenerProtocol

	response, err := s.client.DescribeRules(request)
	if err != nil {
		return nil, fmt.Errorf("describing forwarding rules for listener %s:%d: %w", loadBalancerId, listenerPort, err)
	}

	var rules []ForwardingRuleDetail
	for _, rule := range response.Rules.Rule {
		// Get VServer group name if available
		vsgName := ""
		if rule.VServerGroupId != "" {
			if name, err := s.getVServerGroupName(rule.VServerGroupId); err == nil {
				vsgName = name
			}
		}

		rules = append(rules, ForwardingRuleDetail{
			RuleId:           rule.RuleId,
			RuleName:         rule.RuleName,
			Domain:           rule.Domain,
			Url:              rule.Url,
			VServerGroupId:   rule.VServerGroupId,
			VServerGroupName: vsgName,
		})
	}

	return rules, nil
}

// getECSInstanceDetail retrieves ECS instance details for a given instance ID
func (s *SLBService) getECSInstanceDetail(instanceId string, ecsClient *ecs.Client) *ECSInstanceDetail {
	if ecsClient == nil {
		return nil
	}

	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", instanceId)

	response, err := ecsClient.DescribeInstances(request)
	if err != nil || len(response.Instances.Instance) == 0 {
		return nil
	}

	instance := response.Instances.Instance[0]

	// Get private IP
	privateIP := "N/A"
	if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		privateIP = instance.VpcAttributes.PrivateIpAddress.IpAddress[0]
	} else if len(instance.InnerIpAddress.IpAddress) > 0 {
		privateIP = instance.InnerIpAddress.IpAddress[0]
	}

	// Get public IP or EIP
	publicIP := "N/A"
	if len(instance.PublicIpAddress.IpAddress) > 0 {
		publicIP = instance.PublicIpAddress.IpAddress[0]
	} else if instance.EipAddress.IpAddress != "" {
		publicIP = fmt.Sprintf("EIP: %s", instance.EipAddress.IpAddress)
	}

	return &ECSInstanceDetail{
		InstanceName:     instance.InstanceName,
		PrivateIpAddress: privateIP,
		PublicIpAddress:  publicIP,
	}
}
