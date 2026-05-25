package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	jwtPkg "boread/pkg/jwt"
	"boread/pkg/response"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 2001, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 2001, "invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwtPkg.ParseToken(parts[1])
		if err != nil {
			response.Error(c, 2002, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
