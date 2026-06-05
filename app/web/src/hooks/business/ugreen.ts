import { ref } from 'vue'

/** Promise 超时辅助 */
function withTimeout<T>(promise: Promise<T>, ms: number): Promise<T> {
  return Promise.race([
    promise,
    new Promise<never>((_, reject) => setTimeout(() => reject(new Error('timeout')), ms))
  ]);
}

/** 检测是否运行在 UGOS 环境中（通过实际获取 token 验证环境） */
export function useUgreenEnv() {
  const isUgreenEnv = ref(false);

  async function checkUgreenEnv(): Promise<boolean> {
    try {
      const UGOSCore = await import('@ugreen-nas/core');
      UGOSCore.default.init();

      // 进一步验证：尝试获取 token，确保环境真实可用
      const token = await getUgreenToken(3000);
      if (token) {
        isUgreenEnv.value = true;
        return true;
      }
    } catch {
      // 非UGOS环境或超时，降级
    }
    isUgreenEnv.value = false;
    return false;
  }

  return { isUgreenEnv, checkUgreenEnv };
}

/**
 * 获取 UGREEN third_token
 * @param timeoutMs 超时时间（毫秒），默认 3000
 */
export async function getUgreenToken(timeoutMs = 3000): Promise<string | null> {
  try {
    const cloudWindow = await import('@ugreen-nas/core/cloudWindow');
    const info = await withTimeout(cloudWindow.default.useCapacity('getThirdToken'), timeoutMs);
    return (info?.third_token as string) || null;
  } catch {
    return null;
  }
}
