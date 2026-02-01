package main

import (
	"context"
	"fmt"

	"ghost/application"
	"ghost/models"
	"ghost/remote"
)

// App struct
type App struct {
	ctx     context.Context
	hostApp *application.HostApp
}

// NewApp creates a new App application struct
func NewApp() *App {
	hostApp, err := application.NewHostApp()
	if err != nil {
		fmt.Printf("Failed to initialize host app: %v\n", err)
		return nil
	}

	return &App{
		hostApp: hostApp,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 检查备份目录是否为空，如果为空则备份当前系统hosts文件
	isEmpty, err := a.IsBackupDirEmpty()
	if err != nil {
		fmt.Printf("Warning: failed to check if backup directory is empty: %v\n", err)
		return
	}

	if isEmpty {
		fmt.Println("Backup directory is empty, creating initial system hosts backup")
		err = a.BackupRawSystemHosts()
		if err != nil {
			fmt.Printf("Warning: failed to create initial system hosts backup: %v\n", err)
		} else {
			fmt.Println("Initial system hosts backup created successfully")
		}
	}
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// 创建系统hosts文件的备份
	err := a.hostApp.BackupConfig()
	if err != nil {
		fmt.Printf("Warning: failed to backup config: %v\n", err)
	}

	// 系统hosts文件备份逻辑已移除
}

// GetHostGroups 获取所有Host分组
func (a *App) GetHostGroups() ([]models.HostGroup, error) {
	return a.hostApp.GetHostGroups()
}

// AddHostGroup 添加新的Host分组
func (a *App) AddHostGroup(group models.HostGroup) error {
	return a.hostApp.AddHostGroup(group)
}

// UpdateHostGroup 更新Host分组
func (a *App) UpdateHostGroup(group models.HostGroup) error {
	return a.hostApp.UpdateHostGroup(group)
}

// DeleteHostGroup 删除Host分组
func (a *App) DeleteHostGroup(id string) error {
	return a.hostApp.DeleteHostGroup(id)
}

// ToggleHostGroup 启用或禁用Host分组
func (a *App) ToggleHostGroup(id string, enabled bool) error {
	return a.hostApp.ToggleHostGroup(id, enabled)
}

// ApplyHosts 应用所有启用的Host分组到系统
func (a *App) ApplyHosts() error {
	return a.hostApp.ApplyHosts()
}

// GetSystemHostPath 获取系统hosts文件路径
func (a *App) GetSystemHostPath() string {
	return a.hostApp.GetSystemHostPath()
}

// GetSystemHostsContent 获取系统hosts文件内容
func (a *App) GetSystemHostsContent() (string, error) {
	return a.hostApp.GetSystemHostsContent()
}

// RefreshRemoteGroups 刷新所有远程Host组
func (a *App) RefreshRemoteGroups() error {
	return a.hostApp.RefreshRemoteGroups()
}

// GetConfig 获取应用程序配置
func (a *App) GetConfig() (*models.AppConfig, error) {
	return a.hostApp.GetConfig()
}

// UpdateConfig 更新应用程序配置
func (a *App) UpdateConfig(config models.AppConfig) error {
	return a.hostApp.UpdateConfig(&config)
}

// GetHostGroup 获取指定ID的Host分组
func (a *App) GetHostGroup(id string) (*models.HostGroup, error) {
	return a.hostApp.GetHostGroup(id)
}

// GetRemoteContent 获取指定URL的远程hosts内容
func (a *App) GetRemoteContent(url string) (string, error) {
	remoteFetcher := remote.NewRemoteFetcher()
	content, err := remoteFetcher.FetchRemoteHosts(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch remote content: %w", err)
	}
	return content, nil
}

// RefreshRemoteGroup 刷新指定的远程Host分组
func (a *App) RefreshRemoteGroup(id string) error {
	return a.hostApp.RefreshRemoteGroup(id)
}

// BackupConfig 创建配置备份
func (a *App) BackupConfig() error {
	return a.hostApp.BackupConfig()
}

// BackupData 创建数据文件备份
func (a *App) BackupData() error {
	return a.hostApp.BackupData()
}

// CreateSystemHostsBackup 创建系统hosts文件备份
func (a *App) CreateSystemHostsBackup() (string, error) {
	// Function removed as per requirement - hosts.ghost_backup logic deleted
	return "", nil
}

// BackupAppAndSystemHosts 同时备份应用数据文件和系统hosts文件
func (a *App) BackupAppAndSystemHosts() (string, error) {
	result, err := a.hostApp.BackupAppAndSystemHosts()
	if err != nil {
		return "", err
	}
	return result, nil
}

// ListDataBackups 列出所有数据备份文件
func (a *App) ListDataBackups() ([]string, error) {
	return a.hostApp.ListDataBackups()
}

// RestoreData 从备份文件恢复数据
func (a *App) RestoreData(backupFileName string) error {
	return a.hostApp.RestoreData(backupFileName)
}

// HasRawHostsBackup 检查是否存在原始hosts备份文件
func (a *App) HasRawHostsBackup() (bool, error) {
	return a.hostApp.HasRawHostsBackup()
}

// IsBackupDirEmpty 检查备份目录是否为空
func (a *App) IsBackupDirEmpty() (bool, error) {
	return a.hostApp.IsBackupDirEmpty()
}

// BackupRawSystemHosts 备份当前系统hosts文件
func (a *App) BackupRawSystemHosts() error {
	return a.hostApp.BackupRawSystemHosts()
}

// RestoreRawSystemHosts 从备份恢复系统hosts文件
func (a *App) RestoreRawSystemHosts(backupFileName string) error {
	return a.hostApp.RestoreRawSystemHosts(backupFileName)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
