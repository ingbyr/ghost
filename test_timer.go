package main

import (
	"fmt"
	"log"
	"time"

	"ghost/application"
	"ghost/models"
)

// 测试定时刷新功能的简单脚本
func main() {
	// 创建一个新的HostApp实例
	app, err := application.NewHostApp()
	if err != nil {
		log.Fatal("Failed to create HostApp:", err)
	}

	// 创建一个测试用的远程Host组
	testGroup := models.HostGroup{
		Name:            "Test Remote Group",
		Description:     "A test remote group with timer",
		IsRemote:        true,
		URL:             "https://httpbin.org/html", // 一个简单的测试URL
		RefreshInterval: 10,                         // 10秒刷新一次，用于测试
		Enabled:         true,
	}

	// 添加这个测试组
	err = app.AddHostGroup(testGroup)
	if err != nil {
		log.Fatal("Failed to add test group:", err)
	}

	fmt.Println("Test group added successfully")

	// 启动所有远程组的定时刷新
	err = app.StartAllRemoteGroupRefreshTimers()
	if err != nil {
		log.Printf("Error starting timers: %v", err)
	} else {
		fmt.Println("Started all remote group refresh timers")
	}

	// 等待一段时间观察定时刷新效果
	fmt.Println("Waiting 30 seconds to observe timer behavior...")
	time.Sleep(30 * time.Second)

	// 停止所有定时刷新
	app.StopAllRemoteGroupRefreshTimers()
	fmt.Println("Stopped all remote group refresh timers")

	fmt.Println("Test completed successfully")
}
