# Boread 后端 AI 代码生成规范 (Go + Gin)

> 目标：让 AI 生成符合 Boread 项目规范的后端代码
> 基于项目现有代码模式逆向提取，禁止自行发挥

---

## 1. 项目结构

```
app/server/
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

```go
package model

// SysBookTag 书的标签关联 (book_tag)
type SysBookTag struct {
    BaseModel
    TagName   string `gorm:"column:tag_name;size:64;not null" json:"tagName"`
    TagColor  string `gorm:"column:tag_color;size:32" json:"tagColor"`
    SortOrder int    `gorm:"column:sort_order;default:0" json:"sortOrder"`
}

func (SysBookTag) TableName() string { return "book_tag" }
```

**规则**:
- 必须嵌入 `model.BaseModel`（自动获得 ID/CreateBy/CreateTime/UpdateBy/UpdateTime/DeletedAt）
- `gorm:"column:xxx"` 标签必须与 SQL 列名一致
- `json:"xxx"` 用驼峰命名
- 必须有 `TableName()` 方法
- 枚举放在 model 包常量中

### 4.2 DTO

```go
package dto

// BookTagRequest 创建/更新标签
type BookTagRequest struct {
    TagName   string             `json:"tagName" binding:"required,max=64"`
    TagColor  string             `json:"tagColor"`
    SortOrder int                `json:"sortOrder"`
    Status    model.EnableStatus `json:"status"`
}

// BookTagSearch 标签搜索
type BookTagSearch struct {
    PageRequest
    TagName string             `json:"tagName"`
    Status  model.EnableStatus `json:"status"`
}
```

**规则**:
- 请求体用 `json:"xxx"` 标签，**禁止用 `form:"xxx"`**
- **响应结构体严禁使用 `,omitempty`**（零值被省略后前端 `undefined` 引发 TypeError）
- 搜索结构体嵌入 `PageRequest`（自动获得 current/size/keyword + Normalize/Offset 方法）
- Binding 约束放入 request 结构体（`binding:"required"` / `binding:"required,max=64"`）

### 4.3 Repository

```go
package repository

import (
    "gorm.io/gorm"
)

type SysDeptRepository struct {
    db *gorm.DB
}

func NewSysDeptRepository(db *gorm.DB) *SysDeptRepository {
    return &SysDeptRepository{db: db}
}

// GetByID 单查
func (r *SysDeptRepository) GetByID(ctx context.Context, id uint64) (*model.SysDept, error) {
    var m model.SysDept
    if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
        return nil, err
    }
    return &m, nil
}

// Page 分页查询 (带数据权限)
func (r *SysDeptRepository) Page(ctx context.Context, query *gorm.DB, size, offset int) (rows []T, total int64, err error) {
    err = query.WithContext(ctx).Model(&T{}).Count(&total).Error
    if err != nil || total == 0 {
        return nil, total, err
    }
    err = query.WithContext(ctx).Offset(offset).Limit(size).Find(&rows).Error
    return
}
```

**规则**:
- 每个 Repository 只操作一个 Model
- 方法签名：`(ctx context.Context, 参数...) → (结果, error)`
- 所有查询必须 `r.db.WithContext(ctx)`
- 分页方法的 `query *gorm.DB` 由 Service 注入（含数据权限 scope）
- 不写重复查询（比如同一个条件查了两次）

### 4.4 Service

```go
package service

import (
    "boread/internal/code"
    "boread/internal/dto"
    "boread/internal/model"
    "boread/internal/repository"
    "boread/internal/scope"
)

type DeptService struct {
    repo *repository.SysDeptRepository
}

func NewDeptService(repo *repository.SysDeptRepository) *DeptService {
    return &DeptService{repo: repo}
}

// Create 创建部门
func (s *DeptService) Create(ctx context.Context, req *dto.DeptRequest, userID uint64) (*model.SysDept, error) {
    // 1. 业务校验（返回 sentinel error）
    if _, err := s.repo.GetByCode(ctx, req.DeptCode); err == nil {
        return nil, code.ErrDeptCodeExists
    }
    // 2. 构造 Model
    m := &model.SysDept{...}
    m.CreateBy = &userID
    m.UpdateBy = &userID
    // 3. 入库
    if err := s.repo.Create(ctx, m); err != nil {
        return nil, err
    }
    return m, nil
}

// Page 分页（数据权限流程）
func (s *DeptService) Page(ctx context.Context, search *dto.DeptSearch, userID uint64) (*dto.PageResponse, error) {
    // 1. 计算数据权限上下文
    scopeCtx, err := scope.BuildDataScope(ctx, s.repo.db, userID)
    if err != nil {
        return nil, err
    }
    // 2. 构造基础查询，挂载数据权限 Scope
    query := s.repo.db.Model(&model.SysDept{}).Scopes(scope.ApplyDataScope(scopeCtx))
    // 3. 追加搜索条件
    if search.Keyword != "" {
        query = query.Where("dept_name LIKE ?", "%"+search.Keyword+"%")
    }
    // 4. 调用 Repository 分页
    search.Normalize()
    rows, total, err := s.repo.Page(ctx, query, search.Size, search.Offset())
    if err != nil {
        return nil, err
    }
    return dto.NewPageResponse(rows, total, search), nil
}
```

**规则**:
- Service 可以调用 1 个或多个 Repository
- 业务校验返回 sentinel error（`code.ErrXxx`）
- 写操作带 `userID` 参数，填充 `CreateBy/UpdateBy`
- 数据权限的 `BuildDataScope` + `ApplyDataScope` 在 Service 层组装

### 4.5 事务

涉及多个 Repository 写操作的 Service 方法，必须使用 `db.Transaction`：

```go
// 规约：Service 持有 *gorm.DB 引用，用 s.db.Transaction 包裹多步操作
type BookReaderService struct {
    db            *gorm.DB
    progressRepo  *repository.ReaderReadProgressRepository
    bookshelfRepo *repository.ReaderBookshelfRepository
}

func (s *BookReaderService) ReportProgress(ctx context.Context, userID uint64, req *dto.ReadProgressRequest) error {
    // 事务内调用多个 Repository
    return s.db.Transaction(func(tx *gorm.DB) error {
        if err := s.progressRepo.WithTx(tx).Upsert(ctx, m); err != nil {
            return err   // 自动回滚
        }
        if err := s.bookshelfRepo.WithTx(tx).UpdateTime(ctx, userID, bookID, now); err != nil {
            return err   // 自动回滚
        }
        return nil       // 自动提交
    })
}
```

**规则**:
- Repository 提供 `WithTx(tx *gorm.DB) *XxxRepository` 方法，复用同一事务
- `db.Transaction` 内的 return error 自动回滚，return nil 自动提交
- 只有读操作的 Service 方法不需要事务

### 4.6 Handler

```go
package v1

import (
    "github.com/gin-gonic/gin"
    "boread/internal/code"
    "boread/internal/dto"
    "boread/internal/service"
    "boread/pkg/response"
    "boread/pkg/utils"
)

type DeptHandler struct {
    svc *service.DeptService
}

func NewDeptHandler(svc *service.DeptService) *DeptHandler {
    return &DeptHandler{svc: svc}
}

// Create 创建部门
// @Summary   新增部门
// @Tags      dept
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.DeptRequest  true  "部门"
// @Success  200  {object}  response.Response{data=model.SysDept}
// @Router   /api/manage/dept [post]
func (h *DeptHandler) Create(c *gin.Context) {
    var req dto.DeptRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, code.ParamInvalid, err.Error())
        return
    }
    m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
    if err != nil {
        response.Error(c, code.MapServiceError(err), err.Error())
        return
    }
    response.Success(c, m)
}
```

**规则**:
- 必须有 Swagger 注释（@Summary/@Tags/@Security/@Param/@Success/@Router）
- 路径参数用 `utils.ParseUint64Param(c, "id")`
- 当前用户 id 用 `utils.GetUserID(c)`
- 所有错误统一用 `response.Error(c, httpCode, msg)` 返回
- sentinel error 通过 `code.MapServiceError(err)` 映射为业务错误码

### 4.7 Router — 依赖注入 + 路由注册

**完整装配链（新增一个模块时按此三步走）**：

```go
// 第1步：在 SetupRouter 的 Repository 区域新增
xxxRepo := repository.NewXxxRepository(db)

// 第2步：在 Service 区域新增
xxxSvc := service.NewXxxService(xxxRepo)

// 第3步：在 Handler 区域新增
xxxHandler := v1.NewXxxHandler(xxxSvc)

// 第4步：在路由分组中注册
manage.POST("/xxx", middleware.RequireButton(authSvc, "xxx:create"), xxxHandler.Create)
manage.GET("/xxx/:id", xxxHandler.GetByID)
manage.POST("/xxx/page", xxxHandler.Page)
manage.PUT("/xxx/:id", middleware.RequireButton(authSvc, "xxx:update"), xxxHandler.Update)
manage.DELETE("/xxx/:id", middleware.RequireButton(authSvc, "xxx:delete"), xxxHandler.Delete)
```

**路由分组约定**（见 [`router.go`](file:///c:/Users/zhang/repos/boread/app/server/internal/router/router.go)）：
- `/api/auth/*` — 公开或登录态（无需 manage）
- `/api/manage/*` — 后台管理，写操作需 `middleware.RequireButton(authSvc, "模块:操作")`
- `/api/book/*` — 电子书管理子分组，同样写操作需按钮权限
- `/api/book/note|review|comment/...` — 读者端社交功能
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

项目内置 5 种数据范围，由 `scope` 包统一处理（见 [`scope/data_scope.go`](file:///c:/Users/zhang/repos/boread/app/server/internal/scope/data_scope.go)）：

| 范围 | 代码 | 含义 |
|------|------|------|
| 全部 | `DataScopeAll` | 超管，不看权限 |
| 自定义 | `DataScopeCustom` | 特定部门集合 |
| 本部门 | `DataScopeDept` | 所在部门 |
| 部门及以下 | `DataScopeDeptAndSub` | 部门 + 所有子部门 |
| 仅本人 | `DataScopeSelf` | 只看自己创建 |

**完整用法流程**：

```go
// Service 层 —— 组装 scope 并传递给 Repository
func (s *DeptService) Page(ctx context.Context, search *dto.DeptSearch, userID uint64) (*dto.PageResponse, error) {
    // 第1步：调用 BuildDataScope（查询角色 → 计算数据范围 → 返回上下文）
    scopeCtx, err := scope.BuildDataScope(ctx, s.db, userID)
    if err != nil {
        return nil, err
    }
    // 第2步：用 ApplyDataScope 挂载到 GORM 查询链
    query := s.db.Model(&model.SysBook{}).
        Scopes(scope.ApplyDataScope(scopeCtx))
    // 第3步：追加其他查询条件，然后传给 Repository 的 Page
    search.Normalize()
    rows, total, err := s.repo.Page(ctx, query, search.Size, search.Offset())
    return dto.NewPageResponse(rows, total, search), nil
}
```

```go
// Repository 层 —— 不关心数据权限，只执行查询
func (r *BookRepository) Page(ctx context.Context, query *gorm.DB, size, offset int) ([]model.SysBook, int64, error) {
    var rows []model.SysBook
    var total int64
    err := query.WithContext(ctx).Model(&model.SysBook{}).Count(&total).Error
    if err != nil || total == 0 {
        return nil, total, err
    }
    err = query.WithContext(ctx).Offset(offset).Limit(size).Find(&rows).Error
    return rows, total, err
}
```

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

- **必须遵守** [`project-development.md`](file:///c:/Users/zhang/repos/boread/docs/project-development.md) 中的通用注意事项
- **必须检查** `code/code.go` 中是否已有现成的 `MapServiceError` 映射，新 error 必须加 case
- **禁止** 引入当前 `go.mod` 中没有的依赖
- **禁止** 在 Handler 中写业务逻辑或直接调用 Repository
- **禁止** 在 response 结构体中使用 `,omitempty`
- **必须** 给所有 Handler 方法写完整的 Swagger 注释
- **必须** 给所有写操作路由加 `middleware.RequireButton`
- **多 Repository 写操作必须用 `db.Transaction` 包裹**
- **新建文件** 必须符合本章第 4 节的分层模板框架
