# Ghost Host Manager - 前端组件重构文档

## 项目概述

Ghost是一个基于Go语言和Wails框架的跨平台Host文件切换GUI程序，支持多Host分组管理和远程Host功能。本文档记录了前端组件的重构情况，将原来的单体App.vue文件拆分为多个独立组件，以提高代码可维护性和可读性。

## 重构目标

1. **提高代码可维护性**：将大文件拆分为小的、专注的组件
2. **增强代码复用性**：创建独立的可复用组件
3. **改善开发效率**：团队成员可以并行开发不同组件
4. **遵循Vue最佳实践**：使用props向下传递数据，emits向上触发事件

## 重构后的组件结构

### 1. App.vue (主应用组件)

**功能**：
- 状态管理（groups, selectedGroup, editingGroup等）
- API调用（与Go后端交互）
- 组件协调（连接各个子组件）

**主要职责**：
- 管理全局状态
- 处理业务逻辑
- 协调各组件间的数据流

### 2. Sidebar.vue (侧边栏组件)

**功能**：
- 显示Host群组列表
- 提供搜索功能
- 显示系统Host文件条目
- 控制Host组的启用/禁用和删除

**Props**：
- `groups`: Host群组数组
- `selectedGroup`: 当前选中的群组
- `systemHostPath`: 系统Host文件路径
- `searchQuery`: 搜索查询字符串

**Emits**：
- `select-group`: 选择Host群组
- `select-system-host`: 选择系统Host
- `toggle-status`: 切换群组状态
- `delete-group`: 删除群组
- `open-add-modal`: 打开添加群组模态框
- `update:search-query`: 更新搜索查询

### 3. MainPanel.vue (主面板组件)

**功能**：
- 根据选中的群组类型渲染相应的编辑器
- 包含条件渲染逻辑，决定显示哪种编辑器

**Props**：
- `selectedGroup`: 当前选中的群组
- `editingGroup`: 正在编辑的群组
- `systemHostPath`: 系统Host路径
- `systemHostContent`: 系统Host内容
- `remoteContentPreview`: 远程内容预览
- `isFetchingRemote`: 是否正在获取远程内容

**Emits**：
- `save-group`: 保存群组
- `cancel-edit`: 取消编辑
- `fetch-remote-content`: 获取远程内容
- `mark-as-dirty`: 标记为脏数据
- `refresh-system-host`: 刷新系统Host

### 4. RemoteHostEditor.vue (远程Host编辑器)

**功能**：
- 编辑远程Host群组
- 显示URL字段
- 提供获取远程内容的功能
- 显示远程内容预览

**Props**：
- `group`: 群组对象
- `editingGroup`: 正在编辑的群组
- `remoteContentPreview`: 远程内容预览
- `isFetchingRemote`: 是否正在获取远程内容

**Emits**：
- `save-group`: 保存群组
- `cancel-edit`: 取消编辑
- `fetch-remote-content`: 获取远程内容
- `mark-as-dirty`: 标记为脏数据

### 5. LocalHostEditor.vue (本地Host编辑器)

**功能**：
- 编辑本地Host群组
- 显示内容编辑区域
- 处理本地Host内容的编辑

**Props**：
- `group`: 群组对象
- `editingGroup`: 正在编辑的群组

**Emits**：
- `save-group`: 保存群组
- `cancel-edit`: 取消编辑
- `mark-as-dirty`: 标记为脏数据

### 6. SystemHostPreview.vue (系统Host预览)

**功能**：
- 以只读方式预览系统Host文件
- 提供刷新功能
- 显示系统Host路径和内容

**Props**：
- `systemHostPath`: 系统Host路径
- `systemHostContent`: 系统Host内容

**Emits**：
- `refresh-system-host`: 刷新系统Host

### 7. ActionBar.vue (操作栏组件)

**功能**：
- 提供主要的操作按钮
- 应用Hosts到系统
- 刷新远程群组
- 刷新列表

**Emits**：
- `apply-hosts`: 应用Hosts
- `refresh-remote`: 刷新远程群组
- `refresh-list`: 刷新列表

### 8. AddGroupModal.vue (添加群组模态框)

**功能**：
- 提供添加新Host群组的界面
- 支持本地和远程Host群组类型
- 表单验证和提交

**Props**：
- `showModal`: 是否显示模态框
- `newGroup`: 新群组对象

**Emits**：
- `close-modal`: 关闭模态框
- `add-group`: 添加群组
- `update:newGroup`: 更新新群组数据

## 组件通信模式

### 数据流向（Props Down）
```
App.vue -> Sidebar.vue
App.vue -> MainPanel.vue
App.vue -> ActionBar.vue
App.vue -> AddGroupModal.vue
MainPanel.vue -> RemoteHostEditor.vue
MainPanel.vue -> LocalHostEditor.vue
MainPanel.vue -> SystemHostPreview.vue
```

### 事件流向（Events Up）
```
Sidebar.vue -> App.vue
RemoteHostEditor.vue -> MainPanel.vue -> App.vue
LocalHostEditor.vue -> MainPanel.vue -> App.vue
SystemHostPreview.vue -> MainPanel.vue -> App.vue
ActionBar.vue -> App.vue
AddGroupModal.vue -> App.vue
```

## 设计模式与最佳实践

### 1. 单项数据流
- 所有数据从父组件流向子组件
- 子组件通过事件向父组件发送数据变更请求

### 2. 组件职责分离
- 每个组件有明确的职责边界
- UI展示与业务逻辑适当分离

### 3. 类型安全
- 使用Vue 3的props类型定义
- 确保组件接口的类型安全

### 4. 样式封装
- 使用scoped CSS避免样式冲突
- 保持组件样式的内聚性

## 技术栈

- **前端框架**: Vue 3
- **UI组件库**: Element Plus
- **桌面框架**: Wails
- **构建工具**: Vite
- **编程语言**: JavaScript

## 维护指南

### 添加新功能
1. 确定功能属于哪个现有组件
2. 如需新增组件，遵循命名规范和结构模式
3. 正确定义props和emits接口
4. 保持组件的单一职责

### 修改现有功能
1. 识别受影响的组件
2. 确保组件间的接口兼容性
3. 更新相关文档（如适用）

### 性能优化
1. 利用Vue的组件懒加载
2. 优化不必要的重渲染
3. 合理使用计算属性和侦听器

## 总结

通过此次重构，前端代码结构得到了显著改善：
- 代码可读性提升
- 组件职责更加清晰
- 便于团队协作开发
- 降低了维护成本
- 遵循了Vue 3的最佳实践

这种组件化的架构为未来的功能扩展提供了良好的基础。