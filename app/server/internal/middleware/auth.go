package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"boread/internal/code"
	"boread/internal/repository"
	jwtPkg "boread/pkg/jwt"
	"boread/pkg/response"
)

// Auth JWT 解析中间件: 校验 token 并把 user_id/username 注入 Context
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, code.AuthFailed, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, code.AuthFailed, "invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwtPkg.ParseToken(parts[1])
		if err != nil {
			response.Error(c, code.TokenInvalid, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// FlexAuth 混合认证中间件: 优先 JWT, 其次绿联网关认证
// 当运行在 UGOS 系统网关之后, 网关会注入 Ugreen-User-* 请求头
func FlexAuth(userRepo *repository.SysUserRepository, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 优先尝试 JWT
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := jwtPkg.ParseToken(parts[1])
				if err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Next()
					return
				}
			}
		}

		// 2. 尝试绿联网关认证 (Ugreen-User-ID 由 UGOS 系统网关注入)
		ugreenUserID := c.GetHeader("Ugreen-User-ID")
		if ugreenUserID != "" {
			var user struct {
				ID       uint64
				UserName string
			}
			if err := db.WithContext(c.Request.Context()).
				Table("sys_user").
				Select("id, user_name").
				Where("ugreen_user_id = ?", ugreenUserID).
				First(&user).Error; err == nil {
				c.Set("user_id", user.ID)
				c.Set("username", user.UserName)
				c.Next()
				return
			}
		}

		// 3. 两者都失败
		response.Error(c, code.AuthFailed, "authorization required")
		c.Abort()
	}
}
