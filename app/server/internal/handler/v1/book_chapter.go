package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

type BookChapterHandler struct {
	svc *service.BookChapterService
}

func NewBookChapterHandler(svc *service.BookChapterService) *BookChapterHandler {
	return &BookChapterHandler{svc: svc}
}

// PageChapter 章节分页列表
// @Summary   章节分页列表
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/chapter/page [post]
func (h *BookChapterHandler) PageChapter(c *gin.Context) {
	var req dto.ChapterSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.PageChapter(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// ListChapter 章节列表（不分页）
// @Summary   章节列表
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterListRequest  true  "请求参数"
// @Success  200  {object}  response.Response{data=[]dto.ChapterResponse}
// @Router   /api/book/chapter/list [post]
func (h *BookChapterHandler) ListChapter(c *gin.Context) {
	var req dto.ChapterListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.ListChapter(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetChapterContent 读取章节内容
// @Summary   读取章节内容
// @Tags      book-chapter
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "章节ID"
// @Success  200  {object}  response.Response{data=dto.ChapterContentResponse}
// @Router   /api/book/chapter/{id}/content [get]
func (h *BookChapterHandler) GetChapterContent(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid chapterId")
		return
	}
	resp, err := h.svc.GetChapterContent(c.Request.Context(), id)
	if err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// UpdateChapterTitle 更新章节标题
// @Summary   更新章节标题
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                          true  "章节ID"
// @Param    body  body  dto.ChapterTitleUpdateRequest  true  "标题"
// @Success  200  {object}  response.Response
// @Router   /api/book/chapter/{id}/title [put]
func (h *BookChapterHandler) UpdateChapterTitle(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid chapterId")
		return
	}
	var req dto.ChapterTitleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	if err := h.svc.UpdateChapterTitle(c.Request.Context(), id, req.Title, utils.GetUserID(c)); err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// BatchUpdateChapterTitle 批量更新章节标题
// @Summary   批量更新章节标题
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterTitleBatchRequest  true  "批量标题"
// @Success  200  {object}  response.Response
// @Router   /api/book/chapter/batch-title [put]
func (h *BookChapterHandler) BatchUpdateChapterTitle(c *gin.Context) {
	var req dto.ChapterTitleBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	if err := h.svc.BatchUpdateChapterTitles(c.Request.Context(), req.IDs, req.Title, utils.GetUserID(c)); err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// UpdateChapterStatus 批量更新章节状态
// @Summary   批量更新章节状态
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterStatusBatchRequest  true  "状态"
// @Success  200  {object}  response.Response
// @Router   /api/book/chapter/status [put]
func (h *BookChapterHandler) UpdateChapterStatus(c *gin.Context) {
	var req dto.ChapterStatusBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	if err := h.svc.UpdateChapterStatus(c.Request.Context(), req.IDs, req.Status, utils.GetUserID(c)); err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// DeleteChapter 删除章节（软删除）
// @Summary   删除章节
// @Tags      book-chapter
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "章节ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/chapter/{id} [delete]
func (h *BookChapterHandler) DeleteChapter(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid chapterId")
		return
	}
	if err := h.svc.DeleteChapter(c.Request.Context(), id, utils.GetUserID(c)); err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// MergeChapters 合并章节
// @Summary   合并章节
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterMergeRequest  true  "合并请求"
// @Success  200  {object}  response.Response
// @Router   /api/book/chapter/merge [post]
func (h *BookChapterHandler) MergeChapters(c *gin.Context) {
	var req dto.ChapterMergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	if err := h.svc.MergeChapters(c.Request.Context(), req.BookID, req.TargetID, req.SourceIDs, utils.GetUserID(c)); err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// FormatChapterNumbers 格式化章节编号
// @Summary   格式化章节编号
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterFormatRequest  true  "格式化请求"
// @Success  200  {object}  response.Response
// @Router   /api/book/chapter/format-numbers [post]
func (h *BookChapterHandler) FormatChapterNumbers(c *gin.Context) {
	var req dto.ChapterFormatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	if err := h.svc.FormatChapterNumbers(c.Request.Context(), req.IDs, utils.GetUserID(c)); err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// SaveChapterContent 保存章节内容
// @Summary   保存章节内容
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                           true  "章节ID"
// @Param    body  body  dto.ChapterContentSaveRequest  true  "内容"
// @Success  200  {object}  response.Response
// @Router   /api/book/chapter/{id}/content [put]
func (h *BookChapterHandler) SaveChapterContent(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid chapterId")
		return
	}
	var req dto.ChapterContentSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	req.ChapterID = id
	if err := h.svc.SaveChapterContent(c.Request.Context(), req.BookID, req.ChapterID, req.Content, utils.GetUserID(c)); err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// ReParseChapters 重新识别章节
// @Summary   重新识别章节
// @Tags      book-chapter
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ReParseRequest  true  "请求参数"
// @Success  200  {object}  response.Response{data=dto.ReParseResponse}
// @Router   /api/book/chapter/re-parse [post]
func (h *BookChapterHandler) ReParseChapters(c *gin.Context) {
	var req dto.ReParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.ReParseChapters(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapChapterErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

func mapChapterErr(err error) int {
	return code.MapServiceError(err)
}
