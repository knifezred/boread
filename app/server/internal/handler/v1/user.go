package v1

import (
	"errors"

	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

// 锚定 model 包符号引用, 让 swag 能解析 @Success data=model.SysUser
var _ = model.SysUser{}

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Page 用户分页
// @Summary   用户分页
// @Tags      user
// @Security  BearerAuth
// @Produce   json
// @Param    current     query  int     false  "页码"  default(1)
// @Param    size        query  int     false  "页大小"  default(10)
// @Param    userName    query  string  false  "用户名"
// @Param    nickName    query  string  false  "昵称"
// @Param    userPhone   query  string  false  "手机号"
// @Param    userEmail   query  string  false  "邮箱"
// @Param    userGender  query  string  false  "性别"
// @Param    status      query  string  false  "状态"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/user [get]
func (h *UserHandler) Page(c *gin.Context) {
	var req dto.UserSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.Page(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetByID 用户详情
// @Summary   用户详情
// @Tags      user
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "用户ID"
// @Success  200  {object}  response.Response{data=model.SysUser}
// @Router   /api/manage/user/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	m, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 3001, err.Error())
		return
	}
	response.Success(c, m)
}

// Create 新增用户
// @Summary   新增用户
// @Tags      user
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.UserCreateRequest  true  "用户"
// @Success  200  {object}  response.Response{data=model.SysUser}
// @Router   /api/manage/user [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapUserErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 编辑用户
// @Summary   编辑用户
// @Tags      user
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                     true  "用户ID"
// @Param    body  body  dto.UserUpdateRequest   true  "用户"
// @Success  200  {object}  response.Response{data=model.SysUser}
// @Router   /api/manage/user/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapUserErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除用户
// @Summary   删除用户
// @Tags      user
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "用户ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/user/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, nil)
}

// ResetPassword 重置密码
// @Summary   重置密码
// @Tags      user
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                       true  "用户ID"
// @Param    body  body  dto.UserResetPwdRequest   true  "新密码"
// @Success  200  {object}  response.Response
// @Router   /api/manage/user/{id}/reset-password [put]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.UserResetPwdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	if err := h.svc.ResetPassword(c.Request.Context(), id, req.Password, utils.GetUserID(c)); err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, nil)
}

func mapUserErr(err error) int {
	if errors.Is(err, service.ErrUserExists) {
		return 3001
	}
	return 5001
}
