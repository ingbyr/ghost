package system

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"ghost/storage"
)

const (
	// GhostSectionStart Ghost标记段开始
	GhostSectionStart = "# >>> Ghost Host Entries"
	// GhostSectionEnd Ghost标记段结束
	GhostSectionEnd = "# <<< Ghost Host Entries"
	// GhostBackupExt 备份文件扩展名
	GhostBackupExt = ".ghost_backup"
)

// HostManager 系统hosts文件管理器
type HostManager struct {
	SystemHostPath string
}

// NewHostManager 创建新的系统hosts管理器
func NewHostManager() *HostManager {
	path := GetSystemHostsPath()
	return &HostManager{
		SystemHostPath: path,
	}
}

// GetSystemHostsPath 根据操作系统获取系统hosts文件路径
func GetSystemHostsPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts")
	default: // linux, darwin (macOS), etc
		return "/etc/hosts"
	}
}

// ReadSystemHosts 读取系统hosts文件内容
func (hm *HostManager) ReadSystemHosts() (string, error) {
	content, err := os.ReadFile(hm.SystemHostPath)
	if err != nil {
		return "", fmt.Errorf("failed to read system hosts file: %w", err)
	}

	return string(content), nil
}

// WriteSystemHosts 写入系统hosts文件内容
func (hm *HostManager) WriteSystemHosts(content string) error {
	// 创建备份
	backupPath := hm.createBackup()
	if backupPath != "" {
		fmt.Printf("Created backup: %s\n", backupPath)
	}

	// 写入新内容
	err := os.WriteFile(hm.SystemHostPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write system hosts file: %w", err)
	}

	return nil
}

// getAppDataDir 获取应用程序数据目录
func (hm *HostManager) getAppDataDir() (string, error) {
	usr, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	appDataPath := filepath.Join(usr, storage.AppDataDir)
	err = os.MkdirAll(appDataPath, 0755)
	if err != nil {
		return "", err
	}

	backupPath := filepath.Join(appDataPath, storage.BackupDir)
	err = os.MkdirAll(backupPath, 0755)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

// createBackup 创建系统hosts文件备份
func (hm *HostManager) createBackup() string {
	// 获取应用数据目录用于备份
	backupDir, err := hm.getAppDataDir()
	if err != nil {
		fmt.Printf("Warning: failed to get app data directory for backup: %v\n", err)
		return ""
	}

	// 构建备份文件路径
	backupFileName := filepath.Base(hm.SystemHostPath) + GhostBackupExt + "_" + time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, backupFileName)

	// 读取原文件内容
	content, err := os.ReadFile(hm.SystemHostPath)
	if err != nil {
		fmt.Printf("Warning: failed to read original hosts file for backup: %v\n", err)
		return ""
	}

	// 写入备份文件
	err = os.WriteFile(backupPath, content, 0644)
	if err != nil {
		fmt.Printf("Warning: failed to create backup: %v\n", err)
		return ""
	}

	return backupPath
}

// ApplyHostGroups 将指定的HostGroups应用到系统hosts文件
func (hm *HostManager) ApplyHostGroups(hostGroups []map[string]interface{}) error {
	// 读取当前系统hosts文件内容
	currentContent, err := hm.ReadSystemHosts()
	if err != nil {
		return fmt.Errorf("failed to read current system hosts: %w", err)
	}

	// 移除之前的Ghost段
	contentWithoutGhost, err := hm.removeGhostEntries(currentContent)
	if err != nil {
		return fmt.Errorf("failed to remove previous ghost entries: %w", err)
	}

	// 准备新的Ghost段内容
	var ghostContent strings.Builder
	ghostContent.WriteString(fmt.Sprintf("\n%s\n", GhostSectionStart))
	ghostContent.WriteString("# This section is managed by Ghost - Host Manager\n")
	ghostContent.WriteString("# Changes made outside this section will be preserved\n")
	ghostContent.WriteString("# Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")

	// 添加启用的host组内容
	for _, group := range hostGroups {
		enabled, ok := group["enabled"].(bool)
		if !ok || !enabled {
			continue
		}

		name, _ := group["name"].(string)
		if name == "" {
			name, _ = group["id"].(string)
		}

		content, ok := group["content"].(string)
		if !ok || content == "" {
			continue
		}

		ghostContent.WriteString(fmt.Sprintf("# Start of group: %s\n", name))
		ghostContent.WriteString(content)
		ghostContent.WriteString(fmt.Sprintf("\n# End of group: %s\n\n", name))
	}

	ghostContent.WriteString(GhostSectionEnd)
	ghostContent.WriteString("\n") // 确保文件末尾有换行

	// 确保原内容以换行结尾，避免与Ghost段内容连在一起
	if contentWithoutGhost != "" && !strings.HasSuffix(contentWithoutGhost, "\n") {
		contentWithoutGhost += "\n"
	}
	// 组合最终内容，确保不产生多余空行
	finalContent := strings.TrimRight(contentWithoutGhost, "\n") + "\n" + ghostContent.String()

	// 写入系统hosts文件
	err = hm.WriteSystemHosts(finalContent)
	if err != nil {
		return fmt.Errorf("failed to write updated hosts file: %w", err)
	}

	return nil
}

// removeGhostEntries 从内容中移除现有的Ghost段
func (hm *HostManager) removeGhostEntries(content string) (string, error) {
	lines := strings.Split(content, "\n")
	var result []string
	inGhostSection := false

	for _, line := range lines {
		if strings.Contains(line, GhostSectionStart) {
			inGhostSection = true
			// 保留注释之前的行
			continue
		} else if strings.Contains(line, GhostSectionEnd) {
			inGhostSection = false
			continue
		}

		if !inGhostSection {
			result = append(result, line)
		}
	}

	// 确保最后一行有换行符
	if len(result) > 0 && !strings.HasSuffix(result[len(result)-1], "\n") {
		lastLine := result[len(result)-1]
		if lastLine != "" {
			result[len(result)-1] = lastLine
		}
	}

	return strings.Join(result, "\n"), nil
}

// HasWritePermission 检查是否有写入系统hosts文件的权限
func (hm *HostManager) HasWritePermission() bool {
	// 尝试打开文件进行追加写入，以检查权限
	file, err := os.OpenFile(hm.SystemHostPath, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// RequestAdminPrivileges 提示用户以管理员权限运行（在某些系统上）
func (hm *HostManager) RequestAdminPrivileges() error {
	switch runtime.GOOS {
	case "windows":
		return fmt.Errorf("please run this application as administrator to modify hosts file")
	case "darwin", "linux":
		return fmt.Errorf("please run this application with sudo to modify hosts file")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// GetHostsFileInfo 获取hosts文件信息
func (hm *HostManager) GetHostsFileInfo() (os.FileInfo, error) {
	return os.Stat(hm.SystemHostPath)
}

// RestoreFromBackup 从备份恢复hosts文件
func (hm *HostManager) RestoreFromBackup(backupPath string) error {
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	err = os.WriteFile(hm.SystemHostPath, backupContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}

	return nil
}

// ListBackups 列出所有可用的备份文件
func (hm *HostManager) ListBackups() ([]string, error) {
	// 获取应用数据目录中的备份
	backupDir, err := hm.getAppDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get app data directory: %w", err)
	}

	// 构建匹配模式
	pattern := filepath.Join(backupDir, filepath.Base(hm.SystemHostPath)+GhostBackupExt+"*")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
