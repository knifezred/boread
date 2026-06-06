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

// 锚定 model 包符号引用, 让 swag 能解析 @Success data=model.SysSetting
var _ = model.SysSetting{}

type SettingHandler struct {
	svc *service.SettingService
}

func NewSettingHandler(svc *service.SettingService) *SettingHandler {
	return &SettingHandler{svc: svc}
}

// Page 系统配置分页
// @Summary   系统配置分页
// @Tags      setting
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.SettingSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/setting/page [post]
func (h *SettingHandler) Page(c *gin.Context) {
	var req dto.SettingSearch
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

// GetByID 配置详情
// @Summary   配置详情
// @Tags      setting
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "配置ID"
// @Success  200  {object}  response.Response{data=dto.SettingVO}
// @Router   /api/manage/setting/{id} [get]
func (h *SettingHandler) GetByID(c *gin.Context) {
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

// Create 新增配置
// @Summary   新增配置
// @Tags      setting
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.SettingRequest  true  "配置"
// @Success  200  {object}  response.Response{data=model.SysSetting}
// @Router   /api/manage/setting [post]
func (h *SettingHandler) Create(c *gin.Context) {
	var req dto.SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapSettingErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 编辑配置
// @Summary   编辑配置
// @Tags      setting
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                       true  "配置ID"
// @Param    body  body  dto.SettingUpdateRequest   true  "配置"
// @Success  200  {object}  response.Response{data=model.SysSetting}
// @Router   /api/manage/setting/{id} [put]
func (h *SettingHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	var req dto.SettingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapSettingErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除配置
// @Summary   删除配置
// @Tags      setting
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "配置ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/setting/{id} [delete]
func (h *SettingHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapSettingErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// Categories 获取所有启用配置分类
// @Summary   获取配置分类列表
// @Tags      setting
// @Security  BearerAuth
// @Produce   json
// @Success  200  {object}  response.Response{data=[]string}
// @Router   /api/manage/setting/categories [get]
func (h *SettingHandler) Categories(c *gin.Context) {
	cats, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	if cats == nil {
		cats = []string{}
	}
	response.Success(c, cats)
}

func mapSettingErr(err error) int {
	return code.MapServiceError(err)
}

// GetByCategory 按分类获取所有配置
// @Summary   按分类获取所有配置
// @Tags      setting
// @Security  BearerAuth
// @Produce   json
// @Param    category  path  string  true  "分类"
// @Success  200  {object}  response.Response{data=map[string]dto.SettingVO}
// @Router   /api/manage/setting/by-category/{category} [get]
func (h *SettingHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		response.Error(c, code.ParamInvalid, "category is required")
		return
	}
	result, err := h.svc.GetByCategory(c.Request.Context(), category)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	if result == nil {
		result = map[string]dto.SettingVO{}
	}
	response.Success(c, result)
}

// BatchSave 批量保存配置
// @Summary   批量保存配置 (upsert)
// @Tags      setting
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.BatchSaveRequest  true  "批量保存请求"
// @Success  200  {object}  response.Response{data=dto.BatchSaveResult}
// @Router   /api/manage/setting/batch-save [post]
func (h *SettingHandler) BatchSave(c *gin.Context) {
	var req dto.BatchSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	result, err := h.svc.BatchSave(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, result)
}
