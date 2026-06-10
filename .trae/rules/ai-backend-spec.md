# 后端 AI 代码生成规范 (Go + Gin)

---

## 1. 项目结构

```
server/
├── cmd/api/main.go           # 入口
├── internal/
│   ├── code/code.go           # 错误码常量 + sentinel error + MapServiceError
│   ├── handler/v1/            # 参数校验 + 响应输出
│   ├── service/               # 业务逻辑 + 事务 + 错误码映射
│   ├── repository/            # GORM 数据操作
│   ├── model/                 # GORM Model（嵌入 BaseModel）
│   ├── dto/                   # 请求/响应结构体（仅 json 标签）
│   ├── middleware/            # Auth / CORS / Logger / Permission
│   ├── router/router.go       # 路由注册 + 依赖注入（唯一 DI 点）
│   ├── scope/                 # 数据权限过滤
│   └── seed/                  # 种子数据
├── pkg/                       # 纯工具，不 import internal
│   ├── config/config.go
│   ├── jwt/jwt.go
│   ├── logger/logger.go
│   ├── response/response.go   # 统一响应格式
│   └── utils/                 # getUserID/ParseUint64Param 等
└── configs/config.yaml
```

**关键**：业务代码全放 `internal/` 分层；`pkg/` 只放纯工具，不 import `internal/`。

---

## 2. 调用链

```
Handler (参数校验) → Service (业务逻辑) → Repository (GORM 操作)
                              ↕
                            Model / DTO
```

- **Handler**：`ShouldBindJSON`/`ShouldBindQuery`/`ParseUint64Param` + 调 service + `response.Success/Error`
- **Service**：校验规则 + 数据权限 + 调 Repository + 错误码映射
- **Repository**：GORM 查询 + 返回原生 Model
- **严禁** Handler 直接调 Repository，严禁 Handler 写业务逻辑

---

## 3. 命名规范

| 元素 | 规范 | 示例 |
|------|------|------|
| 文件 | 全小写下划线 | `sys_dept.go`, `book_chapter.go` |
| 包名 | 全小写 | `package service` |
| 导出类型/函数 | PascalCase | `SysDept`, `Create`, `NewDeptService` |
| 非导出 | camelCase | `getUserID`, `buildDeptTree` |
| 常量 | PascalCase | `SerialOngoing`, `StatusEnabled` |
| sentinel error | `Err` 前缀 + PascalCase | `ErrUserNotFound`, `ErrDeptCodeExists` |

---

## 4. 分层规范

### 4.1 Model
- 嵌入 `model.BaseModel`（自动获得 ID/CreateBy/CreateTime/UpdateBy/UpdateTime/DeletedAt）
- `gorm:"column:xxx"` 与 SQL 列名一致，`json:"xxx"` 驼峰
- 必须有 `TableName()`

### 4.2 DTO
- 请求体 `json:"xxx"` 标签，**禁止用 `form:"xxx"`**
- **响应结构体严禁使用 `,omitempty`**（零值省略导致前端 TypeError）
- 搜索结构体嵌入 `PageRequest`（含 current/size/keyword + Normalize/Offset）

### 4.3 Repository
- 一个 Repository 只操作一个 Model
- 签名：`(ctx, 参数...) → (结果, error)`
- 所有查询 `r.db.WithContext(ctx)`
- 分页 `query *gorm.DB` 由 Service 注入（含数据权限 scope）
- Repository 提供 `WithTx(tx *gorm.DB) *XxxRepository` 用于事务复用

### 4.4 Service
- 可调 1+ 个 Repository
- 校验返回 sentinel error（`code.ErrXxx`）
- 写操作带 `userID`，填充 `CreateBy/UpdateBy`
- 数据权限 `BuildDataScope` + `ApplyDataScope` 在 Service 组装
- **多 Repository 写操作必须用 `db.Transaction` 包裹**

### 4.5 Handler
- 完整 Swagger 注释（@Summary/@Tags/@Security/@Param/@Success/@Router）
- 路径参数 `utils.ParseUint64Param(c, "id")`
- 当前用户 `utils.GetUserID(c)`
- 错误统一 `response.Error(c, httpCode, msg)` 返回
- sentinel error 通过 `code.MapServiceError(err)` 映射为业务错误码

### 4.6 Router
- `/api/auth/*` — 公开/登录态（无需 manage）
- 读操作（GET）不加 `RequireButton`，仅需 `FlexAuth`
- 写操作路由必须加 `middleware.RequireButton`

---

## 5. 错误处理

```go
// code.go — 统一错误码分段
// 1000-1999: 参数/文件 | 2000-2999: 认证授权 | 3000-3999: 业务资源
// 4000-4999: 分页/搜索 | 5000-5999: 系统内部

var ErrDeptCodeExists = errors.New("部门编码已存在")

func MapServiceError(err error) int {
    switch {
    case errors.Is(err, ErrDeptCodeExists):
        return ResourceConflict
    }
}
```

**规则**：sentinel error 定义在 `code/code.go` `var` 块中；新 error 必须加 `MapServiceError` case；Handler 层可用 `errors.Is(err, code.ErrXxx)` 精准判断。

---

## 6. 数据权限

| 范围 | 代码 | 含义 |
|------|------|------|
| 全部 | `DataScopeAll` | 超管 |
| 自定义 | `DataScopeCustom` | 特定部门集合 |
| 本部门 | `DataScopeDept` | 所在部门 |
| 部门及以下 | `DataScopeDeptAndSub` | 部门 + 子部门 |
| 仅本人 | `DataScopeSelf` | 只看自己 |

**业务表必须有 `owner_id` 和 `dept_id`** 才能用 `ApplyDataScope`。

---

## 7. 依赖管理

- 新增依赖 `go get` + `go mod tidy`，禁止直接改 `go.mod`
- 常用依赖：`gin`, `gorm`, `golang-jwt`, `zap`, `crypto`, `golang.org/x/text`

---

## 8. AI 编码约束

- **必须检查** `code/code.go` 是否有现成 `MapServiceError` 映射，新 error 必须加 case
- **禁止**引入 `go.mod` 中不存在的依赖
- **禁止**在 Handler 写业务逻辑或直接调 Repository
- **禁止**在 response 结构体用 `,omitempty`
- **必须**给所有 Handler 写完整 Swagger 注释
- **必须**给所有写操作路由加 `middleware.RequireButton`
- **多 Repository 写操作必须用 `db.Transaction`**
- **新建文件** 必须符合第 4 节分层模板
