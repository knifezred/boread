package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserID 从 Gin Context 取登录用户 id, 兼容多种存储类型
// 0 表示未登录或类型异常
func GetUserID(c *gin.Context) uint64 {
	v, ok := c.Get("user_id")
	if !ok {
		return 0
	}
	switch x := v.(type) {
	case uint64:
		return x
	case uint:
		return uint64(x)
	case int64:
		return uint64(x)
	case int:
		return uint64(x)
	case string:
		id, _ := strconv.ParseUint(x, 10, 64)
		return id
	}
	return 0
}

// GetUserName 从 Context 取登录用户名
func GetUserName(c *gin.Context) string {
	v, ok := c.Get("username")
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// ParseUint64Param 解析路径参数为 uint64
func ParseUint64Param(c *gin.Context, key string) (uint64, error) {
	return strconv.ParseUint(c.Param(key), 10, 64)
}
