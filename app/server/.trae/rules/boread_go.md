# Boread 项目AI开发规范

## 一、核心信息

- **项目名**: boread
- **技术栈**: Go 1.26+ / Gin / GORM / MySQL
- **架构**: Handler → Service → Repository 三层架构
- **AI角色**: 资深后端开发，严格遵循本规范生成代码

## 二、目录结构

```
cmd/api/main.go          # 入口
internal/
├── handler/v1/           # HTTP 处理层
├── service/              # 业务逻辑层
├── repository/           # 数据访问层
├── model/                # 数据模型
├── dto/                  # 请求/响应 (request.go, response.go)
├── router/               # 路由注册
└── middleware/           # 中间件
pkg/
├── config/               # 配置
├── jwt/                  # JWT
├── logger/               # 日志
├── response/             # 统一响应
└── utils/                # 工具
configs/config.yaml       # 配置文件
```

## 三、命名规范

| 类型 | 规范 | 示例 |
|:---|:---|:---|
| 目录 | 小写单数 | `handler`, `service` |
| 文件名 | snake_case | `user_service.go` |
| 包名 | 小写单词 | `package service` |
| 导出函数/类型 | PascalCase | `NewUserService`, `UserResponse` |
| 私有函数/变量 | camelCase | `hashPassword`, `userRepo` |
| DTO请求 | `{Action}{Resource}Request` | `CreateUserRequest` |
| DTO响应 | `{Resource}Response` | `UserResponse` |

## 四、分层职责

| 层 | 目录 | 职责 |
|:---|:---|:---|
| Handler | `internal/handler/v1/` | 参数校验、调用Service、返回响应 |
| Service | `internal/service/` | 业务逻辑、事务边界 |
| Repository | `internal/repository/` | GORM 数据库操作 |
| Model | `internal/model/` | 表结构定义 |
| DTO | `internal/dto/` | 请求/响应结构体 |

## 五、统一响应格式

```go
// pkg/response/response.go
type Response struct {
    Code    int         `json:"code"`    // 0=成功
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{})
func Error(c *gin.Context, code int, message string)
```

## 六、错误码

| 范围 | 含义 |
|:---|:---|
| 0 | 成功 |
| 1000-1999 | 参数错误 |
| 2000-2999 | 认证/授权错误 |
| 3000-3999 | 业务错误 |
| 5000-5999 | 系统错误 |

## 七、AI生成指令

```
根据开发规范，生成 {Resource} 模块（如 Product）。
要求：model → dto → repository → service → handler → router 注册
实现：Create, Get, List, Update, Delete
```

## 八、关键约束

1. 数据库操作必须使用 `db.WithContext(ctx)`
2. Handler 不包含业务逻辑，Service 不包含 SQL
3. 依赖注入在 `router.go` 中统一初始化
4. 禁止拼接 SQL，使用 GORM 参数化查询
5. 错误码使用规范中的范围