
import { localStg } from '@/utils/storage'
import UGOSCore from '@ugreen-nas/core'
import cloudWindow from '@ugreen-nas/core/cloudWindow'

/** Setup plugin UGOS */
export async function setupUGOSCore() {
    UGOSCore.init()
    // 获取绿联授权令牌
    const info = await cloudWindow.useCapacity('getThirdToken')
    if (info?.third_token) {
        localStg.set("isUgreenEnv", true)
        localStg.set('ugreenToken', info.third_token as string)
    } else {
        localStg.set("isUgreenEnv", false)
        localStg.remove('ugreenToken')
    }
    window.UGOSCore = UGOSCore
}
