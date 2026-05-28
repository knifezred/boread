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
- **后端新增 DTO** → `json:"xxx"` 标签，不用 `form:"xxx"`
- **后端新增 Model** → 嵌入 `model.BaseModel`，无需手写 `id/createTime/updateTime`
- **后端分页结构** → 返回 `dto.PageResponse`，前端 `Common.PaginatingQueryRecord<T>`
- **前端新增页面** → 更新 4 文件：`routes.ts` + `imports.ts` + `transform.ts` + `elegant-router.d.ts`
- **前端新增 i18n** → 更新 3 文件：`app.d.ts` + `zh-cn.ts` + `en-us.ts`
- **前端 API 函数** → 统一写在 `service/api/system-manage.ts`，用 `request<T>({ url, method, data })`
- **前端表格** → 用 `useNaivePaginatedTable` + `useTableOperate`
- **前端弹窗组件** → `defineProps<Props>()` + `defineModel<boolean>('visible')` + `defineEmits`

---

## 三、开发阶段总览

```
Phase 1: 权限管理模块  ✅ 已完成
  ├── 1.1 建库建表 + 种子数据
  ├── 1.2 认证系统（登录/JWT/用户信息）
  ├── 1.3 部门管理（树形 CRUD + 数据权限）
  ├── 1.4 角色管理（CRUD + 菜单/按钮授权）
  ├── 1.5 用户管理（CRUD + 角色分配 + 密码重置）
  ├── 1.6 菜单管理（树形 CRUD + 按钮配置）
  ├── 1.7 字典管理（分类 CRUD + 字典项 CRUD）
  └── 1.8 日志管理（登录日志 + 操作日志）

Phase 2: 读者与内容管理  ← 当前阶段
  ├── 2.1 读者扩展信息管理 (复用 sys_user + 角色区分)
  ├── 2.2 小说分类管理（树形）
  ├── 2.3 小说标签管理
  ├── 2.4 章节解析规则管理
  └── 2.5 内容过滤规则管理

Phase 3: 小说核心业务
  ├── 3.1 作品管理（多文件上传+聚合）
  ├── 3.2 章节解析引擎 + pread 读取
  ├── 3.3 阅读器 + 阅读进度
  └── 3.4 书架

Phase 4: 社交与统计
  ├── 4.1 笔记/划线
  ├── 4.2 书评 + 评分
  ├── 4.3 章节评论（楼中楼）
  └── 4.4 阅读统计（事件 + 日/周/月/年聚合）
```

---

## 四、Phase 1：权限管理模块（✅ 已完成）

> 以下所有功能已实现，不需要再开发。建表脚本 `system-manage.sql` 已执行，种子数据已写入。

### 4.1 已实现的功能清单

| 功能 | 后端 | 前端页面 | API 路径 |
|:---|:---|:---|:---|
| 登录/登出 | auth.go | login/index.vue | POST /api/auth/login |
| 用户信息 | auth.go | store/auth | GET /api/auth/userInfo |
| 用户菜单 | auth.go | store/route | GET /api/auth/menu |
| 用户按钮 | auth.go | store/auth | GET /api/auth/buttons |
| 部门管理 | dept.go | manage/dept | GET/POST/PUT/DELETE /api/manage/dept/* |
| 角色管理 | role.go | manage/role | GET/POST/PUT/DELETE /api/manage/role/* |
| 用户管理 | user.go | manage/user | GET/POST/PUT/DELETE /api/manage/user/* |
| 菜单管理 | menu.go | manage/menu | GET/POST/PUT/DELETE /api/manage/menu/* |
| 字典管理 | dict.go | manage/dict | GET/POST/PUT/DELETE /api/manage/dict/* |
| 日志管理 | log.go | manage/log | POST /api/manage/log/*/page |

### 4.2 代码规范参考（后续开发严格参照）

**后端 Handler 样板**:

```go
func (h *XxxHandler) Page(c *gin.Context) {
    var req dto.XxxSearch
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 1001, err.Error())
        return
    }
    resp, err := h.svc.Page(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, 5001, err.Error())
        return
    }
    response.Success(c, resp)
}
```

**后端 Service 样板**:

```go
func (s *XxxService) Page(ctx context.Context, req *dto.XxxSearch) (*dto.PageResponse, error) {
    req.Normalize()
    rows, total, err := s.repo.Page(ctx, req)
    if err != nil {
        return nil, err
    }
    return dto.NewPageResponse(rows, total, &req.PageRequest), nil
}
```

**前端页面结构样板**:

```
views/manage/xxx/
├── index.vue                    # 主页面（搜索 + 表格 + 弹窗）
└── modules/
    ├── xxx-search.vue           # 搜索组件
    ├── xxx-operate-modal.vue    # 新增/编辑弹窗（NModal）
    └── xxx-other.vue            # 其他子组件（按需）
```

**前端 API 调用样板**:

```typescript
/** xxx 分页 */
export function fetchGetXxxList(params?: Api.SystemManage.XxxSearchParams) {
  return request<Api.SystemManage.XxxList>({
    url: "/manage/xxx/page",
    method: "post",
    data: params,
  });
}
```

### 4.3 数据权限要点

`data_scope` 5 级，在 `repository` 层用 `scope.WithDataScope` 注入：

```go
func (r *XxxRepository) Page(ctx context.Context, req *dto.XxxSearch) ([]model.Xxx, int64, error) {
    query := r.db.WithContext(ctx).Model(&model.Xxx{})
    query = scope.WithDataScope(query, "owner_id", "dept_id")
    // ... 条件筛选 + 分页
}
```

---

## 五、Phase 2：读者与内容管理（← 当前开发阶段）

### 5.1 读者的认证与管理（跳过）

**设计**: 复用 `sys_user` + 角色管理。用户在 `sys_user` 中创建，通过角色分配 `reader` 角色获得读者身份。无独立认证体系，与管理员共用同一套登录/JWT。

**读者扩展信息**: `sys_user_profile` 表（可选，用户关联此表 = 拥有读者身份）：

```sql
sys_user_profile (
  user_id    PK → sys_user.id,
  nickname  覆盖 sys_user.nick_name,
  avatar    覆盖 sys_user.avatar,
  signature 个性签名
)
```

**后端**: 无需新增认证逻辑，`sys_user` 已有完整的登录/CRUD。只需新增 `sys_user_profile` 的简单 CRUD，管理员管理用户时附带设置。

**前端**:
- 在现有用户管理页面（manage/user）中扩展：新增/编辑用户弹窗增加「读者信息」折叠面板（nickname/avatar/signature）
- 角色选择下拉增加 `reader` 角色选项
- 前端 `store/auth` 中 `userInfo.roles` 已包含角色信息，判断是否含 `reader` 角色即可区分读者/管理员

**路由**: 无需注册新路由，复用 manage/user 的所有路由。

### 5.2 小说分类管理

**建表**: `book_category`（business.sql 已有）

```sql
-- 注意已加 ancestors 字段（上次评审加的）
```

**后端 CRUD**: book_category 的 Handler/Service/Repository + 树形查询（参考 sys_dept）。

**前端**: manage/book-category/index.vue + modules/*（树形表格，参考 manage/dept）

### 5.3 小说标签管理

**建表**: `book_tag`（business.sql 已有）

**后端 CRUD**: book_tag 的 Handler/Service/Repository。

**前端**: manage/book-tag/index.vue（简单表格，含 usage_count 热度排序列）

### 5.4 章节解析规则管理

**建表**: `book_parse_rule`（business.sql 已有）

**后端 CRUD**: book_parse_rule 的 Handler/Service/Repository。

**前端**: manage/book-parse-rule/index.vue（表格 + 正则 pattern 输入框）

### 5.5 内容过滤规则管理

**建表**: `content_filter_rule`（business.sql 已有）

**后端 CRUD**: content_filter_rule 的 Handler/Service/Repository。

**前端**: manage/content-filter-rule/index.vue（表格 + 关键词/正则配置）

---

## 六、Phase 3：小说核心业务

### 6.1 作品管理（核心）

**核心设计**: `book` + `book_file` + `book_chapter` 三层聚合：

```
book（作品，唯一）
  ├── title + author 作为聚合标识
  ├── primary_file_id → 默认使用的文件
  └── aggregate_status（单文件/聚合中/聚合完成）

book_file（物理文件，多个）
  ├── content_path / content_md5 / content_version
  ├── is_primary（是否主版本）
  └── file_status（待处理/处理中/成功/失败）

book_chapter（聚合章节索引）
  ├── UNIQUE(book_id, chapter_no)
  ├── file_id → 来源文件
  └── byte_offset / byte_length → pread 读取
```

**上传流程**:

```
上传 TXT/EPUB
  → book_upload 记录（异步）
  → 解析引擎提取 title + author
  → 查询 book（title+author 聚合）→ 不存在则新建
  → INSERT book_file + 写入 content_path
  → 解析章节 → INSERT book_chapter (ON DUPLICATE KEY UPDATE)
  → 更新 book.total_chapters + aggregate_status
```

**阅读器流程**:

```
请求章节内容
  → 取 book.primary_file_id → book_file.content_path
  → open(file_path) → pread(fd, byte_offset, byte_length)
  → 返回 bytes
```

**后端模块**:

| 模块 | 文件 | 说明 |
|:---|:---|:---|
| 作品 CRUD | handler/service/repository/book.go | 聚合管理 |
| 文件上传 | handler/service/repository/book_file.go | 上传+解析+md5 判重 |
| 章节索引 | handler/service/repository/book_chapter.go | 索引维护+pread 读取 |
| 上传任务 | handler/service/repository/book_upload.go | 异步任务追踪 |

**前端**: manage/book/ 完整页面（含文件列表弹窗、章节管理）

### 6.2 书架

**建表**: `reader_bookshelf`（business.sql 已有）

**接口**: 添加/移除/列表（分组+置顶+最后阅读排序）

**前端**: 读者个人书架页面

### 6.3 阅读进度

**建表**: `reader_read_progress`（business.sql 已有，含 file_id）

**接口**: 上报/获取/文件切换（按 file_id 重新映射 chapter_id）

**前端**: 阅读器组件（与书架联调）

---

## 七、Phase 4：社交与统计

### 7.1 笔记/划线

**建表**: `reader_note`（business.sql 已有，note_type=1/2/3 合并设计）

**接口**: CRUD，支持按 book_id 和 chapter_id 查询

### 7.2 书评 + 评分

**建表**: `book_review`（business.sql 已有，含 CHECK 约束 + owner_id）

**接口**: CRUD，评分后更新 `book.avg_rating` + `book.rating_count`

### 7.3 章节评论

**建表**: `chapter_comment`（business.sql 已有，含 owner_id + parent_id 楼中楼）

**接口**: CRUD，支持按 chapter_id 分页查询

### 7.4 阅读统计

**三表体系**:

```
reader_read_event（原子事件）
    → 客户端 30-60s 心跳上报
    → 服务端回填 word_count（从 book_chapter.word_count）

reader_read_stat_daily（按读者+日聚合）
    → 定时任务：INSERT ... ON DUPLICATE KEY UPDATE（累加）
    → 周/月/年 = SUM GROUP BY 此表

reader_read_stat_book（按读者+书+日聚合）
    → 个人页"我读过的书排行"
```

---

## 八、数据库建表执行计划

| 阶段 | 表 | 脚本位置 | 状态 |
|:---|:---|:---|:---|
| Phase 1 | sys_dept/sys_role/sys_user/sys_menu/sys_dict/sys_log 等全部管理表 | system-manage.sql | ✅ 已执行 |
| Phase 2 | sys_user_profile / book_category / book_tag / book_parse_rule / content_filter_rule | business.sql 上半段 | ⏳ 待执行 |
| Phase 3 | book / book_file / book_chapter / book_upload / reader_bookshelf / reader_read_progress | business.sql 中段 | ⏳ 待执行 |
| Phase 4 | reader_note / book_review / chapter_comment / reader_read_event / reader_read_stat_daily / reader_read_stat_book / book_character / book_character_chapter | business.sql 下半段 | ⏳ 待执行 |

**执行方式**：

```bash
# 分阶段执行，每次执行对应阶段的建表语句
mysql -u root -p boread < app/sql/business.sql
```

---

## 九、通用 Prompt 片段（每次开发前粘贴给 AI）

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

## 十、开发环境快速启动

```bash
# === 后端 ===
cd app/server
cp configs/config.example.yaml configs/config.yaml
# 编辑 config.yaml（数据库连接 + JWT Secret）
make run            # 开发模式，端口 8080
make swag           # 更新 API 文档
# 种子数据（首次）
go run ./cmd/api/ -seed

# === 前端 ===
cd app/web
pnpm install
pnpm dev            # 开发模式，端口 9527
pnpm gen-route      # 新增页面后生成路由
pnpm typecheck      # 类型检查

# === 数据库 ===
mysql -u root -p -e "CREATE DATABASE boread DEFAULT CHARSET utf8mb4;"
mysql -u root -p boread < app/sql/system-manage.sql  # Phase 1 已完成
mysql -u root -p boread < app/sql/business.sql        # 分阶段执行
```

---

## 十一、质量检查清单（提交前过一遍）

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

数据库：
[  ] 建表 SQL 已执行
[  ] 业务表含 deleted_at 软删
[  ] 唯一索引含 IFNULL(deleted_at, '1970-01-01')
```