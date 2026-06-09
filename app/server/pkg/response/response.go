package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// errorCodeHTTPStatus 根据业务错误码范围映射 HTTP 状态码
// 分段规则 (与 code.go 保持一致):
//   0         → 200 (成功)
//   1000-1999 → 400 (参数/文件校验)
//   2000-2999 → 401 (认证授权)
//   3000-3999 → 409 (业务资源冲突)
//   4000-4999 → 400 (搜索/分页参数)
//   5000-5999 → 500 (系统/服务器内部错误)
//   其他      → 400 (默认)
func errorCodeHTTPStatus(code int) int {
	switch {
	case code == 0:
		return http.StatusOK
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	case code >= 2000 && code < 3000:
		return http.StatusUnauthorized
	case code >= 3000 && code < 4000:
		return http.StatusConflict
	case code >= 4000 && code < 5000:
		return http.StatusBadRequest
	case code >= 5000:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(errorCodeHTTPStatus(code), Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(errorCodeHTTPStatus(code), Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
