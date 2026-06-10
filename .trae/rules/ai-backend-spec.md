# 后端 AI 代码生成规范 (Go + Gin)

---

## 1. 项目结构

```
server/
├── cmd/api/main.go           # 入口：配置加载 → 启动逻辑
├── internal/
│   ├── code/code.go           # 错误码常量 + sentinel error + MapServiceError
│   ├── handler/v1/            # Handler 层：参数校验 + 响应输出
│   ├── service/               # Service 层：业务逻辑 + 事务 + 错误码映射
│   ├── repository/            # Repository 层：GORM 数据操作
│   ├── model/                 # GORM Model 定义（嵌入 BaseModel）
│   ├── dto/                   # 请求/响应结构体（仅 json 标签）
│   ├── middleware/            # Auth / CORS / Logger / Permission
│   ├── router/router.go       # 路由注册 + 依赖注入（唯一 DI 点）
│   ├── scope/                 # 数据权限过滤
│   └── seed/                  # 种子数据
├── pkg/
│   ├── config/config.go       # YAML 配置加载
│   ├── jwt/jwt.go             # JWT 签发/验证
│   ├── logger/logger.go       # Zap 日志
│   ├── response/response.go   # 统一响应格式
│   └── utils/                 # 工具函数（ctx/getUserID/encoding）
├── configs/config.yaml        # 配置文件
└── go.mod                     # module boread
```

**关键规则**：
- **所有业务代码必须放在 `internal/` 下分层**（handler → service → repository → model）
- `pkg/` 只放**纯工具**：不 import `internal/` 下的任何包，无状态，无业务逻辑
- 新增包如果 import 了 `internal/` 下的东西，**禁止放 `pkg/`**，必须放 `internal/`

---

## 2. 代码分层与调用链

```
Handler (参数校验) → Service (业务逻辑) → Repository (GORM 操作)
                              ↕
                            Model / DTO
```

- **Handler** 只负责：`ShouldBindJSON` / `ShouldBindQuery` / `ParseUint64Param` + 调用 service + `response.Success/Error`
- **Service** 只负责：校验规则 + 组装数据权限 + 调用 Repository + 错误码映射
- **Repository** 只负责：GORM 查询 + 返回原生 Model
- **严禁** Handler 直接调用 Repository

---

## 3. 命名规范

| 元素 | 规范 | 示例 |
|------|------|------|
| 文件 | 全小写下划线 | `sys_dept.go`, `book_chapter.go` |
| 包名 | 全小写 | `package service` |
| 导出类型/函数 | PascalCase | `SysDept`, `Create`, `NewDeptService` |
| 非导出函数/字段 | camelCase | `getUserID`, `buildDeptTree` |
| 常量 | PascalCase | `SerialOngoing`, `StatusEnabled` |
| 错误码（sentinel） | `Err` 前缀 + PascalCase | `ErrUserNotFound`, `ErrDeptCodeExists` |

---

## 4. 后端分层代码规范

### 4.1 Model

**规则**:
- 必须嵌入 `model.BaseModel`（自动获得 ID/CreateBy/CreateTime/UpdateBy/UpdateTime/DeletedAt）
- `gorm:"column:xxx"` 标签必须与 SQL 列名一致
- `json:"xxx"` 用驼峰命名
- 必须有 `TableName()` 方法

### 4.2 DTO

**规则**:
- 请求体用 `json:"xxx"` 标签，**禁止用 `form:"xxx"`**
- **响应结构体严禁使用 `,omitempty`**（零值被省略后前端 `undefined` 引发 TypeError）
- 搜索结构体嵌入 `PageRequest`（自动获得 current/size/keyword + Normalize/Offset 方法）

### 4.3 Repository

**规则**:
- 每个 Repository 只操作一个 Model
- 方法签名：`(ctx context.Context, 参数...) → (结果, error)`
- 所有查询必须 `r.db.WithContext(ctx)`
- 分页方法的 `query *gorm.DB` 由 Service 注入（含数据权限 scope）
- 不写重复查询（比如同一个条件查了两次）

### 4.4 Service

**规则**:
- Service 可以调用 1 个或多个 Repository
- 业务校验返回 sentinel error（`code.ErrXxx`）
- 写操作带 `userID` 参数，填充 `CreateBy/UpdateBy`
- 数据权限的 `BuildDataScope` + `ApplyDataScope` 在 Service 层组装

### 4.5 事务

涉及多个 Repository 写操作的 Service 方法，必须使用 `db.Transaction`：

**规则**:
- Repository 提供 `WithTx(tx *gorm.DB) *XxxRepository` 方法，复用同一事务
- `db.Transaction` 内的 return error 自动回滚，return nil 自动提交
- 只有读操作的 Service 方法不需要事务

### 4.6 Handler

**规则**:
- 必须有 Swagger 注释（@Summary/@Tags/@Security/@Param/@Success/@Router）
- 路径参数用 `utils.ParseUint64Param(c, "id")`
- 当前用户 id 用 `utils.GetUserID(c)`
- 所有错误统一用 `response.Error(c, httpCode, msg)` 返回
- sentinel error 通过 `code.MapServiceError(err)` 映射为业务错误码

### 4.7 Router — 依赖注入 + 路由注册

**路由分组约定**（见 [`router.go`](app/server/internal/router/router.go)）：
- `/api/auth/*` — 公开或登录态（无需 manage）
- 读操作（GET）不加 `RequireButton`，仅需 `FlexAuth` 认证
- Service 需要 `*gorm.DB` 的，构造函数多传一个 `db` 参数（用于事务）

---

## 5. 错误处理

```go
// code.go — 统一错误码分段
// 1000-1999: 参数/文件校验
// 2000-2999: 认证授权
// 3000-3999: 业务资源/冲突
// 4000-4999: 分页/搜索参数
// 5000-5999: 系统/服务器内部错误

// sentinel error 定义
var ErrDeptCodeExists = errors.New("部门编码已存在")
var ErrDeptHasChildren = errors.New("存在子部门, 不能删除")

// MapServiceError 映射
func MapServiceError(err error) int {
    switch {
    case errors.Is(err, ErrDeptCodeExists):
        return ResourceConflict
    case errors.Is(err, ErrDeptHasChildren):
        return ResourceProtected
    }
}
```

**规则**:
- sentinel error 定义在 `code/code.go` 的 `var` 块中
- `MapServiceError` 必须对新 error 添加 case
- Handler 层优先用 `code.MapServiceError(err)` 映射
- Handler 层可用 `errors.Is(err, code.ErrXxx)` 精准判断

---

## 6. 数据权限

项目内置 5 种数据范围，由 `scope` 包统一处理（见 [`scope/data_scope.go`](app/server/internal/scope/data_scope.go)）：

| 范围 | 代码 | 含义 |
|------|------|------|
| 全部 | `DataScopeAll` | 超管，不看权限 |
| 自定义 | `DataScopeCustom` | 特定部门集合 |
| 本部门 | `DataScopeDept` | 所在部门 |
| 部门及以下 | `DataScopeDeptAndSub` | 部门 + 所有子部门 |
| 仅本人 | `DataScopeSelf` | 只看自己创建 |


**模型要求**：
- 业务表必须有 `owner_id`（创建者）和 `dept_id`（所属部门）两列
- 否则不要套 `ApplyDataScope`

---

## 7. 依赖管理

- 新增依赖必须 `go get`，然后 `go mod tidy`
- 禁止直接修改 `go.mod`
- 常用依赖：`gin`, `gorm`, `golang-jwt`, `zap`, `crypto`, `golang.org/x/text`

---

## 8. AI 编码约束

- **必须检查** `code/code.go` 中是否已有现成的 `MapServiceError` 映射，新 error 必须加 case
- **禁止** 引入当前 `go.mod` 中没有的依赖
- **禁止** 在 Handler 中写业务逻辑或直接调用 Repository
- **禁止** 在 response 结构体中使用 `,omitempty`
- **必须** 给所有 Handler 方法写完整的 Swagger 注释
- **必须** 给所有写操作路由加 `middleware.RequireButton`
- **多 Repository 写操作必须用 `db.Transaction` 包裹**
- **新建文件** 必须符合本章第 4 节的分层模板框架
