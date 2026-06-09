# Boread 小说阅读平台 — AI 辅助开发文档

> 开发者：单人 + AI
> 版本：1.0.0 | 更新：2026-05-27 | 状态：开发中

---

## 一、项目总览

### 1.1 这是什么

Boread 是一个小说阅读平台，含后台运营管理系统 + 前台读者服务。

### 1.2 技术栈速查

| 层 | 技术 | 关键版本 |
|:---|:---|:---|
| 后端语言 | Go | 1.26+ |
| HTTP 框架 | Gin | v1.12 |
| ORM | GORM | v1.31 |
| 数据库 | MySQL | 8.0.13+（函数索引） |
| 前端框架 | Vue 3 | 3.5+ |
| UI 组件 | NaiveUI | 2.44 |
| 状态管理 | Pinia | 3.0 |
| 路由 | Elegant Router / Vue Router | 5.x |
| HTTP 请求 | Axios (@sa/axios) | workspace |
| 表格 Hooks | @sa/hooks (useNaivePaginatedTable) | workspace |
| 认证 | JWT (golang-jwt v5) | — |
| 文档 | Swaggo (gin-swagger) | v1.16 |

### 1.3 项目结构

```
boread/
├── app/
│   ├── server/               # Go 后端
│   │   ├── cmd/api/main.go   # 入口
│   │   ├── internal/
│   │   │   ├── handler/v1/   # Handler 层（参数校验+响应）
│   │   │   ├── service/      # Service 层（业务逻辑）
│   │   │   ├── repository/   # Repository 层（GORM 数据操作）
│   │   │   ├── model/        # GORM Model 定义
│   │   │   ├── dto/          # 请求/响应结构体
│   │   │   ├── middleware/   # Auth/CORS/Logger/Permission
│   │   │   ├── router/       # 路由注册 + 依赖注入
│   │   │   ├── scope/        # 数据权限过滤
│   │   │   └── seed/         # 种子数据
│   │   ├── pkg/              # 公共包 (config/jwt/response/utils)
│   │   ├── configs/          # 配置文件
│   │   ├── docs/             # Swagger 文档
│   │   ├── go.mod
│   │   └── Makefile
│   ├── web/                  # Vue 3 前端
│   │   ├── src/
│   │   │   ├── views/        # 页面组件
│   │   │   ├── service/api/  # API 调用（3 文件：auth/route/system-manage）
│   │   │   ├── store/        # Pinia store（app/auth/route/tab/theme）
│   │   │   ├── hooks/        # 通用 Hooks（table/form/icon）
│   │   │   ├── typings/api/  # API 类型定义
│   │   │   ├── locales/      # i18n（zh-cn + en-us）
│   │   │   ├── router/       # 路由（elegant 自动生成）
│   │   │   └── components/   # 公共组件
│   │   └── package.json
│   └── sql/                  # 数据库建表脚本
│       ├── system-manage.sql # 后台管理表（已执行）
│       └── business.sql      # 小说业务表（待执行）
└── docs/
    └── project-development.md # 本文件
```

---

## 二、开发模式说明（单人 + AI）

### 2.1 工作流

```
1. AI 出代码 → 2. 人工 review → 3. 测试验证 → 4. 提交
```

- **AI 负责**：按本文档的分阶段任务清单，生成后端 Go 代码（Handler/Service/Repository/Model/DTO）、前端 Vue 页面（index.vue + modules/*.vue）、API 调用、类型定义、i18n 翻译、数据库 migration
- **人工负责**：启动开发服务器、执行建表 SQL、运行测试验证、git 提交
- **不需要**：需求评审会议、代码审查流程、团队协作工具、燃尽图

### 2.2 每次给 AI 的 Prompt 模板

```
任务: [功能名称]
后端: 参照 [参考文件] 实现 [Handler/Service/Repository/DTO]
前端: 参照 [参考页面] 实现 [页面/组件]
数据库: [SQL 变更说明]
注意:
- 所有 request 调用要有泛型返回类型
- 前端页面用 $t() 引用 i18n 文本，同步修改 app.d.ts + zh-cn.ts + en-us.ts
- 前端页面参考 manage/user 或 manage/dept 的写法
- 常规字段不要写注释
```

### 2.3 通用注意事项（喂给 AI 前确认）

- **后端新增路由** → 在 `router.go` 中注册，写操作加 `middleware.RequireButton`
- **后端新增 DTO** → `json:"xxx"` 标签，不用 `form:"xxx"`；**响应结构体严禁使用 `,omitempty`**（零值被省略后前端读到 `undefined` 引发 TypeError）
- **后端新增 Model** → 嵌入 `model.BaseModel`，无需手写 `id/createTime/updateTime`
- **后端分页结构** → 返回 `dto.PageResponse`，前端 `Common.PaginatingQueryRecord<T>`，统一使用POST接口
- **前端新增页面** → 无需处理 route 文件,路由会自动生成：`routes.ts` + `imports.ts` + `transform.ts` + `elegant-router.d.ts`
- **前端新增 i18n** → 更新 3 文件：`app.d.ts` + `zh-cn.ts` + `en-us.ts`
- **前端 API 函数** → 统一写在 `service/api/*-manage.ts`，用 `request<T>({ url, method, data })`，禁止任何any定义
- **前端表格** → 用 `useNaivePaginatedTable` + `useTableOperate`
- **前端弹窗组件** → `defineProps<Props>()` + `defineModel<boolean>('visible')` + `defineEmits`

---

## 三、通用 Prompt 片段

### 后端新增模块

```
参照 [已有的 xxx.go] 实现 [yyy.go]：
- Model: 嵌入 BaseModel，字段对应 business.sql 中的 [表名]
- DTO: CreateRequest / UpdateRequest / Search（嵌入 PageRequest）/ VO
- Repository: GetByID / Create / Update / Delete / Page
- Service: 业务校验 + 调用 Repository
- Handler: ShouldBindJSON + svc.xxx + response.Success/Error
- Router: 注册到 manage 分组，写操作加 middleware.RequireButton
```

### 前端新增页面

```
参照 manage/user 实现 manage/[模块名]：
- API: service/api/system-manage.ts 新增函数
- Typing: typings/api/system-manage.d.ts 新增类型
- i18n: typings/app.d.ts + zh-cn.ts + en-us.ts
- 页面: index.vue（useNaivePaginatedTable + useTableOperate）
- 搜索: modules/xxx-search.vue
- 弹窗: modules/xxx-operate-modal.vue（NModal preset="card" w-600px）
- 路由: 更新 elegant 4 文件
```

### 数据权限接入

```
在 Repository 层的 Page 方法中，添加：
query = scope.WithDataScope(query, "owner_id", "dept_id")
确保 Model 有 OwnerID + DeptID 字段
```

---

## 四、质量检查清单（提交前过一遍）

```
后端：
[  ] make lint 通过（go vet）
[  ] 所有 Handler 有 Swagger 注释
[  ] 写操作有 middleware.RequireButton
[  ] 分页接口返回 dto.PageResponse
[  ] 数据权限表加了 scope.WithDataScope

前端：
[  ] pnpm typecheck 无错误
[  ] pnpm lint 无错误
[  ] 新增页面路由 4 文件已更新
[  ] 新增 i18n 3 文件已同步
[  ] request<T> 调用了泛型
[  ] 组件名 PascalCase，defineOptions({ name: 'Xxx' })
[  ] visible 用 defineModel<boolean>('visible')

```