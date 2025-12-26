# 编辑器和分页器配置功能

## 概述

现在 alidash 支持可配置的编辑器和分页器，用户可以通过配置文件或环境变量来自定义使用的编辑器和分页器程序。

## 新增功能

### 按键绑定

- **`e` 键**: 使用配置的编辑器打开JSON数据
- **`v` 键**: 使用配置的分页器查看JSON数据

### 配置优先级

#### 编辑器选择优先级
1. 配置文件顶层的 `editor` 字段
2. `VISUAL` 环境变量
3. `EDITOR` 环境变量  
4. 默认使用 `vim`

#### 分页器选择优先级
1. 配置文件顶层的 `pager` 字段
2. `PAGER` 环境变量
3. 默认使用 `less`

## 配置示例

### 在配置文件中设置

在 `~/.aliyun/config.json` 中在顶层添加 `editor` 和 `pager` 字段：

```json
{
  "current": "default",
  "editor": "code",
  "pager": "less -R",
  "profiles": [
    {
      "name": "default",
      "mode": "AK",
      "access_key_id": "your_access_key",
      "access_key_secret": "your_secret_key",
      "region_id": "cn-hangzhou"
    }
  ]
}
```

### 通过环境变量设置

```bash
# 设置编辑器
export VISUAL="code"
export EDITOR="vim"

# 设置分页器
export PAGER="less -R"
```

## 支持的编辑器示例

- `vim` - 默认编辑器
- `nvim` - Neovim
- `code` - Visual Studio Code
- `nano` - Nano编辑器
- `emacs` - Emacs编辑器
- `subl` - Sublime Text

## 支持的分页器示例

- `less` - 默认分页器
- `less -R` - 支持颜色的less
- `more` - 简单分页器
- `cat` - 直接输出（不分页）

## 使用方法

1. 在任何详情页面（如ECS实例详情、RDS详情等）
2. 按 `e` 键使用编辑器打开JSON数据进行编辑
3. 按 `v` 键使用分页器查看JSON数据
4. 按 `yy` 复制JSON数据到剪贴板

## 注意事项

- 编辑器和分页器命令支持参数，如 `"less -R"` 或 `"code --wait"`
- 确保配置的编辑器和分页器程序已安装在系统中
- 编辑器会在临时文件中打开JSON数据，编辑完成后临时文件会被自动删除
- 分页器用于只读查看，不会修改原始数据
- `editor` 和 `pager` 配置项位于配置文件的顶层，对所有 profiles 生效 