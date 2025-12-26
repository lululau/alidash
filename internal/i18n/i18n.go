package i18n

import (
	"aliyun-tui-viewer/internal/config"
)

// Translation keys
const (
	// App titles
	KeyAppTitle            = "app.title"
	KeyPageMenu            = "page.menu"
	KeyPageECSList         = "page.ecs_list"
	KeyPageECSDetail       = "page.ecs_detail"
	KeyPageECSJSONDetail   = "page.ecs_json_detail"
	KeyPageECSDisks        = "page.ecs_disks"
	KeyPageECSENIs         = "page.ecs_enis"
	KeyPageSecurityGroups  = "page.security_groups"
	KeyPageSGRules         = "page.sg_rules"
	KeyPageSGInstances     = "page.sg_instances"
	KeyPageInstSGs         = "page.inst_sgs"
	KeyPageDNSDomains      = "page.dns_domains"
	KeyPageDNSRecords      = "page.dns_records"
	KeyPageSLBList         = "page.slb_list"
	KeyPageSLBDetail       = "page.slb_detail"
	KeyPageSLBListeners    = "page.slb_listeners"
	KeyPageVServerGroups   = "page.vserver_groups"
	KeyPageBackendServers  = "page.backend_servers"
	KeyPageForwardRules    = "page.forward_rules"
	KeyPageDefaultServers  = "page.default_servers"
	KeyPageOSSBuckets      = "page.oss_buckets"
	KeyPageOSSObjects      = "page.oss_objects"
	KeyPageOSSDetail       = "page.oss_detail"
	KeyPageRDSList         = "page.rds_list"
	KeyPageRDSDetail       = "page.rds_detail"
	KeyPageRDSDatabases    = "page.rds_databases"
	KeyPageRDSAccounts     = "page.rds_accounts"
	KeyPageRedisList       = "page.redis_list"
	KeyPageRedisDetail     = "page.redis_detail"
	KeyPageRedisAccounts   = "page.redis_accounts"
	KeyPageRocketMQList    = "page.rocketmq_list"
	KeyPageRocketMQDetail  = "page.rocketmq_detail"
	KeyPageRocketMQTopics  = "page.rocketmq_topics"
	KeyPageRocketMQGroups  = "page.rocketmq_groups"
	KeyPageResourceFinder  = "page.resource_finder"

	// Menu items
	KeyMenuECS           = "menu.ecs"
	KeyMenuECSDesc       = "menu.ecs_desc"
	KeyMenuSG            = "menu.sg"
	KeyMenuSGDesc        = "menu.sg_desc"
	KeyMenuDNS           = "menu.dns"
	KeyMenuDNSDesc       = "menu.dns_desc"
	KeyMenuSLB           = "menu.slb"
	KeyMenuSLBDesc       = "menu.slb_desc"
	KeyMenuOSS           = "menu.oss"
	KeyMenuOSSDesc       = "menu.oss_desc"
	KeyMenuRDS           = "menu.rds"
	KeyMenuRDSDesc       = "menu.rds_desc"
	KeyMenuRedis         = "menu.redis"
	KeyMenuRedisDesc     = "menu.redis_desc"
	KeyMenuRocketMQ      = "menu.rocketmq"
	KeyMenuRocketMQDesc  = "menu.rocketmq_desc"

	// Header
	KeyHeaderProfile = "header.profile"
	KeyHeaderRegion  = "header.region"

	// Modal
	KeyModalInfo          = "modal.info"
	KeyModalError         = "modal.error"
	KeyModalSuccess       = "modal.success"
	KeyModalOK            = "modal.ok"
	KeyModalConfirm       = "modal.confirm"
	KeyModalCancel        = "modal.cancel"
	KeyModalSelectProfile = "modal.select_profile"
	KeyModalSelectRegion  = "modal.select_region"
	KeyModalLoading       = "modal.loading"
	KeyModalResourceFind  = "modal.resource_find"
	KeyModalInputPrompt   = "modal.input_prompt"
	KeyModalInputExample  = "modal.input_example"
	KeyModalHistory       = "modal.history"
	KeyModalCurrent       = "modal.current"

	// Common columns
	KeyColInstanceID   = "col.instance_id"
	KeyColName         = "col.name"
	KeyColStatus       = "col.status"
	KeyColZone         = "col.zone"
	KeyColCPURAM       = "col.cpu_ram"
	KeyColPrivateIP    = "col.private_ip"
	KeyColPublicIP     = "col.public_ip"
	KeyColExpired      = "col.expired"
	KeyColType         = "col.type"
	KeyColDescription  = "col.description"
	KeyColCreatedAt    = "col.created_at"
	KeyColRegion       = "col.region"

	// ECS specific
	KeyColENIID        = "col.eni_id"
	KeyColENIType      = "col.eni_type"
	KeyColAttachedInst = "col.attached_inst"
	KeyColDiskID       = "col.disk_id"
	KeyColDiskType     = "col.disk_type"
	KeyColDiskSize     = "col.disk_size"
	KeyColDiskCategory = "col.disk_category"
	KeyColDevice       = "col.device"

	// Security Group
	KeyColSGID         = "col.sg_id"
	KeyColSGName       = "col.sg_name"
	KeyColDirection    = "col.direction"
	KeyColProtocol     = "col.protocol"
	KeyColPortRange    = "col.port_range"
	KeyColSource       = "col.source"
	KeyColDestination  = "col.destination"
	KeyColPolicy       = "col.policy"
	KeyColPriority     = "col.priority"

	// DNS
	KeyColDomain       = "col.domain"
	KeyColRecordType   = "col.record_type"
	KeyColRecordValue  = "col.record_value"
	KeyColRR           = "col.rr"
	KeyColTTL          = "col.ttl"
	KeyColLine         = "col.line"

	// SLB
	KeyColSLBID        = "col.slb_id"
	KeyColAddress      = "col.address"
	KeyColAddressType  = "col.address_type"
	KeyColSpec         = "col.spec"
	KeyColPort         = "col.port"
	KeyColBackendPort  = "col.backend_port"
	KeyColWeight       = "col.weight"
	KeyColServerID     = "col.server_id"
	KeyColServerName   = "col.server_name"
	KeyColVServerGroup = "col.vserver_group"

	// OSS
	KeyColBucketName   = "col.bucket_name"
	KeyColObjectKey    = "col.object_key"
	KeyColSize         = "col.size"
	KeyColStorageClass = "col.storage_class"
	KeyColLastModified = "col.last_modified"

	// RDS
	KeyColDBInstanceID = "col.db_instance_id"
	KeyColEngine       = "col.engine"
	KeyColDBName       = "col.db_name"
	KeyColCharset      = "col.charset"
	KeyColAccountName  = "col.account_name"
	KeyColAccountType  = "col.account_type"
	KeyColPrivilege    = "col.privilege"
	KeyColConnString   = "col.conn_string"

	// Redis
	KeyColConnDomain   = "col.conn_domain"
	KeyColArchType     = "col.arch_type"
	KeyColCapacity     = "col.capacity"

	// RocketMQ
	KeyColInstanceName = "col.instance_name"
	KeyColTopicName    = "col.topic_name"
	KeyColGroupID      = "col.group_id"
	KeyColMessageType  = "col.message_type"
	KeyColGroupType    = "col.group_type"

	// Finder section titles
	KeyFinderResult       = "finder.result"
	KeyFinderTotalMatches = "finder.total_matches"
	KeyFinderNoMatch      = "finder.no_match"
	KeyFinderECS          = "finder.ecs"
	KeyFinderENI          = "finder.eni"
	KeyFinderSLB          = "finder.slb"
	KeyFinderDNS          = "finder.dns"
	KeyFinderRDS          = "finder.rds"
	KeyFinderRedis        = "finder.redis"
	KeyFinderRocketMQ     = "finder.rocketmq"

	// ENI Types
	KeyENIPrimary   = "eni.primary"
	KeyENISecondary = "eni.secondary"

	// ENI Status
	KeyENIAvailable = "eni.available"
	KeyENIInUse     = "eni.in_use"
	KeyENIAttaching = "eni.attaching"
	KeyENIDetaching = "eni.detaching"
	KeyENIDeleting  = "eni.deleting"

	// Status values (for display consistency)
	KeyStatusRunning   = "status.running"
	KeyStatusStopped   = "status.stopped"
	KeyStatusCreating  = "status.creating"
	KeyStatusReleased  = "status.released"
	KeyStatusUnknown   = "status.unknown"

	// ECS Detail sections
	KeySectionBasicInfo    = "section.basic_info"
	KeySectionConfigInfo   = "section.config_info"
	KeySectionBoundRes     = "section.bound_resources"
	KeySectionGroupInfo    = "section.group_info"
	KeySectionOtherInfo    = "section.other_info"
	KeySectionUsageOverview = "section.usage_overview"
	KeySectionStorageOverview = "section.storage_overview"

	// ECS Detail labels
	KeyLabelInstanceID     = "label.instance_id"
	KeyLabelInstanceName   = "label.instance_name"
	KeyLabelInstanceStatus = "label.instance_status"
	KeyLabelAvailZone      = "label.avail_zone"
	KeyLabelChargeType     = "label.charge_type"
	KeyLabelExpireTime     = "label.expire_time"
	KeyLabelInstanceSpec   = "label.instance_spec"
	KeyLabelCPUMemory      = "label.cpu_memory"
	KeyLabelImageID        = "label.image_id"
	KeyLabelOSName         = "label.os_name"
	KeyLabelVPC            = "label.vpc"
	KeyLabelVSwitch        = "label.vswitch"
	KeyLabelNetworkType    = "label.network_type"
	KeyLabelBandwidth      = "label.bandwidth"
	KeyLabelBandwidthCharge = "label.bandwidth_charge"
	KeyLabelSecurityGroup  = "label.security_group"
	KeyLabelEIPID          = "label.eip_id"
	KeyLabelSecondaryIP    = "label.secondary_ip"
	KeyLabelResourceGroup  = "label.resource_group"
	KeyLabelTags           = "label.tags"
	KeyLabelHostname       = "label.hostname"
	KeyLabelKeyPair        = "label.key_pair"
	KeyLabelSerialNumber   = "label.serial_number"
	KeyLabelTotalDisks     = "label.total_disks"
	KeyLabelTotalCapacity  = "label.total_capacity"

	// Charge types
	KeyChargePrePaid  = "charge.prepaid"
	KeyChargePostPaid = "charge.postpaid"

	// Network types
	KeyNetworkVPC     = "network.vpc"
	KeyNetworkClassic = "network.classic"

	// Disk related
	KeyDiskCloud           = "disk.cloud"
	KeyDiskCloudEfficiency = "disk.cloud_efficiency"
	KeyDiskCloudSSD        = "disk.cloud_ssd"
	KeyDiskCloudEssd       = "disk.cloud_essd"
	KeyDiskCloudEssdAuto   = "disk.cloud_essd_auto"
	KeyDiskCloudEssdEntry  = "disk.cloud_essd_entry"
	KeyDiskSystem          = "disk.system"
	KeyDiskData            = "disk.data"
	KeyDiskDeleteWithInst  = "disk.delete_with_inst"
	KeyDiskKeepAfterInst   = "disk.keep_after_inst"
	KeyDiskPortable        = "disk.portable"
	KeyDiskNotPortable     = "disk.not_portable"
	KeyColDiskAttribute    = "col.disk_attribute"
	KeyColDiskIOPS         = "col.disk_iops"
	KeyColDiskDeleteBehavior = "col.disk_delete_behavior"
	KeyColDiskPortable     = "col.disk_portable"

	// SLB related
	KeyColServerGroupIDName = "col.server_group_id_name"
	KeyColAssocListeners   = "col.assoc_listeners"
	KeyColAssocRules       = "col.assoc_rules"
	KeyColBackendCount     = "col.backend_count"
	KeyColServerIDName     = "col.server_id_name"
	KeyColPublicPrivateIP  = "col.public_private_ip"
	KeyColRemark           = "col.remark"
	KeyColURL              = "col.url"
	KeySLBVServerGroupList = "slb.vserver_group_list"
	KeySLBForwardingRules  = "slb.forwarding_rules"
	KeySLBDefaultServerGroup = "slb.default_server_group"

	// ENI columns
	KeyColMAC              = "col.mac"

	// Count format
	KeyCountSG             = "count.sg"
	KeyCountENI            = "count.eni"
	KeyCountDisks          = "count.disks"

	// Common actions
	KeyActionCopied   = "action.copied"
	KeyActionLoading  = "action.loading"
	KeyNoData         = "common.no_data"
)

// Region display names (Chinese)
var regionNamesZhCN = map[string]string{
	"cn-hangzhou":    "华东1（杭州）",
	"cn-shanghai":    "华东2（上海）",
	"cn-nanjing":     "华东5（南京）",
	"cn-fuzhou":      "华东6（福州）",
	"cn-wuhan-lr":    "华中1（武汉）",
	"cn-qingdao":     "华北1（青岛）",
	"cn-beijing":     "华北2（北京）",
	"cn-zhangjiakou": "华北3（张家口）",
	"cn-huhehaote":   "华北5（呼和浩特）",
	"cn-wulanchabu":  "华北6（乌兰察布）",
	"cn-shenzhen":    "华南1（深圳）",
	"cn-heyuan":      "华南2（河源）",
	"cn-guangzhou":   "华南3（广州）",
	"cn-chengdu":     "西南1（成都）",
	"cn-hongkong":    "中国（香港）",
	"ap-southeast-1": "新加坡",
	"ap-southeast-2": "澳大利亚（悉尼）",
	"ap-southeast-3": "马来西亚（吉隆坡）",
	"ap-southeast-5": "印度尼西亚（雅加达）",
	"ap-southeast-6": "菲律宾（马尼拉）",
	"ap-southeast-7": "泰国（曼谷）",
	"ap-south-1":     "印度（孟买）",
	"ap-northeast-1": "日本（东京）",
	"ap-northeast-2": "韩国（首尔）",
	"us-west-1":      "美国（硅谷）",
	"us-east-1":      "美国（弗吉尼亚）",
	"eu-central-1":   "德国（法兰克福）",
	"eu-west-1":      "英国（伦敦）",
	"me-east-1":      "阿联酋（迪拜）",
	"me-central-1":   "沙特（利雅得）",
}

// Region display names (English)
var regionNamesEnUS = map[string]string{
	"cn-hangzhou":    "China (Hangzhou)",
	"cn-shanghai":    "China (Shanghai)",
	"cn-nanjing":     "China (Nanjing)",
	"cn-fuzhou":      "China (Fuzhou)",
	"cn-wuhan-lr":    "China (Wuhan)",
	"cn-qingdao":     "China (Qingdao)",
	"cn-beijing":     "China (Beijing)",
	"cn-zhangjiakou": "China (Zhangjiakou)",
	"cn-huhehaote":   "China (Hohhot)",
	"cn-wulanchabu":  "China (Ulanqab)",
	"cn-shenzhen":    "China (Shenzhen)",
	"cn-heyuan":      "China (Heyuan)",
	"cn-guangzhou":   "China (Guangzhou)",
	"cn-chengdu":     "China (Chengdu)",
	"cn-hongkong":    "China (Hong Kong)",
	"ap-southeast-1": "Singapore",
	"ap-southeast-2": "Australia (Sydney)",
	"ap-southeast-3": "Malaysia (Kuala Lumpur)",
	"ap-southeast-5": "Indonesia (Jakarta)",
	"ap-southeast-6": "Philippines (Manila)",
	"ap-southeast-7": "Thailand (Bangkok)",
	"ap-south-1":     "India (Mumbai)",
	"ap-northeast-1": "Japan (Tokyo)",
	"ap-northeast-2": "South Korea (Seoul)",
	"us-west-1":      "US (Silicon Valley)",
	"us-east-1":      "US (Virginia)",
	"eu-central-1":   "Germany (Frankfurt)",
	"eu-west-1":      "UK (London)",
	"me-east-1":      "UAE (Dubai)",
	"me-central-1":   "Saudi Arabia (Riyadh)",
}

// GetRegionDisplayName returns a localized human-readable name for a region ID
func GetRegionDisplayName(regionID string) string {
	var regionNames map[string]string
	if IsChinese() {
		regionNames = regionNamesZhCN
	} else {
		regionNames = regionNamesEnUS
	}

	if name, ok := regionNames[regionID]; ok {
		return name + " (" + regionID + ")"
	}
	return regionID
}

// translations holds all translations
var translations = map[string]map[string]string{
	config.LocaleEnUS: enUS,
	config.LocaleZhCN: zhCN,
}

// English translations
var enUS = map[string]string{
	// App titles
	KeyAppTitle:            "Aliyun TUI Dashboard",
	KeyPageMenu:            "Aliyun TUI Dashboard",
	KeyPageECSList:         "ECS Instances",
	KeyPageECSDetail:       "ECS Detail",
	KeyPageECSJSONDetail:   "ECS JSON Detail",
	KeyPageECSDisks:        "ECS Disks",
	KeyPageECSENIs:         "ECS Network Interfaces",
	KeyPageSecurityGroups:  "Security Groups",
	KeyPageSGRules:         "Security Group Rules",
	KeyPageSGInstances:     "Security Group Instances",
	KeyPageInstSGs:         "Instance Security Groups",
	KeyPageDNSDomains:      "DNS Domains",
	KeyPageDNSRecords:      "DNS Records",
	KeyPageSLBList:         "SLB Instances",
	KeyPageSLBDetail:       "SLB Detail",
	KeyPageSLBListeners:    "SLB Listeners",
	KeyPageVServerGroups:   "VServer Groups",
	KeyPageBackendServers:  "Backend Servers",
	KeyPageForwardRules:    "Forwarding Rules",
	KeyPageDefaultServers:  "Default Servers",
	KeyPageOSSBuckets:      "OSS Buckets",
	KeyPageOSSObjects:      "OSS Objects",
	KeyPageOSSDetail:       "OSS Object Detail",
	KeyPageRDSList:         "RDS Instances",
	KeyPageRDSDetail:       "RDS Detail",
	KeyPageRDSDatabases:    "RDS Databases",
	KeyPageRDSAccounts:     "RDS Accounts",
	KeyPageRedisList:       "Redis Instances",
	KeyPageRedisDetail:     "Redis Detail",
	KeyPageRedisAccounts:   "Redis Accounts",
	KeyPageRocketMQList:    "RocketMQ Instances",
	KeyPageRocketMQDetail:  "RocketMQ Detail",
	KeyPageRocketMQTopics:  "RocketMQ Topics",
	KeyPageRocketMQGroups:  "RocketMQ Groups",
	KeyPageResourceFinder:  "Resource Finder",

	// Menu items
	KeyMenuECS:          "(s) ECS Instances",
	KeyMenuECSDesc:      "View ECS instances",
	KeyMenuSG:           "(g) Security Groups",
	KeyMenuSGDesc:       "View ECS security groups",
	KeyMenuDNS:          "(d) DNS Management",
	KeyMenuDNSDesc:      "View AliDNS domains and records",
	KeyMenuSLB:          "(b) SLB Instances",
	KeyMenuSLBDesc:      "View SLB instances",
	KeyMenuOSS:          "(o) OSS Management",
	KeyMenuOSSDesc:      "Browse OSS buckets and objects",
	KeyMenuRDS:          "(r) RDS Instances",
	KeyMenuRDSDesc:      "View RDS instances",
	KeyMenuRedis:        "(i) Redis Instances",
	KeyMenuRedisDesc:    "View Redis instances",
	KeyMenuRocketMQ:     "(m) RocketMQ Instances",
	KeyMenuRocketMQDesc: "View RocketMQ instances",

	// Header
	KeyHeaderProfile: "Profile",
	KeyHeaderRegion:  "Region",

	// Modal
	KeyModalInfo:          "Info",
	KeyModalError:         "Error",
	KeyModalSuccess:       "Success",
	KeyModalOK:            "OK",
	KeyModalConfirm:       "Confirm",
	KeyModalCancel:        "Cancel",
	KeyModalSelectProfile: "Select Profile",
	KeyModalSelectRegion:  "Select Region",
	KeyModalLoading:       "Loading regions with resources...",
	KeyModalResourceFind:  "Resource Finder",
	KeyModalInputPrompt:   "Enter IP address or domain:",
	KeyModalInputExample:  "e.g.: 192.168.1.1 or example.com",
	KeyModalHistory:       "History",
	KeyModalCurrent:       "current",

	// Common columns
	KeyColInstanceID:   "Instance ID",
	KeyColName:         "Name",
	KeyColStatus:       "Status",
	KeyColZone:         "Zone",
	KeyColCPURAM:       "CPU/RAM",
	KeyColPrivateIP:    "Private IP",
	KeyColPublicIP:     "Public IP",
	KeyColExpired:      "Expired",
	KeyColType:         "Type",
	KeyColDescription:  "Description",
	KeyColCreatedAt:    "Created At",
	KeyColRegion:       "Region",

	// ECS specific
	KeyColENIID:        "ENI ID",
	KeyColENIType:      "ENI Type",
	KeyColAttachedInst: "Attached Instance",
	KeyColDiskID:       "Disk ID",
	KeyColDiskType:     "Disk Type",
	KeyColDiskSize:     "Size",
	KeyColDiskCategory: "Category",
	KeyColDevice:       "Device",

	// Security Group
	KeyColSGID:        "Security Group ID",
	KeyColSGName:      "Security Group Name",
	KeyColDirection:   "Direction",
	KeyColProtocol:    "Protocol",
	KeyColPortRange:   "Port Range",
	KeyColSource:      "Source",
	KeyColDestination: "Destination",
	KeyColPolicy:      "Policy",
	KeyColPriority:    "Priority",

	// DNS
	KeyColDomain:      "Domain",
	KeyColRecordType:  "Record Type",
	KeyColRecordValue: "Value",
	KeyColRR:          "Host Record",
	KeyColTTL:         "TTL",
	KeyColLine:        "Line",

	// SLB
	KeyColSLBID:        "SLB ID",
	KeyColAddress:      "Address",
	KeyColAddressType:  "Address Type",
	KeyColSpec:         "Spec",
	KeyColPort:         "Port",
	KeyColBackendPort:  "Backend Port",
	KeyColWeight:       "Weight",
	KeyColServerID:     "Server ID",
	KeyColServerName:   "Server Name",
	KeyColVServerGroup: "VServer Group",

	// OSS
	KeyColBucketName:   "Bucket Name",
	KeyColObjectKey:    "Object Key",
	KeyColSize:         "Size",
	KeyColStorageClass: "Storage Class",
	KeyColLastModified: "Last Modified",

	// RDS
	KeyColDBInstanceID: "DB Instance ID",
	KeyColEngine:       "Engine",
	KeyColDBName:       "Database Name",
	KeyColCharset:      "Charset",
	KeyColAccountName:  "Account Name",
	KeyColAccountType:  "Account Type",
	KeyColPrivilege:    "Privilege",
	KeyColConnString:   "Connection String",

	// Redis
	KeyColConnDomain: "Connection Domain",
	KeyColArchType:   "Architecture",
	KeyColCapacity:   "Capacity",

	// RocketMQ
	KeyColInstanceName: "Instance Name",
	KeyColTopicName:    "Topic Name",
	KeyColGroupID:      "Group ID",
	KeyColMessageType:  "Message Type",
	KeyColGroupType:    "Group Type",

	// Finder
	KeyFinderResult:       "Resource Search Result",
	KeyFinderTotalMatches: "Found %d matching resources",
	KeyFinderNoMatch:      "No matching data",
	KeyFinderECS:          "ECS Instances",
	KeyFinderENI:          "Elastic Network Interfaces",
	KeyFinderSLB:          "Load Balancers",
	KeyFinderDNS:          "DNS Records",
	KeyFinderRDS:          "RDS Instances",
	KeyFinderRedis:        "Redis Instances",
	KeyFinderRocketMQ:     "RocketMQ Instances",

	// ENI Types
	KeyENIPrimary:   "Primary",
	KeyENISecondary: "Secondary",

	// ENI Status
	KeyENIAvailable: "Available",
	KeyENIInUse:     "In Use",
	KeyENIAttaching: "Attaching",
	KeyENIDetaching: "Detaching",
	KeyENIDeleting:  "Deleting",

	// Status
	KeyStatusRunning:  "Running",
	KeyStatusStopped:  "Stopped",
	KeyStatusCreating: "Creating",
	KeyStatusReleased: "Released",
	KeyStatusUnknown:  "Unknown",

	// ECS Detail sections
	KeySectionBasicInfo:      "Basic Information",
	KeySectionConfigInfo:     "Configuration",
	KeySectionBoundRes:       "Bound Resources",
	KeySectionGroupInfo:      "Group Information",
	KeySectionOtherInfo:      "Other Information",
	KeySectionUsageOverview:  "Usage Overview",
	KeySectionStorageOverview: "Storage Overview",

	// ECS Detail labels
	KeyLabelInstanceID:      "Instance ID",
	KeyLabelInstanceName:    "Instance Name",
	KeyLabelInstanceStatus:  "Instance Status",
	KeyLabelAvailZone:       "Availability Zone",
	KeyLabelChargeType:      "Billing Method",
	KeyLabelExpireTime:      "Expiration Time",
	KeyLabelInstanceSpec:    "Instance Type",
	KeyLabelCPUMemory:       "CPU & Memory",
	KeyLabelImageID:         "Image ID",
	KeyLabelOSName:          "Operating System",
	KeyLabelVPC:             "VPC",
	KeyLabelVSwitch:         "VSwitch",
	KeyLabelNetworkType:     "Network Type",
	KeyLabelBandwidth:       "Bandwidth",
	KeyLabelBandwidthCharge: "Bandwidth Billing",
	KeyLabelSecurityGroup:   "Security Groups",
	KeyLabelEIPID:           "EIP ID",
	KeyLabelSecondaryIP:     "Secondary Private IPs",
	KeyLabelResourceGroup:   "Resource Group",
	KeyLabelTags:            "Tags",
	KeyLabelHostname:        "Hostname",
	KeyLabelKeyPair:         "Key Pair",
	KeyLabelSerialNumber:    "Serial Number",
	KeyLabelTotalDisks:      "Total Disks",
	KeyLabelTotalCapacity:   "Total Capacity",

	// Charge types
	KeyChargePrePaid:  "Subscription",
	KeyChargePostPaid: "Pay-As-You-Go",

	// Network types
	KeyNetworkVPC:     "VPC",
	KeyNetworkClassic: "Classic",

	// Disk related
	KeyDiskCloud:           "Basic Cloud Disk",
	KeyDiskCloudEfficiency: "Ultra Cloud Disk",
	KeyDiskCloudSSD:        "SSD Cloud Disk",
	KeyDiskCloudEssd:       "ESSD",
	KeyDiskCloudEssdAuto:   "ESSD AutoPL",
	KeyDiskCloudEssdEntry:  "ESSD Entry",
	KeyDiskSystem:          "System Disk",
	KeyDiskData:            "Data Disk",
	KeyDiskDeleteWithInst:  "Release with Instance",
	KeyDiskKeepAfterInst:   "Retain after Release",
	KeyDiskPortable:        "Supported",
	KeyDiskNotPortable:     "Not Supported",
	KeyColDiskAttribute:    "Attribute",
	KeyColDiskIOPS:         "IOPS",
	KeyColDiskDeleteBehavior: "Release Behavior",
	KeyColDiskPortable:     "Portable",

	// SLB related
	KeyColServerGroupIDName: "Server Group ID/Name",
	KeyColAssocListeners:    "Associated Listeners",
	KeyColAssocRules:        "Associated Rules",
	KeyColBackendCount:      "Backend Count",
	KeyColServerIDName:      "Server ID/Name",
	KeyColPublicPrivateIP:   "Public/Private IP",
	KeyColRemark:            "Remark",
	KeyColURL:               "URL",
	KeySLBVServerGroupList:  "VServer Groups",
	KeySLBForwardingRules:   "Forwarding Rules",
	KeySLBDefaultServerGroup: "Default Server Group",

	// ENI columns
	KeyColMAC: "MAC Address",

	// Count format
	KeyCountSG:    "%d security groups",
	KeyCountENI:   "%d ENIs",
	KeyCountDisks: "%d disks",

	// Common
	KeyActionCopied:  "Copied to clipboard!",
	KeyActionLoading: "Loading...",
	KeyNoData:        "N/A",
}

// Chinese translations
var zhCN = map[string]string{
	// App titles
	KeyAppTitle:            "阿里云 TUI 控制台",
	KeyPageMenu:            "阿里云 TUI 控制台",
	KeyPageECSList:         "ECS 云服务器",
	KeyPageECSDetail:       "ECS 实例详情",
	KeyPageECSJSONDetail:   "ECS JSON 详情",
	KeyPageECSDisks:        "ECS 云盘",
	KeyPageECSENIs:         "ECS 弹性网卡",
	KeyPageSecurityGroups:  "安全组",
	KeyPageSGRules:         "安全组规则",
	KeyPageSGInstances:     "安全组关联实例",
	KeyPageInstSGs:         "实例安全组",
	KeyPageDNSDomains:      "DNS 域名",
	KeyPageDNSRecords:      "DNS 解析记录",
	KeyPageSLBList:         "负载均衡",
	KeyPageSLBDetail:       "SLB 详情",
	KeyPageSLBListeners:    "SLB 监听",
	KeyPageVServerGroups:   "虚拟服务器组",
	KeyPageBackendServers:  "后端服务器",
	KeyPageForwardRules:    "转发规则",
	KeyPageDefaultServers:  "默认服务器组",
	KeyPageOSSBuckets:      "OSS 存储桶",
	KeyPageOSSObjects:      "OSS 对象",
	KeyPageOSSDetail:       "OSS 对象详情",
	KeyPageRDSList:         "RDS 云数据库",
	KeyPageRDSDetail:       "RDS 详情",
	KeyPageRDSDatabases:    "RDS 数据库",
	KeyPageRDSAccounts:     "RDS 账号",
	KeyPageRedisList:       "Redis 实例",
	KeyPageRedisDetail:     "Redis 详情",
	KeyPageRedisAccounts:   "Redis 账号",
	KeyPageRocketMQList:    "RocketMQ 实例",
	KeyPageRocketMQDetail:  "RocketMQ 详情",
	KeyPageRocketMQTopics:  "RocketMQ 主题",
	KeyPageRocketMQGroups:  "RocketMQ 消费组",
	KeyPageResourceFinder:  "资源查找",

	// Menu items
	KeyMenuECS:          "(s) ECS 云服务器",
	KeyMenuECSDesc:      "查看 ECS 云服务器实例",
	KeyMenuSG:           "(g) 安全组",
	KeyMenuSGDesc:       "查看 ECS 安全组",
	KeyMenuDNS:          "(d) DNS 解析",
	KeyMenuDNSDesc:      "查看云解析 DNS 域名和记录",
	KeyMenuSLB:          "(b) 负载均衡",
	KeyMenuSLBDesc:      "查看 SLB 负载均衡实例",
	KeyMenuOSS:          "(o) OSS 存储",
	KeyMenuOSSDesc:      "浏览 OSS 存储桶和对象",
	KeyMenuRDS:          "(r) RDS 数据库",
	KeyMenuRDSDesc:      "查看 RDS 云数据库实例",
	KeyMenuRedis:        "(i) Redis 缓存",
	KeyMenuRedisDesc:    "查看 Redis 实例",
	KeyMenuRocketMQ:     "(m) RocketMQ 消息队列",
	KeyMenuRocketMQDesc: "查看 RocketMQ 实例",

	// Header
	KeyHeaderProfile: "配置",
	KeyHeaderRegion:  "地域",

	// Modal
	KeyModalInfo:          "信息",
	KeyModalError:         "错误",
	KeyModalSuccess:       "成功",
	KeyModalOK:            "确定",
	KeyModalConfirm:       "确认",
	KeyModalCancel:        "取消",
	KeyModalSelectProfile: "选择配置",
	KeyModalSelectRegion:  "选择地域",
	KeyModalLoading:       "正在加载有资源的地域...",
	KeyModalResourceFind:  "资源查找",
	KeyModalInputPrompt:   "请输入 IP 地址或域名:",
	KeyModalInputExample:  "例如: 192.168.1.1 或 example.com",
	KeyModalHistory:       "历史",
	KeyModalCurrent:       "当前",

	// Common columns
	KeyColInstanceID:   "实例 ID",
	KeyColName:         "名称",
	KeyColStatus:       "状态",
	KeyColZone:         "可用区",
	KeyColCPURAM:       "CPU/内存",
	KeyColPrivateIP:    "私网 IP",
	KeyColPublicIP:     "公网 IP",
	KeyColExpired:      "到期时间",
	KeyColType:         "类型",
	KeyColDescription:  "描述",
	KeyColCreatedAt:    "创建时间",
	KeyColRegion:       "地域",

	// ECS specific
	KeyColENIID:        "网卡 ID",
	KeyColENIType:      "网卡类型",
	KeyColAttachedInst: "绑定实例",
	KeyColDiskID:       "云盘 ID",
	KeyColDiskType:     "云盘类型",
	KeyColDiskSize:     "容量",
	KeyColDiskCategory: "云盘种类",
	KeyColDevice:       "设备名",

	// Security Group
	KeyColSGID:        "安全组 ID",
	KeyColSGName:      "安全组名称",
	KeyColDirection:   "方向",
	KeyColProtocol:    "协议",
	KeyColPortRange:   "端口范围",
	KeyColSource:      "源地址",
	KeyColDestination: "目标地址",
	KeyColPolicy:      "策略",
	KeyColPriority:    "优先级",

	// DNS
	KeyColDomain:      "域名",
	KeyColRecordType:  "记录类型",
	KeyColRecordValue: "记录值",
	KeyColRR:          "主机记录",
	KeyColTTL:         "TTL",
	KeyColLine:        "解析线路",

	// SLB
	KeyColSLBID:        "SLB ID",
	KeyColAddress:      "IP 地址",
	KeyColAddressType:  "网络类型",
	KeyColSpec:         "规格",
	KeyColPort:         "端口",
	KeyColBackendPort:  "后端端口",
	KeyColWeight:       "权重",
	KeyColServerID:     "服务器 ID",
	KeyColServerName:   "服务器名称",
	KeyColVServerGroup: "虚拟服务器组",

	// OSS
	KeyColBucketName:   "存储桶名称",
	KeyColObjectKey:    "对象键",
	KeyColSize:         "大小",
	KeyColStorageClass: "存储类型",
	KeyColLastModified: "最后修改",

	// RDS
	KeyColDBInstanceID: "实例 ID",
	KeyColEngine:       "数据库类型",
	KeyColDBName:       "数据库名",
	KeyColCharset:      "字符集",
	KeyColAccountName:  "账号名",
	KeyColAccountType:  "账号类型",
	KeyColPrivilege:    "权限",
	KeyColConnString:   "连接地址",

	// Redis
	KeyColConnDomain: "连接地址",
	KeyColArchType:   "架构类型",
	KeyColCapacity:   "容量",

	// RocketMQ
	KeyColInstanceName: "实例名称",
	KeyColTopicName:    "主题名称",
	KeyColGroupID:      "消费组 ID",
	KeyColMessageType:  "消息类型",
	KeyColGroupType:    "消费组类型",

	// Finder
	KeyFinderResult:       "资源查找结果",
	KeyFinderTotalMatches: "共找到 %d 个匹配资源",
	KeyFinderNoMatch:      "暂无匹配数据",
	KeyFinderECS:          "ECS 实例",
	KeyFinderENI:          "弹性网卡",
	KeyFinderSLB:          "负载均衡",
	KeyFinderDNS:          "DNS 记录",
	KeyFinderRDS:          "RDS 实例",
	KeyFinderRedis:        "Redis 实例",
	KeyFinderRocketMQ:     "RocketMQ 实例",

	// ENI Types
	KeyENIPrimary:   "主网卡",
	KeyENISecondary: "辅助网卡",

	// ENI Status
	KeyENIAvailable: "可用",
	KeyENIInUse:     "已绑定",
	KeyENIAttaching: "绑定中",
	KeyENIDetaching: "解绑中",
	KeyENIDeleting:  "删除中",

	// Status
	KeyStatusRunning:  "运行中",
	KeyStatusStopped:  "已停止",
	KeyStatusCreating: "创建中",
	KeyStatusReleased: "已释放",
	KeyStatusUnknown:  "未知",

	// ECS Detail sections
	KeySectionBasicInfo:      "基本信息",
	KeySectionConfigInfo:     "配置信息",
	KeySectionBoundRes:       "绑定资源",
	KeySectionGroupInfo:      "分组信息",
	KeySectionOtherInfo:      "其他信息",
	KeySectionUsageOverview:  "使用率概览",
	KeySectionStorageOverview: "存储概览",

	// ECS Detail labels
	KeyLabelInstanceID:      "实例 ID",
	KeyLabelInstanceName:    "实例名称",
	KeyLabelInstanceStatus:  "实例状态",
	KeyLabelAvailZone:       "所在可用区",
	KeyLabelChargeType:      "付费类型",
	KeyLabelExpireTime:      "到期时间",
	KeyLabelInstanceSpec:    "实例规格",
	KeyLabelCPUMemory:       "CPU & 内存",
	KeyLabelImageID:         "镜像 ID",
	KeyLabelOSName:          "操作系统",
	KeyLabelVPC:             "专有网络",
	KeyLabelVSwitch:         "虚拟交换机",
	KeyLabelNetworkType:     "网络类型",
	KeyLabelBandwidth:       "公网带宽",
	KeyLabelBandwidthCharge: "带宽计费方式",
	KeyLabelSecurityGroup:   "安全组",
	KeyLabelEIPID:           "弹性公网 IP ID",
	KeyLabelSecondaryIP:     "辅助私网 IP",
	KeyLabelResourceGroup:   "资源组",
	KeyLabelTags:            "标签",
	KeyLabelHostname:        "主机名",
	KeyLabelKeyPair:         "密钥对",
	KeyLabelSerialNumber:    "序列号",
	KeyLabelTotalDisks:      "云盘总数",
	KeyLabelTotalCapacity:   "总存储容量",

	// Charge types
	KeyChargePrePaid:  "包年包月",
	KeyChargePostPaid: "按量付费",

	// Network types
	KeyNetworkVPC:     "专有网络",
	KeyNetworkClassic: "经典网络",

	// Disk related
	KeyDiskCloud:           "普通云盘",
	KeyDiskCloudEfficiency: "高效云盘",
	KeyDiskCloudSSD:        "SSD 云盘",
	KeyDiskCloudEssd:       "ESSD 云盘",
	KeyDiskCloudEssdAuto:   "ESSD AutoPL 云盘",
	KeyDiskCloudEssdEntry:  "ESSD Entry 云盘",
	KeyDiskSystem:          "系统盘",
	KeyDiskData:            "数据盘",
	KeyDiskDeleteWithInst:  "随实例释放",
	KeyDiskKeepAfterInst:   "不随盘释放",
	KeyDiskPortable:        "支持",
	KeyDiskNotPortable:     "不支持",
	KeyColDiskAttribute:    "属性",
	KeyColDiskIOPS:         "IOPS",
	KeyColDiskDeleteBehavior: "云盘释放行为",
	KeyColDiskPortable:     "可卸载",

	// SLB related
	KeyColServerGroupIDName: "服务器组ID/名称",
	KeyColAssocListeners:    "关联监听",
	KeyColAssocRules:        "关联转发策略",
	KeyColBackendCount:      "后端服务器数量",
	KeyColServerIDName:      "云服务器ID/名称",
	KeyColPublicPrivateIP:   "公网/内网IP地址",
	KeyColRemark:            "备注",
	KeyColURL:               "URL",
	KeySLBVServerGroupList:  "虚拟服务器组",
	KeySLBForwardingRules:   "转发策略列表",
	KeySLBDefaultServerGroup: "默认服务器组",

	// ENI columns
	KeyColMAC: "MAC 地址",

	// Count format
	KeyCountSG:    "%d 个安全组",
	KeyCountENI:   "%d 个",
	KeyCountDisks: "%d 个",

	// Common
	KeyActionCopied:  "已复制到剪贴板!",
	KeyActionLoading: "加载中...",
	KeyNoData:        "无",
}

// currentLocale caches the current locale
var currentLocale string

// T returns the translation for the given key
// Falls back to English if the key is not found
func T(key string) string {
	locale := GetLocale()
	if trans, ok := translations[locale]; ok {
		if value, ok := trans[key]; ok {
			return value
		}
	}
	// Fallback to English
	if value, ok := enUS[key]; ok {
		return value
	}
	return key
}

// GetLocale returns the current locale
func GetLocale() string {
	if currentLocale == "" {
		currentLocale = config.GetLocale()
	}
	return currentLocale
}

// RefreshLocale forces a refresh of the cached locale
func RefreshLocale() {
	currentLocale = config.GetLocale()
}

// SetLocale sets the locale (for testing purposes)
func SetLocale(locale string) {
	currentLocale = locale
}

// IsChinese returns true if current locale is Chinese
func IsChinese() bool {
	return GetLocale() == config.LocaleZhCN
}

