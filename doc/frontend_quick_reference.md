# Ghost Host Manager - 前端快速参考

## 文件结构

```
frontend/
├── src/
│   ├── App.vue                 # 主应用组件
│   ├── main.js                 # 应用入口
│   ├── style.css              # 全局样式
│   └── components/            # 组件目录
│       ├── Sidebar.vue        # 侧边栏
│       ├── MainPanel.vue      # 主面板
│       ├── RemoteHostEditor.vue # 远程Host编辑器
│       ├── LocalHostEditor.vue  # 本地Host编辑器
│       ├── SystemHostPreview.vue # 系统Host预览
│       ├── ActionBar.vue      # 操作栏
│       └── AddGroupModal.vue  # 添加群组模态框
```

## 组件功能速查

### Sidebar.vue
**用途**: 显示Host群组列表和搜索功能
**主要功能**:
- Host群组列表展示
- 搜索过滤
- 群组状态切换
- 群组删除
- 系统Host展示

### MainPanel.vue
**用途**: 根据选中的群组类型渲染对应编辑器
**主要功能**:
- 条件渲染不同编辑器
- 协调编辑状态

### RemoteHostEditor.vue
**用途**: 编辑远程Host群组
**主要功能**:
- 编辑远程Host名称、描述、URL
- 获取远程内容预览
- 保存远程Host配置
- "获取host内容"按钮用于手动更新远程内容
- 当组启用时，内容显示为只读预览模式

### LocalHostEditor.vue
**用途**: 编辑本地Host群组
**主要功能**:
- 编辑本地Host名称、描述
- 编辑Host内容
- 保存本地Host配置
- 当组启用时，内容显示为只读预览模式

### SystemHostPreview.vue
**用途**: 预览系统Host文件（只读）
**主要功能**:
- 显示系统Host路径
- 显示系统Host内容
- 刷新系统Host内容

### ActionBar.vue
**用途**: 提供主要操作按钮
**主要功能**:
- 刷新远程群组
- 备份数据
- 刷新列表

### AddGroupModal.vue
**用途**: 添加新的Host群组
**主要功能**:
- 选择群组类型（本地/远程）
- 输入群组信息
- 创建新群组

## 常见操作示例

### 添加新组件
```javascript
// 1. 在components目录创建新组件
// 2. 在父组件中import并注册
import NewComponent from './components/NewComponent.vue'

// 3. 在components选项中注册
components: {
  NewComponent
}

// 4. 在模板中使用
<NewComponent :prop-data="data" @event="handler" />
```

### 修改现有组件
``javascript
// 示例：修改Sidebar.vue中的搜索功能
// 1. 在Sidebar.vue中找到相关代码
<input 
  :value="searchQuery"
  @input="$emit('update:search-query', $event.target.value)"
  placeholder="Search groups..." 
  class="search-input"
/>

// 2. 在props中确保定义了searchQuery
props: {
  searchQuery: {
    type: String,
    default: ''
  }
}
```

### 组件间通信
```javascript
// 父组件向子组件传递数据 (Props Down)
<Sidebar :groups="groups" :search-query="searchQuery" />

// 子组件向父组件发送事件 (Events Up)
this.$emit('update:search-query', newValue)
```

### 远程Host行为说明
```
1. 启用/禁用远程Host: 仅改变其启用状态，不会自动获取最新内容
2. 更新远程内容: 点击"获取host内容"按钮手动获取最新内容
3. 内容应用: 当远程Host启用且内容更新后，会自动应用到系统Hosts文件
```

## 调试技巧

### 检查组件状态
在浏览器控制台中：
```javascript
// 检查Vue实例
app.__vue__ // Vue 3应用实例

// 检查特定组件数据
// 在组件内部使用console.log(this.$props, this.$data)
```

### 组件生命周期
- **created**: 初始化数据
- **mounted**: DOM挂载完成，可访问DOM元素
- **updated**: 数据更新后
- **unmounted**: 组件卸载前

## 构建与运行

### 开发模式
```bash
cd frontend
npm run dev
```

### 生产构建
```bash
wails build
```

### 依赖管理
```bash
# 安装新依赖
npm install package-name

# 检查依赖更新
npm outdated
```

## 最佳实践

### 组件命名
- 使用PascalCase命名组件文件
- 组件名应体现其功能

### Props定义
- 明确定义类型
- 提供默认值（如适用）
- 使用验证函数（如需要）

### 事件命名
- 使用kebab-case命名自定义事件
- 事件名应描述发生了什么

### 样式管理
- 使用scoped样式避免冲突
- 遵循BEM命名规范
- 优先使用Element Plus样式系统

## 故障排除

### 构建错误
1. 检查组件导入/导出路径
2. 验证props和emits定义
3. 确认template语法正确

### 运行时错误
1. 检查数据类型匹配
2. 验证异步操作处理
3. 确认事件处理器定义

### 组件不更新
1. 检查响应式数据定义
2. 验证key属性使用
3. 确认computed/watch使用恰当

### 启用组的只读模式
```
1. 当本地或远程Host组启用时，编辑器将内容区域设为只读状态（但仍显示在textarea中）
2. 在只读模式下，不能编辑组的内容、名称或描述
3. 浮动保存按钮将被隐藏
4. 启用的组显示内容为只读，方便查看当前应用的内容
5. 要编辑内容，需要先禁用该组，此时内容区域变为可编辑状态
```
