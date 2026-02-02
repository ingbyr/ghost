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

// UpdateRemoteHostGroup 更新远程Host组
func (rf *RemoteFetcher) UpdateRemoteHostGroup(group *models.HostGroup) error {
	if !group.IsRemote || group.URL == "" {
		return fmt.Errorf("not a remote host group or URL is empty")
	}

	content, err := rf.FetchRemoteHosts(group.URL)
	if err != nil {
		return fmt.Errorf("failed to fetch remote hosts: %w", err)
	}

	// 不再验证获取到的内容是否符合hosts文件格式，直接更新
	// 更新组内容
	group.Content = content
	group.LastUpdated = time.Now().Format(time.RFC3339)

	// 更新RemoteConfig中的LastFetched字段
	// 如果group是远程组，我们也可以更新它的最后获取时间
	// 注意：这里假设HostGroup结构中没有嵌入RemoteConfig，所以只能更新LastUpdated字段

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
