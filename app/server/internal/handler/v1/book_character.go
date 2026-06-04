package v1

import (
	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
	"boread/pkg/utils"
)

// ======================== CharacterHandler ========================

type CharacterHandler struct {
	svc *service.BookCharacterService
}

func NewCharacterHandler(svc *service.BookCharacterService) *CharacterHandler {
	return &CharacterHandler{svc: svc}
}

// CreateCharacter 创建角色
// @Summary   创建角色
// @Tags      book-character
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.CharacterRequest  true  "角色参数"
// @Success  200  {object}  response.Response{data=dto.CharacterResponse}
// @Router   /api/book/character [post]
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	var req dto.CharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, mapCharErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// UpdateCharacter 更新角色
// @Summary   更新角色
// @Tags      book-character
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    id    path  int                true  "角色ID"
// @Param    body  body  dto.CharacterRequest  true  "角色参数"
// @Success  200  {object}  response.Response{data=dto.CharacterResponse}
// @Router   /api/book/character/{id} [put]
func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	var req dto.CharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), utils.GetUserID(c), id, &req)
	if err != nil {
		response.Error(c, mapCharErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// DeleteCharacter 删除角色
// @Summary   删除角色
// @Tags      book-character
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "角色ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/character/{id} [delete]
func (h *CharacterHandler) DeleteCharacter(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapCharErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// GetCharacter 获取角色详情
// @Summary   获取角色详情
// @Tags      book-character
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "角色ID"
// @Success  200  {object}  response.Response{data=dto.CharacterResponse}
// @Router   /api/book/character/{id} [get]
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	resp, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, mapCharErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// PageCharacter 分页查询角色
// @Summary   分页查询角色
// @Tags      book-character
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.CharacterSearch  true  "搜索参数"
// @Success  200  {object}  response.Response{data=dto.PageResponse}
// @Router   /api/book/character/page [post]
func (h *CharacterHandler) PageCharacter(c *gin.Context) {
	var req dto.CharacterSearch
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

// ListByCharacterBook 按书查询角色列表
// @Summary   按书查询角色列表
// @Tags      book-character
// @Security  BearerAuth
// @Produce   json
// @Param    bookId  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response{data=[]dto.CharacterResponse}
// @Router   /api/book/character/book/{bookId} [get]
func (h *CharacterHandler) ListByCharacterBook(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
		return
	}
	resp, err := h.svc.ListByBook(c.Request.Context(), bookID)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// ======================== CharacterRelHandler ========================

type CharacterRelHandler struct {
	svc *service.BookCharacterRelService
}

func NewCharacterRelHandler(svc *service.BookCharacterRelService) *CharacterRelHandler {
	return &CharacterRelHandler{svc: svc}
}

// CreateRelation 创建角色关系
// @Summary   创建角色关系
// @Tags      book-character-rel
// @Security  BearerAuth
// @Accept    json
// @Produce   json
// @Param    body  body  dto.CharacterRelRequest  true  "关系参数"
// @Success  200  {object}  response.Response{data=dto.CharacterRelResponse}
// @Router   /api/book/character/rel [post]
func (h *CharacterRelHandler) CreateRelation(c *gin.Context) {
	var req dto.CharacterRelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), utils.GetUserID(c), &req)
	if err != nil {
		response.Error(c, mapCharRelErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// DeleteRelation 删除角色关系
// @Summary   删除角色关系
// @Tags      book-character-rel
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "关系ID"
// @Success  200  {object}  response.Response
// @Router   /api/book/character/rel/{id} [delete]
func (h *CharacterRelHandler) DeleteRelation(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, mapCharRelErr(err), err.Error())
		return
	}
	response.Success(c, nil)
}

// GetRelation 获取角色关系详情
// @Summary   获取角色关系详情
// @Tags      book-character-rel
// @Security  BearerAuth
// @Produce   json
// @Param    id  path  int  true  "关系ID"
// @Success  200  {object}  response.Response{data=dto.CharacterRelResponse}
// @Router   /api/book/character/rel/{id} [get]
func (h *CharacterRelHandler) GetRelation(c *gin.Context) {
	id, err := utils.ParseUint64Param(c, "id")
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}
	resp, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, mapCharRelErr(err), err.Error())
		return
	}
	response.Success(c, resp)
}

// ListRelationsByCharacter 查询某个角色的所有关系
// @Summary   查询角色所有关系
// @Tags      book-character-rel
// @Security  BearerAuth
// @Produce   json
// @Param    characterId  path  int  true  "角色ID"
// @Success  200  {object}  response.Response{data=[]dto.CharacterRelResponse}
// @Router   /api/book/character/rel/character/{characterId} [get]
func (h *CharacterRelHandler) ListRelationsByCharacter(c *gin.Context) {
	characterID, err := utils.ParseUint64Param(c, "characterId")
	if err != nil {
		response.Error(c, 1001, "invalid characterId")
		return
	}
	resp, err := h.svc.ListByCharacter(c.Request.Context(), characterID)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// ListRelationsByBook 查询某本书的所有角色关系
// @Summary   查询某本书所有角色关系
// @Tags      book-character-rel
// @Security  BearerAuth
// @Produce   json
// @Param    bookId  path  int  true  "书籍ID"
// @Success  200  {object}  response.Response{data=[]dto.CharacterRelResponse}
// @Router   /api/book/character/rel/book/{bookId} [get]
func (h *CharacterRelHandler) ListRelationsByBook(c *gin.Context) {
	bookID, err := utils.ParseUint64Param(c, "bookId")
	if err != nil {
		response.Error(c, 1001, "invalid bookId")
		return
	}
	resp, err := h.svc.ListByBook(c.Request.Context(), bookID)
	if err != nil {
		response.Error(c, 5001, err.Error())
		return
	}
	response.Success(c, resp)
}

// ======================== 错误码映射 ========================

func mapCharErr(err error) int {
	switch {
	case service.ErrCharacterNotFound == err:
		return 3601
	case service.ErrBookExists == err:
		return 3001
	default:
		return 5001
	}
}

func mapCharRelErr(err error) int {
	switch {
	case service.ErrCharacterRelNotFound == err:
		return 3701
	case service.ErrCharacterNotFound == err:
		return 3601
	default:
		return 5001
	}
}
