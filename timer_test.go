package main

import (
	"fmt"
	"testing"
	"time"

	"ghost/application"
	"ghost/models"
)

// TestRemoteHostRefreshTimer 测试定时刷新功能的完整测试
func TestRemoteHostRefreshTimer(t *testing.T) {
	// 创建一个新的HostApp实例
	app, err := application.NewHostApp()
	if err != nil {
		t.Fatalf("Failed to create HostApp: %v", err)
	}

	// 创建一个测试用的远程Host组
	testGroup := models.HostGroup{
		Name:            "Test Remote Group Timer",
		Description:     "A test remote group with timer",
		IsRemote:        true,
		URL:             "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts", // 一个有效的hosts文件URL
		RefreshInterval: 5,                                                                  // 5秒刷新一次，用于测试
		Enabled:         true,
	}

	// 添加这个测试组
	err = app.AddHostGroup(testGroup)
	if err != nil {
		t.Fatalf("Failed to add test group: %v", err)
	}

	fmt.Println("Test group added successfully")

	// 获取刚添加的组ID
	groups, err := app.GetHostGroups()
	if err != nil {
		t.Fatalf("Failed to get host groups: %v", err)
	}

	var testGroupID string
	for _, group := range groups {
		if group.Name == testGroup.Name {
			testGroupID = group.ID
			break
		}
	}

	if testGroupID == "" {
		t.Fatalf("Failed to find test group after adding")
	}

	// 测试单个定时器启动
	err = app.StartRemoteGroupRefreshTimer(testGroupID)
	if err != nil {
		t.Errorf("Failed to start refresh timer for single group: %v", err)
	} else {
		fmt.Println("Started refresh timer for single group successfully")
	}

	// 测试启动所有定时器
	err = app.StartAllRemoteGroupRefreshTimers()
	if err != nil {
		t.Errorf("Error starting all timers: %v", err)
	} else {
		fmt.Println("Started all remote group refresh timers successfully")
	}

	// 等待一段时间观察定时刷新效果
	fmt.Println("Waiting 15 seconds to observe timer behavior...")
	time.Sleep(15 * time.Second)

	// 停止特定组的定时器
	app.StopRemoteGroupRefreshTimer(testGroupID)
	fmt.Printf("Stopped refresh timer for group %s\n", testGroupID)

	// 停止所有定时刷新
	app.StopAllRemoteGroupRefreshTimers()
	fmt.Println("Stopped all remote group refresh timers")

	// 测试删除组是否会自动清理定时器
	err = app.DeleteHostGroup(testGroupID)
	if err != nil {
		t.Errorf("Failed to delete test group: %v", err)
	} else {
		fmt.Printf("Deleted test group %s successfully\n", testGroupID)
	}

	fmt.Println("Test completed successfully")
}

// TestRemoteHostRefreshTimerWithToggle 测试启用/禁用对定时器的影响
func TestRemoteHostRefreshTimerWithToggle(t *testing.T) {
	// 创建一个新的HostApp实例
	app, err := application.NewHostApp()
	if err != nil {
		t.Fatalf("Failed to create HostApp: %v", err)
	}

	// 先清理可能存在的定时器
	app.StopAllRemoteGroupRefreshTimers()

	// 创建一个测试用的远程Host组
	testGroup := models.HostGroup{
		Name:            "Test Toggle Group Timer",
		Description:     "A test remote group for toggle timer test",
		IsRemote:        true,
		URL:             "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews/hosts", // 一个有效的hosts文件URL
		RefreshInterval: 10,                                                                                     // 10秒刷新一次，用于测试
		Enabled:         true,
	}

	// 添加这个测试组
	err = app.AddHostGroup(testGroup)
	if err != nil {
		t.Fatalf("Failed to add test group: %v", err)
	}

	// 获取刚添加的组ID
	groups, err := app.GetHostGroups()
	if err != nil {
		t.Fatalf("Failed to get host groups: %v", err)
	}

	var testGroupID string
	for _, group := range groups {
		if group.Name == testGroup.Name {
			testGroupID = group.ID
			break
		}
	}

	if testGroupID == "" {
		t.Fatalf("Failed to find test group after adding")
	}

	// 启动定时器
	err = app.StartRemoteGroupRefreshTimer(testGroupID)
	if err != nil {
		t.Errorf("Failed to start refresh timer: %v", err)
	} else {
		fmt.Println("Started refresh timer for toggle test")
	}

	// 禁用组，应该停止定时器
	err = app.ToggleHostGroup(testGroupID, false)
	if err != nil {
		t.Errorf("Failed to disable group: %v", err)
	} else {
		fmt.Printf("Disabled group %s, timer should stop\n", testGroupID)
	}

	// 重新启用组，应该启动定时器
	err = app.ToggleHostGroup(testGroupID, true)
	if err != nil {
		t.Errorf("Failed to enable group: %v", err)
	} else {
		fmt.Printf("Re-enabled group %s, timer should start\n", testGroupID)
	}

	// 等待观察行为
	time.Sleep(5 * time.Second)

	// 清理：停止所有定时器
	app.StopAllRemoteGroupRefreshTimers()

	// 删除测试组
	err = app.DeleteHostGroup(testGroupID)
	if err != nil {
		t.Errorf("Failed to delete test group: %v", err)
	}

	fmt.Println("Toggle test completed successfully")
}

// TestRemoteHostRefreshTimerEdgeCases 测试边缘情况
func TestRemoteHostRefreshTimerEdgeCases(t *testing.T) {
	// 创建一个新的HostApp实例
	app, err := application.NewHostApp()
	if err != nil {
		t.Fatalf("Failed to create HostApp: %v", err)
	}

	// 测试不存在的组ID
	err = app.StartRemoteGroupRefreshTimer("non-existent-id")
	if err == nil {
		t.Error("Expected error when starting timer for non-existent group")
	} else {
		fmt.Printf("Correctly got error for non-existent group: %v\n", err)
	}

	// 测试停止不存在的组的定时器不应报错
	app.StopRemoteGroupRefreshTimer("non-existent-id")
	fmt.Println("Successfully handled stopping timer for non-existent group (no error)")

	// 测试带有0刷新间隔的组（不应该启动定时器）
	testGroup := models.HostGroup{
		Name:            "Test Zero Interval Group",
		Description:     "A test remote group with zero refresh interval",
		IsRemote:        true,
		URL:             "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
		RefreshInterval: 0, // 0表示不启用定时刷新
		Enabled:         true,
	}

	err = app.AddHostGroup(testGroup)
	if err != nil {
		t.Fatalf("Failed to add test group: %v", err)
	}

	// 获取刚添加的组ID
	groups, err := app.GetHostGroups()
	if err != nil {
		t.Fatalf("Failed to get host groups: %v", err)
	}

	var testGroupID string
	for _, group := range groups {
		if group.Name == testGroup.Name {
			testGroupID = group.ID
			break
		}
	}

	if testGroupID == "" {
		t.Fatalf("Failed to find test group after adding")
	}

	// 尝试为0间隔的组启动定时器，应该失败
	err = app.StartRemoteGroupRefreshTimer(testGroupID)
	if err != nil {
		fmt.Printf("Correctly got error for zero interval group: %v\n", err)
	} else {
		t.Error("Expected error when starting timer for group with zero refresh interval")
	}

	// 删除测试组
	err = app.DeleteHostGroup(testGroupID)
	if err != nil {
		t.Errorf("Failed to delete test group: %v", err)
	}

	fmt.Println("Edge cases test completed successfully")
}
