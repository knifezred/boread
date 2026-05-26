package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
)

type LogHandler struct {
	svc *service.LogService
}

func NewLogHandler(svc *service.LogService) *LogHandler {
	return &LogHandler{svc: svc}
}

// PageLogin 登录日志分页
// @Summary   登录日志分页
// @Tags      log
// @Security  BearerAuth
// @Produce   json
// @Param    current      query  int     false  "页码"  default(1)
// @Param    size         query  int     false  "页大小"  default(10)
// @Param    userName     query  string  false  "用户名"
// @Param    loginIp      query  string  false  "登录IP"
// @Param    loginType    query  string  false  "类型: 1登录 2登出"
// @Param    loginResult  query  string  false  "结果: 1成功 2失败"
// @Param    startTime    query  string  false  "开始时间"
// @Param    endTime      query  string  false  "结束时间"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/log/login [get]
func (h *LogHandler) PageLogin(c *gin.Context) {
	var req dto.LoginLogSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageLogin(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// PageOperation 操作日志分页
// @Summary   操作日志分页
// @Tags      log
// @Security  BearerAuth
// @Produce   json
// @Param    current    query  int     false  "页码"  default(1)
// @Param    size       query  int     false  "页大小"  default(10)
// @Param    userName   query  string  false  "用户名"
// @Param    module     query  string  false  "模块"
// @Param    action     query  string  false  "动作"
// @Param    clientIp   query  string  false  "客户端IP"
// @Param    startTime  query  string  false  "开始时间"
// @Param    endTime    query  string  false  "结束时间"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/log/operation [get]
func (h *LogHandler) PageOperation(c *gin.Context) {
	var req dto.OperationLogSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageOperation(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}
