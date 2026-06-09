import { createApp } from 'vue'
import './plugins/assets'
import { setupVueRootValidator } from 'vite-plugin-vue-transition-root-validator/client'
import { setupAppVersionNotification, setupDayjs, setupIconifyOffline, setupLoading, setupNProgress, setupUGOSCore } from './plugins'
import { setupStore } from './store'
import { setupRouter } from './router'
import { getLocale, setupI18n } from './locales'
import App from './App.vue'

async function setupApp() {

  setupLoading()

  setupNProgress()

  setupIconifyOffline()

  setupDayjs()


  const app = createApp(App)

  setupStore(app)

  // 启动 UGOS 初始化（不阻塞路由装配）
  // 已登录用户的请求在 SDK 初始化完成前可能缺少 Ugreen-Ttk，绿联环境网关会容错处理
  setupUGOSCore()

  await setupRouter(app)

  setupI18n(app)

  setupAppVersionNotification()

  setupVueRootValidator(app, {
    lang: getLocale() === 'zh-CN' ? 'zh' : 'en'
  })

  app.mount('#app')
}

setupApp()
