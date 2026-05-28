package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

var _ = model.BookCategory{}

type BookCategoryHandler struct {
	svc *service.BookCategoryService
}

func NewBookCategoryHandler(svc *service.BookCategoryService) *BookCategoryHandler {
	return &BookCategoryHandler{svc: svc}
}

// Tree 分类树
// @Summary   分类树
// @Tags      book-category
// @Security  BearerAuth
// @Produce   json
// @Success  200  {object}  response.Response{data=[]dto.CategoryNode}
// @Router   /api/manage/book-category/tree [get]
func (h *BookCategoryHandler) Tree(c *gin.Context) {
	tree, err := h.svc.Tree(c.Request.Context())
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, tree)
}

// GetByID 分类详情
// @Summary   分类详情
// @Tags      book-category
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "分类ID"
// @Success  200  {object}  response.Response{data=model.BookCategory}
// @Router   /api/manage/book-category/{id} [get]
func (h *BookCategoryHandler) GetByID(c *gin.Context) {
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

// Create 新增分类
// @Summary   新增分类
// @Tags      book-category
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.CategoryRequest  true  "分类"
// @Success  200  {object}  response.Response{data=model.BookCategory}
// @Router   /api/manage/book-category [post]
func (h *BookCategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Create(c.Request.Context(), &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapCategoryErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Update 编辑分类
// @Summary   编辑分类
// @Tags      book-category
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                true  "分类ID"
// @Param    body  body  dto.CategoryRequest  true  "分类"
// @Success  200  {object}  response.Response{data=model.BookCategory}
// @Router   /api/manage/book-category/{id} [put]
func (h *BookCategoryHandler) Update(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	m, err := h.svc.Update(c.Request.Context(), id, &req, utils.GetUserID(c))
	if err != nil {
		response.Error(c, mapCategoryErr(err), err.Error())
		return
	}
	response.Success(c, m)
}

// Delete 删除分类
// @Summary   删除分类
// @Tags      book-category
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "分类ID"
// @Success  200  {object}  response.Response
// @Router   /api/manage/book-category/{id} [delete]
func (h *BookCategoryHandler) Delete(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapCategoryErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// Page 分类分页列表 (树形)
// @Summary   分类分页列表 (树形)
// @Tags      book-category
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param     body  body  dto.CategorySearch  true  "搜索参数"
// @Success   200      {object}  response.Response{data=dto.PageResponse{records=[]dto.CategoryNode}}
// @Router   /api/manage/book-category/page [post]
func (h *BookCategoryHandler) Page(c *gin.Context) {
	var req dto.CategorySearch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 4001, err.Error())
		return
	}
	res, err := h.svc.Page(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, res)
}

func mapCategoryErr(err error) int {
	switch {
	case service.ErrCategoryCodeExists == err,
		service.ErrCategoryParentNotFound == err:
		return 3001
	case service.ErrCategoryHasChildren == err:
		return 3002
	default:
		return 5001
	}
}
