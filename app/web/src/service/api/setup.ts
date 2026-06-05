import { request } from '../request';

/** 系统初始化 - 数据库配置请求 */
interface DbConfigRequest {
  host: string;
  port: number;
  username: string;
  password: string;
  dbname: string;
}

/** 检查系统是否已配置数据库 */
export function fetchGetSetupStatus() {
  return request<{ configured: boolean }>({
    url: '/setup/status',
    method: 'get',
  });
}

/** 保存数据库配置 */
export function fetchSaveDatabaseConfig(data: DbConfigRequest) {
  return request<null>({
    url: '/setup/database',
    method: 'post',
    data,
  });
}
