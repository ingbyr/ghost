# Ghost Host Manager - 前端组件关系图

## 组件层级结构

```
App.vue (根组件)
├── Sidebar.vue (侧边栏)
├── MainPanel.vue (主面板)
│   ├── RemoteHostEditor.vue (远程Host编辑器)
│   ├── LocalHostEditor.vue (本地Host编辑器)
│   └── SystemHostPreview.vue (系统Host预览)
├── ActionBar.vue (操作栏)
└── AddGroupModal.vue (添加群组模态框)
```

## 组件依赖关系

### 1. App.vue 依赖
- Element Plus (全局注册)
- Sidebar.vue
- MainPanel.vue
- ActionBar.vue
- AddGroupModal.vue

### 2. Sidebar.vue 依赖
- Element Plus 组件（ElRow, ElCol - 如果需要）

### 3. MainPanel.vue 依赖
- RemoteHostEditor.vue
- LocalHostEditor.vue
- SystemHostPreview.vue

### 4. RemoteHostEditor.vue 依赖
- Element Plus 组件（ElRow, ElCol）

### 5. LocalHostEditor.vue 依赖
- Element Plus 组件（ElRow, ElCol）

### 6. AddGroupModal.vue 依赖
- Element Plus 组件（ElRow, ElCol）

## 数据流图

```
                    +------------------+
                    |     App.vue      |
                    | (State Manager)  |
                    +--------+---------+
                             |
         +-------------------+-------------------+
         |                   |                   |
+--------v---------+ +-------v---------+ +-------v---------+
|   Sidebar.vue    | |  MainPanel.vue  | |  ActionBar.vue  |
| (Host List)      | | (Editor Panel)  | | (Actions)       |
+--------+---------+ +--------+--------+ +-----------------+
         |                   |
         |                   +-------------------+
         |                   |                   |
+--------v---------+ +-------v---------+ +-------v---------+
| AddGroupModal.vue| |RemoteHostEditor| |LocalHostEditor |
| (Add Group)      | | (Remote Edit)   | | (Local Edit)    |
+------------------+ +----------------+ +-----------------+
                                    |
                            +-------v---------+
                            |SystemHostPreview|
                            | (System View)   |
                            +-----------------+
```

## 事件流图

```
组件事件流向（从子到父）:

Sidebar.vue
├── select-group → App.vue
├── toggle-status → App.vue
├── delete-group → App.vue
├── open-add-modal → App.vue
└── update:search-query → App.vue

MainPanel.vue
├── save-group → App.vue
├── cancel-edit → App.vue
├── fetch-remote-content → App.vue
├── mark-as-dirty → App.vue
└── refresh-system-host → App.vue

RemoteHostEditor.vue
├── save-group → MainPanel.vue
├── cancel-edit → MainPanel.vue
├── fetch-remote-content → MainPanel.vue
└── mark-as-dirty → MainPanel.vue

LocalHostEditor.vue
├── save-group → MainPanel.vue
├── cancel-edit → MainPanel.vue
└── mark-as-dirty → MainPanel.vue

SystemHostPreview.vue
└── refresh-system-host → MainPanel.vue

ActionBar.vue
├── refresh-remote → App.vue
└── refresh-list → App.vue

AddGroupModal.vue
├── close-modal → App.vue
├── add-group → App.vue
└── update:newGroup → App.vue
```

## Props 数据流向

```
App.vue → Sidebar.vue
├── groups
├── selectedGroup
├── systemHostPath
└── searchQuery

App.vue → MainPanel.vue
├── selectedGroup
├── editingGroup
├── systemHostPath
├── systemHostContent
├── remoteContentPreview
└── isFetchingRemote

MainPanel.vue → RemoteHostEditor.vue
├── group
├── editingGroup
├── remoteContentPreview
└── isFetchingRemote

MainPanel.vue → LocalHostEditor.vue
├── group
└── editingGroup

MainPanel.vue → SystemHostPreview.vue
├── systemHostPath
└── systemHostContent

App.vue → ActionBar.vue (无props，纯事件驱动)

App.vue → AddGroupModal.vue
├── showModal
└── newGroup
```

## 组件状态管理

### App.vue 状态
- groups: Host群组列表
- selectedGroup: 当前选中群组
- editingGroup: 正在编辑的群组
- newGroup: 新建群组数据
- showAddGroupModal: 模态框显示状态
- searchQuery: 搜索查询
- message/messageType: 消息提示
- isDirty: 脏数据标记
- systemHostPath/systemHostContent: 系统Host信息
- remoteContentPreview: 远程内容预览
- isRefreshingRemote/isFetchingRemote: 加载状态

### 局部组件状态
- RemoteHostEditor.vue: localEditingGroup
- LocalHostEditor.vue: localEditingGroup
- AddGroupModal.vue: localNewGroup

## 组件通信协议

### 通用事件命名约定
- `update:*` - 更新特定值的事件
- `*changed` - 值改变事件
- `*requested` - 请求执行某个操作
- `*completed` - 操作完成事件

### 错误处理模式
- 通过消息提示组件统一处理错误显示
- 异步操作使用try/catch处理错误
- 用户确认对话框用于重要操作

## 扩展建议

### 新增组件时的考虑
1. 确定组件在层级中的位置
2. 定义清晰的props接口
3. 规划合理的事件发射
4. 考虑样式封装需求
5. 验证与其他组件的兼容性

### 性能优化点
1. 合理使用v-show/v-if控制渲染
2. 组件懒加载对于大型组件
3. 合理使用计算属性缓存
4. 避免深层嵌套的数据传递