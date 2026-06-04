package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

// ======================== NoteHandler ========================

type NoteHandler struct {
	svc *service.BookSocialService
}

func NewNoteHandler(svc *service.BookSocialService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

// CreateNote 创建笔记/划线
// @Summary   创建笔记/划线
// @Tags      reader-note
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.NoteRequest  true  "笔记参数"
// @Success  200  {object}  response.Response{data=dto.NoteResponse}
// @Router   /api/book/note [post]
func (h *NoteHandler) CreateNote(c *gin.Context) {
	var req dto.NoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.CreateNote(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, mapNoteErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// UpdateNote 更新笔记/划线
// @Summary   更新笔记/划线
// @Tags      reader-note
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int              true  "笔记ID"
// @Param    body  body  dto.NoteRequest  true  "笔记参数"
// @Success  200  {object}  response.Response{data=dto.NoteResponse}
// @Router   /api/book/note/{id} [put]
func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.NoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.UpdateNote(c.Request.Context(), utils.GetUserID(c), id, &req)
	if err != nil {
		response.Error(c, mapNoteErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// DeleteNote 删除笔记/划线
// @Summary   删除笔记/划线
// @Tags      reader-note
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "笔记ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/note/{id} [delete]
func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.DeleteNote(c.Request.Context(), utils.GetUserID(c), id); err != nil {
		response.Error(c, mapNoteErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// GetNote 获取笔记详情
// @Summary   获取笔记详情
// @Tags      reader-note
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "笔记ID"
// @Success  200  {object}  response.Response{data=dto.NoteResponse}
// @Router   /api/book/note/{id} [get]
func (h *NoteHandler) GetNote(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	resp, err := h.svc.GetNote(c.Request.Context(), id, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapNoteErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// PageNote 分页查询笔记
// @Summary   分页查询笔记
// @Tags      reader-note
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.NoteSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/note/page [post]
func (h *NoteHandler) PageNote(c *gin.Context) {
	var req dto.NoteSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageNote(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// ListNotesByBook 获取某本书的公开笔记
// @Summary   获取某本书的公开笔记
// @Tags      reader-note
// @Security  BearerAuth
// @Produce   json
// @Param    bookId  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response{data=[]dto.NoteResponse}
// @Router   /api/book/note/book/{bookId} [get]
func (h *NoteHandler) ListNotesByBook(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
		return
	}
	noteType := c.Query("noteType")
	var pageReq dto.PageRequest
	pageReq.Current = 1
	pageReq.Size = 50
	resp, err := h.svc.ListNotesByBook(c.Request.Context(), bookID, noteType, &pageReq)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

func mapNoteErr(err error) int {
	switch {
	case service.ErrNoteNotFound == err:
		return 3301
	case service.ErrNoteNotOwner == err:
		return 3302
	case service.ErrBookNotExists == err:
		return 3001
	default:
		return 5001
	}
}

// ======================== ReviewHandler ========================

type ReviewHandler struct {
	svc *service.BookSocialService
}

func NewReviewHandler(svc *service.BookSocialService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

// CreateReview 创建书评
// @Summary   创建书评
// @Tags      book-review
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ReviewRequest  true  "书评参数"
// @Success  200  {object}  response.Response{data=dto.ReviewResponse}
// @Router   /api/book/review [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.CreateReview(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, mapReviewErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// UpdateReview 更新书评
// @Summary   更新书评
// @Tags      book-review
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int               true  "书评ID"
// @Param    body  body  dto.ReviewRequest  true  "书评参数"
// @Success  200  {object}  response.Response{data=dto.ReviewResponse}
// @Router   /api/book/review/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.UpdateReview(c.Request.Context(), utils.GetUserID(c), id, &req)
	if err != nil {
		response.Error(c, mapReviewErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// DeleteReview 删除书评
// @Summary   删除书评
// @Tags      book-review
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "书评ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/review/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.DeleteReview(c.Request.Context(), utils.GetUserID(c), id); err != nil {
		response.Error(c, mapReviewErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// GetReview 获取书评详情
// @Summary   获取书评详情
// @Tags      book-review
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "书评ID"
// @Success  200  {object}  response.Response{data=dto.ReviewResponse}
// @Router   /api/book/review/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	resp, err := h.svc.GetReview(c.Request.Context(), id, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapReviewErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// PageReview 分页查询书评
// @Summary   分页查询书评
// @Tags      book-review
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.ReviewSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/review/page [post]
func (h *ReviewHandler) PageReview(c *gin.Context) {
	var req dto.ReviewSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageReview(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

func mapReviewErr(err error) int {
	switch {
	case service.ErrReviewNotFound == err:
		return 3401
	case service.ErrNoteNotOwner == err:
		return 3402
	case service.ErrBookNotExists == err:
		return 3001
	default:
		return 5001
	}
}

// ======================== CommentHandler ========================

type CommentHandler struct {
	svc *service.BookSocialService
}

func NewCommentHandler(svc *service.BookSocialService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// CreateComment 创建章节评论
// @Summary   创建章节评论
// @Tags      book-comment
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.CommentRequest  true  "评论参数"
// @Success  200  {object}  response.Response{data=dto.CommentResponse}
// @Router   /api/book/comment [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.CreateComment(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, mapCommentErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// DeleteComment 删除章节评论
// @Summary   删除章节评论
// @Tags      book-comment
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "评论ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/comment/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.DeleteComment(c.Request.Context(), utils.GetUserID(c), id); err != nil {
		response.Error(c, mapCommentErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// GetComment 获取评论详情
// @Summary   获取评论详情
// @Tags      book-comment
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "评论ID"
// @Success  200  {object}  response.Response{data=dto.CommentResponse}
// @Router   /api/book/comment/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	resp, err := h.svc.GetComment(c.Request.Context(), id, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapCommentErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// PageComment 分页查询评论
// @Summary   分页查询评论
// @Tags      book-comment
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.CommentSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/comment/page [post]
func (h *CommentHandler) PageComment(c *gin.Context) {
	var req dto.CommentSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.PageComment(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

func mapCommentErr(err error) int {
	switch {
	case service.ErrCommentNotFound == err:
		return 3501
	case service.ErrCommentNotOwner == err:
		return 3502
	case service.ErrParentCommentNotExist == err:
		return 3503
	case service.ErrChapterNotExists == err:
		return 3041
	default:
		return 5001
	}
}

// ======================== LikeHandler ========================

type LikeHandler struct {
	svc *service.BookSocialService
}

func NewLikeHandler(svc *service.BookSocialService) *LikeHandler {
	return &LikeHandler{svc: svc}
}

// ToggleLike 切换点赞
// @Summary   切换点赞
// @Tags      like
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.LikeRequest  true  "点赞参数"
// @Success  200  {object}  response.Response{data=dto.LikeResponse}
// @Router   /api/book/like/toggle [post]
func (h *LikeHandler) ToggleLike(c *gin.Context) {
	var req dto.LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.ToggleLike(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetLikeStatus 批量查询点赞状态
// @Summary   批量查询点赞状态
// @Tags      like
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.LikeStatusRequest  true  "查询参数"
// @Success  200  {object}  response.Response{data=map[string]bool}
// @Router   /api/book/like/status [post]
func (h *LikeHandler) GetLikeStatus(c *gin.Context) {
	var req dto.LikeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.GetLikeStatus(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// CountLikes 查询点赞数
// @Summary   查询点赞数
// @Tags      like
// @Security  BearerAuth
// @Produce   json
// @Param    targetType  path  string  true  "目标类型"
// @Param    targetId    path  int     true  "目标ID"
// @Success  200  {object}  response.Response{data=int64}
// @Router   /api/book/like/count/{targetType}/{targetId} [get]
func (h *LikeHandler) CountLikes(c *gin.Context) {
	targetType := c.Param("targetType")
	targetID, err := utils.ParseUint64Param(c, "targetId")
	if err != nil {
		response.Error(c, 1001, "invalid targetId")
		return
	}
	count, err := h.svc.CountLikes(c.Request.Context(), targetType, targetID)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, count)
}
