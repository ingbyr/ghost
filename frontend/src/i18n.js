import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN.json'
import enUS from './locales/en-US.json'

// 从localStorage获取用户首选语言，如果没有则默认使用中文
const getLocale = () => {
  const savedLocale = localStorage.getItem('locale')
  if (savedLocale && ['zh-CN', 'en-US'].includes(savedLocale)) {
    return savedLocale
  }
  // 检测浏览器语言
  const browserLang = navigator.language || navigator.userLanguage
  if (browserLang.startsWith('zh')) {
    return 'zh-CN'
  }
  return 'en-US'
}

const i18n = createI18n({
  locale: getLocale(), // 设置默认语言
  fallbackLocale: 'en-US', // 设置备用语言
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS
  },
  legacy: false // 使用 Composition API 模式
})

export default i18n