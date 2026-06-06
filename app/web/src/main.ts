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

  // 先初始化绿联网关（获取 ugreenToken），路由守卫中所有请求需要携带 Ugreen-Ttk
  await setupUGOSCore()

  await setupRouter(app)

  setupI18n(app)

  setupAppVersionNotification()

  setupVueRootValidator(app, {
    lang: getLocale() === 'zh-CN' ? 'zh' : 'en'
  })

  app.mount('#app')
}

setupApp()
