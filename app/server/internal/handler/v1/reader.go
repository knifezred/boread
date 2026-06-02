package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

var _ = dto.BookshelfRequest{}

type ReaderHandler struct {
	svc *service.ReaderService
}

func NewReaderHandler(svc *service.ReaderService) *ReaderHandler {
	return &ReaderHandler{svc: svc}
}

// AddToBookshelf 添加到书架
// @Summary   添加到书架
// @Tags      reader-bookshelf
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.BookshelfRequest  true  "书架参数"
// @Success  200  {object}  response.Response{data=dto.BookshelfResponse}
// @Router   /api/manage/bookshelf [post]
func (h *ReaderHandler) AddToBookshelf(c *gin.Context) {
	var req dto.BookshelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
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
// @Router   /api/manage/bookshelf/{bookId} [delete]
func (h *ReaderHandler) RemoveFromBookshelf(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
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
// @Router   /api/manage/bookshelf/{bookId} [put]
func (h *ReaderHandler) UpdateBookshelf(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
		return
	}
	var req dto.BookshelfUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
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
// @Router   /api/manage/bookshelf/page [post]
func (h *ReaderHandler) GetBookshelfPage(c *gin.Context) {
	var req dto.BookshelfSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.GetBookshelfPage(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
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
// @Router   /api/manage/bookshelf/groups [get]
func (h *ReaderHandler) ListGroups(c *gin.Context) {
	groups, err := h.svc.ListGroups(c.Request.Context(), utils.GetUserID(c))
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, groups)
}

// ReportProgress 上报阅读进度
// @Summary   上报阅读进度
// @Tags      reader-progress
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    bookId  path  int                         true  "书籍ID"
// @Param    body    body  dto.ReadProgressRequest      true  "进度参数"
// @Success  200  {object}  response.Response{data=dto.ReadProgressResponse}
// @Router   /api/manage/reader/progress/{bookId} [put]
func (h *ReaderHandler) ReportProgress(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
		return
	}
	var req dto.ReadProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.ReportProgress(c.Request.Context(), utils.GetUserID(c), bookID, &req)
	if err != nil {
		response.Error(c, mapProgressErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// GetProgress 获取阅读进度
// @Summary   获取阅读进度
// @Tags      reader-progress
// @Security  BearerAuth
// @Produce   json
// @Param    bookId  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response{data=dto.ReadProgressResponse}
// @Router   /api/manage/reader/progress/{bookId} [get]
func (h *ReaderHandler) GetProgress(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
		return
	}
	resp, err := h.svc.GetProgress(c.Request.Context(), utils.GetUserID(c), bookID)
	if err != nil {
		response.Error(c, mapProgressErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

func mapBookshelfErr(err error) int {
	switch {
	case service.ErrBookshelfNotFound == err:
		return 3101
	case service.ErrBookNotExist == err:
		return 3001
	case service.ErrAlreadyInBookshelf == err:
		return 3102
	default:
		return 5001
	}
}

func mapProgressErr(err error) int {
	switch {
	case service.ErrProgressNotFound == err:
		return 3201
	case service.ErrBookNotExist == err:
		return 3001
	default:
		return 5001
	}
}
