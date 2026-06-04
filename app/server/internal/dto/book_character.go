package dto

import "boread/internal/model"

// ==================== 小说角色 ====================

// CharacterRequest 创建/更新角色请求
type CharacterRequest struct {
	BookID    uint64  `json:"bookId" binding:"required"`
	Name      string  `json:"name" binding:"required,max=128"`
	Alias     *string `json:"alias"`
	RoleType  string  `json:"roleType" binding:"required,oneof=1 2 3 4"`
	Avatar    *string `json:"avatar"`
	Intro     *string `json:"intro"`
	Extra     *string `json:"extra"`
	SortOrder int     `json:"sortOrder"`
}

// CharacterResponse 角色响应
type CharacterResponse struct {
	model.BookCharacter
}

// CharacterSearch 角色搜索
type CharacterSearch struct {
	PageRequest
	BookID   uint64 `json:"bookId" binding:"required"`
	RoleType string `json:"roleType"`
	Keyword  string `json:"keyword"`
}

// ==================== 角色关系 ====================

// CharacterRelRequest 创建角色关系请求
type CharacterRelRequest struct {
	BookID       uint64  `json:"bookId" binding:"required"`
	CharacterAID uint64  `json:"characterAId" binding:"required"`
	CharacterBID uint64  `json:"characterBId" binding:"required"`
	RelationType string  `json:"relationType" binding:"required,max=32"`
	RelationDesc *string `json:"relationDesc"`
	SortOrder    int     `json:"sortOrder"`
}

// CharacterRelResponse 角色关系响应
type CharacterRelResponse struct {
	model.BookCharacterRel
	CharacterAName *string `json:"characterAName,omitempty"`
	CharacterBName *string `json:"characterBName,omitempty"`
}
