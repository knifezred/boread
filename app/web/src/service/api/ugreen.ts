import { request } from '../request';

/** 获取绿联用户信息（登录前确认用） */
export function fetchUgreenProfile() {
  return request<Api.Ugreen.UgreenProfile>({
    url: '/auth/ugreen-profile',
    method: 'get'
  });
}

/** 绿联NAS登录 */
export function fetchUgreenLogin() {
  return request<Api.Auth.LoginToken>({
    url: '/auth/ugreen-login',
    method: 'post'
  });
}
