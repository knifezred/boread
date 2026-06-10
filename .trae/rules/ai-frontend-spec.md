# 前端 AI 代码生成规范 (Vue 3 + TypeScript)

---

## 1. 项目结构

```
web/
├── src/
│   ├── views/                        # 页面组件（自动生成路由）
│   │   ├── admin/                    # 后台管理页面
│   │   │   ├── system/               # 系统管理（user/dept/role/menu/dict）
│   │   │   │   ├── dept/
│   │   │   │   │   ├── index.vue                 # 列表页
│   │   │   │   │   └── modules/
│   │   │   │   │       ├── dept-search.vue       # 搜索区
│   │   │   │   │       └── dept-operate-modal.vue # 弹窗
│   │   │   │   ├── user/
│   │   │   │   ├── role/
│   │   │   │   └── dict/
│   │   │   └── library/             # 电子书管理
│   │   ├── _builtin/                # 内置页面（login/setup/403/404/500）
│   │   ├── bookshelf/               # 前台页面
│   │   ├── book-detail/
│   │   ├── book-reader/
│   │   ├── home/
│   │   └── admin/.../               # 其他管理页面
│   ├── service/
│   │   ├── api/                     # API 调用
│   │   │   ├── index.ts             # 统一 re-export
│   │   │   ├── auth.ts
│   │   │   ├── system-manage.ts
│   │   │   ├── book-manage.ts
│   │   │   └── chapter-manage.ts
│   │   └── request/index.ts         # createFlatRequest 实例化
│   ├── typings/
│   │   ├── api/                     # API 类型定义（declare namespace）
│   │   │   ├── common.d.ts
│   │   │   ├── system-manage.d.ts
│   │   │   ├── book-manage.d.ts
│   │   │   └── auth.d.ts
│   │   └── app.d.ts                 # App 全局类型（含 I18n.Schema）
│   ├── locales/
│   │   ├── langs/zh-cn.ts           # 中文 i18n（App.I18n.Schema 类型）
│   │   ├── langs/en-us.ts           # 英文 i18n
│   │   └── locale.ts                # 合并导出
│   ├── constants/
│   │   └── business.ts              # 业务常量（enableStatusRecord 等）
│   ├── utils/
│   │   ├── service.ts               # createServiceConfig / getServiceBaseURL
│   │   └── common.ts                # transformRecordToOption 等
│   ├── hooks/common/                # 通用 Hooks
│   └── store/modules/               # Pinia store
```

---

## 2. 页面结构规范

每个管理模块遵循以下三文件结构：

```
admin/system/模块名/
├── index.vue                           # 列表页（useNaivePaginatedTable + useTableOperate）
└── modules/
    ├── 模块名-search.vue                # 搜索组件
    └── 模块名-operate-modal.vue          # 新增/编辑弹窗
```

---

## 3. 命名规范

| 元素 | 规范 | 示例 |
|------|------|------|
| 目录 | 全小写 kebab-case | `book-tag`, `user-detail` |
| .vue 文件 | kebab-case | `dept-operate-modal.vue` |
| 组件名（defineOptions） | PascalCase | `defineOptions({ name: 'DeptOperateModal' })` |
| 函数名 | camelCase | `handleAdd`, `fetchGetDeptList` |
| TS/TSX 文件 | kebab-case | `system-manage.ts` |
| API 函数 | `fetch` 前缀 | `fetchGetDeptList`, `fetchCreateDept` |
| Typing namespace | 嵌套 namespace | `Api.SystemManage.Dept` |
| i18n key 路径模式 | `page.admin.system.模块.字段` | `page.admin.system.dept.deptName` |

---

## 4. TypeScript 类型定义规范

**规则**:
- 统一在 `declare namespace Api { namespace Xxx { ... } }` 下
- 继承 `Common.CommonRecord` 获得基础字段（`id`/`createBy`/`createTime`/`updateBy`/`updateTime`/`status`）
- 搜索参数用 `CommonType.RecordNullable<Pick<...> & CommonSearchParams>`
- 分页列表用 `Common.PaginatingQueryRecord<T>`

---

## 5. 路由

**规则**: AI **无需手动处理**路由文件。`elegant-router` 会自动从 `views/` 目录结构生成 4 个文件（`routes.ts` + `imports.ts` + `transform.ts` + `elegant-router.d.ts`）。

---

## 6. 状态常量

```ts
// src/constants/business.ts
export const enableStatusRecord: Record<Api.Common.EnableStatus, App.I18n.I18nKey> = {
  '1': 'common.enable',
  '2': 'common.disable'
};
export const enableStatusOptions = transformRecordToOption(enableStatusRecord);
```

---

## 7. 前端错误处理

后端统一返回 `{ code, message, data }` 结构，前端 `request` 层自动处理（见 [`request/index.ts`](app/web/src/service/request/index.ts)）：

```ts
// 后端响应格式
interface ServiceResponse<T> {
  code: string;    // "0000"=成功，其余=失败
  message: string;
  data: T;
}
```

**前端调用约定**：
- `isBackendSuccess` 判断：`code === VITE_SERVICE_SUCCESS_CODE`（默认 `"0000"`）
- `onBackendFail` 自动处理：401/403 弹登录弹窗、过期 token 自动刷新
- 登录过期 codes：`VITE_SERVICE_LOGOUT_CODES`（如 `"2001,2002,2003"`）
- modal 弹窗 codes：`VITE_SERVICE_MODAL_LOGOUT_CODES`（如 `"2008"`）

**组件中错误处理**：
```ts
const { error } = await fetchCreateDept(data);
if (error) {
  window.$message?.error(error.message);
} else {
  window.$message?.success($t('common.updateSuccess'));
}
```

---

## 8. AI 编码约束

- **必须** 使用 `$t()` 引用所有展示文本
- **必须** 用 `request<T>()` 调用 API，必须带泛型参数（返回 `{ data, error }` 结构）
- **禁止** 使用 `any` 类型
- **必须** 给 `.vue` 组件写 `defineOptions({ name: 'Xxx' })`
- **必须** 用 `defineModel<boolean>('visible')` 控制弹窗/抽屉显隐
- **禁止** 直接修改路由文件（elegant-router 自动生成）
- **新建业务模块** 必须遵循 `三文件结构`（index + search + operate-modal/drawer）
- **新建 API 函数** 必须在 `src/service/api/index.ts` 中 re-export
- **列渲染** 用 TSX（`lang="tsx"`），页面逻辑用 TS（`lang="ts"`）
- **弹窗组件** 用 `<NModal preset="card" class="w-600px">`，抽屉用 `<NDrawer :width="520">`
- **`.d.ts` 文件用 `declare namespace`**，全局生效，不需要 export
- **后端请求失败**统一用 `error.message` 显示，不直接显示原始错误
- **消息提示/通知/对话框** 必须用全局挂载的 `window.$message`、`window.$notification`、`window.$dialog`、`window.$loadingBar`（在 `app-provider.vue` 中注册），**禁止** import `useMessage`/`useNotification`/`useDialog`/`useLoadingBar` 从 naive-ui
