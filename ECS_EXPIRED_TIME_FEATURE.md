# ECS 表格 ExpiredTime 列功能

## 概述

为 ECS 实例列表表格添加了 "Expired Time" (过期时间) 列，用于显示 ECS 实例的过期时间信息。

## 修改内容

### 1. 修改的文件

#### `internal/tui/pages/ecs.go`
- **ECS 实例列表**:
  - 在表头中添加 "Expired Time" 列
  - 从 ECS 实例的 `ExpiredTime` 字段获取过期时间数据
  - 如果过期时间为空，显示 "N/A"

- **安全组关联实例列表**:
  - 同样添加 "Expired Time" 列，保持与主 ECS 列表的一致性
  - 适用于查看使用特定安全组的实例列表

### 2. 表格列结构

**修改前：**
```
Instance ID | Status | Zone | CPU/RAM | Private IP | Public IP | Name
```

**修改后：**
```
Instance ID | Status | Zone | CPU/RAM | Private IP | Public IP | Name | Expired Time
```

### 3. 数据源

- 使用阿里云 ECS API 中实例对象的 `ExpiredTime` 字段
- 该字段包含实例的过期时间信息（主要用于按量付费实例设置的自动释放时间）
- 如果实例没有设置过期时间，字段为空，显示为 "N/A"

## 技术实现

### 字段映射
```go
// Expired Time
expiredTime := "N/A"
if instance.ExpiredTime != "" {
    expiredTime = instance.ExpiredTime
}
```

### 表格行设置
```go
// 使用 Bubble Tea 的 table 组件设置行数据
row := table.Row{instanceId, status, zone, cpuRam, privateIp, publicIp, name, expiredTime}
```

## 影响范围

1. **主 ECS 实例列表**: 所有 ECS 实例现在都显示过期时间
2. **安全组关联实例列表**: 查看使用特定安全组的实例时也显示过期时间
3. **文档更新**: 
   - `README.md` - 更新了 ECS 实例功能描述
   - `USAGE_EXAMPLE.md` - 更新了中文功能说明

## 用户价值

1. **运维管理**: 方便查看实例的过期时间，避免意外的实例自动释放
2. **成本控制**: 帮助识别设置了自动释放时间的按量付费实例
3. **一致性**: 所有显示 ECS 实例的列表都包含相同的列信息

## 测试状态

- ✅ 编译测试通过
- ✅ 代码语法正确
- ✅ 向后兼容，不影响现有功能
- ✅ 文档已更新

## 使用示例

启动应用程序后：
1. 进入 ECS 实例列表
2. 可以看到新增的 "Expired Time" 列
3. 对于设置了自动释放时间的实例，会显示具体的过期时间
4. 对于没有设置过期时间的实例，显示 "N/A"

## 注意事项

- 过期时间主要适用于按量付费实例，包年包月实例通常不设置此字段
- 显示格式为阿里云 API 返回的原始时间格式
- 如果需要更友好的时间显示格式，可在后续版本中进行改进 