package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

// 锚定 model 包符号引用, 让 swag 能解析 @Success data=model.SysXxx
var _ = model.SysDept{}

type DeptHandler struct {
	svc *service.DeptService
}

func NewDeptHandler(svc *service.DeptService) *DeptHandler {
	return &DeptHandler{svc: svc}
}

// Page 部门分页列表
// @Summary   部门分页列表 (平级)
// @Tags      dept
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     body  body  dto.DeptSearch  true  "搜索参数"
// @Success   200      {object}  response.Response{data=dto.PageResponse{records=[]dto.DeptNode}}
// @Router    /api/manage/dept/page [post]
func (h *DeptHandler) Page(c *gin.Context) {
	var req dto.DeptSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.SearchParamInvalid, err.Error())
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
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, res)
}

// Tree 部门树
// @Summary   部门树
// @Tags      dept
// @Security  BearerAuth
// @Produce   json
// @Param    deptName  query  string  false  "部门名称"
// @Param    deptCode  query  string  false  "部门编码"
// @Param    status    query  string  false  "状态: 1启用 2禁用"
// @Success  200  {object}  response.Response{data=[]dto.DeptNode}
// @Router   /api/manage/dept/tree [get]
func (h *DeptHandler) Tree(c *gin.Context) {
	var search dto.DeptSearch
	if err := c.ShouldBindQuery(&search); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	tree, err := h.svc.Tree(c.Request.Context(), &search)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, tree)
}

// GetByID 部门详情
// @Summary   部门详情
// @Tags      dept
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "部门ID"
// @Success  200  {object}  response.Response{data=model.SysDept}
// @Router   /api/manage/dept/{id} [get]
func (h *DeptHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	m, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, code.ResourceConflict, err.Error())
		return
	}
	response.Success(c, m)
}

// Create 创建
// @Summary   新增部门
// @Tags      dept
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.DeptRequest  true  "部门"
// @Success  200  {object}  response.Response{data=model.SysDept}
// @Router   /api/manage/dept [post]
func (h *DeptHandler) Create(c *gin.Context) {
	var req dto.DeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapDeptErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 更新
// @Summary   编辑部门
// @Tags      dept
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int              true  "部门ID"
// @Param    body  body  dto.DeptRequest  true  "部门"
// @Success  200  {object}  response.Response{data=model.SysDept}
// @Router   /api/manage/dept/{id} [put]
func (h *DeptHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	var req dto.DeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapDeptErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除
// @Summary   删除部门
// @Tags      dept
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "部门ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/dept/{id} [delete]
func (h *DeptHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapDeptErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

func mapDeptErr(err error) int {
	return code.MapServiceError(err)
}
