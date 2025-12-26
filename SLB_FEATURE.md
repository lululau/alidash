# SLB (Server Load Balancer) Feature Implementation

## Overview

This document describes the implementation of the SLB feature in the tali application, which allows users to navigate through SLB instances, their listeners, virtual server groups, and backend servers.

## Features Implemented

### 1. SLB Instance List View
- Displays all SLB instances with columns: SLB ID, Name, IP Address, Type, Status
- Supports search functionality with vim-like navigation
- Supports yank (copy) functionality with `yy` key combination
- Press Enter to view detailed JSON information of the selected SLB instance

### 2. Key Navigation from SLB List
- **Press `l` (lowercase L)**: Navigate to listeners list for the selected SLB instance
- **Press `v`**: Navigate to virtual server groups list for the selected SLB instance
- **Press Enter**: View detailed JSON information of the selected SLB instance

### 3. SLB Listeners View
- Shows listeners for a specific SLB instance with detailed information
- Displays Protocol, Port, Backend Port, Status, Health Check, Scheduler, 服务器组 columns
- Automatically detects listener type (HTTP, HTTPS, TCP, UDP) and fetches detailed configuration
- Shows associated virtual server group information

### 4. Virtual Server Groups View
- Lists all virtual server groups for a specific SLB instance with detailed information
- Shows VServer Group ID, VServer Group Name, Backend Server Count, 关联监听 columns
- Displays actual backend server count by querying each virtual server group
- Shows which listeners are associated with each virtual server group
- Press Enter to navigate to backend servers list for the selected virtual server group

### 5. Backend Servers View
- Shows backend servers for a specific virtual server group with detailed ECS information
- Displays Server ID, ECS名称, Port, Weight, Type, 内网IP, 公网IP/EIP, Description columns
- Automatically fetches ECS instance details including instance name and IP addresses
- Shows both private IP addresses and public IP addresses (or EIP if available)
- Provides comprehensive view of backend server configuration and network details

## Files Modified/Added

### 1. Service Layer (`internal/service/slb.go`)
- **Added `FetchListeners(loadBalancerId string)`**: Retrieves basic listeners for a specific SLB instance
- **Added `FetchDetailedListeners(loadBalancerId string)`**: Retrieves detailed listener information with protocol-specific details
- **Added `FetchVServerGroups(loadBalancerId string)`**: Retrieves virtual server groups for a specific SLB instance  
- **Added `FetchDetailedVServerGroups(loadBalancerId string)`**: Retrieves detailed virtual server group information with backend counts and associated listeners
- **Added `FetchVServerGroupBackendServers(vServerGroupId string)`**: Retrieves backend servers for a specific virtual server group
- **Added `FetchDetailedBackendServers(vServerGroupId string, ecsClient)`**: Retrieves backend servers with ECS instance details
- **Added helper methods**: `fetchHTTPListenerDetail`, `fetchHTTPSListenerDetail`, `fetchTCPListenerDetail`, `fetchUDPListenerDetail`, `getECSInstanceDetail`

### 2. Page Types (`internal/tui/types/pages.go`)
- Added `PageSLBListeners`
- Added `PageSLBVServerGroups`
- Added `PageSLBBackendServers`

### 3. SLB Pages (`internal/tui/pages/slb.go`)
- **SLBListModel**: SLB instances list with navigation
- **SLBListenersModel**: Listeners table with protocol, backend port, status, health check, scheduler, and server group information
- **SLBVServerGroupsModel**: Virtual server groups table with backend count and associated listeners
- **SLBBackendServersModel**: Backend servers table with ECS instance name, private IP, and public IP/EIP information

### 4. Application State (`internal/tui/app.go`)
- Uses Bubble Tea's Model-Update-View architecture
- Page models stored in the main application model
- Navigation handled via message passing

### 5. Navigation (`internal/tui/app.go`)
- **Added `setupSlbKeyHandlers()`**: Sets up key handlers for SLB-specific actions
  - `l` key: Navigate to listeners view
  - `v` key: Navigate to virtual server groups view
- **Added `switchToSlbListenersView()`**: Navigation method for listeners view
- **Added `switchToSlbVServerGroupsView()`**: Navigation method for virtual server groups view
- **Added `switchToSlbVServerGroupBackendServersView()`**: Navigation method for backend servers view
- **Updated escape and back key handlers** to support navigation between SLB pages

## Navigation Flow

```
Main Menu
    ↓ (Select SLB Instances)
SLB Instances List
    ↓ (Press 'l')          ↓ (Press 'v')           ↓ (Press Enter)
SLB Listeners View    Virtual Server Groups    SLB Detail View
                           ↓ (Press Enter)
                    Backend Servers View
```

## Key Bindings

### In SLB Instances List:
- `l`: View listeners for selected SLB instance
- `v`: View virtual server groups for selected SLB instance
- `Enter`: View detailed JSON information
- `yy`: Copy selected row data to clipboard
- `/`: Search within the table
- `q` or `Escape`: Go back to main menu

### In Listeners/VServer Groups/Backend Servers Views:
- `yy`: Copy selected row data to clipboard
- `/`: Search within the table
- `q` or `Escape`: Go back to SLB instances list
- `Enter`: (In VServer Groups) Navigate to backend servers

## API Integration

The implementation uses the Alibaba Cloud SLB and ECS SDKs with the following API calls:

### SLB API Calls:
1. **DescribeLoadBalancerAttribute**: Retrieves basic listener information for an SLB instance
2. **DescribeLoadBalancerHTTPListenerAttribute**: Retrieves detailed HTTP listener configuration
3. **DescribeLoadBalancerHTTPSListenerAttribute**: Retrieves detailed HTTPS listener configuration
4. **DescribeLoadBalancerTCPListenerAttribute**: Retrieves detailed TCP listener configuration
5. **DescribeLoadBalancerUDPListenerAttribute**: Retrieves detailed UDP listener configuration
6. **DescribeVServerGroups**: Retrieves virtual server groups for an SLB instance
7. **DescribeVServerGroupAttribute**: Retrieves backend servers and details for a virtual server group

### ECS API Calls:
1. **DescribeInstances**: Retrieves ECS instance details including name, private IP, and public IP/EIP information for backend servers

## Error Handling

- All API calls include proper error handling with user-friendly error messages
- Network errors and API failures are displayed in modal dialogs
- Invalid selections are handled gracefully

## Future Enhancements

1. **Enhanced Listener Details**: Implement specific API calls for different listener types (HTTP, HTTPS, TCP, UDP) to show detailed configuration
2. **Backend Server Health Status**: Add real-time health check status for backend servers
3. **Load Balancer Rules**: Add support for viewing and managing forwarding rules
4. **SSL Certificate Information**: Display SSL certificate details for HTTPS listeners
5. **Real-time Metrics**: Integrate with CloudMonitor to show performance metrics

## Testing

The implementation has been tested with:
- Compilation verification
- Basic application startup
- Navigation flow testing
- Error handling verification

## Usage Example

1. Start the application: `./tali`
2. Select "SLB Instances" from the main menu
3. Navigate through the SLB list using arrow keys
4. Press `l` to view listeners or `v` to view virtual server groups
5. In virtual server groups view, press Enter to see backend servers
6. Use `q` or Escape to navigate back through the hierarchy 