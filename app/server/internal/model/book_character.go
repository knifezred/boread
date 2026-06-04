package model

import "gorm.io/gorm"

// BookCharacter 小说角色表 (book_character)
type BookCharacter struct {
	BaseModel
	BookID    uint64  `gorm:"column:book_id;not null" json:"bookId"`
	Name      string  `gorm:"column:name;size:128;not null" json:"name"`
	Alias     *string `gorm:"column:alias;size:255" json:"alias"`
	RoleType  string  `gorm:"column:role_type;type:char(1);not null;default:'2'" json:"roleType"`
	Avatar    *string `gorm:"column:avatar;size:512" json:"avatar"`
	Intro     *string `gorm:"column:intro;type:text" json:"intro"`
	Extra     *string `gorm:"column:extra;type:json" json:"extra"`
	SortOrder int     `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`
}

func (BookCharacter) TableName() string { return "book_character" }

// BookCharacterRel 角色关系表 (book_character_rel)
type BookCharacterRel struct {
	BaseModel
	BookID       uint64  `gorm:"column:book_id;not null" json:"bookId"`
	CharacterAID uint64  `gorm:"column:character_a_id;not null" json:"characterAId"`
	CharacterBID uint64  `gorm:"column:character_b_id;not null" json:"characterBId"`
	RelationType string  `gorm:"column:relation_type;size:32;not null" json:"relationType"`
	RelationDesc *string `gorm:"column:relation_desc;size:255" json:"relationDesc"`
	SortOrder    int     `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`
}

func (BookCharacterRel) TableName() string { return "book_character_rel" }

var (
	_ = BookCharacter{}
	_ = BookCharacterRel{}
	_ = gorm.DeletedAt{}
)
