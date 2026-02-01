# Ghost Host Manager - 开发者国际化指南

## 概述

本文档为开发者提供关于如何在 Ghost Host Manager 项目中添加和管理国际化内容的详细指导。

## 技术架构

### 前端国际化
- 使用 `vue-i18n` 库 (v9.x 版本) 实现 Vue 3 应用的国际化
- 采用 Composition API 方式集成
- 支持动态语言切换

### 语言包结构
语言包位于 `frontend/src/locales/` 目录下：
- `zh-CN.json` - 中文语言包
- `en-US.json` - 英文语言包
- 可扩展其他语言包

## 语言包组织结构

语言包按功能模块分类组织：

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
    },
    "actionBar": {
      "refreshRemoteGroups": "刷新远程分组"
    }
  },
  "messages": {
    "confirmDelete": "您确定要删除此分组吗？"
  },
  "tags": {
    "remote": "远程",
    "local": "本地"
  },
  "status": {
    "enabled": "已启用",
    "disabled": "已禁用"
  },
  "types": {
    "remote": "远程",
    "local": "本地"
  }
}
```

## 开发实践

### 1. 在组件中使用国际化

#### Setup 语法
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

#### 模板中使用
```html
<template>
  <div>
    <h1>{{ t('components.sidebar.title') }}</h1>
    <p>{{ t('messages.confirmDelete') }}</p>
  </div>
</template>
```

### 2. 添加新的国际化文本

当需要添加新的国际化文本时：

1. **更新语言包**：在所有支持的语言包中添加新的键值对
```json
// zh-CN.json
{
  "newFeature": {
    "buttonLabel": "新功能按钮"
  }
}

// en-US.json
{
  "newFeature": {
    "buttonLabel": "New Feature Button"
  }
}
```

2. **在组件中使用**：
```javascript
<template>
  <button>{{ t('newFeature.buttonLabel') }}</button>
</template>
```

### 3. 带参数的国际化文本

使用参数化消息：
```javascript
// 语言包中定义
{
  "messages": {
    "restoreConfirmation": "您确定要从{backup}恢复吗？\\n\\n这将覆盖您当前的数据！"
  }
}

// 在组件中使用
const message = t('messages.restoreConfirmation', { backup: 'backup_20231201.json' });
```

## 添加新语言支持

### 1. 创建语言包
在 `frontend/src/locales/` 目录下创建新语言包文件，例如 `ja-JP.json`（日语）：

```json
{
  "common": {
    "save": "保存",
    "cancel": "キャンセル",
    "confirm": "確認"
  },
  // ... 其他翻译
}
```

### 2. 更新 i18n 配置
修改 `frontend/src/i18n.js`：

```javascript
import jaJP from './locales/ja-JP.json'

const i18n = createI18n({
  locale: getLocale(),
  fallbackLocale: 'en-US',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
    'ja-JP': jaJP  // 添加新语言
  },
  legacy: false
});
```

### 3. 更新语言切换组件
修改 `frontend/src/components/LanguageSwitcher.vue` 添加新语言选项。

## 最佳实践

### 1. 语言包维护
- 保持所有语言包的键结构一致
- 定期审查未使用的国际化键
- 使用一致的键命名约定

### 2. 文本长度考虑
- 不同语言的文本长度差异很大
- 确保UI能够适应不同长度的文本
- 避免固定宽度导致文本截断

### 3. 日期和数字格式
- 对于日期、时间、数字等，使用适当的本地化格式
- 考虑使用专门的国际化库处理复杂格式化

### 4. 测试策略
- 测试所有支持的语言
- 验证动态内容的国际化
- 检查特殊字符的显示

## 故障排除

### 常见问题

**Q: 新添加的文本没有国际化？**
A: 检查是否在所有语言包中都添加了对应的键值对，以及是否正确使用了 `t()` 函数。

**Q: 语言切换后部分文本未更新？**
A: 确保组件正确使用了响应式的国际化函数，检查是否存在硬编码文本。

**Q: 特殊字符显示异常？**
A: 确保语言包文件使用UTF-8编码保存。

## 维护指南

### 1. 审查国际化覆盖率
定期检查代码中是否还有未国际化的硬编码文本：
```bash
grep -r "['\"][^'\"]\{3,\}[^'\"]*['\"]" frontend/src/ --include="*.vue" --include="*.js"
```

### 2. 同步语言包
当添加新功能时，确保同步更新所有语言包：
- 在默认语言包（如en-US）中添加新键
- 将新键同步到其他语言包
- 标记待翻译的键以便后续处理

### 3. 性能考虑
- 避免在循环中使用复杂的国际化处理
- 考虑懒加载大型语言包
- 使用合适的缓存策略

## 扩展阅读

- [vue-i18n 官方文档](https://vue-i18n.intlify.dev/)
- [国际化最佳实践](https://vue-i18n.intlify.dev/guide/best-practices.html)
- [Wails 与国际化集成](https://wails.io/)