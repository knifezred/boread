package middleware

import (
	"github.com/gin-gonic/gin"

	"boread/internal/service"
	"boread/pkg/response"
)

// RequireButton 按钮级权限校验中间件
// 用法: protected.POST("/role", middleware.RequireButton(authSvc, "role:create"), handler.CreateRole)
// 流程: 从 Context 取 user_id (由 Auth 中间件注入)
//
//	查该用户拥有的所有按钮 code 集合
//	检查目标 code 是否在集合内, 否则 403
//
// 性能注意: 当前每次请求查 DB. 后续可加 sync.Map / Redis 缓存,
//
//	但生产前先确认真有性能问题再优化, 别过早优化.
func RequireButton(authSvc *service.AuthService, requiredCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user_id")
		if !ok {
			response.Error(c, 2001, "unauthorized")
			c.Abort()
			return
		}
		userID, ok := v.(uint64)
		if !ok {
			response.Error(c, 2001, "invalid user_id type")
			c.Abort()
			return
		}

		codes, err := authSvc.GetButtons(c.Request.Context(), userID)
		if err != nil {
			response.Error(c, 5001, "failed to load permissions")
			c.Abort()
			return
		}

		for _, code := range codes {
			if code == requiredCode {
				c.Next()
				return
			}
		}

		response.Error(c, 2005, "permission denied: "+requiredCode)
		c.Abort()
	}
}
