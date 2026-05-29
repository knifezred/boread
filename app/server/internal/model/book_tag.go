package model

type BookTag struct {
	BaseModel
	TagName     string `gorm:"column:tag_name;size:64;not null" json:"tagName"`
	Description string `gorm:"column:description;size:255" json:"description"`
	UsageCount  uint32 `gorm:"column:usage_count;default:0" json:"usageCount"`
}

func (BookTag) TableName() string { return "book_tag" }
