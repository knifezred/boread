package model

type BookCategory struct {
	BaseModel
	ParentID     uint64       `gorm:"column:parent_id;default:0;index" json:"parentId"`
	Ancestors    string       `gorm:"column:ancestors;size:1024;default:''" json:"ancestors"`
	CategoryName string       `gorm:"column:category_name;size:64;not null" json:"categoryName"`
	CategoryCode string       `gorm:"column:category_code;size:64;not null" json:"categoryCode"`
	Description  string       `gorm:"column:description;size:255" json:"description"`
	SortOrder    int          `gorm:"column:sort_order;default:0" json:"sortOrder"`
	Status       EnableStatus `gorm:"column:status;type:char(1);default:'1'" json:"status"`
}

func (BookCategory) TableName() string { return "book_category" }