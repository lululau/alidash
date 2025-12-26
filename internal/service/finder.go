package service

import (
	"net"
	"regexp"
	"strings"
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

// FinderService provides resource finding functionality
type FinderService struct {
	ecsService      *ECSService
	dnsService      *DNSService
	slbService      *SLBService
	rdsService      *RDSService
	redisService    *RedisService
	rocketMQService *RocketMQService
}

// NewFinderService creates a new finder service
func NewFinderService(
	ecs *ECSService,
	dns *DNSService,
	slb *SLBService,
	rds *RDSService,
	redis *RedisService,
	rocketMQ *RocketMQService,
) *FinderService {
	return &FinderService{
		ecsService:      ecs,
		dnsService:      dns,
		slbService:      slb,
		rdsService:      rds,
		redisService:    redis,
		rocketMQService: rocketMQ,
	}
}

// FindResult contains all matching resources
type FindResult struct {
	Query              string
	ResolvedIPs        []string
	ECSInstances       []ecs.Instance
	ENIs               []ecs.NetworkInterfaceSet
	SLBInstances       []slb.LoadBalancer
	DNSRecords         []DNSRecordMatch
	RDSInstances       []RDSInstanceDetail // Changed to detailed instances
	RedisInstances     []r_kvstore.KVStoreInstance
	RocketMQInstances  []RocketMQInstance
}

// DNSRecordMatch contains a matched DNS record with its domain
type DNSRecordMatch struct {
	DomainName string
	Record     alidns.Record
}

// IsIP checks if the input is an IP address
func IsIP(input string) bool {
	return net.ParseIP(input) != nil
}

// IsDomain checks if the input looks like a domain name
func IsDomain(input string) bool {
	// Simple domain pattern check
	domainPattern := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)
	return domainPattern.MatchString(input)
}

// ResolveToIPs resolves the input to IP addresses
// If input is an IP, returns it directly
// If input is a domain, tries to resolve via Aliyun DNS first, then system DNS
func (s *FinderService) ResolveToIPs(input string) ([]string, string, error) {
	input = strings.TrimSpace(input)
	
	// If already an IP, return it
	if IsIP(input) {
		return []string{input}, input, nil
	}

	// Try to resolve via Aliyun DNS records first
	if s.dnsService != nil {
		ips, err := s.resolveViaAliyunDNS(input)
		if err == nil && len(ips) > 0 {
			return ips, input, nil
		}
	}

	// Fall back to system DNS
	addrs, err := net.LookupHost(input)
	if err != nil {
		// Return empty IPs but keep the domain for DNS record search
		return nil, input, nil
	}

	return addrs, input, nil
}

// resolveViaAliyunDNS looks up the domain in Aliyun DNS records
func (s *FinderService) resolveViaAliyunDNS(domain string) ([]string, error) {
	if s.dnsService == nil {
		return nil, nil
	}

	// Get all domains
	domains, err := s.dnsService.FetchDomains()
	if err != nil {
		return nil, err
	}

	var ips []string

	// Find matching records
	for _, d := range domains {
		// Check if the input domain is under this DNS domain
		if strings.HasSuffix(domain, d.DomainName) || domain == d.DomainName {
			records, err := s.dnsService.FetchDomainRecords(d.DomainName)
			if err != nil {
				continue
			}

			// Calculate the expected RR (subdomain part)
			expectedRR := strings.TrimSuffix(domain, "."+d.DomainName)
			if domain == d.DomainName {
				expectedRR = "@"
			}

			for _, r := range records {
				// Match A records with the subdomain
				if r.Type == "A" && (r.RR == expectedRR || r.RR+"." +d.DomainName == domain) {
					ips = append(ips, r.Value)
				}
			}
		}
	}

	return ips, nil
}

// FindResources searches for resources matching the given IPs and domain
func (s *FinderService) FindResources(ips []string, domain string) (*FindResult, error) {
	result := &FindResult{
		Query:       domain,
		ResolvedIPs: ips,
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Search ECS instances
	if s.ecsService != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instances, err := s.ecsService.FetchInstances()
			if err != nil {
				return
			}
			matched := s.matchECSInstances(instances, ips)
			mu.Lock()
			result.ECSInstances = matched
			mu.Unlock()
		}()
	}

	// Search ENIs
	if s.ecsService != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			enis, err := s.fetchAllENIs()
			if err != nil {
				return
			}
			matched := s.matchENIs(enis, ips)
			mu.Lock()
			result.ENIs = matched
			mu.Unlock()
		}()
	}

	// Search SLB instances
	if s.slbService != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lbs, err := s.slbService.FetchInstances()
			if err != nil {
				return
			}
			matched := s.matchSLBInstances(lbs, ips)
			mu.Lock()
			result.SLBInstances = matched
			mu.Unlock()
		}()
	}

	// Search DNS records
	if s.dnsService != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			matched, _ := s.matchDNSRecords(ips, domain)
			mu.Lock()
			result.DNSRecords = matched
			mu.Unlock()
		}()
	}

	// Search RDS instances (with network info for public address matching)
	if s.rdsService != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instances, err := s.rdsService.FetchDetailedInstances()
			if err != nil {
				return
			}
			matched := s.matchRDSDetailedInstances(instances, ips, domain)
			mu.Lock()
			result.RDSInstances = matched
			mu.Unlock()
		}()
	}

	// Search Redis instances
	if s.redisService != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instances, err := s.redisService.FetchInstances()
			if err != nil {
				return
			}
			matched := s.matchRedisInstances(instances, ips, domain)
			mu.Lock()
			result.RedisInstances = matched
			mu.Unlock()
		}()
	}

	// Search RocketMQ instances
	if s.rocketMQService != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instances, err := s.rocketMQService.FetchInstances()
			if err != nil {
				return
			}
			matched := s.matchRocketMQInstances(instances, ips, domain)
			mu.Lock()
			result.RocketMQInstances = matched
			mu.Unlock()
		}()
	}

	wg.Wait()
	return result, nil
}

// fetchAllENIs fetches all ENIs (this is a simplified version)
func (s *FinderService) fetchAllENIs() ([]ecs.NetworkInterfaceSet, error) {
	// Note: This requires iterating through all instances or using DescribeNetworkInterfaces API
	// For simplicity, we'll fetch instances first and then their ENIs
	instances, err := s.ecsService.FetchInstances()
	if err != nil {
		return nil, err
	}

	var allENIs []ecs.NetworkInterfaceSet
	for _, inst := range instances {
		enis, err := s.ecsService.FetchNetworkInterfaces(inst.InstanceId)
		if err != nil {
			continue
		}
		allENIs = append(allENIs, enis...)
	}
	return allENIs, nil
}

// containsAny checks if target contains any of the search strings (case-insensitive)
func containsAny(target string, searches []string) bool {
	targetLower := strings.ToLower(target)
	for _, s := range searches {
		if s != "" && strings.Contains(targetLower, strings.ToLower(s)) {
			return true
		}
	}
	return false
}

// matchECSInstances finds ECS instances matching the given IPs (using contains matching)
func (s *FinderService) matchECSInstances(instances []ecs.Instance, ips []string) []ecs.Instance {
	if len(ips) == 0 {
		return nil
	}

	var matched []ecs.Instance
	for _, inst := range instances {
		// Check public IP
		for _, pip := range inst.PublicIpAddress.IpAddress {
			if containsAny(pip, ips) {
				matched = append(matched, inst)
				goto next
			}
		}
		// Check private IP (VPC)
		for _, pip := range inst.VpcAttributes.PrivateIpAddress.IpAddress {
			if containsAny(pip, ips) {
				matched = append(matched, inst)
				goto next
			}
		}
		// Check inner IP (classic)
		for _, pip := range inst.InnerIpAddress.IpAddress {
			if containsAny(pip, ips) {
				matched = append(matched, inst)
				goto next
			}
		}
		// Check EIP
		if containsAny(inst.EipAddress.IpAddress, ips) {
			matched = append(matched, inst)
			goto next
		}
	next:
	}
	return matched
}

// matchENIs finds ENIs matching the given IPs (using contains matching)
func (s *FinderService) matchENIs(enis []ecs.NetworkInterfaceSet, ips []string) []ecs.NetworkInterfaceSet {
	if len(ips) == 0 {
		return nil
	}

	var matched []ecs.NetworkInterfaceSet
	for _, eni := range enis {
		// Check primary private IP
		if containsAny(eni.PrivateIpAddress, ips) {
			matched = append(matched, eni)
			continue
		}
		// Check private IP sets
		for _, pip := range eni.PrivateIpSets.PrivateIpSet {
			if containsAny(pip.PrivateIpAddress, ips) {
				matched = append(matched, eni)
				break
			}
		}
	}
	return matched
}

// matchSLBInstances finds SLB instances matching the given IPs (using contains matching)
func (s *FinderService) matchSLBInstances(lbs []slb.LoadBalancer, ips []string) []slb.LoadBalancer {
	if len(ips) == 0 {
		return nil
	}

	var matched []slb.LoadBalancer
	for _, lb := range lbs {
		if containsAny(lb.Address, ips) {
			matched = append(matched, lb)
		}
	}
	return matched
}

// matchDNSRecords finds DNS records matching the given IPs or domain (using contains matching)
func (s *FinderService) matchDNSRecords(ips []string, domain string) ([]DNSRecordMatch, error) {
	domains, err := s.dnsService.FetchDomains()
	if err != nil {
		return nil, err
	}

	var matched []DNSRecordMatch
	for _, d := range domains {
		records, err := s.dnsService.FetchDomainRecords(d.DomainName)
		if err != nil {
			continue
		}

		for _, r := range records {
			// Match by IP value (A records) using contains matching
			if r.Type == "A" && containsAny(r.Value, ips) {
				matched = append(matched, DNSRecordMatch{
					DomainName: d.DomainName,
					Record:     r,
				})
				continue
			}
			// Match by domain in record value (CNAME, etc.)
			if domain != "" && strings.Contains(strings.ToLower(r.Value), strings.ToLower(domain)) {
				matched = append(matched, DNSRecordMatch{
					DomainName: d.DomainName,
					Record:     r,
				})
				continue
			}
			// Match by subdomain (using contains matching)
			if domain != "" {
				fullRecord := r.RR + "." + d.DomainName
				if strings.Contains(strings.ToLower(fullRecord), strings.ToLower(domain)) ||
					strings.Contains(strings.ToLower(r.RR), strings.ToLower(domain)) {
					matched = append(matched, DNSRecordMatch{
						DomainName: d.DomainName,
						Record:     r,
					})
				}
			}
		}
	}
	return matched, nil
}

// matchRDSDetailedInstances finds RDS instances matching the given IPs or domain (using contains matching)
func (s *FinderService) matchRDSDetailedInstances(instances []RDSInstanceDetail, ips []string, domain string) []RDSInstanceDetail {
	var matched []RDSInstanceDetail
	for _, inst := range instances {
		// Check internal connection string
		if domain != "" && strings.Contains(strings.ToLower(inst.InternalConnectionStr), strings.ToLower(domain)) {
			matched = append(matched, inst)
			continue
		}
		// Check public connection string
		if domain != "" && strings.Contains(strings.ToLower(inst.PublicConnectionStr), strings.ToLower(domain)) {
			matched = append(matched, inst)
			continue
		}
		// Check internal IP using contains matching
		if containsAny(inst.InternalIP, ips) {
			matched = append(matched, inst)
			continue
		}
		// Check public IP using contains matching
		if containsAny(inst.PublicIP, ips) {
			matched = append(matched, inst)
			continue
		}
		// Check if connection strings contain IP
		if containsAny(inst.InternalConnectionStr, ips) || containsAny(inst.PublicConnectionStr, ips) {
			matched = append(matched, inst)
		}
	}
	return matched
}

// matchRedisInstances finds Redis instances matching the given IPs or domain (using contains matching)
func (s *FinderService) matchRedisInstances(instances []r_kvstore.KVStoreInstance, ips []string, domain string) []r_kvstore.KVStoreInstance {
	var matched []r_kvstore.KVStoreInstance
	for _, inst := range instances {
		// Check connection domain contains the domain
		if domain != "" && strings.Contains(strings.ToLower(inst.ConnectionDomain), strings.ToLower(domain)) {
			matched = append(matched, inst)
			continue
		}
		// Check private IP using contains matching
		if containsAny(inst.PrivateIp, ips) {
			matched = append(matched, inst)
		}
	}
	return matched
}

// matchRocketMQInstances finds RocketMQ instances matching the given IPs or domain
// Note: RocketMQ instances don't have direct IP/endpoint info in the basic list API
// This is a placeholder that matches by instance name containing the domain
func (s *FinderService) matchRocketMQInstances(instances []RocketMQInstance, ips []string, domain string) []RocketMQInstance {
	if domain == "" {
		return nil
	}

	var matched []RocketMQInstance
	for _, inst := range instances {
		// Check if instance name contains the domain
		if strings.Contains(strings.ToLower(inst.InstanceName), strings.ToLower(domain)) {
			matched = append(matched, inst)
		}
	}
	return matched
}

// HasResults checks if the find result has any matches
func (r *FindResult) HasResults() bool {
	return len(r.ECSInstances) > 0 ||
		len(r.ENIs) > 0 ||
		len(r.SLBInstances) > 0 ||
		len(r.DNSRecords) > 0 ||
		len(r.RDSInstances) > 0 ||
		len(r.RedisInstances) > 0 ||
		len(r.RocketMQInstances) > 0
}

// TotalCount returns the total number of matched resources
func (r *FindResult) TotalCount() int {
	return len(r.ECSInstances) +
		len(r.ENIs) +
		len(r.SLBInstances) +
		len(r.DNSRecords) +
		len(r.RDSInstances) +
		len(r.RedisInstances) +
		len(r.RocketMQInstances)
}

