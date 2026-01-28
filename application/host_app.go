package application

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"ghost/models"
	"ghost/remote"
	"ghost/storage"
	"ghost/system"
)

// HostApp 主应用程序逻辑
type HostApp struct {
	configStorage *storage.ConfigStorage
	hostManager   *system.HostManager
	mu            sync.RWMutex
	autoRefresh   bool
	stopChan      chan struct{}
}

// NewHostApp 创建新的Host应用程序实例
func NewHostApp() (*HostApp, error) {
	configStorage, err := storage.NewConfigStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config storage: %w", err)
	}

	hostManager := system.NewHostManager()

	app := &HostApp{
		configStorage: configStorage,
		hostManager:   hostManager,
		autoRefresh:   false,
		stopChan:      make(chan struct{}),
	}

	return app, nil
}

// GetHostGroups 获取所有Host分组
func (app *HostApp) GetHostGroups() ([]models.HostGroup, error) {
	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return nil, fmt.Errorf("failed to load host manager: %w", err)
	}

	return manager.Groups, nil
}

// AddHostGroup 添加新的Host分组
func (app *HostApp) AddHostGroup(group models.HostGroup) error {
	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return fmt.Errorf("failed to load host manager: %w", err)
	}

	// 自动生成UUID作为ID
	group.ID = uuid.New().String()

	// 检查ID是否已存在
	for _, existingGroup := range manager.Groups {
		if existingGroup.ID == group.ID {
			return fmt.Errorf("host group with ID %s already exists", group.ID)
		}
	}

	// 验证必要字段（不再需要验证ID，因为是自动生成的）
	if group.Name == "" {
		return fmt.Errorf("group name cannot be empty")
	}

	// 如果是远程组，验证URL
	if group.IsRemote && strings.TrimSpace(group.URL) == "" {
		return fmt.Errorf("remote group URL cannot be empty")
	}

	// 设置默认值
	if group.CreatedAt == "" {
		group.CreatedAt = time.Now().Format(time.RFC3339)
	}
	group.UpdatedAt = time.Now().Format(time.RFC3339)

	manager.Groups = append(manager.Groups, group)
	manager.UpdatedAt = time.Now().Format(time.RFC3339)

	err = app.configStorage.SaveHostManager(manager)
	if err != nil {
		return fmt.Errorf("failed to save host manager: %w", err)
	}

	return nil
}

// UpdateHostGroup 更新Host分组
func (app *HostApp) UpdateHostGroup(group models.HostGroup) error {
	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return fmt.Errorf("failed to load host manager: %w", err)
	}

	updated := false
	for i, existingGroup := range manager.Groups {
		if existingGroup.ID == group.ID {
			// 验证必要字段
			if group.Name == "" {
				return fmt.Errorf("group name cannot be empty")
			}

			// 如果是远程组，验证URL
			if group.IsRemote && strings.TrimSpace(group.URL) == "" {
				return fmt.Errorf("remote group URL cannot be empty")
			}

			// 保留创建时间
			group.CreatedAt = existingGroup.CreatedAt
			group.UpdatedAt = time.Now().Format(time.RFC3339)

			manager.Groups[i] = group
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("host group with ID %s not found", group.ID)
	}

	manager.UpdatedAt = time.Now().Format(time.RFC3339)
	err = app.configStorage.SaveHostManager(manager)
	if err != nil {
		return fmt.Errorf("failed to save host manager: %w", err)
	}

	return nil
}

// DeleteHostGroup 删除Host分组
func (app *HostApp) DeleteHostGroup(id string) error {
	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return fmt.Errorf("failed to load host manager: %w", err)
	}

	updatedGroups := make([]models.HostGroup, 0, len(manager.Groups))
	found := false

	for _, group := range manager.Groups {
		if group.ID == id {
			found = true
			continue
		}
		updatedGroups = append(updatedGroups, group)
	}

	if !found {
		return fmt.Errorf("host group with ID %s not found", id)
	}

	manager.Groups = updatedGroups
	manager.UpdatedAt = time.Now().Format(time.RFC3339)

	err = app.configStorage.SaveHostManager(manager)
	if err != nil {
		return fmt.Errorf("failed to save host manager: %w", err)
	}

	return nil
}

// ToggleHostGroup 启用或禁用Host分组
func (app *HostApp) ToggleHostGroup(id string, enabled bool) error {
	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return fmt.Errorf("failed to load host manager: %w", err)
	}

	updated := false
	for i, group := range manager.Groups {
		if group.ID == id {
			manager.Groups[i].Enabled = enabled
			manager.Groups[i].UpdatedAt = time.Now().Format(time.RFC3339)
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("host group with ID %s not found", id)
	}

	manager.UpdatedAt = time.Now().Format(time.RFC3339)
	err = app.configStorage.SaveHostManager(manager)
	if err != nil {
		return fmt.Errorf("failed to save host manager: %w", err)
	}

	return nil
}

// ApplyHosts 应用所有启用的Host分组到系统
func (app *HostApp) ApplyHosts() error {
	// 检查权限
	if !app.hostManager.HasWritePermission() {
		// 尝试以管理员权限重新启动（仅在必要时）
		err := app.requestAdminPrivileges()
		if err != nil {
			return err
		}
		// 如果重新启动成功，这里不会执行到
		return app.hostManager.RequestAdminPrivileges()
	}

	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return fmt.Errorf("failed to load host manager: %w", err)
	}

	// 构建hostGroups数据结构用于应用到系统
	var hostGroups []map[string]interface{}
	for _, group := range manager.Groups {
		if group.Enabled {
			hostGroup := map[string]interface{}{
				"id":       group.ID,
				"name":     group.Name,
				"content":  group.Content,
				"enabled":  group.Enabled,
				"isRemote": group.IsRemote,
			}
			hostGroups = append(hostGroups, hostGroup)
			log.Printf("Applying group: %s (Remote: %t, Enabled: %t)", group.Name, group.IsRemote, group.Enabled)
		}
	}

	// 使用HostManager的ApplyHostGroups方法，该方法会保留系统原有内容
	err = app.hostManager.ApplyHostGroups(hostGroups)
	if err != nil {
		return fmt.Errorf("failed to apply host groups to system: %w", err)
	}

	log.Printf("Applied %d enabled host groups to system hosts file", len(hostGroups))

	return nil
}

// GetSystemHostsContent 获取系统hosts文件内容
func (app *HostApp) GetSystemHostsContent() (string, error) {
	content, err := app.hostManager.ReadSystemHosts()
	if err != nil {
		return "", fmt.Errorf("failed to read system hosts file: %w", err)
	}
	return content, nil
}

// GetSystemHostPath 获取系统hosts文件路径
func (app *HostApp) GetSystemHostPath() string {
	return app.hostManager.SystemHostPath
}

// RefreshRemoteGroups 刷新所有远程Host组
func (app *HostApp) RefreshRemoteGroups() error {
	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return fmt.Errorf("failed to load host manager: %w", err)
	}

	remoteFetcher := remote.NewRemoteFetcher()
	updated := false

	for i := range manager.Groups {
		group := &manager.Groups[i]
		if group.IsRemote && group.URL != "" {
			log.Printf("Fetching remote content from URL: %s for group: %s", group.URL, group.Name)
			oldContent := group.Content
			err := remoteFetcher.UpdateRemoteHostGroup(group)
			if err != nil {
				log.Printf("Error updating remote group %s from URL %s: %v", group.Name, group.URL, err)
				continue
			}

			// 检查内容是否有变化
			if oldContent != group.Content {
				log.Printf("Remote group %s updated with new content", group.Name)
				updated = true
			} else {
				log.Printf("Remote group %s content unchanged", group.Name)
			}
		}
	}

	if updated {
		manager.UpdatedAt = time.Now().Format(time.RFC3339)
		err = app.configStorage.SaveHostManager(manager)
		if err != nil {
			return fmt.Errorf("failed to save updated host manager: %w", err)
		}
		log.Println("Successfully saved updated host manager with new remote content")
	}

	return nil
}

// StartAutoRefresh 启动自动刷新功能
func (app *HostApp) StartAutoRefresh() error {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.autoRefresh {
		return fmt.Errorf("auto refresh is already running")
	}

	config, err := app.configStorage.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !config.AutoRefresh {
		return fmt.Errorf("auto refresh is disabled in config")
	}

	app.autoRefresh = true
	interval := config.RefreshInterval
	if interval == 0 {
		interval = 3600 // 默认1小时（秒）
	}

	go func() {
		// 将秒数转换为Duration
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Refreshing remote host groups...")
				err := app.RefreshRemoteGroups()
				if err != nil {
					log.Printf("Error refreshing remote groups: %v", err)
				}
			case <-app.stopChan:
				log.Println("Stopping auto refresh...")
				return
			}
		}
	}()

	return nil
}

// StopAutoRefresh 停止自动刷新功能
func (app *HostApp) StopAutoRefresh() {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.autoRefresh {
		close(app.stopChan)
		app.autoRefresh = false
		app.stopChan = make(chan struct{})
	}
}

// GetConfig 获取应用程序配置
func (app *HostApp) GetConfig() (*models.AppConfig, error) {
	return app.configStorage.LoadConfig()
}

// UpdateConfig 更新应用程序配置
func (app *HostApp) UpdateConfig(config *models.AppConfig) error {
	config.UpdatedAt = time.Now().Format(time.RFC3339)
	return app.configStorage.SaveConfig(config)
}

// BackupConfig 创建配置备份
func (app *HostApp) BackupConfig() error {
	return app.configStorage.BackupConfig()
}

// BackupData 创建数据备份
func (app *HostApp) BackupData() error {
	return app.configStorage.BackupData()
}

// GetHostGroup 获取指定ID的Host分组
func (app *HostApp) GetHostGroup(id string) (*models.HostGroup, error) {
	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return nil, fmt.Errorf("failed to load host manager: %w", err)
	}

	for _, group := range manager.Groups {
		if group.ID == id {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("host group with ID %s not found", id)
}

// RefreshRemoteGroup 刷新指定的远程Host分组
func (app *HostApp) RefreshRemoteGroup(id string) error {
	manager, err := app.configStorage.LoadHostManager()
	if err != nil {
		return fmt.Errorf("failed to load host manager: %w", err)
	}

	var targetGroup *models.HostGroup
	for i := range manager.Groups {
		if manager.Groups[i].ID == id {
			targetGroup = &manager.Groups[i]
			break
		}
	}

	if targetGroup == nil {
		return fmt.Errorf("host group with ID %s not found", id)
	}

	if !targetGroup.IsRemote {
		return fmt.Errorf("host group is not a remote group")
	}

	remoteFetcher := remote.NewRemoteFetcher()
	err = remoteFetcher.UpdateRemoteHostGroup(targetGroup)
	if err != nil {
		return fmt.Errorf("failed to update remote group: %w", err)
	}

	// 更新组的更新时间
	targetGroup.UpdatedAt = time.Now().Format(time.RFC3339)
	manager.UpdatedAt = time.Now().Format(time.RFC3339)

	err = app.configStorage.SaveHostManager(manager)
	if err != nil {
		return fmt.Errorf("failed to save host manager: %w", err)
	}

	return nil
}

// requestAdminPrivileges 尝试以管理员权限重新启动应用
func (app *HostApp) requestAdminPrivileges() error {
	switch runtime.GOOS {
	case "windows":
		// 在Windows上尝试以管理员身份运行
		cmd := exec.Command("powershell", "/C", "Start-Process", "cmd", "/C", "pause")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("could not elevate privileges: %w", err)
		}
	case "darwin", "linux":
		// 在Unix系统上提示使用sudo
		return fmt.Errorf("please run this application with sudo to modify hosts file")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return nil
}

// HostManager 返回HostManager实例
func (app *HostApp) HostManager() *system.HostManager {
	return app.hostManager
}

// BackupAppAndSystemHosts 同时备份应用数据文件和系统hosts文件
func (app *HostApp) BackupAppAndSystemHosts() (string, error) {
	// 创建系统hosts文件备份
	hostsBackupPath := app.hostManager.CreateBackup()
	if hostsBackupPath == "" {
		return "", fmt.Errorf("failed to create system hosts backup")
	}

	// 创建应用数据备份（即data.json）
	err := app.configStorage.BackupData()
	if err != nil {
		return hostsBackupPath, fmt.Errorf("created hosts backup at %s, but failed to create data backup: %w", hostsBackupPath, err)
	}

	return fmt.Sprintf("System hosts backed up to: %s, App data backed up to: data backup created", hostsBackupPath), nil
}