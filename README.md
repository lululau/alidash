# alidash

A Terminal User Interface (TUI) application for managing Alibaba Cloud resources. Navigate and inspect your cloud infrastructure directly from the command line with vim-style navigation and powerful search capabilities.

## Features

### Supported Services
- **ECS Instances**: View instance details with zone, CPU/RAM configuration, private/public IPs, and full JSON details
- **Security Groups**: Browse security groups, view rules, and see associated instances
- **DNS Management**: Browse AliDNS domains and their DNS records
- **SLB (Server Load Balancer)**: Monitor SLB instances, listeners, VServer groups, and backend servers
- **OSS (Object Storage)**: Browse OSS buckets and objects with pagination
- **RDS (Relational Database)**: Inspect RDS instances, databases, and accounts
- **Redis**: View Redis instances and accounts
- **RocketMQ**: Browse RocketMQ instances, topics, and consumer groups

### Interactive Features
- **Vim-style Navigation**: Use j/k keys for navigation, Enter to select
- **Powerful Search**: Search across all data with `/` key, navigate results with n/N
- **Data Export**: Copy any data as JSON to clipboard with `yy` (double-y)
- **External Editing**: Edit JSON data in nvim with `e` key
- **Mouse Support**: Text selection in detail views
- **Profile Management**: Switch between multiple Alibaba Cloud profiles
- **Real-time Mode Line**: Shows current profile and contextual shortcuts
- **Pagination**: Navigate large datasets with intuitive controls

## Prerequisites

- Go 1.22.2 or later
- Valid Alibaba Cloud account with API access
- Access Key ID and Secret with appropriate permissions
- (Optional) nvim for external JSON editing

## Installation

### Using Homebrew (macOS)

The easiest way to install alidash on macOS is using Homebrew:

```bash
# Add the tap (replace with your actual repository)
brew tap lululau/alidash

# Install alidash
brew install alidash
```

Or install directly from the formula URL:

```bash
brew install https://raw.githubusercontent.com/lululau/alidash/main/alidash.rb
```

For detailed Homebrew setup instructions, including creating your own tap, see [HOMEBREW.md](HOMEBREW.md).

**For maintainers**: Use the provided script to generate the formula:
```bash
./scripts/generate-formula.sh 1.0.0 lululau
```

### From Source

1. Clone the repository:
```bash
git clone <repository-url>
cd alidash
```

2. Build the application:
```bash
go build -o alidash cmd/main.go
```

3. (Optional) Install globally:
```bash
go install ./cmd
```

## Configuration

### Alibaba Cloud CLI Configuration

The application uses the standard Alibaba Cloud CLI configuration format. Create a configuration file at `~/.aliyun/config.json`:

```json
{
  "current": "default",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your-access-key-id",
      "access_key_secret": "your-access-key-secret",
      "region_id": "cn-hangzhou",
      "oss_endpoint": "oss-cn-hangzhou.aliyuncs.com"
    },
    {
      "name": "production",
      "mode": "AK",
      "access_key_id": "prod-access-key-id",
      "access_key_secret": "prod-access-key-secret",
      "region_id": "cn-shanghai",
      "oss_endpoint": "oss-cn-shanghai.aliyuncs.com"
    }
  ]
}
```

### Configuration Fields

- **name**: Profile name (used for identification)
- **mode**: Authentication mode (use "AK" for Access Key)
- **access_key_id**: Your Alibaba Cloud Access Key ID
- **access_key_secret**: Your Alibaba Cloud Access Key Secret
- **region_id**: Target region ID
- **oss_endpoint**: OSS endpoint (optional, auto-generated if not specified)

### Common Region IDs
- `cn-hangzhou` - China (Hangzhou)
- `cn-shanghai` - China (Shanghai)
- `cn-beijing` - China (Beijing)
- `cn-shenzhen` - China (Shenzhen)
- `us-west-1` - US West (Silicon Valley)
- `ap-southeast-1` - Asia Pacific (Singapore)

## Usage

### Running the Application

```bash
./alidash
```

Or if installed globally:
```bash
alidash
```

### Navigation and Controls

The application uses vim-style keyboard navigation with contextual shortcuts displayed in the mode line at the bottom.

#### Global Controls
- `Q` - Quit application (uppercase Q)
- `q` or `Esc` - Go back to previous screen/menu
- `P` - Open profile selection dialog (uppercase P)
- `R` - Open region selection dialog (uppercase R)
- `Ctrl+C` - Force quit

#### Main Menu Navigation
- `j/k` or `↑/↓` - Navigate up/down
- `Enter` - Select current service
- Service shortcuts:
  - `s` - ECS Instances
  - `g` - Security Groups  
  - `d` - DNS Management
  - `b` - SLB Instances
  - `o` - OSS Management
  - `r` - RDS Instances
  - `i` - Redis Instances
  - `m` - RocketMQ Instances

#### List Navigation
- `j/k` or `↑/↓` - Move up/down in lists
- `Enter` - Select item for detailed view or sub-navigation
- `/` - Enter search mode
- `n/N` - Navigate to next/previous search result
- `yy` - Copy current row data as JSON to clipboard

#### Service-Specific Shortcuts

**ECS Instances:**
- `g` - View security groups for selected instance

**Security Groups:**
- `Enter` - View security group rules
- `s` - View instances using this security group

**SLB Instances:**
- `l` - View listeners for selected SLB
- `v` - View VServer groups for selected SLB

**RDS Instances:**
- `D` - View databases for selected RDS instance
- `A` - View accounts for selected RDS instance

**Redis Instances:**
- `A` - View accounts for selected Redis instance

**RocketMQ Instances:**
- `T` - View topics for selected RocketMQ instance
- `G` - View consumer groups for selected RocketMQ instance

#### Detail View Controls
- `q/Esc` - Go back to list view
- `yy` - Copy complete JSON data to clipboard
- `e` - Open JSON data in nvim for editing
- `/` - Search within JSON data
- `n/N` - Navigate search results within JSON
- Mouse selection supported for copying text

#### Search Functionality
- `/` - Enter search mode (vim-style search bar appears at bottom)
- `Enter` - Execute search and highlight matches
- `Esc` - Exit search mode
- `n` - Go to next search result
- `N` - Go to previous search result
- Search is case-insensitive by default
- Works in all table views and JSON detail views

#### Profile Management
- Press `P` to open profile selection dialog
- Use `j/k` to navigate available profiles
- Use `/` to filter profiles by name
- Press `Enter` to select a profile
- Press `q` or `Esc` to cancel
- After switching profiles:
  - All client connections are recreated
  - All cached data is cleared
  - Region resets to the profile's default region
  - Application returns to main menu
  - New credentials take effect immediately

#### Region Management
- Press `R` to open region selection dialog
- The region list only shows regions where your account has resources
- Region list is cached per profile (7-day expiration, stored in `~/.aliyun/region_cache.json`)
- Use `j/k` to navigate available regions
- Use `/` to filter regions
- Press `Enter` to select a region
- Press `q` or `Esc` to cancel
- After switching regions:
  - All client connections are recreated with the new region
  - All cached data is cleared
  - Application returns to main menu
  - New region takes effect immediately

#### OSS Object Pagination
- `[` - Previous page
- `]` - Next page
- `0` - Go to first page
- Page information displayed in mode line

### Service Details

#### ECS Instances
- Lists all ECS instances with ID, status, zone, CPU/RAM configuration, private IP, public IP, name, and expired time
- Press `g` on any instance to view its security groups
- Select an instance to view complete JSON details including:
  - Instance specifications and configuration
  - Network configuration and IP addresses
  - Security groups and network interfaces
  - Storage details and disk information
  - All available metadata

#### Security Groups
- Lists all ECS security groups with ID, name, description, VPC ID, type, and creation time
- Press `Enter` to view security group rules (ingress/egress)
- Press `s` to view instances using this security group
- Select for complete JSON configuration including:
  - Security group rules and policies
  - Associated instances and network interfaces
  - VPC and network configuration
  - All available metadata

#### DNS Management
- Browse all domains in your account
- View record count and version information
- Select a domain to view all DNS records
- See record types (A, CNAME, MX, etc.), values, TTL, and status
- Full JSON details for domains and records

#### SLB (Server Load Balancer)
- List all SLB instances with ID, name, IP address, type, and status
- Press `l` to view listeners for selected SLB
- Press `v` to view VServer groups for selected SLB
- Navigate to backend servers from VServer groups
- Complete JSON configuration including:
  - Load balancer specifications
  - Network configuration and IP addresses
  - Health check settings
  - All available metadata

#### OSS (Object Storage)
- Browse all OSS buckets with name, location, creation date, and storage class
- Select a bucket to view all objects with pagination
- Object details include key, size, last modified date, storage class, and ETag
- Navigate large object lists with `[`, `]`, and `0` keys
- Select an object to view complete JSON metadata

#### RDS (Relational Database)
- Browse all RDS database instances
- View engine type, version, instance class, and status
- Press `D` to view databases for selected RDS instance
- Press `A` to view accounts for selected RDS instance
- Complete JSON configuration including:
  - Connection strings and ports
  - Storage and backup information
  - Network configuration
  - Maintenance windows and settings
  - All available metadata

#### Redis
- Browse all Redis instances with version, class, and status information
- Press `A` to view accounts for selected Redis instance
- Complete JSON configuration including:
  - Connection information
  - Memory and performance settings
  - Network configuration
  - All available metadata

#### RocketMQ
- Browse all RocketMQ instances
- Press `T` to view topics for selected instance
- Press `G` to view consumer groups for selected instance
- Complete JSON configuration including:
  - Instance specifications
  - Network configuration
  - Topic and group management details
  - All available metadata

## Required Permissions

Your Alibaba Cloud Access Key needs the following permissions:

- **ECS**: `ecs:DescribeInstances`, `ecs:DescribeSecurityGroups`, `ecs:DescribeSecurityGroupAttribute`
- **DNS**: `alidns:DescribeDomains`, `alidns:DescribeDomainRecords`
- **SLB**: `slb:DescribeLoadBalancers`, `slb:DescribeLoadBalancerAttribute`, `slb:DescribeVServerGroups`, `slb:DescribeVServerGroupAttribute`
- **RDS**: `rds:DescribeDBInstances`, `rds:DescribeDatabases`, `rds:DescribeAccounts`
- **Redis**: `r-kvstore:DescribeInstances`, `r-kvstore:DescribeAccounts`
- **RocketMQ**: `ons:OnsInstanceInServiceList`, `ons:OnsTopicList`, `ons:OnsGroupList`
- **OSS**: `oss:ListBuckets`, `oss:ListObjects`, `oss:GetObjectMeta`

## Troubleshooting

### Common Issues

1. **Authentication Error**
   - Verify your Access Key ID and Secret are correct in `~/.aliyun/config.json`
   - Ensure the keys have the required permissions
   - Check that the region ID is valid and accessible

2. **Configuration File Not Found**
   - Ensure `~/.aliyun/config.json` exists and is properly formatted
   - Check file permissions (should be readable by your user)
   - Verify JSON syntax is correct

3. **No Resources Found**
   - Verify you're using the correct region ID
   - Ensure resources exist in the specified region
   - Check that your account has access to the resources
   - Verify the profile has the correct permissions

4. **Profile Switching Issues**
   - Ensure all profiles in config.json have required fields
   - Check that profile names are unique
   - Verify the "current" field points to an existing profile

5. **Network Issues**
   - Ensure you have internet connectivity
   - Check if your firewall allows HTTPS traffic
   - Verify the OSS endpoint is correct for your region

6. **nvim Editor Issues**
   - Ensure nvim is installed and in your PATH
   - Check that temporary file creation works in your system
   - Verify nvim can be launched from the terminal

### Debug Mode

For additional debugging information, you can check the application logs. The application will display error messages in modal dialogs for most issues.

## Development

### Dependencies

The project uses Go modules. Key dependencies include:

- `github.com/aliyun/alibaba-cloud-sdk-go` - Alibaba Cloud SDK
- `github.com/aliyun/aliyun-oss-go-sdk` - OSS SDK
- `github.com/charmbracelet/bubbletea` - Terminal UI framework (Elm architecture)
- `github.com/charmbracelet/bubbles` - TUI components (table, list, viewport, etc.)
- `github.com/charmbracelet/lipgloss` - Style definitions for terminal rendering

### Building from Source

```bash
# Download dependencies
go mod download

# Run tests (if any)
go test ./...

# Build
go build -o alidash cmd/main.go
```

### Project Structure

```
alidash/
├── cmd/                    # Application entry point
├── internal/
│   ├── client/            # Alibaba Cloud client management
│   ├── config/            # Configuration loading and management
│   ├── service/           # Service layer for API calls (including RegionService)
│   └── tui/               # Terminal user interface (Bubble Tea)
│       ├── components/    # Reusable UI components (table, modal, header, etc.)
│       ├── pages/         # Page models for each service
│       ├── types/         # Shared types and messages
│       ├── app.go         # Main application model
│       ├── keys.go        # Key bindings
│       ├── messages.go    # Message types
│       └── styles.go      # Style definitions
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md             # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[Add your license information here]

## Support

[Add support information here]