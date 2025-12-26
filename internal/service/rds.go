package service

import (
	"fmt"
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
)

// RDSService handles RDS operations
type RDSService struct {
	client *rds.Client
}

// RDSInstanceDetail contains RDS instance with network info
type RDSInstanceDetail struct {
	Instance              rds.DBInstance
	InternalConnectionStr string // 内网连接地址
	InternalIP            string // 内网IP
	PublicConnectionStr   string // 外网连接地址
	PublicIP              string // 外网IP
}

// NewRDSService creates a new RDS service
func NewRDSService(client *rds.Client) *RDSService {
	return &RDSService{client: client}
}

// FetchInstances retrieves all RDS instances using pagination
func (s *RDSService) FetchInstances() ([]rds.DBInstance, error) {
	var allInstances []rds.DBInstance
	pageNumber := 1
	pageSize := 100 // 使用最大页面大小以减少请求次数

	for {
		request := rds.CreateDescribeDBInstancesRequest()
		request.Scheme = "https"
		request.PageNumber = requests.NewInteger(pageNumber)
		request.PageSize = requests.NewInteger(pageSize)

		response, err := s.client.DescribeDBInstances(request)
		if err != nil {
			return nil, fmt.Errorf("describing RDS instances (page %d): %w", pageNumber, err)
		}

		// 添加当前页的实例到总列表
		allInstances = append(allInstances, response.Items.DBInstance...)

		// 检查是否还有更多页面
		// 如果当前页的实例数量小于页面大小，说明这是最后一页
		if len(response.Items.DBInstance) < pageSize {
			break
		}

		// 也可以通过TotalRecordCount来判断是否获取完所有数据
		if len(allInstances) >= response.TotalRecordCount {
			break
		}

		pageNumber++
	}

	return allInstances, nil
}

// FetchDatabases retrieves all databases for a specific RDS instance
func (s *RDSService) FetchDatabases(dbInstanceId string) ([]rds.Database, error) {
	request := rds.CreateDescribeDatabasesRequest()
	request.Scheme = "https"
	request.DBInstanceId = dbInstanceId

	response, err := s.client.DescribeDatabases(request)
	if err != nil {
		return nil, fmt.Errorf("describing databases for instance %s: %w", dbInstanceId, err)
	}

	return response.Databases.Database, nil
}

// FetchAccounts retrieves all accounts for a specific RDS instance
func (s *RDSService) FetchAccounts(dbInstanceId string) ([]rds.DBInstanceAccount, error) {
	request := rds.CreateDescribeAccountsRequest()
	request.Scheme = "https"
	request.DBInstanceId = dbInstanceId

	response, err := s.client.DescribeAccounts(request)
	if err != nil {
		return nil, fmt.Errorf("describing accounts for instance %s: %w", dbInstanceId, err)
	}

	return response.Accounts.DBInstanceAccount, nil
}

// FetchInstanceNetInfo retrieves network info for a specific RDS instance
func (s *RDSService) FetchInstanceNetInfo(dbInstanceId string) ([]rds.DBInstanceNetInfo, error) {
	request := rds.CreateDescribeDBInstanceNetInfoRequest()
	request.Scheme = "https"
	request.DBInstanceId = dbInstanceId

	response, err := s.client.DescribeDBInstanceNetInfo(request)
	if err != nil {
		return nil, fmt.Errorf("describing network info for instance %s: %w", dbInstanceId, err)
	}

	return response.DBInstanceNetInfos.DBInstanceNetInfo, nil
}

// FetchDetailedInstances retrieves all RDS instances with their network info
func (s *RDSService) FetchDetailedInstances() ([]RDSInstanceDetail, error) {
	// First fetch all instances
	instances, err := s.FetchInstances()
	if err != nil {
		return nil, err
	}

	// Fetch network info in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	detailedInstances := make([]RDSInstanceDetail, len(instances))

	for i, inst := range instances {
		detailedInstances[i] = RDSInstanceDetail{
			Instance:              inst,
			InternalConnectionStr: inst.ConnectionString, // Default from instance list
		}

		wg.Add(1)
		go func(idx int, instanceId string) {
			defer wg.Done()

			netInfos, err := s.FetchInstanceNetInfo(instanceId)
			if err != nil {
				return // Skip if error
			}

			mu.Lock()
			defer mu.Unlock()

			for _, netInfo := range netInfos {
				// IPType: "Private" for internal, "Public" for external
				// Also check ConnectionStringType: "Normal" for standard connection
				switch netInfo.IPType {
				case "Private":
					detailedInstances[idx].InternalConnectionStr = netInfo.ConnectionString
					detailedInstances[idx].InternalIP = netInfo.IPAddress
				case "Public":
					detailedInstances[idx].PublicConnectionStr = netInfo.ConnectionString
					detailedInstances[idx].PublicIP = netInfo.IPAddress
				}
			}
		}(i, inst.DBInstanceId)
	}

	wg.Wait()
	return detailedInstances, nil
}
