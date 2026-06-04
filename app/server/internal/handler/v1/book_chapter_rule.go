package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

type BookChapterRuleHandler struct {
	svc *service.BookChapterRuleService
}

func NewBookChapterRuleHandler(svc *service.BookChapterRuleService) *BookChapterRuleHandler {
	return &BookChapterRuleHandler{svc: svc}
}

// CreateChapterRule 创建章节识别规则
// @Summary   创建章节识别规则
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterRuleRequest  true  "规则"
// @Success  200  {object}  response.Response{data=model.BookChapterRule}
// @Router   /api/book/chapter-rule [post]
func (h *BookChapterRuleHandler) CreateChapterRule(c *gin.Context) {
	var req dto.ChapterRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.CreateChapterRule(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
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
// @Router   /api/book/chapter-rule/{id} [put]
func (h *BookChapterRuleHandler) UpdateChapterRule(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	var req dto.ChapterRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.UpdateChapterRule(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
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
// @Router   /api/book/chapter-rule/{id} [delete]
func (h *BookChapterRuleHandler) DeleteChapterRule(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	if err := h.svc.DeleteChapterRule(c.Request.Context(), id); err != nil {
		response.Error(c, code.ServerError, err.Error())
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
// @Router   /api/book/chapter-rule/{id} [get]
func (h *BookChapterRuleHandler) GetChapterRuleByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	m, err := h.svc.GetChapterRuleByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
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
// @Router   /api/book/chapter-rule/page [post]
func (h *BookChapterRuleHandler) PageChapterRule(c *gin.Context) {
	var req dto.ChapterRuleSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.PageChapterRule(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// BindChapterRule 为书籍绑定章节识别规则
// @Summary   为书籍绑定章节识别规则
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ChapterRuleBindRequest  true  "绑定参数"
// @Success  200  {object}  response.Response{data=dto.ChapterRuleBindResponse}
// @Router   /api/book/chapter-rule/bind [post]
func (h *BookChapterRuleHandler) BindChapterRule(c *gin.Context) {
	var req dto.ChapterRuleBindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.BindChapterRule(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// UnbindChapterRule 解绑书籍的章节识别规则
// @Summary   解绑书籍的章节识别规则
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Produce   json
// @Param    bookId  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/chapter-rule/bind/{bookId} [delete]
func (h *BookChapterRuleHandler) UnbindChapterRule(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid bookId")
		return
	}
	if err := h.svc.UnbindChapterRule(c.Request.Context(), bookID, utils.GetUserID(c)); err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetBoundChapterRule 获取书籍绑定的章节识别规则
// @Summary   获取书籍绑定的章节识别规则
// @Tags      book-chapter-rule
// @Security  BearerAuth
// @Produce   json
// @Param    bookId  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response{data=dto.ChapterRuleBindResponse}
// @Router   /api/book/chapter-rule/bind/{bookId} [get]
func (h *BookChapterRuleHandler) GetBoundChapterRule(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid bookId")
		return
	}
	resp, err := h.svc.GetBoundChapterRule(c.Request.Context(), bookID, utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}
