# Ghost Host Manager - 备份机制概览

## 概述

Ghost Host Manager 实现了一套完整的备份机制，涵盖应用程序配置和Host分组数据的备份和恢复功能。该机制确保用户数据安全，并在出现问题时能够快速恢复。

## 备份机制架构

### 1. 备份目录结构

```
~/.ghost/                    # 应用程序数据目录
├── config.json             # 应用程序配置文件
├── data.json               # Host分组数据文件
└── backups/                # 备份文件目录
    ├── config_YYYYMMDD_HHMMSS.json  # 配置文件备份
    └── YYYY-MM-DD_HH-MM-SS.json     # 数据文件备份
```

### 2. 备份类型

#### 2.1 应用程序配置备份 (Config Backup)
- **目标文件**: `config.json`
- **备份路径**: `~/.ghost/backups/config_YYYYMMDD_HHMMSS.json`
- **触发时机**:
  - 手动调用 [BackupConfig()](file:///D:/ghost/app.go#L127-L129)
  - 配置更改时（可选）

#### 2.2 Host分组数据备份 (Data Backup)
- **目标文件**: `data.json`
- **备份路径**: `~/.ghost/backups/YYYY-MM-DD_HH-MM-SS.json`
- **触发时机**:
  - 手动调用 [BackupData()](file:///D:/ghost/application/host_app.go#L382-L384)
  - 数据更改前的自动备份
  - 恢复操作前的安全备份

## 核心备份功能

### 1. 单一备份功能

#### 配置备份
```go
func (cs *ConfigStorage) BackupConfig() error
```
- 为 `config.json` 创建带时间戳的备份
- 文件名格式: `config_YYYYMMDD_HHMMSS.json`

#### 数据备份
```go
func (cs *ConfigStorage) BackupData() error
```
- 为 `data.json` 创建带时间戳的备份
- 文件名格式: `YYYY-MM-DD_HH-MM-SS.json`

### 2. 组合备份功能

#### 完整备份
```go
func (app *HostApp) BackupAppAndSystemHosts() (string, error)
```
- 仅备份应用数据
- 返回备份状态信息

## 备份清理策略

### 旧备份清理
- **触发时机**: 每次创建新备份后
- **保留数量**: 默认保留最近10个备份（可通过配置调整）
- **清理算法**: 按修改时间排序，删除超出数量限制的最旧备份

### 特殊备份文件
- `pre_restore_*.json`: 恢复操作前的临时备份，用于防止恢复失败时的数据丢失
- 这些文件不会在常规备份列表中显示

## 用户界面集成

### 前端备份操作
- **立即备份按钮**: 调用 [BackupAppAndSystemHosts()](file:///D:/ghost/frontend/wailsjs/go/main/App.d.ts#L8-L8)，创建应用数据备份
- **恢复功能**: 允许从备份文件恢复数据
- **备份列表**: 显示可用的数据备份文件

### 前端备份API
```javascript
// 创建应用数据备份
await BackupAppAndSystemHosts()

// 列出数据备份
await ListDataBackups()

// 从备份恢复
await RestoreData(backupFileName)

// 仅备份配置
await BackupConfig()
```

## 安全措施

### 恢复安全
- 恢复操作前会创建当前数据的临时备份
- 验证备份文件格式有效性
- 防止无效数据写入配置文件

## 配置选项

### 备份配置参数
- `BackupEnabled`: 启用/禁用自动备份功能
- `MaxBackups`: 最大备份文件数量（默认10个）
- 自动备份会在应用Hosts更改前自动触发

## 备份文件格式

### 配置备份 (JSON)
```json
{
  "autoRefresh": false,
  "refreshInterval": 3600,
  "activeGroups": [],
  "backupEnabled": true,
  "maxBackups": 10,
  "createdAt": "...",
  "updatedAt": "..."
}
```

### 数据备份 (JSON)
包含完整的Host分组信息，包括ID、名称、描述、内容、启用状态、远程URL等。