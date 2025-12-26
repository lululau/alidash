# Security Groups Feature Implementation

## Overview

This document describes the implementation of the Security Groups feature in the alidash application, following the same pattern as the existing ECS module.

## Files Modified/Added

### 1. Service Layer (`internal/service/ecs.go`)
- Added `FetchSecurityGroups()` method to the ECSService
- Implements pagination to retrieve all security groups
- Uses the same error handling pattern as `FetchInstances()`

### 2. Page Types (`internal/tui/types/pages.go`)
- Added `PageSecurityGroups`
- Added `PageSecurityGroupRules`
- Added `PageSecurityGroupInstances`
- Added `PageInstanceSecurityGroups`

### 3. Security Groups Pages (`internal/tui/pages/securitygroups.go`)
- **SecurityGroupsModel**: Security groups list with columns: Security Group ID, Name, Description, VPC ID, Type, Creation Time
- **SecurityGroupRulesModel**: Security group rules list
- Supports navigation to view instances using a security group
- Follows the Bubble Tea Model-Update-View pattern

### 4. Main Menu (`internal/tui/pages/menu.go`)
- Added security groups option to the main menu
- Assigned shortcut key 'g' for security groups

### 5. Application State (`internal/tui/app.go`)
- Page models managed in the main Model struct
- Navigation via message passing
- Data caching with automatic clear on profile/region switch

### 6. Navigation
- Navigation handled via `NavigateMsg` and `GoBackMsg` messages
- Back navigation using page stack
- Automatic data clearing on profile switch

### 7. Documentation (`README.md`)
- Updated main menu options to include security groups
- Added security groups feature description
- Updated required permissions to include `ecs:DescribeSecurityGroups`

## Features Implemented

### Security Groups List View
- **Table Columns**: Security Group ID, Name, Description, VPC ID, Type, Creation Time
- **Navigation**: vim-style navigation (j/k keys)
- **Search**: Full search functionality with `/` key
- **Selection**: Enter key to view details
- **Copy**: Double-y (yy) to copy row data to clipboard

### Security Groups Detail View
- **JSON Display**: Complete security group data in formatted JSON
- **Search**: Search within JSON data
- **Copy**: Double-y (yy) to copy complete JSON to clipboard
- **Edit**: 'e' key to open in nvim editor
- **Mouse Support**: Text selection with mouse

### Navigation
- **From Main Menu**: Press 'g' or select "Security Groups"
- **Back Navigation**: 'q' or Escape to go back
- **Detail Navigation**: Enter on any security group to view details

## Integration Points

The security groups feature integrates seamlessly with existing functionality:

1. **Profile Switching**: Security groups cache is cleared when switching profiles
2. **Search System**: Uses the same search infrastructure as other modules
3. **Copy/Edit System**: Integrates with the existing clipboard and nvim editing features
4. **Error Handling**: Uses the same error modal system
5. **UI Consistency**: Follows the same visual and interaction patterns

## API Permissions Required

The feature requires the following additional permission:
- `ecs:DescribeSecurityGroups`

## Testing

The implementation has been tested for:
- ✅ Compilation without errors
- ✅ Integration with existing codebase
- ✅ Consistent code patterns
- ✅ Documentation updates

## Usage

1. Start the application: `./alidash`
2. From the main menu, select "Security Groups" or press 'g'
3. Navigate the list using j/k keys or arrow keys
4. Press Enter on any security group to view details
5. Use search functionality with '/' key
6. Copy data with 'yy' (double-y)
7. Edit JSON with 'e' key in detail view
8. Navigate back with 'q' or Escape

The security groups feature is now fully integrated and ready for use! 