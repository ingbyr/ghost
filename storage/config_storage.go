package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"ghost/models"
)

const (
	AppDataDir = ".ghost"
	ConfigFile = "config.json"
	DataFile   = "data.json"
	BackupDir  = "backups"
)

// ConfigStorage 处理配置文件的读写
type ConfigStorage struct {
	configPath string
	dataPath   string
	mutex      sync.RWMutex
}

// NewConfigStorage 创建新的配置存储实例
func NewConfigStorage() (*ConfigStorage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appDataPath := filepath.Join(homeDir, AppDataDir)
	err = os.MkdirAll(appDataPath, 0755)
	if err != nil {
		return nil, err
	}

	backupPath := filepath.Join(appDataPath, BackupDir)
	err = os.MkdirAll(backupPath, 0755)
	if err != nil {
		return nil, err
	}

	return &ConfigStorage{
		configPath: filepath.Join(appDataPath, ConfigFile),
		dataPath:   filepath.Join(appDataPath, DataFile),
	}, nil
}

// LoadConfig 加载应用程序配置
func (cs *ConfigStorage) LoadConfig() (*models.AppConfig, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	config := &models.AppConfig{}

	// 设置默认值
	config.AutoRefresh = false
	config.RefreshInterval = 3600 // 1 hour
	config.ActiveGroups = []string{}
	config.BackupEnabled = true
	config.MaxBackups = 10

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(cs.configPath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(cs.configPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig 保存应用程序配置
func (cs *ConfigStorage) SaveConfig(config *models.AppConfig) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cs.configPath, data, 0644)
}

// LoadHostManager 加载Host管理器数据
func (cs *ConfigStorage) LoadHostManager() (*models.HostManager, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	// 如果数据文件不存在，返回空的HostManager
	if _, err := os.Stat(cs.dataPath); os.IsNotExist(err) {
		return &models.HostManager{
			Config:    models.AppConfig{},
			Groups:    []models.HostGroup{},
			Version:   "1.0.0",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}, nil
	}

	data, err := os.ReadFile(cs.dataPath)
	if err != nil {
		return nil, err
	}

	var manager models.HostManager
	err = json.Unmarshal(data, &manager)
	if err != nil {
		return nil, err
	}

	return &manager, nil
}

// SaveHostManager 保存Host管理器数据
func (cs *ConfigStorage) SaveHostManager(manager *models.HostManager) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	manager.UpdatedAt = time.Now().Format(time.RFC3339)

	data, err := json.MarshalIndent(manager, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cs.dataPath, data, 0644)
}

// BackupData 创建数据文件备份
func (cs *ConfigStorage) BackupData() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	backupPath := filepath.Join(homeDir, AppDataDir, BackupDir)
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFile := filepath.Join(backupPath, timestamp+".json")

	// 读取当前数据
	data, err := os.ReadFile(cs.dataPath)
	if err != nil {
		return err
	}

	// 写入备份
	err = os.WriteFile(backupFile, data, 0644)
	if err != nil {
		return err
	}

	// 清理旧备份
	return cs.cleanupOldBackups()
}

// BackupConfig 创建配置备份
func (cs *ConfigStorage) BackupConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	backupPath := filepath.Join(homeDir, AppDataDir, BackupDir)
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupPath, "config_"+timestamp+".json")

	// 读取当前配置
	data, err := os.ReadFile(cs.configPath)
	if err != nil {
		return err
	}

	// 写入备份
	err = os.WriteFile(backupFile, data, 0644)
	if err != nil {
		return err
	}

	// 清理旧备份
	return cs.cleanupOldBackups()
}

// ListDataBackups 列出所有数据备份文件
func (cs *ConfigStorage) ListDataBackups() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	backupPath := filepath.Join(homeDir, AppDataDir, BackupDir)
	files, err := os.ReadDir(backupPath)
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" && !strings.HasPrefix(file.Name(), "pre_restore_") {
			backups = append(backups, file.Name())
		}
	}

	return backups, nil
}

// RestoreData 从备份文件恢复数据
func (cs *ConfigStorage) RestoreData(backupFileName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	backupPath := filepath.Join(homeDir, AppDataDir, BackupDir, backupFileName)

	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFileName)
	}

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	// 先创建当前数据的备份（安全备份，计入总计限制）
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	safetyBackup := filepath.Join(homeDir, AppDataDir, BackupDir, "pre_restore_"+timestamp+".json")

	currentData, err := os.ReadFile(cs.dataPath)
	if err == nil {
		os.WriteFile(safetyBackup, currentData, 0644)
	}

	// 读取备份数据
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// 验证 JSON 格式
	var manager models.HostManager
	err = json.Unmarshal(backupData, &manager)
	if err != nil {
		return fmt.Errorf("invalid backup file format: %w", err)
	}

	// 写入到 data.json
	err = os.WriteFile(cs.dataPath, backupData, 0644)
	if err != nil {
		return fmt.Errorf("failed to restore data: %w", err)
	}

	// 注意：不在这里调用 cleanupOldBackups，避免与 LoadConfig 的死锁
	// 安全备份由下次自动备份时清理
	return nil
}

// HasRawHostsBackup 检查是否存在原始hosts备份文件
func (cs *ConfigStorage) HasRawHostsBackup() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	backupPath := filepath.Join(homeDir, AppDataDir, BackupDir, "raw_hosts_backup.txt")
	_, err = os.Stat(backupPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// IsBackupDirEmpty 检查备份目录是否为空
func (cs *ConfigStorage) IsBackupDirEmpty() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	backupPath := filepath.Join(homeDir, AppDataDir, BackupDir)
	files, err := os.ReadDir(backupPath)
	if err != nil {
		return false, err
	}

	// 检查是否有非临时文件（忽略pre_restore_开头的临时备份）
	for _, file := range files {
		if !file.IsDir() && !strings.HasPrefix(file.Name(), "pre_restore_") {
			return false, nil
		}
	}

	return true, nil
}

// BackupRawSystemHosts 创建原始系统hosts文件的备份
func (cs *ConfigStorage) BackupRawSystemHosts(systemHostPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	backupPath := filepath.Join(homeDir, AppDataDir, BackupDir)
	backupFile := filepath.Join(backupPath, "raw_hosts_backup.txt")

	// 读取系统hosts文件内容
	content, err := os.ReadFile(systemHostPath)
	if err != nil {
		return fmt.Errorf("failed to read system hosts file: %w", err)
	}

	// 写入备份
	err = os.WriteFile(backupFile, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

// cleanupOldBackups 清理旧备份文件
func (cs *ConfigStorage) cleanupOldBackups() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	backupPath := filepath.Join(homeDir, AppDataDir, BackupDir)
	files, err := os.ReadDir(backupPath)
	if err != nil {
		return err
	}

	// 获取所有备份文件并排序
	var backupFiles []os.FileInfo
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			info, err := file.Info()
			if err != nil {
				continue
			}
			backupFiles = append(backupFiles, info)
		}
	}

	// 按修改时间排序（最新的在前）
	sort.Slice(backupFiles, func(i, j int) bool {
		return backupFiles[i].ModTime().After(backupFiles[j].ModTime())
	})

	// 删除超出最大数量的备份
	config, err := cs.LoadConfig()
	if err != nil {
		return err
	}

	maxBackups := config.MaxBackups
	if maxBackups <= 0 {
		maxBackups = 10
	}

	if len(backupFiles) > maxBackups {
		for i := maxBackups; i < len(backupFiles); i++ {
			err := os.Remove(filepath.Join(backupPath, backupFiles[i].Name()))
			if err != nil {
				// 记录错误但不停止整个过程
				continue
			}
		}
	}

	return nil
}
