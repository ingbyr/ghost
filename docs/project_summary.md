# Ghost - Host管理器项目总结

## 项目概述
Ghost是一个基于Go语言和Wails框架的跨平台Host文件切换GUI程序，支持多Host分组管理和远程Host功能。

## 核心功能
1. **多Host分组管理** - 支持创建、编辑、删除和管理多个Host文件分组
2. **远程Host支持** - 可配置URL定期拉取远程Host文件，启用后应用到本地Host
3. **开关控制** - 可选择性地启用或禁用各个Host分组
4. **一键应用** - 将所有启用的Host分组合并并应用到系统Hosts文件
5. **跨平台支持** - 自动识别Windows、macOS和Linux系统的Hosts文件位置及权限要求
6. **国际化支持** - 支持中英文界面切换，自动检测语言偏好

## 技术架构
- **前端** - Vue.js + Wails框架构建桌面界面
- **后端** - Go语言实现业务逻辑
- **数据存储** - 基于JSON文件的持久化存储方案
- **远程获取** - 支持从URL获取远程Host内容并验证格式
- **国际化** - 使用vue-i18n实现多语言支持

## 关键文件结构
- `models/host_models.go` - 数据模型定义
- `storage/config_storage.go` - 配置存储管理
- `system/host_manager.go` - 系统Hosts文件操作
- `remote/remote_fetcher.go` - 远程Host获取功能
- `application/host_app.go` - 应用主逻辑
- `frontend/src/App.vue` - Wails前端界面
- `frontend/src/i18n.js` - 国际化配置
- `frontend/src/locales/` - 语言包文件

## 跨平台兼容性
- **Windows**: `C:\Windows\System32\drivers\etc\hosts`
- **Unix/Linux/macOS**: `/etc/hosts`
- **权限管理**: 需要管理员权限来修改系统Hosts文件
  - **Windows**: 通过应用程序清单文件配置，实现永久性管理员权限，用户确认后整个会话期间都具有所需权限
  - **Linux/macOS**: 通过图形化sudo工具请求临时权限，根据系统配置可能需要周期性重新认证
- **国际化**: 支持中英文界面切换，自动检测浏览器语言偏好

## UI界面功能
- 左侧树形结构显示Host分组列表
- 系统Host文件条目置顶显示，点击后在右侧显示其内容的只读预览界面
- Host分组列表显示（支持本地与远程Host分类管理）
- 新增/编辑/删除Host分组
- 启用/禁用分组的Switch类型UI控件
- 自动UUID生成，取消手动填充分组ID
- 刷新远程Host按钮
- 备份数据功能
- 分组内容预览
- 远程Host内容预览功能，可实时查看远程内容
- 应用启动时默认显示系统Host文件
- 启用远程Host时不再自动更新远程内容，仅在点击"获取host内容"或"Refresh Remote Groups"按钮时更新内容
- 启用的Host组在编辑器中显示为只读预览模式，以防止对活跃配置的意外修改
- LocalHost编辑器实现自动保存功能，当内容编辑区域失去焦点时自动保存到data.json
- **国际化界面**: 支持中英文界面切换，用户可选择偏好的语言
- **定时刷新功能**: 远程Host支持设置1小时、4小时、8小时、24小时固定间隔的自动刷新
- **自动保存功能**: 编辑Host分组时，当输入框失去焦点时自动保存更改，无需手动点击保存按钮

## 错误处理
- 完善的错误处理和用户提示
- 自动备份系统Hosts文件
- 远程Host内容格式验证

## 开发环境与调试
- Go语言 Wails框架
- JSON文件存储
- 前端技术使用Vue和ElementUI前端组件
- 使用 Windows PowerShell 命令格式调试
- 国际化使用 vue-i18n 库实现