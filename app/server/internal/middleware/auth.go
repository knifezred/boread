package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"boread/internal/code"
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
