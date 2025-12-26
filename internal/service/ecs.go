package service

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// ECSService handles ECS operations
type ECSService struct {
	client *ecs.Client
}

// NewECSService creates a new ECS service
func NewECSService(client *ecs.Client) *ECSService {
	return &ECSService{client: client}
}

// FetchInstances retrieves all ECS instances using pagination
func (s *ECSService) FetchInstances() ([]ecs.Instance, error) {
	var allInstances []ecs.Instance
	pageNumber := 1
	pageSize := 100 // 使用最大页面大小以减少请求次数

	for {
		request := ecs.CreateDescribeInstancesRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(pageNumber)
		request.PageSize = requests.NewInteger(pageSize)

		response, err := s.client.DescribeInstances(request)
		if err != nil {
			return nil, fmt.Errorf("describing ECS instances (page %d): %w", pageNumber, err)
		}

		// 添加当前页的实例到总列表
		allInstances = append(allInstances, response.Instances.Instance...)

		// 检查是否还有更多页面
		// 如果当前页的实例数量小于页面大小，说明这是最后一页
		if len(response.Instances.Instance) < pageSize {
			break
		}

		// 也可以通过TotalCount来判断是否获取完所有数据
		if len(allInstances) >= response.TotalCount {
			break
		}

		pageNumber++
	}

	return allInstances, nil
}

// FetchSecurityGroups retrieves all security groups using pagination
func (s *ECSService) FetchSecurityGroups() ([]ecs.SecurityGroup, error) {
	var allSecurityGroups []ecs.SecurityGroup
	pageNumber := 1
	pageSize := 100 // 使用最大页面大小以减少请求次数

	for {
		request := ecs.CreateDescribeSecurityGroupsRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(pageNumber)
		request.PageSize = requests.NewInteger(pageSize)

		response, err := s.client.DescribeSecurityGroups(request)
		if err != nil {
			return nil, fmt.Errorf("describing security groups (page %d): %w", pageNumber, err)
		}

		// 添加当前页的安全组到总列表
		allSecurityGroups = append(allSecurityGroups, response.SecurityGroups.SecurityGroup...)

		// 检查是否还有更多页面
		// 如果当前页的安全组数量小于页面大小，说明这是最后一页
		if len(response.SecurityGroups.SecurityGroup) < pageSize {
			break
		}

		// 也可以通过TotalCount来判断是否获取完所有数据
		if len(allSecurityGroups) >= response.TotalCount {
			break
		}

		pageNumber++
	}

	return allSecurityGroups, nil
}

// FetchSecurityGroupRules retrieves security group rules for a specific security group
func (s *ECSService) FetchSecurityGroupRules(securityGroupId string) (*ecs.DescribeSecurityGroupAttributeResponse, error) {
	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.Scheme = "https"
	request.SecurityGroupId = securityGroupId

	response, err := s.client.DescribeSecurityGroupAttribute(request)
	if err != nil {
		return nil, fmt.Errorf("describing security group rules for %s: %w", securityGroupId, err)
	}

	return response, nil
}

// FetchInstancesBySecurityGroup retrieves ECS instances that use a specific security group
func (s *ECSService) FetchInstancesBySecurityGroup(securityGroupId string) ([]ecs.Instance, error) {
	var allInstances []ecs.Instance
	pageNumber := 1
	pageSize := 100

	for {
		request := ecs.CreateDescribeInstancesRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(pageNumber)
		request.PageSize = requests.NewInteger(pageSize)
		request.SecurityGroupId = securityGroupId

		response, err := s.client.DescribeInstances(request)
		if err != nil {
			return nil, fmt.Errorf("describing instances for security group %s (page %d): %w", securityGroupId, pageNumber, err)
		}

		// 添加当前页的实例到总列表
		allInstances = append(allInstances, response.Instances.Instance...)

		// 检查是否还有更多页面
		if len(response.Instances.Instance) < pageSize {
			break
		}

		if len(allInstances) >= response.TotalCount {
			break
		}

		pageNumber++
	}

	return allInstances, nil
}

// FetchDisks retrieves all disks for a specific ECS instance
func (s *ECSService) FetchDisks(instanceId string) ([]ecs.Disk, error) {
	var allDisks []ecs.Disk
	pageNumber := 1
	pageSize := 100

	for {
		request := ecs.CreateDescribeDisksRequest()
		request.Scheme = "https"
		request.InstanceId = instanceId
		request.PageNumber = requests.NewInteger(pageNumber)
		request.PageSize = requests.NewInteger(pageSize)

		response, err := s.client.DescribeDisks(request)
		if err != nil {
			return nil, fmt.Errorf("describing disks for instance %s (page %d): %w", instanceId, pageNumber, err)
		}

		allDisks = append(allDisks, response.Disks.Disk...)

		if len(response.Disks.Disk) < pageSize {
			break
		}

		if len(allDisks) >= response.TotalCount {
			break
		}

		pageNumber++
	}

	return allDisks, nil
}

// FetchSecurityGroupsByInstance retrieves security groups for a specific ECS instance
func (s *ECSService) FetchSecurityGroupsByInstance(instanceId string) ([]ecs.SecurityGroup, error) {
	// 首先获取实例详情以获取安全组ID列表
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", instanceId)

	response, err := s.client.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("describing instance %s: %w", instanceId, err)
	}

	if len(response.Instances.Instance) == 0 {
		return []ecs.SecurityGroup{}, nil
	}

	instance := response.Instances.Instance[0]
	securityGroupIds := instance.SecurityGroupIds.SecurityGroupId

	if len(securityGroupIds) == 0 {
		return []ecs.SecurityGroup{}, nil
	}

	// 获取安全组详情
	var securityGroups []ecs.SecurityGroup
	for _, sgId := range securityGroupIds {
		sgRequest := ecs.CreateDescribeSecurityGroupsRequest()
		sgRequest.Scheme = "https"
		sgRequest.SecurityGroupIds = fmt.Sprintf("[\"%s\"]", sgId)

		sgResponse, err := s.client.DescribeSecurityGroups(sgRequest)
		if err != nil {
			// 记录错误但继续处理其他安全组
			fmt.Printf("Warning: failed to describe security group %s: %v\n", sgId, err)
			continue
		}

		securityGroups = append(securityGroups, sgResponse.SecurityGroups.SecurityGroup...)
	}

	return securityGroups, nil
}
