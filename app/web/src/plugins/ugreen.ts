
import { localStg } from '@/utils/storage'
import UGOSCore from '@ugreen-nas/core'
import cloudWindow from '@ugreen-nas/core/cloudWindow'

/** 整体超时（毫秒），防止非绿联环境一直加载 */
const UGOS_INIT_TIMEOUT = 3000

/** Setup plugin UGOS
 *  整体超时不超过 3 秒：init + getThirdToken 共用一个 deadline，
 *  任一环节卡住均视为非 UGOS 环境，静默降级。
 */
export async function setupUGOSCore() {
    try {
      const info = await Promise.race([
        initAndGetToken(),
        new Promise<null>((_, reject) =>
          setTimeout(() => reject(new Error(`UGOS setup timeout (${UGOS_INIT_TIMEOUT}ms)`)), UGOS_INIT_TIMEOUT)
        ),
      ])
      if (info?.third_token) {
        localStg.set('isUgreenEnv', true)
        localStg.set('ugreenToken', info.third_token as string)
      } else {
        localStg.set('isUgreenEnv', false)
        localStg.remove('ugreenToken')
      }
    } catch {
      // 超时或异常 → 非 UGOS 环境，静默降级
      localStg.set('isUgreenEnv', false)
      localStg.remove('ugreenToken')
    }
    window.UGOSCore = UGOSCore
}

/** init bridge 后获取 third_token */
async function initAndGetToken() {
    // init() 返回 Promise，必须 await 等待 bridge 就绪，否则调用 getThirdToken 会失败
    await UGOSCore.init().catch(() => {})
    return await cloudWindow.useCapacity('getThirdToken') as { third_token?: string } | null
}
