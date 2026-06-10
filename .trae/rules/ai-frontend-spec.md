# 前端 AI 代码生成规范 (Vue 3 + TypeScript)

---

## 1. 项目结构

```
web/src/
├── views/                        # 页面组件（elegant-router 自动生成路由）
│   └── admin/system/{模块}/
│       ├── index.vue                       # 列表页
│       └── modules/
│           ├── {模块}-search.vue            # 搜索组件
│           └── {模块}-operate-{drawer|modal}.vue  # 新增/编辑弹窗
├── service/api/                  # API 调用，fetchXxx 命名
│   ├── index.ts                  # 统一 re-export
│   └── *.ts
├── typings/api/                  # declare namespace Api { ... }
│   ├── common.d.ts               # Common.CommonRecord, PaginatingQueryRecord 等
│   └── *.d.ts
├── constants/business.ts         # 业务常量（enableStatusRecord 等）
├── locales/langs/                # i18n
├── utils/                        # 工具函数
├── hooks/common/                 # 通用 Hooks
└── store/modules/                # Pinia store
```

---

## 2. 页面结构规范

每个业务模块必须是三文件结构：

```
views/admin/system/{模块}/
├── index.vue                           # useNaivePaginatedTable + useTableOperate
└── modules/
    ├── {模块}-search.vue                # 搜索区
    └── {模块}-operate-{drawer|modal}.vue  # 新增/编辑弹窗
```

**弹窗选型**：根据表单字段数量决定 —— 字段多（≥6 项，如 user）用 `NDrawer`，字段少（<6 项，如 dept/dict）用 `NModal`。文件命名为 `{模块}-operate-drawer.vue` 或 `{模块}-operate-modal.vue`。

**参考模块**：`views/admin/system/user/` 是标准实现，新建模块须严格对齐其代码风格。

---

## 3. 命名规范

| 元素 | 规范 | 示例 |
|------|------|------|
| 目录 | kebab-case | `user-detail` |
| .vue 文件 | kebab-case | `user-operate-modal.vue` |
| 组件名（defineOptions） | PascalCase | `defineOptions({ name: 'UserOperateModal' })` |
| API 函数 | `fetch` 前缀 | `fetchGetUserList`, `fetchCreateUser` |
| i18n key | `page.admin.system.模块.字段` | `page.admin.system.user.userName` |
| Typing namespace | 嵌套 namespace | `Api.SystemManage.User` |

---

## 4. TypeScript 类型定义

- 统一在 `declare namespace Api { namespace Xxx { ... } }` 下，全局生效不 export
- 实体类型继承 `Common.CommonRecord` 获得 `id`/`createBy`/`createTime`/`updateBy`/`updateTime`/`status`
- 搜索参数用 `CommonType.RecordNullable<Pick<...> & CommonSearchParams>`
- 分页列表定义 `xxxList = Common.PaginatingQueryRecord<T>`
- **禁止使用 `any`**

---

## 5. 路由

**无需手动处理**路由文件。`elegant-router` 根据 `views/` 目录结构自动生成。

---

## 6. 状态常量

定义在 `constants/business.ts`，参考 `enableStatusRecord` 模式：

```ts
export const enableStatusRecord: Record<Api.Common.EnableStatus, App.I18n.I18nKey> = {
  '1': 'common.enable',
  '2': 'common.disable'
};
export const enableStatusOptions = transformRecordToOption(enableStatusRecord);
```

---

## 7. 错误处理

后端统一 `{ code, message, data }` 结构。`request` 层已封装 `{ data, error }` 返回。

**组件中调用**：
```ts
const { error } = await fetchCreateDept(data);
if (error) {
  window.$message?.error(error.message);   // 不直接显示原始错误
} else {
  window.$message?.success($t('common.updateSuccess'));
}
```

**消息提示**必须用全局挂载的 `window.$message` / `window.$notification` / `window.$dialog` / `window.$loadingBar`，**禁止** import `useMessage` 等。

---

## 8. 样式

- **必须**用 naive-ui 组件库
- **推荐**用 unocss 修改样式
- 禁止引入新的组件库

---

## 9. AI 编码约束

- **必须**用 `$t()` 引用所有展示文本
- **必须**用 `request<T>()` 调用 API，必须带泛型参数
- **禁止**使用 `any`
- **必须**给 `.vue` 组件写 `defineOptions({ name: 'Xxx' })`
- **必须**用 `defineModel<boolean>('visible')` 控制弹窗/抽屉显隐
- **禁止**直接修改路由文件
- **新建业务模块**必须遵循三文件结构，参考 `views/admin/system/user/`
- **新建 API 函数**必须在 `service/api/index.ts` 中 re-export
- **列渲染**用 TSX（`lang="tsx"`），页面逻辑用 TS（`lang="ts"`）
- **后端请求失败**统一用 `error.message` 显示
- **禁止** import `useMessage`/`useNotification`/`useDialog`/`useLoadingBar`，改用 `window.$message` 等

---

## 10. 参考实现

`views/admin/system/user/` 是标准模块模板，包含：

| 文件 | 作用 | 关键模式 |
|------|------|----------|
| `index.vue` | 列表页 | `useNaivePaginatedTable` + `useTableOperate`，TSX 列渲染，`NDataTable` + `TableHeaderOperation` |
| `modules/user-search.vue` | 搜索 | `defineModel('model')` 双向绑定，`NCard + NCollapse + NForm + NGrid` 布局 |
| `modules/user-operate-drawer.vue` | 弹窗（多字段） | `defineModel('visible')`，`watch(visible, ...)` 初始化，字段多（≥6）用 `NDrawer` |
| `modules/{模块}-operate-modal.vue` | 弹窗（少字段） | 同上模式，字段少（<6）用 `NModal` |
