package models

// HostGroup 表示一个Host分组
type HostGroup struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	Content         string `json:"content"`         // Host内容
	Enabled         bool   `json:"enabled"`         // 是否启用
	IsRemote        bool   `json:"isRemote"`        // 是否为远程Host
	URL             string `json:"url,omitempty"`   // 远程URL（仅当IsRemote=true时有效）
	RefreshInterval int64  `json:"refreshInterval"` // 刷新间隔（秒），0表示不启用定时刷新
	LastUpdated     string `json:"lastUpdated"`     // 最后更新时间
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// RemoteConfig 远程Host配置
type RemoteConfig struct {
	URL         string `json:"url"`
	Interval    int64  `json:"interval"`    // 更新间隔 (以纳秒为单位)
	LastFetched string `json:"lastFetched"` // 最后获取时间
	ETag        string `json:"etag"`        // 用于HTTP缓存验证
}

// AppConfig 应用程序配置
type AppConfig struct {
	AutoRefresh     bool     `json:"autoRefresh"`     // 是否自动刷新远程Host
	RefreshInterval int64    `json:"refreshInterval"` // 刷新间隔（纳秒）
	ActiveGroups    []string `json:"activeGroups"`    // 当前激活的分组ID列表
	BackupEnabled   bool     `json:"backupEnabled"`   // 是否启用备份
	MaxBackups      int      `json:"maxBackups"`      // 最大备份数量
	SystemHostPath  string   `json:"systemHostPath"`  // 系统Host文件路径
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

// HostManager 管理所有Host分组
type HostManager struct {
	Config    AppConfig   `json:"config"`
	Groups    []HostGroup `json:"groups"`
	Version   string      `json:"version"` // 配置版本
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
}
