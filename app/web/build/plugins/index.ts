import type { PluginOption } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import progress from 'vite-plugin-progress'
import vueRootValidator from 'vite-plugin-vue-transition-root-validator'
import { setupElegantRouter } from './router'
import { setupUnocss } from './unocss'
import { setupUnplugin } from './unplugin'
import { setupHtmlPlugin } from './html'
import { setupDevtoolsPlugin } from './devtools'
import { UgosViteBuilder } from '@ugreen-nas/builder-open'

export function setupVitePlugins(viteEnv: Env.ImportMeta, buildTime: string) {

  const builder = new UgosViteBuilder({
    windowConfig: {
      width: 1200,
      height: 800,
      hideTitle: true,
    },
    getIgnoreFolder: (c: unknown) => c,
  })
  const plugins: PluginOption = [
    vue(),
    vueJsx(),
    setupDevtoolsPlugin(viteEnv),
    setupElegantRouter(),
    setupUnocss(viteEnv),
    ...setupUnplugin(viteEnv),
    progress(),
    setupHtmlPlugin(buildTime),
    vueRootValidator(),
    // builder.pluginEntry()
  ]

  return plugins
}
