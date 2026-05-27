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

// 锚定 model 包符号引用, 让 swag 能解析 @Success data=model.SysRole
var _ = model.SysRole{}

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// Page 角色分页
// @Summary   角色分页
// @Tags      role
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.RoleSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/role/page [post]
func (h *RoleHandler) Page(c *gin.Context) {
	var req dto.RoleSearch
	if err := c.ShouldBindJSON(&req); err != nil {
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

// GetByID 角色详情
// @Summary   角色详情
// @Tags      role
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "角色ID"
// @Success  200  {object}  response.Response{data=model.SysRole}
// @Router   /api/manage/role/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
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

// AllBrief 全量角色 (下拉)
// @Summary   全量角色 (下拉)
// @Tags      role
// @Security  BearerAuth
// @Produce   json
// @Success  200  {object}  response.Response{data=[]dto.RoleBrief}
// @Router   /api/manage/role/all [get]
func (h *RoleHandler) AllBrief(c *gin.Context) {
	rows, err := h.svc.AllBrief(c.Request.Context())
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, rows)
}

// Create 新增角色
// @Summary   新增角色
// @Tags      role
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.RoleRequest  true  "角色"
// @Success  200  {object}  response.Response{data=model.SysRole}
// @Router   /api/manage/role [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapRoleErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 编辑角色
// @Summary   编辑角色
// @Tags      role
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int              true  "角色ID"
// @Param    body  body  dto.RoleRequest  true  "角色"
// @Success  200  {object}  response.Response{data=model.SysRole}
// @Router   /api/manage/role/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapRoleErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除角色
// @Summary   删除角色
// @Tags      role
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "角色ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/role/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapRoleErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// GrantMenus 角色授权菜单
// @Summary   角色授权菜单
// @Tags      role
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                   true  "角色ID"
// @Param    body  body  dto.RoleMenuRequest   true  "菜单 IDs"
// @Success  200  {object}  response.Response
// @Router   /api/manage/role/{id}/menus [put]
func (h *RoleHandler) GrantMenus(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.RoleMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	if err := h.svc.GrantMenus(c.Request.Context(), id, req.MenuIDs); err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, nil)
}

// GrantButtons 角色授权按钮
// @Summary   角色授权按钮
// @Tags      role
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                     true  "角色ID"
// @Param    body  body  dto.RoleButtonRequest   true  "按钮 IDs"
// @Success  200  {object}  response.Response
// @Router   /api/manage/role/{id}/buttons [put]
func (h *RoleHandler) GrantButtons(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.RoleButtonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	if err := h.svc.GrantButtons(c.Request.Context(), id, req.ButtonIDs); err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetMenuIDs 角色已授权的菜单 IDs
// @Summary   角色已授权的菜单 IDs
// @Tags      role
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "角色ID"
// @Success  200  {object}  response.Response{data=[]int}
// @Router   /api/manage/role/{id}/menus [get]
func (h *RoleHandler) GetMenuIDs(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	ids, err := h.svc.GetMenuIDs(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, ids)
}

// GetButtonIDs 角色已授权的按钮 IDs
// @Summary   角色已授权的按钮 IDs
// @Tags      role
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "角色ID"
// @Success  200  {object}  response.Response{data=[]int}
// @Router   /api/manage/role/{id}/buttons [get]
func (h *RoleHandler) GetButtonIDs(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	ids, err := h.svc.GetButtonIDs(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, ids)
}

func mapRoleErr(err error) int {
	switch {
	case errors.Is(err, service.ErrRoleCodeExists):
		return 3001
	case errors.Is(err, service.ErrRoleSystem):
		return 3002
	case errors.Is(err, service.ErrRoleHasUsers):
		return 3003
	default:
		return 5001
	}
}
