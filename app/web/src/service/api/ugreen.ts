import { request } from '../request';

/** 获取绿联用户信息（登录前确认用） */
export function fetchUgreenProfile() {
  return request<Api.Ugreen.UgreenProfile>({
    url: '/auth/ugreen-profile',
    method: 'get'
  });
}

/** 绿联NAS登录
 * @param allowCreate 是否允许自动创建新用户（确认按钮时传 true，路由守卫自动登录传 false）
 */
export function fetchUgreenLogin(allowCreate = false) {
  return request<Api.Auth.LoginToken>({
    url: '/auth/ugreen-login',
    method: 'post',
    params: { create: allowCreate }
  });
}
