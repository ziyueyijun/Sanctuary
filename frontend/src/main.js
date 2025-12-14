/**
 * Sanctuary - 前端入口文件
 * 
 * 初始化 Vue 3 应用并注册 Arco Design Vue 组件库
 */
import { createApp } from 'vue'
import ArcoVue from '@arco-design/web-vue'       // Arco Design Vue 组件库
import '@arco-design/web-vue/dist/arco.css'     // Arco Design 样式
import App from './App.vue'                      // 根组件
import './style.css'                             // 全局样式

// 创建并挂载 Vue 应用
const app = createApp(App)
app.use(ArcoVue)  // 注册 Arco Design 组件
app.mount('#app') // 挂载到 DOM
