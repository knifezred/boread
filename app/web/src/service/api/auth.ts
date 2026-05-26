import { request } from '../request';

/**
 * Login
 *
 * @param userName User name
 * @param password Password
 */
export function fetchLogin(userName: string, password: string) {
  return request<Api.Auth.LoginToken>({
    url: '/auth/login',
    method: 'post',
    data: {
      // boread 后端接收字段名为 username (小写), 此处保留前端命名一致性, 后端 JSON 标签匹配
      username: userName,
      password
    }
  });
}

/** Get user info */
export function fetchGetUserInfo() {
  return request<Api.Auth.UserInfo>({ url: '/auth/userInfo' });
}

/**
 * Refresh token
 *
 * 当前后端 MVP 阶段尚未实现独立 refresh 流程, 临时复用 login token 结构,
 * 此函数走 userInfo 触发重新登录态校验 (token 过期会触发 401, 由全局 logout 处理)
 *
 * @param _refreshToken 暂未使用, 保留签名以兼容 axios 框架
 */
export function fetchRefreshToken(_refreshToken: string) {
  return request<Api.Auth.LoginToken>({
    url: '/auth/refresh-token',
    method: 'post',
    data: { refreshToken: _refreshToken }
  });
}

/**
 * return custom backend error (调试用, boread 未实现)
 *
 * @param code error code
 * @param msg error message
 */
export function fetchCustomBackendError(code: string, msg: string) {
  return request({ url: '/auth/error', params: { code, msg } });
}
