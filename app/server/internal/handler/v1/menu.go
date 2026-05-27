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

// 锚定 model 包符号引用, 让 swag 能解析 @Success data=model.SysMenu / SysMenuButton
var (
	_ = model.SysMenu{}
	_ = model.SysMenuButton{}
)

type MenuHandler struct {
	svc *service.MenuService
}

func NewMenuHandler(svc *service.MenuService) *MenuHandler {
	return &MenuHandler{svc: svc}
}

// Page 菜单分页列表
// @Summary   菜单分页列表 (平级, 含按钮)
// @Tags      menu
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     body  body  dto.MenuSearch  true  "搜索参数"
// @Success  200      {object}  response.Response{data=dto.PageResponse[dto.MenuNode]}
// @Router   /api/manage/menu/page [post]
func (h *MenuHandler) Page(c *gin.Context) {
	var req dto.MenuSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 4001, err.Error())
		return
	}
	// 分页参数默认值
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 10
	}
	res, err := h.svc.Page(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, res)
}

// Tree 全量菜单树
// @Summary   全量菜单树 (含按钮)
// @Tags      menu
// @Security  BearerAuth
// @Produce   json
// @Success  200  {object}  response.Response{data=[]dto.MenuNode}
// @Router   /api/manage/menu/tree [get]
func (h *MenuHandler) Tree(c *gin.Context) {
	tree, err := h.svc.Tree(c.Request.Context())
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, tree)
}

// GetByID 菜单详情
// @Summary   菜单详情
// @Tags      menu
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "菜单ID"
// @Success  200  {object}  response.Response{data=model.SysMenu}
// @Router   /api/manage/menu/{id} [get]
func (h *MenuHandler) GetByID(c *gin.Context) {
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

// Create 新增菜单
// @Summary   新增菜单
// @Tags      menu
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.MenuRequest  true  "菜单"
// @Success  200  {object}  response.Response{data=model.SysMenu}
// @Router   /api/manage/menu [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req dto.MenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapMenuErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 编辑菜单
// @Summary   编辑菜单
// @Tags      menu
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int              true  "菜单ID"
// @Param    body  body  dto.MenuRequest  true  "菜单"
// @Success  200  {object}  response.Response{data=model.SysMenu}
// @Router   /api/manage/menu/{id} [put]
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.MenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapMenuErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除菜单
// @Summary   删除菜单
// @Tags      menu
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "菜单ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/menu/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapMenuErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// CreateButton 新增菜单按钮
// @Summary   新增菜单按钮
// @Tags      menu
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.MenuButtonRequest  true  "按钮"
// @Success  200  {object}  response.Response{data=model.SysMenuButton}
// @Router   /api/manage/menu/button [post]
func (h *MenuHandler) CreateButton(c *gin.Context) {
	var req dto.MenuButtonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	b, err := h.svc.CreateButton(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, b)
}

// DeleteButton 删除菜单按钮
// @Summary   删除菜单按钮
// @Tags      menu
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "按钮ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/menu/button/{id} [delete]
func (h *MenuHandler) DeleteButton(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.DeleteButton(c.Request.Context(), id); err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListButtonsByMenu 按菜单查按钮
// @Summary   按菜单查按钮
// @Tags      menu
// @Security  BearerAuth
// @Produce   json
// @Param    menuId  path  int  true  "菜单ID"
// @Success  200  {object}  response.Response{data=[]model.SysMenuButton}
// @Router   /api/manage/menu/buttons/{menuId} [get]
func (h *MenuHandler) ListButtonsByMenu(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "menuId")
	if err != nil {
		response.Error(c, 1001, "invalid menuId")
		return
	}
	rows, err := h.svc.ListButtonsByMenu(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, rows)
}

func mapMenuErr(err error) int {
	switch {
	case errors.Is(err, service.ErrMenuRouteExists):
		return 3001
	case errors.Is(err, service.ErrMenuSystem):
		return 3002
	case errors.Is(err, service.ErrMenuHasChildren):
		return 3003
	default:
		return 5001
	}
}
