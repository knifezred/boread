package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

// BookReaderHandler 读者阅读处理器 (阅读进度 + 阅读事件)
type BookReaderHandler struct {
	svc *service.BookReaderService
}

func NewBookReaderHandler(svc *service.BookReaderService) *BookReaderHandler {
	return &BookReaderHandler{svc: svc}
}

// ==================== 阅读进度 ====================

// ReportProgress 上报阅读进度
// @Summary   上报阅读进度
// @Tags      book-reader
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    bookId  path  int                         true  "书籍ID"
// @Param    body    body  dto.ReadProgressRequest      true  "进度参数"
// @Success  200  {object}  response.Response{data=dto.ReadProgressResponse}
// @Router   /api/book/reader/progress/{bookId} [put]
func (h *BookReaderHandler) ReportProgress(c *gin.Context) {
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
		response.Error(c, mapReaderErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// GetProgress 获取阅读进度
// @Summary   获取阅读进度
// @Tags      book-reader
// @Security  BearerAuth
// @Produce   json
// @Param    bookId  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response{data=dto.ReadProgressResponse}
// @Router   /api/book/reader/progress/{bookId} [get]
func (h *BookReaderHandler) GetProgress(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
		return
	}
	resp, err := h.svc.GetProgress(c.Request.Context(), utils.GetUserID(c), bookID)
	if err != nil {
		response.Error(c, mapReaderErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// ==================== 阅读事件 ====================

// ReportEvent 上报阅读事件
// @Summary   上报阅读事件
// @Tags      book-reader
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ReadEventRequest  true  "阅读事件参数"
// @Success  200  {object}  response.Response
// @Router   /api/book/reader/read-event [post]
func (h *BookReaderHandler) ReportEvent(c *gin.Context) {
	var req dto.ReadEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	if err := h.svc.ReportEvent(c.Request.Context(), utils.GetUserID(c), &req); err != nil {
		response.Error(c, mapReaderErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// ======================== 错误码映射 ========================

func mapReaderErr(err error) int {
	switch {
	case service.ErrProgressNotFound == err:
		return 3201
	case service.ErrBookNotExist == err:
		return 3001
	default:
		return 5001
	}
}
