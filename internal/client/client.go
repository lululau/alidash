package client

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ons20190214 "github.com/alibabacloud-go/ons-20190214/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliyunClients holds all Aliyun service clients
type AliyunClients struct {
	ECS      *ecs.Client
	DNS      *alidns.Client
	SLB      *slb.Client
	RDS      *rds.Client
	OSS      *oss.Client
	Redis    *r_kvstore.Client
	RocketMQ *ons20190214.Client
	config   *Config
}

// Config represents the client configuration
type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	OssEndpoint     string
}

// NewAliyunClients creates and initializes all Aliyun service clients
func NewAliyunClients(cfg *Config) (*AliyunClients, error) {
	clients := &AliyunClients{config: cfg}

	// Initialize ECS client
	ecsClient, err := ecs.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating ECS client: %w", err)
	}
	clients.ECS = ecsClient

	// Initialize DNS client
	dnsClient, err := alidns.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating DNS client: %w", err)
	}
	clients.DNS = dnsClient

	// Initialize SLB client
	slbClient, err := slb.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating SLB client: %w", err)
	}
	clients.SLB = slbClient

	// Initialize RDS client
	rdsClient, err := rds.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating RDS client: %w", err)
	}
	clients.RDS = rdsClient

	// Initialize OSS client
	ossClient, err := oss.New(cfg.OssEndpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating OSS client: %w", err)
	}
	clients.OSS = ossClient

	// Initialize Redis client
	redisClient, err := r_kvstore.NewClientWithAccessKey(cfg.RegionID, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("creating Redis client: %w", err)
	}
	clients.Redis = redisClient

	// Initialize RocketMQ client using V2.0 SDK
	rocketmqConfig := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyID),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		RegionId:        tea.String(cfg.RegionID),
		Endpoint:        tea.String(fmt.Sprintf("ons.%s.aliyuncs.com", cfg.RegionID)),
	}
	rocketmqClient, err := ons20190214.NewClient(rocketmqConfig)
	if err != nil {
		return nil, fmt.Errorf("creating RocketMQ client: %w", err)
	}
	clients.RocketMQ = rocketmqClient

	return clients, nil
}

// GetConfig returns the client configuration
func (c *AliyunClients) GetConfig() *Config {
	return c.config
}

// UpdateRegion creates new clients with a different region
func (c *AliyunClients) UpdateRegion(regionID string) (*AliyunClients, error) {
	newConfig := &Config{
		AccessKeyID:     c.config.AccessKeyID,
		AccessKeySecret: c.config.AccessKeySecret,
		RegionID:        regionID,
		OssEndpoint:     fmt.Sprintf("oss-%s.aliyuncs.com", regionID),
	}
	return NewAliyunClients(newConfig)
}
