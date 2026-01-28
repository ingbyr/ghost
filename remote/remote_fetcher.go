package remote

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

// UpdateRemoteHostGroup 更新远程Host组
func (rf *RemoteFetcher) UpdateRemoteHostGroup(group *models.HostGroup) error {
	if !group.IsRemote || group.URL == "" {
		return fmt.Errorf("not a remote host group or URL is empty")
	}

	content, err := rf.FetchRemoteHosts(group.URL)
	if err != nil {
		return fmt.Errorf("failed to fetch remote hosts: %w", err)
	}

	// 验证获取到的内容是否符合hosts文件格式
	if !rf.isValidHostsContent(content) {
		return fmt.Errorf("fetched content does not appear to be a valid hosts file format")
	}

	// 更新组内容
	group.Content = content
	group.LastUpdated = time.Now().Format(time.RFC3339)

	// 更新RemoteConfig中的LastFetched字段
	// 如果group是远程组，我们也可以更新它的最后获取时间
	// 注意：这里假设HostGroup结构中没有嵌入RemoteConfig，所以只能更新LastUpdated字段

	return nil
}

// isValidHostsContent 验证内容是否符合基本的hosts文件格式
func (rf *RemoteFetcher) isValidHostsContent(content string) bool {
	// 简单的验证：检查是否包含IP地址和域名的基本格式
	lines := strings.Split(content, "\n")
	validLines := 0

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 跳过注释和空行
		if strings.HasPrefix(trimmedLine, "#") || trimmedLine == "" {
			continue
		}

		// 检查是否包含至少一个IP地址和域名的组合
		parts := strings.Fields(trimmedLine)
		if len(parts) >= 2 {
			ipPart := parts[0]
			domainPart := parts[1]

			// 简单检查IP地址格式（支持IPv4）
			if isValidIP(ipPart) && isValidDomain(domainPart) {
				validLines++
			}
		}
	}

	// 如果至少有一行符合格式，则认为内容有效
	return validLines > 0
}

// isValidIP 检查字符串是否为有效的IP地址格式
func isValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
	}

	return true
}

// isValidDomain 检查字符串是否为有效的域名格式
func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// 检查是否只包含允许的字符
	for _, char := range domain {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '.' && char != '_' {
			return false
		}
	}

	// 检查域名部分
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false
		}
		if part[0] == '-' || part[len(part)-1] == '-' {
			return false
		}
	}

	return true
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