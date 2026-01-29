package permissions

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ElevateIfNeeded 检查当前权限并在需要时提升权限
func ElevateIfNeeded() error {
	if IsAdmin() {
		return nil // 已经具有管理员权限
	}

	return RequestElevation()
}

// IsAdmin 检查当前进程是否以管理员/root权限运行
func IsAdmin() bool {
	switch runtime.GOOS {
	case "windows":
		// 在Windows上，尝试访问受保护资源来检查权限
		_, err := os.OpenFile(`\\.\PHYSICALDRIVE0`, os.O_RDONLY, 0)
		return err == nil
	case "linux", "darwin":
		// 在Unix系统上，检查是否为root用户
		return os.Geteuid() == 0
	default:
		return false
	}
}

// RequestElevation 请求提升到管理员权限
func RequestElevation() error {
	switch runtime.GOOS {
	case "windows":
		return elevateWindows()
	case "linux":
		return elevateLinux()
	case "darwin":
		return elevateMacOS()
	default:
		return fmt.Errorf("platform %s does not support privilege elevation", runtime.GOOS)
	}
}

// CanSudoWithoutPassword 检查是否可以无需密码执行sudo命令（Linux/macOS）
func CanSudoWithoutPassword() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	// 尝试执行一个无害的sudo命令来检查是否需要密码
	cmd := exec.Command("sudo", "-vn")
	// -v 选项会刷新时间戳而不执行命令
	// -n 选项表示不提示输入密码
	err := cmd.Run()
	return err == nil
}

// HasSudoAccess 检查当前用户是否有sudo权限（Linux/macOS）
func HasSudoAccess() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	// 检查是否能运行sudo命令
	cmd := exec.Command("sudo", "-l")
	// 这会列出当前用户可以执行的sudo命令
	err := cmd.Run()
	return err == nil
}

// elevateWindows 使用ShellExecute以管理员身份运行当前程序
func elevateWindows() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// 使用PowerShell调用Start-Process以管理员身份运行
	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command",
		fmt.Sprintf("Start-Process '%s' -Verb RunAs", executable))

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start elevated process: %w", err)
	}

	// 退出当前非特权进程
	os.Exit(0)
	return nil // 不会到达这里
}

// elevateLinux 尝试使用图形化sudo工具提升权限
func elevateLinux() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// 按优先级尝试不同的图形化sudo工具
	tools := []string{"pkexec", "gksudo", "kdesudo", "lxqt-sudo"}

	var lastErr error
	for _, tool := range tools {
		// 检查工具是否存在
		if _, err := exec.LookPath(tool); err == nil {
			var cmd *exec.Cmd
			switch tool {
			case "pkexec":
				cmd = exec.Command(tool, "--disable-internal-agent", executable)
			case "gksudo":
				cmd = exec.Command(tool, "--", executable)
			case "kdesudo":
				cmd = exec.Command(tool, "--", executable)
			case "lxqt-sudo":
				cmd = exec.Command(tool, "--", executable)
			default:
				cmd = exec.Command(tool, executable)
			}

			err := cmd.Start()
			if err == nil {
				// 成功启动提权进程，退出当前进程
				os.Exit(0)
				return nil // 不会到达这里
			}
			lastErr = err
		}
	}

	if lastErr != nil {
		return fmt.Errorf("failed to elevate privileges using any available tool: %w", lastErr)
	}

	return fmt.Errorf("no graphical sudo tools found (tried: %s)", strings.Join(tools, ", "))
}

// elevateMacOS 使用osascript以管理员权限运行当前程序
func elevateMacOS() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	script := fmt.Sprintf(
		"do shell script \"%s\" with administrator privileges",
		executable,
	)

	cmd := exec.Command("osascript", "-e", script)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start elevated process on macOS: %w", err)
	}

	// 退出当前非特权进程
	os.Exit(0)
	return nil // 不会到达这里
}
