package remote

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"ghost/models"
)

// RemoteFetcher 用于从远程URL获取Host内容
type RemoteFetcher struct {
	httpClient *http.Client
	userAgent  string
}

// NewRemoteFetcher 创建新的远程获取器
func NewRemoteFetcher() *RemoteFetcher {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &RemoteFetcher{
		httpClient: client,
		userAgent:  "Ghost Host Manager/1.0",
	}
}

// FetchRemoteHosts 从指定URL获取远程Hosts内容
func (rf *RemoteFetcher) FetchRemoteHosts(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", rf.userAgent)

	resp, err := rf.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch from URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code %d from URL %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// ValidateHostsContent 验证Hosts内容格式
func (rf *RemoteFetcher) ValidateHostsContent(content string) error {
	// 检查内容是否符合基本的hosts文件格式
	lines := splitLines(content)

	for i, line := range lines {
		line = trimSpace(line)

		// 跳过注释和空行
		if line == "" || line[0] == '#' {
			continue
		}

		// 检查是否至少包含IP和域名
		fields := splitFields(line)
		if len(fields) < 2 {
			return fmt.Errorf("invalid hosts entry at line %d: %s", i+1, line)
		}

		ip := fields[0]
		domain := fields[1]

		// 简单验证IP格式
		if !isValidIP(ip) {
			return fmt.Errorf("invalid IP format at line %d: %s", i+1, ip)
		}

		// 简单验证域名格式
		if !isValidDomain(domain) {
			return fmt.Errorf("invalid domain format at line %d: %s", i+1, domain)
		}
	}

	return nil
}

// UpdateRemoteHostGroup 更新远程Host组
func (rf *RemoteFetcher) UpdateRemoteHostGroup(group *models.HostGroup) error {
	if !group.IsRemote || group.URL == "" {
		return fmt.Errorf("not a remote host group or URL is empty")
	}

	content, err := rf.FetchRemoteHosts(group.URL)
	if err != nil {
		return fmt.Errorf("failed to fetch remote hosts: %w", err)
	}

	// 验证内容
	err = rf.ValidateHostsContent(content)
	if err != nil {
		return fmt.Errorf("invalid hosts content: %w", err)
	}

	// 更新组内容
	group.Content = content
	group.LastUpdated = time.Now()

	return nil
}

// DownloadToFile 下载远程内容到临时文件
func (rf *RemoteFetcher) DownloadToFile(url, filePath string) error {
	content, err := rf.FetchRemoteHosts(url)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}

	return nil
}

// splitLines 分割字符串为行
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, r := range s {
		if r == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// trimSpace 去除字符串首尾空白字符
func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r' || s[start] == '\n') {
		start++
	}
	for start < end && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}

// splitFields 按空白字符分割字段
func splitFields(s string) []string {
	var fields []string
	start := -1
	for i, r := range s {
		if r == ' ' || r == '\t' {
			if start >= 0 {
				fields = append(fields, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		fields = append(fields, s[start:])
	}
	return fields
}

// isValidIP 简单验证IP地址格式
func isValidIP(ip string) bool {
	if len(ip) < 7 || len(ip) > 15 { // 最短: "0.0.0.0", 最长: "255.255.255.255"
		return false
	}

	parts := splitByChar(ip, '.')
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if !isValidIPPart(part) {
			return false
		}
	}

	return true
}

// splitByChar 按指定字符分割字符串
func splitByChar(s string, sep byte) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// isValidIPPart 验证IP地址的一个部分
func isValidIPPart(part string) bool {
	if len(part) == 0 || len(part) > 3 {
		return false
	}

	for _, r := range part {
		if r < '0' || r > '9' {
			return false
		}
	}

	// 检查前导零（除了单独的"0"）
	if len(part) > 1 && part[0] == '0' {
		return false
	}

	// 转换为数字并检查范围
	num := 0
	for _, r := range part {
		num = num*10 + int(r-'0')
	}

	return num >= 0 && num <= 255
}

// isValidDomain 简单验证域名格式
func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// 检查域名是否只包含有效的字符
	for _, r := range domain {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '.' || r == '_' {
			continue
		}
		return false
	}

	// 检查域名格式（简单验证）
	parts := splitByChar(domain, '.')
	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false
		}
		// 检查每个标签不能以连字符开头或结尾
		if len(part) > 1 && (part[0] == '-' || part[len(part)-1] == '-') {
			return false
		}
		if len(part) == 1 && part[0] == '-' {
			return false
		}
	}

	return true
}
