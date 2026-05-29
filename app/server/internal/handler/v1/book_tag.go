package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

var _ = model.BookTag{}

type BookTagHandler struct {
	svc *service.BookTagService
}

func NewBookTagHandler(svc *service.BookTagService) *BookTagHandler {
	return &BookTagHandler{svc: svc}
}

// Page 标签分页
// @Summary   标签分页
// @Tags      book-tag
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.TagSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/manage/book-tag/page [post]
func (h *BookTagHandler) Page(c *gin.Context) {
	var req dto.TagSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.Page(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetByID 标签详情
// @Summary   标签详情
// @Tags      book-tag
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "标签ID"
// @Success  200  {object}  response.Response{data=model.BookTag}
// @Router   /api/manage/book-tag/{id} [get]
func (h *BookTagHandler) GetByID(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	m, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 3001, err.Error())
		return
	}
	response.Success(c, m)
}

// Create 新增标签
// @Summary   新增标签
// @Tags      book-tag
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.TagRequest  true  "标签"
// @Success  200  {object}  response.Response{data=model.BookTag}
// @Router   /api/manage/book-tag [post]
func (h *BookTagHandler) Create(c *gin.Context) {
	var req dto.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapTagErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 编辑标签
// @Summary   编辑标签
// @Tags      book-tag
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int              true  "标签ID"
// @Param    body  body  dto.TagRequest   true  "标签"
// @Success  200  {object}  response.Response{data=model.BookTag}
// @Router   /api/manage/book-tag/{id} [put]
func (h *BookTagHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapTagErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除标签
// @Summary   删除标签
// @Tags      book-tag
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "标签ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/book-tag/{id} [delete]
func (h *BookTagHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapTagErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

func mapTagErr(err error) int {
	if service.ErrTagNameExists == err {
		return 3001
	}
	return 5001
}
