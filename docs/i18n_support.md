# Ghost Host Manager - 国际化(i18n)功能说明

## 概述

Ghost Host Manager 支持多语言界面，当前支持中文和英文两种语言。系统会自动检测用户浏览器语言偏好，并允许用户手动切换界面语言。

## 支持的语言

- 中文 (zh-CN)
- 英文 (en-US)

## 语言切换功能

### 自动检测
- 系统会自动检测用户的浏览器语言偏好
- 如果检测到中文环境，默认使用中文界面
- 否则使用英文界面

### 手动切换
- 界面右上角提供语言切换下拉菜单
- 用户可随时切换中英文界面
- 选择的语言偏好会持久化保存到localStorage

## 技术实现

### 前端技术栈
- 使用 `vue-i18n` 库实现国际化功能
- Vue 3 Composition API 集成
- Element Plus 组件库兼容

### 语言资源配置
- 语言文件存储在 `frontend/src/locales/` 目录
- 中文语言包: `zh-CN.json`
- 英文语言包: `en-US.json`
- 采用模块化结构组织翻译文本

### 语言包结构
语言包按照组件进行分类组织：

```json
{
  "common": {
    "save": "保存",
    "cancel": "取消",
    "confirm": "确认"
  },
  "components": {
    "sidebar": {
      "title": "主机分组",
      "searchPlaceholder": "搜索分组..."
    }
  },
  "messages": {
    "confirmDelete": "您确定要删除此分组吗？"
  }
}
```

## 国际化范围

### 已国际化的界面元素

#### 侧边栏 (Sidebar)
- 标题文本："Host Groups" / "主机分组"
- 搜索框占位符
- 系统Host文件条目
- 状态标签："ON"/"OFF" / "开启"/"关闭"
- 操作按钮："Delete" / "删除"

#### 操作栏 (ActionBar)
- "Refresh Remote Groups" / "刷新远程分组"
- "Backup Now" / "立即备份"
- "Restore Backup" / "恢复备份"

#### 添加群组模态框 (AddGroupModal)
- 标题和按钮文本
- 类型选择："Local Host" / "Remote Host"
- 字段标签

#### 编辑器组件 (RemoteHostEditor, LocalHostEditor)
- 表单字段标签
- 按钮文本："Save Changes" / "保存更改"
- 状态显示文本

#### 系统Host预览 (SystemHostPreview)
- 标题和字段标签

### 后端国际化考虑

当前版本中，后端Go代码保持英文错误消息，主要原因如下：
1. 错误消息主要用于日志记录和调试
2. 用户看到的消息由前端处理和国际化
3. 保持后端错误消息的一致性便于开发维护

## 扩展新语言

如需添加新语言支持，请按以下步骤操作：

1. 在 `frontend/src/locales/` 目录下创建新的语言包文件
2. 参考现有语言包的结构和分类
3. 在 `frontend/src/i18n.js` 中注册新语言
4. 在语言切换组件中添加新语言选项

## 用户体验

- 语言切换即时生效，无需刷新页面
- 所有界面元素（包括提示消息）均得到本地化
- 保持一致的用户体验和交互逻辑
- 键盘快捷键和操作流程不受语言影响

## 开发指南

### 添加新的国际化文本

1. 在对应语言包中添加新的翻译键值对
2. 在组件中使用 `{{ t('key.path') }}` 语法引用
3. 确保所有语言包都包含对应的翻译

### 组件中使用国际化

```javascript
import { useI18n } from 'vue-i18n';

export default {
  setup() {
    const { t } = useI18n();
    return { t };
  },
  template: `
    <button>{{ t('common.save') }}</button>
  `
}
```

## 故障排除

### 常见问题

**Q: 切换语言后部分文本未更新？**
A: 检查对应语言包是否包含了缺失的翻译键，或者尝试刷新页面。

**Q: 某些特殊字符显示异常？**
A: 确保语言包文件使用UTF-8编码保存。

**Q: 如何测试不同语言环境？**
A: 可以在浏览器开发者工具中临时修改navigator.language值进行测试。