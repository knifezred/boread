import { request } from "../request";

// =====================================================================
// 后端: boread (Go + Gin)
// 路径约定: /api 作为 baseURL, 此处只写 /manage/* 子路径
// 响应分页结构: { records, current, size, total } 对齐前端 Common.PaginatingQueryRecord
// =====================================================================

// -------- Role --------

/** 角色分页 */
export function fetchGetRoleList(params?: Api.SystemManage.RoleSearchParams) {
  return request<Api.SystemManage.RoleList>({
    url: "/manage/role/page",
    method: "post",
    data: params,
  });
}

/** 全量角色 (启用, 下拉用) */
export function fetchGetAllRoles() {
  return request<Api.SystemManage.AllRole[]>({
    url: "/manage/role/all",
    method: "get",
  });
}

/** 角色详情 */
export function fetchGetRole(id: string | number) {
  return request<Api.SystemManage.Role>({
    url: `/manage/role/${id}`,
    method: "get",
  });
}

/** 新增角色 */
export function fetchCreateRole(
  data: Partial<Api.SystemManage.Role> & {
    dataScope: string;
    deptIds?: number[];
  },
) {
  return request<Api.SystemManage.Role>({
    url: "/manage/role",
    method: "post",
    data,
  });
}

/** 编辑角色 */
export function fetchUpdateRole(
  id: string | number,
  data: Partial<Api.SystemManage.Role> & {
    dataScope: string;
    deptIds?: number[];
  },
) {
  return request<Api.SystemManage.Role>({
    url: `/manage/role/${id}`,
    method: "put",
    data,
  });
}

/** 删除角色 */
export function fetchDeleteRole(id: string | number) {
  return request<null>({
    url: `/manage/role/${id}`,
    method: "delete",
  });
}

/** 查角色已授权菜单 IDs */
export function fetchGetRoleMenuIds(id: string | number) {
  return request<number[]>({
    url: `/manage/role/${id}/menus`,
    method: "get",
  });
}

/** 授权菜单 */
export function fetchGrantRoleMenus(id: string | number, menuIds: number[]) {
  return request<null>({
    url: `/manage/role/${id}/menus`,
    method: "put",
    data: { menuIds },
  });
}

/** 查角色已授权按钮 IDs */
export function fetchGetRoleButtonIds(id: string | number) {
  return request<number[]>({
    url: `/manage/role/${id}/buttons`,
    method: "get",
  });
}

/** 授权按钮 */
export function fetchGrantRoleButtons(
  id: string | number,
  buttonIds: number[],
) {
  return request<null>({
    url: `/manage/role/${id}/buttons`,
    method: "put",
    data: { buttonIds },
  });
}

// -------- User --------

/** 用户分页 */
export function fetchGetUserList(params?: Api.SystemManage.UserSearchParams) {
  return request<Api.SystemManage.UserList>({
    url: "/manage/user/page",
    method: "post",
    data: params,
  });
}

/** 用户详情 */
export function fetchGetUser(id: string | number) {
  return request<Api.SystemManage.User>({
    url: `/manage/user/${id}`,
    method: "get",
  });
}

/** 新增用户 */
export function fetchCreateUser(
  data: Partial<Api.SystemManage.User> & {
    password: string;
    roleIds?: number[];
    deptId?: number | null;
  },
) {
  return request<Api.SystemManage.User>({
    url: "/manage/user",
    method: "post",
    data,
  });
}

/** 编辑用户 */
export function fetchUpdateUser(
  id: string | number,
  data: Partial<Api.SystemManage.User> & {
    roleIds?: number[];
    deptId?: number | null;
  },
) {
  return request<Api.SystemManage.User>({
    url: `/manage/user/${id}`,
    method: "put",
    data,
  });
}

/** 删除用户 */
export function fetchDeleteUser(id: string | number) {
  return request<null>({
    url: `/manage/user/${id}`,
    method: "delete",
  });
}

/** 批量删除用户 */
export function fetchBatchDeleteUser(ids: (string | number)[]) {
  return request<null>({
    url: "/manage/user/batch-delete",
    method: "post",
    data: { ids },
  });
}

/** 重置密码 */
export function fetchResetUserPassword(id: string | number, password: string) {
  return request<null>({
    url: `/manage/user/${id}/reset-password`,
    method: "put",
    data: { password },
  });
}

// -------- Menu --------

/** 菜单分页列表 */
export function fetchGetMenuList(params?: Api.SystemManage.MenuSearchParams) {
  return request<Api.SystemManage.MenuList>({
    url: "/manage/menu/page",
    method: "post",
    data: params,
  });
}

/** 全量菜单树 (含按钮, 用于权限分配) */
export function fetchGetMenuTree() {
  return request<Api.SystemManage.MenuTreeNode[]>({
    url: "/manage/menu/tree",
    method: "get",
  });
}

/**
 * 全部页面组件 key (用于菜单编辑下拉 component 字段)
 *
 * boread 后端当前未提供该接口, 返回空数组. 实际项目可由 elegant-router 在前端枚举生成.
 */
export async function fetchGetAllPages() {
  return { data: [] as string[], error: null };
}

/** 菜单详情 */
export function fetchGetMenu(id: string | number) {
  return request<Api.SystemManage.Menu>({
    url: `/manage/menu/${id}`,
    method: "get",
  });
}

/** 新增菜单 */
export function fetchCreateMenu(data: Partial<Api.SystemManage.Menu>) {
  return request<Api.SystemManage.Menu>({
    url: "/manage/menu",
    method: "post",
    data,
  });
}

/** 编辑菜单 */
export function fetchUpdateMenu(
  id: string | number,
  data: Partial<Api.SystemManage.Menu>,
) {
  return request<Api.SystemManage.Menu>({
    url: `/manage/menu/${id}`,
    method: "put",
    data,
  });
}

/** 删除菜单 */
export function fetchDeleteMenu(id: string | number) {
  return request<null>({
    url: `/manage/menu/${id}`,
    method: "delete",
  });
}

/** 按菜单查按钮 */
export function fetchGetMenuButtons(menuId: string | number) {
  return request<Api.SystemManage.MenuButton[]>({
    url: `/manage/menu/buttons/${menuId}`,
    method: "get",
  });
}

/** 新增菜单按钮 */
export function fetchCreateMenuButton(data: {
  menuId: number;
  buttonCode: string;
  buttonDesc?: string;
}) {
  return request<Api.SystemManage.MenuButton>({
    url: "/manage/menu/button",
    method: "post",
    data,
  });
}

/** 删除菜单按钮 */
export function fetchDeleteMenuButton(id: string | number) {
  return request<null>({
    url: `/manage/menu/button/${id}`,
    method: "delete",
  });
}

// -------- Dept --------

/** 部门树 */
export function fetchGetDeptTree(params?: {
  deptName?: string;
  deptCode?: string;
  status?: string;
}) {
  return request<Api.SystemManage.Dept[]>({
    url: "/manage/dept/tree",
    method: "get",
    params,
  });
}

/** 部门分页列表 */
export function fetchGetDeptList(params?: Api.SystemManage.DeptSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.Dept>>({
    url: "/manage/dept/page",
    method: "post",
    data: params,
  });
}

/** 部门详情 */
export function fetchGetDept(id: string | number) {
  return request<Api.SystemManage.Dept>({
    url: `/manage/dept/${id}`,
    method: "get",
  });
}

/** 新增部门 */
export function fetchCreateDept(data: Record<string, any>) {
  return request<Api.SystemManage.Dept>({
    url: "/manage/dept",
    method: "post",
    data,
  });
}

/** 编辑部门 */
export function fetchUpdateDept(
  id: string | number,
  data: Record<string, any>,
) {
  return request<Api.SystemManage.Dept>({
    url: `/manage/dept/${id}`,
    method: "put",
    data,
  });
}

/** 删除部门 */
export function fetchDeleteDept(id: string | number) {
  return request<null>({
    url: `/manage/dept/${id}`,
    method: "delete",
  });
}

// -------- Dict --------

/** 字典分页 */
export function fetchGetDictList(params?: Api.SystemManage.DictSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.Dict>>({
    url: "/manage/dict/page",
    method: "post",
    data: params,
  });
}

/** 字典详情 */
export function fetchGetDict(id: string | number) {
  return request<Api.SystemManage.Dict>({
    url: `/manage/dict/${id}`,
    method: "get",
  });
}

/** 新增字典 */
export function fetchCreateDict(data: Record<string, any>) {
  return request<Api.SystemManage.Dict>({
    url: "/manage/dict",
    method: "post",
    data,
  });
}

/** 编辑字典 */
export function fetchUpdateDict(
  id: string | number,
  data: Record<string, any>,
) {
  return request<Api.SystemManage.Dict>({
    url: `/manage/dict/${id}`,
    method: "put",
    data,
  });
}

/** 删除字典 */
export function fetchDeleteDict(id: string | number) {
  return request<null>({
    url: `/manage/dict/${id}`,
    method: "delete",
  });
}

/** 按字典 ID 查项 */
export function fetchGetDictItems(dictId: string | number) {
  return request<Api.SystemManage.DictItem[]>({
    url: `/manage/dict/items/${dictId}`,
    method: "get",
  });
}

/** 按字典 code 拉项 (前端下拉高频) */
export function fetchGetDictItemsByCode(code: string) {
  return request<Api.SystemManage.DictItem[]>({
    url: `/manage/dict/code/${code}`,
    method: "get",
  });
}

/** 新增字典项 */
export function fetchCreateDictItem(data: {
  dictId: number;
  itemLabel: string;
  itemValue: string;
  itemDesc?: string;
  sortOrder?: number;
  status?: string;
}) {
  return request<Api.SystemManage.DictItem>({
    url: "/manage/dict/item",
    method: "post",
    data,
  });
}

/** 编辑字典项 */
export function fetchUpdateDictItem(
  id: string | number,
  data: {
    dictId: number;
    itemLabel: string;
    itemValue: string;
    itemDesc: string;
    sortOrder: number;
    status: string;
  },
) {
  return request<Api.SystemManage.DictItem>({
    url: `/manage/dict/item/${id}`,
    method: "put",
    data,
  });
}

/** 删除字典项 */
export function fetchDeleteDictItem(id: string | number) {
  return request<null>({
    url: `/manage/dict/item/${id}`,
    method: "delete",
  });
}

// -------- Log --------

/** 登录日志分页 */
export function fetchGetLoginLogList(params?: Api.SystemManage.LoginLogSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<any>>({
    url: "/manage/log/login/page",
    method: "post",
    data: params,
  });
}

/** 操作日志分页 */
export function fetchGetOperationLogList(params?: Api.SystemManage.OperationLogSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<any>>({
    url: "/manage/log/operation/page",
    method: "post",
    data: params,
  });
}

// -------- Setting --------

/** 系统配置分页 */
export function fetchGetSettingList(params?: Api.SystemManage.SettingSearchParams) {
  return request<Api.Common.PaginatingQueryRecord<Api.SystemManage.Setting>>({
    url: "/manage/setting/page",
    method: "post",
    data: params,
  });
}

/** 配置详情 */
export function fetchGetSettingById(id: string | number) {
  return request<Api.SystemManage.Setting>({
    url: `/manage/setting/${id}`,
    method: "get",
  });
}

/** 新增配置 */
export function fetchCreateSetting(data: {
  category: string;
  key: string;
  value: string;
  valueType: string;
  description?: string;
  editable?: boolean;
  status?: string;
}) {
  return request<Api.SystemManage.Setting>({
    url: "/manage/setting",
    method: "post",
    data,
  });
}

/** 编辑配置 */
export function fetchUpdateSetting(
  id: string | number,
  data: {
    value: string;
    valueType: string;
    description?: string;
    status?: string;
  },
) {
  return request<Api.SystemManage.Setting>({
    url: `/manage/setting/${id}`,
    method: "put",
    data,
  });
}

/** 删除配置 */
export function fetchDeleteSetting(id: string | number) {
  return request<null>({
    url: `/manage/setting/${id}`,
    method: "delete",
  });
}

/** 获取配置分类列表 */
export function fetchGetSettingCategories() {
  return request<string[]>({
    url: "/manage/setting/categories",
    method: "get",
  });
}

/** 按分类获取所有配置 */
export function fetchGetSettingsByCategory(category: string) {
  return request<Api.SystemManage.SettingCategoryMap>({
    url: `/manage/setting/by-category/${encodeURIComponent(category)}`,
    method: "get",
  });
}

/** 批量保存配置 (upsert) */
export function fetchBatchSaveSettings(data: Api.SystemManage.SettingBatchSaveRequest) {
  return request<Api.SystemManage.SettingBatchSaveResult>({
    url: "/manage/setting/batch-save",
    method: "post",
    data,
  });
}
