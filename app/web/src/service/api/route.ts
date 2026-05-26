import { request } from '../request';

/**
 * get constant routes
 *
 * 常量路由 = 无需登录即可访问的路由(login/403/404 等), 通常由前端 router/elegant 直接定义.
 * 后端 boread 当前未提供该接口, 返回空数组让 route store 直接走前端静态路由.
 */
export async function fetchGetConstantRoutes() {
  return { data: [] as Api.Route.MenuRoute[], error: null };
}

/** get user routes */
export function fetchGetUserRoutes() {
  return request<Api.Route.UserRoute>({ url: '/auth/menu' });
}

/**
 * whether the route is exist
 *
 * boread 后端当前未实现此查询, 简单返回 true (路由是否存在由前端 router 自己判断)
 *
 * @param _routeName route name
 */
export async function fetchIsRouteExist(_routeName: string) {
  return { data: true, error: null };
}
