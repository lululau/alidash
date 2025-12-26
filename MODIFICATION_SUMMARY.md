# 编辑器和分页器配置功能 - 修改总结

## 修改概述

根据用户需求，实现了可配置的编辑器和分页器功能，替换了原来硬编码的nvim编辑器。配置项位于配置文件的顶层，对所有 profiles 生效。

## 修改的文件

### 1. `internal/config/config.go`

**结构体修改:**
- 从 `ConfigProfile` 结构体中移除了 `Editor` 和 `Pager` 字段
- 在 `AliyunConfig` 结构体顶层添加了 `Editor` 和 `Pager` 字段
- `Config` 结构体保持 `Editor` 和 `Pager` 字段不变

**新增函数:**
- `GetEditor()` - 获取编辑器命令，按优先级：配置文件顶层 → VISUAL环境变量 → EDITOR环境变量 → vim
- `GetPager()` - 获取分页器命令，按优先级：配置文件顶层 → PAGER环境变量 → less

**加载逻辑修改:**
- `LoadAliyunConfig()` 函数现在从配置文件顶层读取 `editor` 和 `pager` 字段

### 2. `internal/tui/` - Bubble Tea UI 框架

**UI 架构:**
- 使用 Bubble Tea (Elm architecture) 重构整个 UI 层
- `internal/tui/components/` - 可复用 UI 组件（表格、模态框、视口等）
- `internal/tui/pages/` - 各服务的页面模型
- `internal/tui/app.go` - 主应用程序模型

**编辑器和分页器功能:**
- `OpenInEditor()` - 使用配置的编辑器打开 JSON 数据
- `OpenInPager()` - 使用配置的分页器查看 JSON 数据
- 支持带参数的命令（如 `"less -R"`）
- 自动创建和清理临时文件

**详情页面快捷键:**
- `e` 键: 在外部编辑器中编辑 JSON
- `v` 键: 在分页器中查看 JSON
- `yy`: 复制 JSON 到剪贴板

## 新增按键绑定

- **`e` 键**: 使用配置的编辑器打开JSON数据
- **`v` 键**: 使用配置的分页器查看JSON数据

## 配置结构

### 新的配置文件结构
```json
{
  "current": "default",
  "editor": "code",
  "pager": "less -R",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your_key",
      "access_key_secret": "your_secret",
      "region_id": "cn-hangzhou"
    }
  ]
}
```

### 配置优先级

#### 编辑器选择
1. 配置文件顶层的 `editor` 字段
2. `VISUAL` 环境变量
3. `EDITOR` 环境变量
4. 默认 `vim`

#### 分页器选择
1. 配置文件顶层的 `pager` 字段
2. `PAGER` 环境变量
3. 默认 `less`

## 向后兼容性

- 如果配置文件中没有顶层的 `editor` 和 `pager` 字段，功能仍然正常工作
- 现有的配置文件无需修改即可继续使用
- 保持了所有原有的功能和按键绑定
- 旧的 profile 级别的 `editor` 和 `pager` 字段会被忽略

## 测试状态

- ✅ 编译成功
- ✅ 所有linter错误已修复
- ✅ 保持了原有功能的完整性
- ✅ 新功能按预期工作

## 使用示例

```json
{
  "current": "default",
  "editor": "code",
  "pager": "less -R",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your_key",
      "access_key_secret": "your_secret",
      "region_id": "cn-hangzhou"
    }
  ]
}
```

用户现在可以：
1. 在详情页面按 `e` 使用自定义编辑器编辑JSON
2. 在详情页面按 `v` 使用自定义分页器查看JSON
3. 通过配置文件顶层或环境变量自定义编辑器和分页器
4. 编辑器和分页器配置对所有 profiles 生效 