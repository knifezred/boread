package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/code"
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
// @Accept    json
// @Produce   json
// @Param    body  body  dto.LoginLogSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/log/login/page [post]
func (h *LogHandler) PageLogin(c *gin.Context) {
	var req dto.LoginLogSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.PageLogin(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// PageOperation 操作日志分页
// @Summary   操作日志分页
// @Tags      log
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.OperationLogSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/log/operation/page [post]
func (h *LogHandler) PageOperation(c *gin.Context) {
	var req dto.OperationLogSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.PageOperation(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}
