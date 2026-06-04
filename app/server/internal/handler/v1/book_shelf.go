package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

type BookshelfHandler struct {
	svc *service.ReaderBookshelfService
}

func NewBookshelfHandler(svc *service.ReaderBookshelfService) *BookshelfHandler {
	return &BookshelfHandler{svc: svc}
}

// AddToBookshelf 添加到书架
// @Summary   添加到书架
// @Tags      reader-bookshelf
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.BookshelfRequest  true  "书架参数"
// @Success  200  {object}  response.Response{data=dto.BookshelfResponse}
// @Router   /api/book/shelf [post]
func (h *BookshelfHandler) AddToBookshelf(c *gin.Context) {
	var req dto.BookshelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.AddToBookshelf(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, mapBookshelfErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// RemoveFromBookshelf 从书架移除
// @Summary   从书架移除
// @Tags      reader-bookshelf
// @Security  BearerAuth
// @Produce   json
// @Param    bookId  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/shelf/{bookId} [delete]
func (h *BookshelfHandler) RemoveFromBookshelf(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid bookId")
		return
	}
	if err := h.svc.RemoveFromBookshelf(c.Request.Context(), utils.GetUserID(c), bookID); err != nil {
		response.Error(c, mapBookshelfErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// UpdateBookshelf 更新书架 (分组/置顶)
// @Summary   更新书架
// @Tags      reader-bookshelf
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    bookId  path  int                         true  "书籍ID"
// @Param    body    body  dto.BookshelfUpdateRequest  true  "更新参数"
// @Success  200  {object}  response.Response{data=dto.BookshelfResponse}
// @Router   /api/book/shelf/{bookId} [put]
func (h *BookshelfHandler) UpdateBookshelf(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid bookId")
		return
	}
	var req dto.BookshelfUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.UpdateBookshelf(c.Request.Context(), utils.GetUserID(c), bookID, &req)
	if err != nil {
		response.Error(c, mapBookshelfErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// GetBookshelfPage 书架分页列表
// @Summary   书架分页列表
// @Tags      reader-bookshelf
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.BookshelfSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.BookshelfPageResponse}
// @Router   /api/book/shelf/page [post]
func (h *BookshelfHandler) GetBookshelfPage(c *gin.Context) {
	var req dto.BookshelfSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	resp, err := h.svc.GetBookshelfPage(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// ListGroups 获取分组列表
// @Summary   获取书架分组列表
// @Tags      reader-bookshelf
// @Security  BearerAuth
// @Produce   json
// @Success  200  {object}  response.Response{data=[]dto.BookshelfGroupItem}
// @Router   /api/book/shelf/groups [get]
func (h *BookshelfHandler) ListGroups(c *gin.Context) {
	groups, err := h.svc.ListGroups(c.Request.Context(), utils.GetUserID(c))
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, groups)
}

func mapBookshelfErr(err error) int {
	return code.MapServiceError(err)
}
