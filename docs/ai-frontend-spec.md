# Boread 前端 AI 代码生成规范 (Vue 3 + TypeScript)

> 目标：让 AI 生成符合 Boread 项目规范的前端代码
> 基于项目现有代码模式逆向提取，禁止自行发挥

---

## 1. 项目结构

```
app/web/
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
└── packages/                        # 工作空间包（utils/axios/hooks/materials...）
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

## 4. 页面模板 — index.vue

### 4.1 组合式 API 约定

```vue
<script setup lang="tsx">
import { ref } from 'vue';
import { NButton, NPopconfirm, NTag } from 'naive-ui';
import { useBoolean } from '@sa/hooks';
import { enableStatusRecord } from '@/constants/business';
import { fetchGetDeptList, fetchDeleteDept } from '@/service/api';
import { useAppStore } from '@/store/modules/app';
import { defaultTransform, useNaivePaginatedTable, useTableOperate } from '@/hooks/common/table';
import { $t } from '@/locales';
import DeptOperateModal, { type OperateType } from './modules/dept-operate-modal.vue';

const appStore = useAppStore();

// 搜索参数
const searchParams = ref<Api.SystemManage.DeptSearchParams>({
  current: 1,
  size: 10,
  deptName: null,
  status: null,
});

// 表格
const { columns, columnChecks, data, loading, pagination, getData, getDataByPage } = useNaivePaginatedTable({
  api: () => fetchGetDeptList(searchParams.value),
  onPaginationParamsChange: params => {
    searchParams.value.current = params.page;
    searchParams.value.size = params.pageSize;
  },
  transform: response => defaultTransform(response),
  columns: () => [
    // 列定义...
  ],
});

// 增删改操作
const { checkedRowKeys, onBatchDeleted, onDeleted } = useTableOperate(data, 'id', getData);
</script>
```

### 4.2 列渲染规则

- 状态字段用 `<NTag>` 配合 `enableStatusRecord` 常量
- JSX 列用 `() => (...)` 箭头函数
- 操作列用 `NButton` + `NPopconfirm`（delete 确认）
- 索引列：`render: (_, index) => index + 1`

```tsx
{
  key: 'status',
  title: $t('page.admin.system.dept.status'),
  align: 'center',
  width: 80,
  render: (row: Api.SystemManage.Dept) => {
    const tagMap: Record<Api.Common.EnableStatus, NaiveUI.ThemeColor> = {
      '1': 'success',
      '2': 'warning'
    };
    const label = $t(enableStatusRecord[row.status]);
    return <NTag type={tagMap[row.status]}>{label}</NTag>;
  }
}
```

### 4.3 操作列

```tsx
{
  key: 'operate',
  title: $t('common.operate'),
  align: 'center',
  width: 230,
  render: (row: Api.SystemManage.Dept) => (
    <div class="flex-center justify-end gap-8px">
      <NButton type="primary" ghost size="small" onClick={() => handleEdit(row)}>
        {$t('common.edit')}
      </NButton>
      <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
        {{
          default: () => $t('common.confirmDelete'),
          trigger: () => (
            <NButton type="error" ghost size="small">
              {$t('common.delete')}
            </NButton>
          )
        }}
      </NPopconfirm>
    </div>
  )
}
```

### 4.4 Template 模板

```vue
<template>
  <div class="flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <SearchComponent v-model:model="searchParams" @search="getDataByPage" />
    <NCard :title="$t('page.admin.system.dept.title')" :bordered="false" size="small" class="card-wrapper sm:flex-1-hidden">
      <template #header-extra>
        <TableHeaderOperation v-model:columns="columnChecks" :disabled-delete="checkedRowKeys.length === 0"
          :loading="loading" @add="handleAdd" @delete="handleBatchDelete" @refresh="getData" />
      </template>
      <NDataTable v-model:checked-row-keys="checkedRowKeys" :columns="columns" :data="data" size="small"
        :flex-height="!appStore.isMobile" :scroll-x="1088" :loading="loading" :row-key="row => row.id" remote
        :pagination="pagination" class="sm:h-full" />
      <OperateModal v-model:visible="visible" :operate-type="operateType" :row-data="editingData"
        @submitted="getDataByPage" />
    </NCard>
  </div>
</template>
```

---

## 5. 弹窗组件模板 — operate-modal.vue

```vue
<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { NButton, NForm, NFormItem, NInput, NInputNumber, NModal, NScrollbar, NRadioGroup, NRadio, NSpace } from 'naive-ui';
import { enableStatusOptions } from '@/constants/business';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { $t } from '@/locales';
import { fetchCreateDept, fetchUpdateDept } from '@/service/api';

defineOptions({ name: 'DeptOperateModal' });

export type OperateType = NaiveUI.TableOperateType | 'addChild';

interface Props {
  operateType: OperateType;
  rowData?: Api.SystemManage.Dept | null;
}

const props = defineProps<Props>();
interface Emits { (e: 'submitted'): void }
const emit = defineEmits<Emits>();
const visible = defineModel<boolean>('visible', { default: false });

const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<string, string> = {
    add: $t('page.admin.system.dept.addDept'),
    edit: $t('page.admin.system.dept.editDept')
  };
  return titles[props.operateType];
});

type Model = {
  id: number;
  parentId: number;
  deptName: string;
  deptCode: string;
  leader: string;
  sortOrder: number;
  status: Api.Common.EnableStatus;
};

const model = ref<Model>(createDefaultModel());

function createDefaultModel(): Model {
  return {
    id: 0,
    parentId: 0,
    deptName: '',
    deptCode: '',
    leader: '',
    sortOrder: 0,
    status: '1'
  };
}

type RuleKey = Extract<keyof Model, 'deptName' | 'deptCode' | 'status'>;
const rules: Record<RuleKey, App.Global.FormRule> = {
  deptName: defaultRequiredRule,
  deptCode: defaultRequiredRule,
  status: defaultRequiredRule
};

function handleInitModel() {
  model.value = createDefaultModel();
  if (!props.rowData) return;
  if (props.operateType === 'addChild') {
    const { id } = props.rowData;
    Object.assign(model.value, { parentId: id });
  }
  if (props.operateType === 'edit') {
    Object.assign(model.value, props.rowData);
  }
}

function closeModal() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();
  const { id, ...restData } = model.value;
  if (props.operateType === 'edit') {
    await fetchUpdateDept(id, restData);
  } else {
    await fetchCreateDept(restData);
  }
  window.$message?.success($t('common.updateSuccess'));
  closeModal();
  emit('submitted');
}

watch(visible, () => {
  if (visible.value) {
    handleInitModel();
    restoreValidation();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-600px">
    <NScrollbar class="h-360px pr-20px">
      <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" :label-width="100">
        <NFormItem :label="$t('page.admin.system.dept.deptName')" path="deptName">
          <NInput v-model:value="model.deptName" :placeholder="$t('page.admin.system.dept.form.deptName')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dept.deptCode')" path="deptCode">
          <NInput v-model:value="model.deptCode" :placeholder="$t('page.admin.system.dept.form.deptCode')" />
        </NFormItem>
        <NFormItem :label="$t('page.admin.system.dept.status')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio v-for="item in enableStatusOptions" :key="item.value" :value="item.value" :label="$t(item.label)" />
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </NScrollbar>
    <template #footer>
      <NSpace justify="end" :size="16">
        <NButton @click="closeModal">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
      </NSpace>
    </template>
  </NModal>
</template>
```

**规则**:
- `visible` 必须用 `defineModel<boolean>('visible', { default: false })`
- `submit` emit 名为 `submitted`
- 窗口默认宽度 `w-600px`
- `defineOptions({ name: 'Xxx' })` 必须写
- export `type OperateType` 供父组件使用
- `handleInitModel` 中根据 `operateType` 决定初始值来源
- `handleSubmit` 区分 add/edit 调用不同 API
- `watch(visible)` 在弹窗打开时初始化模型并重置表单校验状态

---

## 6. 搜索组件模板

```vue
<script setup lang="ts">
import { toRaw } from 'vue';
import { jsonClone } from '@sa/utils';
import { enableStatusOptions } from '@/constants/business';
import { translateOptions } from '@/utils/common';
import { $t } from '@/locales';

defineOptions({ name: 'RoleSearch' });

interface Emits { (e: 'search'): void }
const emit = defineEmits<Emits>();
const model = defineModel<Api.SystemManage.RoleSearchParams>('model', { required: true });

const defaultModel = jsonClone(toRaw(model.value));
function resetModel() { Object.assign(model.value, defaultModel); }
function search() { emit('search'); }
</script>

<template>
  <NCard :bordered="false" size="small" class="card-wrapper">
    <NCollapse :default-expanded-names="['role-search']">
      <NCollapseItem :title="$t('common.search')" name="role-search">
        <NForm :model="model" label-placement="left" :label-width="80">
          <NGrid responsive="screen" item-responsive>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.admin.system.role.roleName')" path="roleName" class="pr-24px">
              <NInput v-model:value="model.roleName" :placeholder="$t('page.admin.system.role.form.roleName')" />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6">
              <NSpace class="w-full" justify="end">
                <NButton @click="resetModel">{{ $t('common.reset') }}</NButton>
                <NButton type="primary" ghost @click="search">{{ $t('common.search') }}</NButton>
              </NSpace>
            </NFormItemGi>
          </NGrid>
        </NForm>
      </NCollapseItem>
    </NCollapse>
  </NCard>
</template>
```

**规则**:
- 使用 `defineModel<T>('model', { required: true })` 双向绑定搜索参数
- 搜索按钮 emit `search` 事件
- 重置调用 `jsonClone(toRaw(model.value))` 保存深拷贝副本

---

## 7. API 文件规范

```ts
// src/service/api/system-manage.ts
import { request } from '../request';

/** 部门分页 */
export function fetchGetDeptList(params?: Api.SystemManage.DeptSearchParams) {
  return request<Api.SystemManage.DeptList>({
    url: '/manage/dept/page',
    method: 'post',
    data: params,
  });
}

/** 新增部门 */
export function fetchCreateDept(data: Partial<Api.SystemManage.Dept>) {
  return request<Api.SystemManage.Dept>({
    url: '/manage/dept',
    method: 'post',
    data,
  });
}

/** 删除部门 */
export function fetchDeleteDept(id: string | number) {
  return request<null>({
    url: `/manage/dept/${id}`,
    method: 'delete',
  });
}
```

**规则**:
- 函数命名前缀 `fetch`（如 `fetchGetXxx`, `fetchCreateXxx`, `fetchDeleteXxx`）
- 必须使用 `request<T>({ url, method, data })` 调用
- **禁止任何 `any` 类型定义**
- 必须导出（`export function`）
- API 函数统一在 `src/service/api/index.ts` 中 re-export（`export * from './system-manage'`）
- url 路径只需写 `/manage/xxx`（baseURL 为 `/api`）

---

## 8. TypeScript 类型定义规范

```ts
// src/typings/api/system-manage.d.ts
declare namespace Api {
  namespace SystemManage {
    /** 部门 */
    type Dept = Common.CommonRecord<{
      parentId: number;
      deptName: string;
      deptCode: string;
      leader: string | null;
      sortOrder: number;
      children?: Dept[] | null;
    }>;

    /** 部门搜索参数 */
    type DeptSearchParams = CommonType.RecordNullable<
      Pick<Api.SystemManage.Dept, 'deptName' | 'status'> & CommonSearchParams
    >;

    /** 部门列表 */
    type DeptList = Common.PaginatingQueryRecord<Dept>;
  }
}
```

**规则**:
- 统一在 `declare namespace Api { namespace Xxx { ... } }` 下
- 继承 `Common.CommonRecord` 获得基础字段（`id`/`createBy`/`createTime`/`updateBy`/`updateTime`/`status`）
- 搜索参数用 `CommonType.RecordNullable<Pick<...> & CommonSearchParams>`
- 分页列表用 `Common.PaginatingQueryRecord<T>`
- **`.d.ts` 用 `declare namespace` 全局生效，无需在 `index.ts` 中 export**

---

## 9. i18n 规范

新增页面/字段需修改 **3 个文件**：

```ts
// 第1步：src/locales/langs/zh-cn.ts — 往 Schema 对象加字段（I18nKey 类型自动推导）
const locale: App.I18n.Schema = {
  common: {
    enable: '启用',
    disable: '禁用',
    // ...
  },
  page: {
    admin: {
      system: {
        dept: {
          title: '部门管理',
          deptName: '部门名称',
          deptCode: '部门编码',
          status: '状态',
        }
      }
    }
  }
};
```

```ts
// 第2步：src/locales/langs/en-us.ts — 同样结构，英文翻译
const locale: App.I18n.Schema = {
  page: {
    admin: {
      system: {
        dept: {
          title: 'Department Management',
          deptName: 'Department Name',
          deptCode: 'Department Code',
          status: 'Status',
        }
      }
    }
  }
};
```

**重要**：**不需要** 手写 `I18nKey` union type。`I18nKey` 由 `GetI18nKey<Schema>` 递归类型自动推导（见 [`app.d.ts`](file:///c:/Users/zhang/repos/boread/app/web/src/typings/app.d.ts#L1206-L1215)）。只需往 `Schema` 对象里加字段，类型系统自动补全。

---

## 10. 路由

**规则**: AI **无需手动处理**路由文件。`elegant-router` 会自动从 `views/` 目录结构生成 4 个文件（`routes.ts` + `imports.ts` + `transform.ts` + `elegant-router.d.ts`）。

---

## 11. 状态常量

```ts
// src/constants/business.ts
export const enableStatusRecord: Record<Api.Common.EnableStatus, App.I18n.I18nKey> = {
  '1': 'common.enable',
  '2': 'common.disable'
};
export const enableStatusOptions = transformRecordToOption(enableStatusRecord);
```

---

## 12. 前端错误处理

后端统一返回 `{ code, message, data }` 结构，前端 `request` 层自动处理（见 [`request/index.ts`](file:///c:/Users/zhang/repos/boread/app/web/src/service/request/index.ts)）：

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

## 13. AI 编码约束

- **必须遵守** [`project-development.md`](file:///c:/Users/zhang/repos/boread/docs/project-development.md) 中的通用注意事项
- **必须** 使用 `$t()` 引用所有展示文本
- **必须** 用 `request<T>()` 调用 API，必须带泛型参数（返回 `{ data, error }` 结构）
- **禁止** 使用 `any` 类型
- **必须** 给 `.vue` 组件写 `defineOptions({ name: 'Xxx' })`
- **必须** 用 `defineModel<boolean>('visible')` 控制弹窗/抽屉显隐
- **禁止** 直接修改路由文件（elegant-router 自动生成，新建`.vue`需要重启 dev server）
- **新建业务模块** 必须遵循 `三文件结构`（index + search + operate-modal）
- **新建 API 函数** 必须在 `src/service/api/index.ts` 中 re-export
- **列渲染** 用 TSX（`lang="tsx"`），页面逻辑用 TS（`lang="ts"`）
- **弹窗组件** 用 `<NModal preset="card" class="w-600px">`，抽屉用 `<NDrawer :width="520">`
- **`.d.ts` 文件用 `declare namespace`**，全局生效，不需要 export
- **i18n 新增字段**只改 `zh-cn.ts` + `en-us.ts` 里的 `Schema` 对象，不改 `app.d.ts` 类型
- **后端请求失败**统一用 `error.message` 显示，不直接显示原始错误
