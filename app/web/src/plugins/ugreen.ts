
import { localStg } from '@/utils/storage'
import UGOSCore from '@ugreen-nas/core'
import cloudWindow from '@ugreen-nas/core/cloudWindow'

/** Setup plugin UGOS */
export async function setupUGOSCore() {
    UGOSCore.init()
    // 获取绿联授权令牌（非UGOS环境3秒超时退出，避免Windows等平台一直加载）
    const info = await Promise.race([
      cloudWindow.useCapacity('getThirdToken'),
      new Promise<null>((_, reject) =>
        setTimeout(() => reject(new Error('UGOS getThirdToken timeout (3s)')), 3000)
      ),
    ]).catch(() => null)
    if (info?.third_token) {
        localStg.set("isUgreenEnv", true)
        localStg.set('ugreenToken', info.third_token as string)
    } else {
        localStg.set("isUgreenEnv", false)
        localStg.remove('ugreenToken')
    }
    window.UGOSCore = UGOSCore
}
