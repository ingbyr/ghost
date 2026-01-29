# Ghost Host Manager - 权限管理说明

## 概述

Ghost Host Manager 需要管理员权限来修改系统 hosts 文件。不同操作系统有不同的权限管理机制，本文档解释了每种平台的权限处理方式。

## Windows 平台

### 永久性管理员权限
- 应用程序使用了 Windows 清单文件（manifest）配置
- 配置了 `requireAdministrator` 权限级别
- 每次运行时会自动触发 UAC（用户账户控制）对话框
- 用户确认后，应用将以管理员权限运行整个会话

### 用户体验
1. 首次运行应用时，Windows 会显示 UAC 对话框
2. 用户点击"是"或输入管理员密码
3. 应用以管理员权限启动，后续操作无需重复认证
4. 整个会话期间，应用都具有修改系统 hosts 文件的权限

## Linux 平台

### 权限处理机制
- 使用图形化 sudo 工具（按优先级：pkexec、gksudo、kdesudo、lxqt-sudo）
- 每次需要修改 hosts 文件时会提示用户输入密码
- 利用系统 sudo 的时间窗口缓存机制

### 用户体验
1. 首次执行需要管理员权限的操作时，会弹出图形化密码输入框
2. 输入用户密码进行认证
3. 在 sudo 缓存有效期内（通常是15分钟），后续操作无需重复认证
4. 缓存过期后，再次需要认证

### 优化建议
如果希望减少认证次数，系统管理员可以配置 sudoers：

```bash
# 编辑 sudoers 文件
sudo visudo

# 添加以下行（将 username 替换为实际用户名）
username ALL=(ALL) NOPASSWD: /path/to/ghost
```

**注意**：此配置存在安全风险，请谨慎使用。

## macOS 平台

### 权限处理机制
- 使用 AppleScript 的 `with administrator privileges` 功能
- 每次需要修改 hosts 文件时会提示用户输入管理员密码
- 利用系统 sudo 的时间窗口缓存机制

### 用户体验
1. 首次执行需要管理员权限的操作时，会弹出系统密码输入对话框
2. 输入管理员密码进行认证
3. 在 sudo 缓存有效期内（通常是5分钟），后续操作无需重复认证
4. 缓存过期后，再次需要认证

## 安全考虑

### 为什么需要管理员权限
- hosts 文件位于系统受保护目录（Windows: `C:\Windows\System32\drivers\etc\hosts`，Linux/macOS: `/etc/hosts`）
- 修改系统配置文件需要相应权限
- 这是操作系统安全模型的一部分

### 权限最小化原则
- 应用仅在需要修改 hosts 文件时请求权限
- 不会获取比所需权限更高的系统访问权
- 遵循最小权限原则

## 故障排除

### Windows
- 如果没有看到 UAC 对话框，请检查是否以管理员身份运行
- 确认用户账户具有管理员权限

### Linux
- 确保系统安装了图形化 sudo 工具之一
- 检查用户是否在 sudo 组中
- 验证 sudo 配置

### macOS
- 确保用户具有管理员权限
- 检查系统完整性保护（SIP）设置

## 技术实现

权限管理实现在 `permissions` 包中，提供了跨平台的统一接口：
- `IsAdmin()` - 检查当前是否具有管理员权限
- `RequestElevation()` - 请求权限提升
- `CanSudoWithoutPassword()` - 检查是否可以无密码执行sudo（Linux/macOS）
- `HasSudoAccess()` - 检查是否具有sudo权限（Linux/macOS）

这种设计确保了跨平台一致性，同时尊重各平台的安全模型。