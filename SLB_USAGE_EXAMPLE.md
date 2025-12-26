# SLB Feature Usage Example

## Quick Start Guide

### 1. Launch the Application
```bash
./alidash
```

### 2. Navigate to SLB Instances
- From the main menu, select "SLB Instances" (press 'b' or use arrow keys + Enter)
- You'll see a list of all your SLB instances with columns:
  - SLB ID
  - Name  
  - IP Address
  - Type
  - Status

### 3. Explore SLB Features

#### View Listeners
1. Select an SLB instance using arrow keys
2. Press `l` (lowercase L) to view listeners
3. You'll see a table showing:
   - Protocol
   - Port
   - Backend Port
   - Status
   - Health Check
   - Scheduler

#### View Virtual Server Groups
1. From the SLB instances list, select an instance
2. Press `v` to view virtual server groups
3. You'll see a table showing:
   - VServer Group ID
   - VServer Group Name
   - Backend Server Count

#### View Backend Servers
1. From the virtual server groups view
2. Select a virtual server group using arrow keys
3. Press `Enter` to view backend servers
4. You'll see a table showing:
   - Server ID (ECS Instance ID)
   - Port
   - Weight
   - Type
   - Description

### 4. Navigation Tips

#### Key Bindings
- `l`: View listeners (from SLB list)
- `v`: View virtual server groups (from SLB list)
- `Enter`: Navigate deeper or view details
- `q` or `Escape`: Go back to previous view
- `/`: Search within current table
- `yy`: Copy current row data to clipboard
- `j/k`: Navigate up/down (vim-style)
- Arrow keys: Standard navigation

#### Navigation Flow
```
Main Menu → SLB Instances
                ↓
    ┌───────────┼───────────┐
    ↓           ↓           ↓
Listeners   VServer Groups  Details
                ↓
         Backend Servers
```

### 5. Search Functionality
- Press `/` in any table view to start searching
- Type your search term and press Enter
- Use `n` to find next match, `N` for previous match
- Press Escape to exit search mode

### 6. Copy Data
- Navigate to any row in any table
- Press `y` twice quickly (`yy`) to copy the row data to clipboard
- The data is copied in JSON format for easy use in scripts or documentation

## Example Workflow

### Scenario: Investigating Load Balancer Configuration

1. **Start with SLB list**: `./alidash` → Select "SLB Instances"
2. **Find your SLB**: Use `/` to search for "web-lb-prod"
3. **Check listeners**: Press `l` to see what ports are configured
4. **Check backend distribution**: Press `q` to go back, then `v` to view virtual server groups
5. **Inspect backend servers**: Select a VServer group and press Enter
6. **Copy configuration**: Use `yy` to copy server details for documentation

### Scenario: Troubleshooting Backend Server Issues

1. Navigate to SLB instances
2. Find the problematic load balancer
3. Press `v` to view virtual server groups
4. Select the relevant VServer group
5. Press Enter to see backend servers
6. Check server IDs, ports, and weights
7. Use `yy` to copy server information for further investigation

## Advanced Features

### Profile and Region Switching
- Press `P` (uppercase P) from anywhere to switch between different Alibaba Cloud profiles
- Press `R` (uppercase R) from anywhere to switch between different regions
- Region list only shows regions where your account has resources
- Useful when managing multiple accounts or regions

### Error Handling
- If API calls fail, you'll see user-friendly error messages
- Network issues are handled gracefully
- Invalid selections are prevented

## Tips and Best Practices

1. **Use search frequently**: Large SLB deployments can have many instances - use `/` to quickly find what you need

2. **Copy before modifying**: Use `yy` to copy current configurations before making changes elsewhere

3. **Navigate efficiently**: Learn the key bindings (`l`, `v`, `q`) for quick navigation

4. **Check all levels**: Don't forget to check both listeners (`l`) and virtual server groups (`v`) for complete understanding

5. **Use profiles**: If managing multiple environments, set up different profiles and use `P` to switch between them, or `R` to switch regions

## Troubleshooting

### Common Issues

**"No SLB instances found"**
- Check your Alibaba Cloud credentials
- Verify you're in the correct region
- Ensure your account has SLB permissions

**"Failed to fetch listeners/groups"**
- The SLB instance might not have listeners or virtual server groups configured
- Check API permissions for SLB operations

**Navigation not working**
- Make sure you're using lowercase `l` and `v` keys
- Ensure you've selected a row before pressing navigation keys

### Getting Help
- All views support the standard help keys
- Press `?` in most views for context-sensitive help
- Check the main documentation for detailed API information 