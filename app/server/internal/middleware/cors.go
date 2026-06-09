package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 返回 CORS 中间件
// allowedOrigins: 允许的来源列表; 为空表示允许所有源 (*)
func Cors(allowedOrigins ...string) gin.HandlerFunc {
	allowAll := len(allowedOrigins) == 0 || len(allowedOrigins[0]) == 0

	return func(c *gin.Context) {
		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			origin := c.Request.Header.Get("Origin")
			for _, allowed := range allowedOrigins {
				if allowed == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
