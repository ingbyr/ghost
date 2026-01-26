package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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

	manager.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(manager, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cs.dataPath, data, 0644)
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
	// TODO: 实际实现时需要按名称或时间戳排序

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
			os.Remove(filepath.Join(backupPath, backupFiles[i].Name()))
		}
	}

	return nil
}
