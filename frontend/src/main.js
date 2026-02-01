import { createApp } from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import './style.css';
import i18n from './i18n'; // 导入国际化配置

const app = createApp(App);

app.use(ElementPlus);
app.use(i18n); // 使用国际化插件

app.mount('#app');
