package middleware

import (
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"boread/internal/model"
	"boread/internal/repository"
)

// OperationLog 操作日志中间件
// 记录 POST/PUT/DELETE 请求到 sys_operation_log 表，异步写入不阻塞响应
func OperationLog(logRepo *repository.SysLogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// 只记录写操作
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			c.Next()
			return
		}

		start := time.Now()

		// 缓存请求体（handler 读取前需要先复制）
		var bodyStr string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				bodyStr = sanitizeBody(string(bodyBytes))
				// 将原始 body 放回，供后续 handler 读取
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 执行 handler
		c.Next()

		// handler 完成后异步写日志
		costMs := uint32(time.Since(start).Milliseconds())
		statusCode := c.Writer.Status()
		url := c.Request.URL.Path
		clientIP := c.ClientIP()
		ua := c.Request.UserAgent()

		// 从 FlexAuth 中间件提取的用户信息
		userID, _ := c.Get("user_id")
		userName, _ := c.Get("username")

		// 解析 module / action / targetID
		module, action, targetID := parseRoute(url, method)

		uid := toUint64(userID)
		uname, _ := userName.(string)

		go func() {
			log := &model.SysOperationLog{
				UserID:        uid,
				UserName:      uname,
				Module:        module,
				Action:        action,
				TargetID:      targetID,
				RequestURL:    &url,
				RequestMethod: &method,
				RequestBody:   &bodyStr,
				ResponseCode:  &statusCode,
				ClientIP:      &clientIP,
				UserAgent:     &ua,
				CostMs:        costMs,
			}
			_ = logRepo.CreateOperationLog(log)
		}()
	}
}

// parseRoute 从 URL 路径提取 module、action、targetID
// /api/manage/dept         POST   → module=dept, action=create
// /api/manage/dept/:id     PUT    → module=dept, action=update, targetID=id
// /api/manage/dept/:id     DELETE → module=dept, action=delete, targetID=id
// /api/manage/role/:id/menus PUT  → module=role, action=grant_menus
func parseRoute(url, method string) (module, action string, targetID *uint64) {
	// 去掉 /api/manage/ 前缀，取剩余路径段
	path := strings.TrimPrefix(url, "/api/manage/")
	segments := strings.Split(path, "/")

	if len(segments) == 0 || segments[0] == "" {
		return "unknown", strings.ToLower(method), nil
	}

	module = segments[0] // dept, role, book, etc.

	switch len(segments) {
	case 1:
		// /api/manage/dept → 集合操作
		action = methodToAction(method)
	case 2:
		// /api/manage/dept/:id 或 /api/manage/dept/page
		if segments[1] == "page" {
			// 分页查询虽是 POST 但不属于写操作，不过既然走到这里也记录
			action = "page"
		} else {
			action = methodToAction(method)
			if id, err := strconv.ParseUint(segments[1], 10, 64); err == nil {
				targetID = &id
			}
		}
	default:
		// /api/manage/role/:id/menus → grant_menus
		if _, err := strconv.ParseUint(segments[1], 10, 64); err == nil {
			targetIDVal, _ := strconv.ParseUint(segments[1], 10, 64)
			targetID = &targetIDVal
		}
		action = strings.ToLower(method) + "_" + segments[len(segments)-1]
	}

	return module, action, targetID
}

func methodToAction(method string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "DELETE":
		return "delete"
	case "PATCH":
		return "patch"
	default:
		return strings.ToLower(method)
	}
}

// sensitiveFieldReplacements 预编译的正则 + 替换模板
// 匹配 JSON 中的敏感字段值并替换为 ***
var sensitiveFieldReplacements = func() []struct {
	re          *regexp.Regexp
	replacement string
} {
	fields := []string{"password", "token", "refreshToken", "secret"}
	out := make([]struct {
		re          *regexp.Regexp
		replacement string
	}, len(fields))
	for i, field := range fields {
		out[i].re = regexp.MustCompile(`"` + regexp.QuoteMeta(field) + `"\s*:\s*"[^"]*"`)
		out[i].replacement = `"` + field + `":"***"`
	}
	return out
}()

// sanitizeBody 脱敏请求体中的敏感字段
func sanitizeBody(body string) string {
	if body == "" {
		return ""
	}
	// 截断过长 body（超过 2KB）
	if len(body) > 2048 {
		body = body[:2048] + "...(truncated)"
	}
	for _, sr := range sensitiveFieldReplacements {
		body = sr.re.ReplaceAllString(body, sr.replacement)
	}
	return body
}

// toUint64 将 gin context 中的 user_id 转为 uint64
func toUint64(v interface{}) uint64 {
	switch x := v.(type) {
	case uint64:
		return x
	case uint:
		return uint64(x)
	case int64:
		return uint64(x)
	case string:
		id, _ := strconv.ParseUint(x, 10, 64)
		return id
	}
	return 0
}
