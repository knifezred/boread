package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

var (
	_ = model.BookFile{}
	_ = model.BookChapter{}
	_ = model.BookUpload{}
	_ = model.BookChapterRule{}
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
// @Router   /api/manage/book/upload [post]
func (h *BookFileHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, 1001, "请选择要上传的文件")
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
// @Router   /api/manage/book/confirm-import [post]
func (h *BookFileHandler) ConfirmImport(c *gin.Context) {
	var req dto.ConfirmImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
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
// @Router   /api/manage/book/scan [post]
func (h *BookFileHandler) ScanAll(c *gin.Context) {
	resp, err := h.svc.ScanPending(c.Request.Context())
	if err != nil {
		response.Error(c, 5001, err.Error())
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
// @Router   /api/manage/book/scan-path [post]
func (h *BookFileHandler) ScanPath(c *gin.Context) {
	var req dto.ScanPathRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.ScanPath(c.Request.Context(), req.Path, utils.GetUserID(c))
	if err != nil {
		response.Error(c, 5001, err.Error())
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
// @Router   /api/manage/book/scan/{id} [post]
func (h *BookFileHandler) ScanByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	resp, err := h.svc.ScanByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// GetChapterContent 读取章节内容
// @Summary   读取章节内容
// @Tags      book-file
// @Security  BearerAuth
// @Produce   json
// @Param    id     path  int  true  "书籍ID"
// @Param    chapterNo  path  int  true  "章节序号"
// @Success  200  {object}  response.Response{data=dto.ChapterContentResponse}
// @Router   /api/manage/book/{id}/chapter/{chapterNo} [get]
func (h *BookFileHandler) GetChapterContent(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
		return
	}
	chapterNo, err := utils.ParseUint64Param(c, "chapterNo")
	if err != nil {
		response.Error(c, 1001, "invalid chapterNo")
		return
	}
	resp, err := h.svc.GetChapterContent(c.Request.Context(), bookID, uint32(chapterNo))
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// ReParseChapters 重新识别章节
// @Summary   重新识别章节（删除旧章节索引，重新解析并创建）
// @Tags      book-file
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ReParseRequest  true  "请求参数"
// @Success  200  {object}  response.Response{data=dto.ReParseResponse}
// @Router   /api/manage/book/re-parse [post]
func (h *BookFileHandler) ReParseChapters(c *gin.Context) {
	var req dto.ReParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.ReParseChapters(c.Request.Context(), &req)
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
// @Router   /api/manage/book/upload/page [post]
func (h *BookFileHandler) PageUpload(c *gin.Context) {
	var req dto.UploadSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageUpload(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
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
// @Router   /api/manage/book/file/page [post]
func (h *BookFileHandler) PageFile(c *gin.Context) {
	var req dto.FileSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageFile(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// ==================== 章节 ====================

// PageChapter 章节分页
// @Summary   章节分页列表
// @Tags      book-file
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/book/chapter/page [post]
func (h *BookFileHandler) PageChapter(c *gin.Context) {
	var req dto.ChapterSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageChapter(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// ListChapter 章节列表（不分页）
// @Summary   章节列表
// @Tags      book-file
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterListRequest  true  "请求参数"
// @Success  200  {object}  response.Response{data=[]dto.ChapterResponse}
// @Router   /api/manage/book/chapter/list [post]
func (h *BookFileHandler) ListChapter(c *gin.Context) {
	var req dto.ChapterListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.ListChapter(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// ==================== 章节识别规则 ====================

// CreateChapterRule 创建章节识别规则
// @Summary   创建章节识别规则
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterRuleRequest  true  "规则"
// @Success  200  {object}  response.Response{data=model.BookChapterRule}
// @Router   /api/manage/book/chapter-rule [post]
func (h *BookFileHandler) CreateChapterRule(c *gin.Context) {
	var req dto.ChapterRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.CreateChapterRule(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// UpdateChapterRule 更新章节识别规则
// @Summary   更新章节识别规则
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                    true  "规则ID"
// @Param    body  body  dto.ChapterRuleRequest  true  "规则"
// @Success  200  {object}  response.Response{data=model.BookChapterRule}
// @Router   /api/manage/book/chapter-rule/{id} [put]
func (h *BookFileHandler) UpdateChapterRule(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.ChapterRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.UpdateChapterRule(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// DeleteChapterRule 删除章节识别规则
// @Summary   删除章节识别规则
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "规则ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/book/chapter-rule/{id} [delete]
func (h *BookFileHandler) DeleteChapterRule(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.DeleteChapterRule(c.Request.Context(), id); err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// GetChapterRuleByID 章节识别规则详情
// @Summary   章节识别规则详情
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "规则ID"
// @Success  200  {object}  response.Response{data=model.BookChapterRule}
// @Router   /api/manage/book/chapter-rule/{id} [get]
func (h *BookFileHandler) GetChapterRuleByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	m, err := h.svc.GetChapterRuleByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, mapFileErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// PageChapterRule 章节识别规则分页
// @Summary   章节识别规则分页
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterRuleSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/book/chapter-rule/page [post]
func (h *BookFileHandler) PageChapterRule(c *gin.Context) {
	var req dto.ChapterRuleSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageChapterRule(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
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
// @Router   /api/manage/book/filter-rule [post]
func (h *BookFileHandler) CreateFilterRule(c *gin.Context) {
	var req dto.FilterRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
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
// @Router   /api/manage/book/filter-rule/{id} [put]
func (h *BookFileHandler) UpdateFilterRule(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.FilterRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
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
// @Router   /api/manage/book/filter-rule/{id} [delete]
func (h *BookFileHandler) DeleteFilterRule(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
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
// @Router   /api/manage/book/filter-rule/{id} [get]
func (h *BookFileHandler) GetFilterRuleByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
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
// @Router   /api/manage/book/filter-rule/page [post]
func (h *BookFileHandler) PageFilterRule(c *gin.Context) {
	var req dto.FilterRuleSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageFilterRule(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

func mapFileErr(err error) int {
	switch {
	case service.ErrFileTooLarge == err,
		service.ErrFileTypeUnsupported == err,
		service.ErrFileEmpty == err,
		service.ErrFileMD5Exists == err:
		return 1002
	case service.ErrUploadNotFound == err,
		service.ErrChapterNotFound == err,
		service.ErrBookFileNotFound == err,
		service.ErrRuleNotFound == err,
		service.ErrFilterRuleNotFound == err:
		return 3001
	default:
		return 5001
	}
}
