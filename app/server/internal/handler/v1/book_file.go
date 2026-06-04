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

var (
	_ = model.BookFile{}
	_ = model.BookUpload{}
	_ = model.BookContentFilterRule{}
)

type BookFileHandler struct {
	svc *service.BookFileService
}

func NewBookFileHandler(svc *service.BookFileService) *BookFileHandler {
	return &BookFileHandler{svc: svc}
}

// Upload 上传小说文件
// @Summary   上传小说文件
// @Tags      book-file
// @Security  BearerAuth
// @Accept    multipart/form-data
// @Produce   json
// @Param    file  formData  file    true  "小说文件 (txt/epub/mobi/pdf)"
// @Success  200   {object}  response.Response{data=dto.FileUploadResponse}
// @Router   /api/book/upload [post]
func (h *BookFileHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, code.ParamInvalid, "请选择要上传的文件")
		return
	}
	defer file.Close()

	resp, err := h.svc.Upload(c.Request.Context(), file, header.Filename, uint64(header.Size), utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// ConfirmImport 确认入库
// @Summary   确认入库（匹配或创建书籍、写入章节索引）
// @Tags      book-file
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ConfirmImportRequest  true  "确认入库参数"
// @Success  200  {object}  response.Response{data=dto.ConfirmImportResponse}
// @Router   /api/book/confirm-import [post]
func (h *BookFileHandler) ConfirmImport(c *gin.Context) {
	var req dto.ConfirmImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.ConfirmImport(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// ScanAll 扫描所有待处理上传任务
// @Summary   扫描入库
// @Tags      book-file
// @Security  BearerAuth
// @Produce   json
// @Success  200  {object}  response.Response{data=dto.ScanAllResponse}
// @Router   /api/book/scan [post]
func (h *BookFileHandler) ScanAll(c *gin.Context) {
	resp, err := h.svc.ScanPending(c.Request.Context())
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// ScanPath 扫描本地目录
// @Summary   扫描本地目录中的小说文件
// @Tags      book-file
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ScanPathRequest  true  "扫描路径"
// @Success  200  {object}  response.Response{data=dto.ScanPathResponse}
// @Router   /api/book/scan-path [post]
func (h *BookFileHandler) ScanPath(c *gin.Context) {
	var req dto.ScanPathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.ScanPath(c.Request.Context(), req.Path, utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// ScanByID 扫描单个上传任务
// @Summary   扫描单个上传任务
// @Tags      book-file
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "上传记录ID"
// @Success  200  {object}  response.Response{data=dto.ScanResult}
// @Router   /api/book/scan/{id} [post]
func (h *BookFileHandler) ScanByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	resp, err := h.svc.ScanByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// ==================== 上传记录 ====================

// PageUpload 上传记录分页
// @Summary   上传记录分页
// @Tags      book-file
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.UploadSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/upload/page [post]
func (h *BookFileHandler) PageUpload(c *gin.Context) {
	var req dto.UploadSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.PageUpload(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// ==================== 文件记录 ====================

// PageFile 文件记录分页
// @Summary   文件记录分页
// @Tags      book-file
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.FileSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/file/page [post]
func (h *BookFileHandler) PageFile(c *gin.Context) {
	var req dto.FileSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.PageFile(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// ==================== 内容净化规则 ====================

// CreateFilterRule 创建内容净化规则
// @Summary   创建内容净化规则
// @Tags      book-filter-rule
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.FilterRuleRequest  true  "规则"
// @Success  200  {object}  response.Response{data=model.BookContentFilterRule}
// @Router   /api/book/filter-rule [post]
func (h *BookFileHandler) CreateFilterRule(c *gin.Context) {
	var req dto.FilterRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.CreateFilterRule(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// UpdateFilterRule 更新内容净化规则
// @Summary   更新内容净化规则
// @Tags      book-filter-rule
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                   true  "规则ID"
// @Param    body  body  dto.FilterRuleRequest  true  "规则"
// @Success  200  {object}  response.Response{data=model.BookContentFilterRule}
// @Router   /api/book/filter-rule/{id} [put]
func (h *BookFileHandler) UpdateFilterRule(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	var req dto.FilterRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.UpdateFilterRule(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// DeleteFilterRule 删除内容净化规则
// @Summary   删除内容净化规则
// @Tags      book-filter-rule
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "规则ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/filter-rule/{id} [delete]
func (h *BookFileHandler) DeleteFilterRule(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	if err := h.svc.DeleteFilterRule(c.Request.Context(), id); err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// GetFilterRuleByID 内容净化规则详情
// @Summary   内容净化规则详情
// @Tags      book-filter-rule
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "规则ID"
// @Success  200  {object}  response.Response{data=model.BookContentFilterRule}
// @Router   /api/book/filter-rule/{id} [get]
func (h *BookFileHandler) GetFilterRuleByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	m, err := h.svc.GetFilterRuleByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// PageFilterRule 内容净化规则分页
// @Summary   内容净化规则分页
// @Tags      book-filter-rule
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.FilterRuleSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/filter-rule/page [post]
func (h *BookFileHandler) PageFilterRule(c *gin.Context) {
	var req dto.FilterRuleSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.PageFilterRule(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

func mapFileErr(err error) int {
	return code.MapServiceError(err)
}
