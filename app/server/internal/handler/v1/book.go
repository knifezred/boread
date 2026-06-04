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

var _ = model.Book{}

type BookHandler struct {
	svc *service.BookService
}

func NewBookHandler(svc *service.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

// GetByID 书籍详情
// @Summary   书籍详情
// @Tags      book
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response{data=dto.BookResponse}
// @Router   /api/book/{id} [get]
func (h *BookHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	m, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, mapBookErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Create 新增书籍
// @Summary   新增书籍
// @Tags      book
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.BookRequest  true  "书籍"
// @Success  200  {object}  response.Response{data=dto.BookResponse}
// @Router   /api/book [post]
func (h *BookHandler) Create(c *gin.Context) {
	var req dto.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapBookErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 编辑书籍
// @Summary   编辑书籍
// @Tags      book
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int            true  "书籍ID"
// @Param    body  body  dto.BookRequest  true  "书籍"
// @Success  200  {object}  response.Response{data=dto.BookResponse}
// @Router   /api/book/{id} [put]
func (h *BookHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	var req dto.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapBookErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除书籍
// @Summary   删除书籍
// @Tags      book
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/{id} [delete]
func (h *BookHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapBookErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// Page 书籍分页列表
// @Summary   书籍分页列表
// @Tags      book
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.BookSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/page [post]
func (h *BookHandler) Page(c *gin.Context) {
	var req dto.BookSearch
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

// UpdateStatus 更新上架状态
// @Summary   更新上架状态
// @Tags      book
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                         true  "书籍ID"
// @Param    body  body  dto.BookUpdateStatusRequest  true  "状态"
// @Success  200  {object}  response.Response
// @Router   /api/book/{id}/status [put]
func (h *BookHandler) UpdateStatus(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, code.ParamInvalid, "invalid id")
		return
	}
	var req dto.BookUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status, utils.GetUserID(c)); err != nil {
		response.Error(c, mapBookErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

func mapBookErr(err error) int {
	return code.MapServiceError(err)
}
