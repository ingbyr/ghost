# Ghost - Host管理器

基于Go语言和Wails框架的跨平台Host文件切换GUI程序，支持多Host分组管理和远程Host功能。

## 新增功能：树形结构UI

最新版本引入了全新的左右分栏界面设计：

- **左侧树形结构**：清晰展示所有Host分组
- **右侧编辑面板**：提供详细的内容编辑功能
- **搜索功能**：快速查找特定Host分组
- **状态管理**：实时显示各分组启用/禁用状态
- **防丢失保护**：检测未保存更改，防止意外丢失

## 核心功能

1. **多Host分组管理** - 支持创建、编辑、删除和管理多个Host文件分组
2. **远程Host支持** - 可配置URL定期拉取远程Host文件，启用后应用到本地Host
3. **开关控制** - 可选择性地启用或禁用各个Host分组
4. **一键应用** - 将所有启用的Host分组合并并应用到系统Hosts文件
5. **跨平台支持** - 自动识别Windows、macOS和Linux系统的Hosts文件位置及权限要求

## 技术架构

- **前端** - Vue.js + Wails框架构建桌面界面
- **后端** - Go语言实现业务逻辑
- **数据存储** - 基于JSON文件的持久化存储方案
- **远程获取** - 支持从URL获取远程Host内容并验证格式

## 跨平台兼容性

- **Windows**: `C:\Windows\System32\drivers\etc\hosts`
- **Unix/Linux/macOS**: `/etc/hosts`
- **权限管理**: 需要管理员权限来修改系统Hosts文件

## 开发说明

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## 构建

To build a redistributable, production mode package, use `wails build`.