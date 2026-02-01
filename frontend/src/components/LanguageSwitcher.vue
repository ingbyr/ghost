<template>
  <div class="language-switcher">
    <el-dropdown @command="handleLanguageChange" placement="bottom-end">
      <span class="el-dropdown-link">
        {{ currentLanguage }}
        <el-icon><arrow-down /></el-icon>
      </span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="zh-CN" :class="{ active: currentLocale === 'zh-CN' }">
            中文
          </el-dropdown-item>
          <el-dropdown-item command="en-US" :class="{ active: currentLocale === 'en-US' }">
            English
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';
import { ArrowDown } from '@element-plus/icons-vue';
import { computed } from 'vue';

export default {
  name: 'LanguageSwitcher',
  components: {
    ArrowDown
  },
  setup() {
    const { locale, t } = useI18n();
    
    const handleLanguageChange = (lang) => {
      locale.value = lang;
      // 保存用户选择的语言到localStorage
      localStorage.setItem('locale', lang);
    };
    
    // 获取当前语言显示名称
    const currentLanguage = computed(() => {
      return locale.value === 'zh-CN' ? '中文' : 'English';
    });
    
    return {
      locale,
      currentLocale: locale,
      currentLanguage,
      handleLanguageChange,
      t
    };
  }
};
</script>

<style scoped>
.language-switcher {
  margin-left: 10px;
}

.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 4px;
  background: #f0f2f5;
  transition: background-color 0.3s;
}

.el-dropdown-link:hover {
  background: #e0e2e6;
}

.el-dropdown-menu {
  min-width: 100px;
}

.el-dropdown-item.active {
  background-color: #f0f9ff;
  color: #409eff;
}
</style>