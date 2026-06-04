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

// 锚定 model 包符号引用, 让 swag 能解析 @Success data=model.SysDict / SysDictItem
var (
	_ = model.SysDict{}
	_ = model.SysDictItem{}
)

type DictHandler struct {
	svc *service.DictService
}

func NewDictHandler(svc *service.DictService) *DictHandler {
	return &DictHandler{svc: svc}
}

// Page 字典分页
// @Summary   字典分页
// @Tags      dict
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.DictSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/dict/page [post]
func (h *DictHandler) Page(c *gin.Context) {
	var req dto.DictSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.Page(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetByID 字典详情
// @Summary   字典详情
// @Tags      dict
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "字典ID"
// @Success  200  {object}  response.Response{data=model.SysDict}
// @Router   /api/manage/dict/{id} [get]
func (h *DictHandler) GetByID(c *gin.Context) {
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

// Create 新增字典
// @Summary   新增字典
// @Tags      dict
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.DictRequest  true  "字典"
// @Success  200  {object}  response.Response{data=model.SysDict}
// @Router   /api/manage/dict [post]
func (h *DictHandler) Create(c *gin.Context) {
	var req dto.DictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapDictErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 编辑字典
// @Summary   编辑字典
// @Tags      dict
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int              true  "字典ID"
// @Param    body  body  dto.DictRequest  true  "字典"
// @Success  200  {object}  response.Response{data=model.SysDict}
// @Router   /api/manage/dict/{id} [put]
func (h *DictHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	var req dto.DictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapDictErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除字典
// @Summary   删除字典
// @Tags      dict
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "字典ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/dict/{id} [delete]
func (h *DictHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapDictErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// ItemsByDictID 按字典ID查项
// @Summary   按字典ID查项
// @Tags      dict
// @Security  BearerAuth
// @Produce   json
// @Param    dictId  path  int  true  "字典ID"
// @Success  200  {object}  response.Response{data=[]model.SysDictItem}
// @Router   /api/manage/dict/items/{dictId} [get]
func (h *DictHandler) ItemsByDictID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "dictId")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid dictId")
		return
	}
	rows, err := h.svc.ListItemsByDictID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, rows)
}

// ItemsByCode 按字典 code 拉取字典项
// @Summary   按字典 code 拉取字典项 (前端高频)
// @Tags      dict
// @Security  BearerAuth
// @Produce   json
// @Param    code  path  string  true  "字典编码"
// @Success  200  {object}  response.Response{data=[]model.SysDictItem}
// @Router   /api/manage/dict/code/{code} [get]
func (h *DictHandler) ItemsByCode(c *gin.Context) {
	dictCode := c.Param("code")
	if dictCode == "" {
		response.Error(c, code.ParamInvalid, "invalid code")
		return
	}
	rows, err := h.svc.ListItemsByCode(c.Request.Context(), dictCode)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, rows)
}

// CreateItem 新增字典项
// @Summary   新增字典项
// @Tags      dict
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.DictItemRequest  true  "字典项"
// @Success  200  {object}  response.Response{data=model.SysDictItem}
// @Router   /api/manage/dict/item [post]
func (h *DictHandler) CreateItem(c *gin.Context) {
	var req dto.DictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.CreateItem(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, m)
}

// UpdateItem 编辑字典项
// @Summary   编辑字典项
// @Tags      dict
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                  true  "字典项ID"
// @Param    body  body  dto.DictItemRequest  true  "字典项"
// @Success  200  {object}  response.Response{data=model.SysDictItem}
// @Router   /api/manage/dict/item/{id} [put]
func (h *DictHandler) UpdateItem(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	var req dto.DictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.UpdateItem(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, m)
}

// DeleteItem 删除字典项
// @Summary   删除字典项
// @Tags      dict
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "字典项ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/dict/item/{id} [delete]
func (h *DictHandler) DeleteItem(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	if err := h.svc.DeleteItem(c.Request.Context(), id); err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

func mapDictErr(err error) int {
	return code.MapServiceError(err)
}
