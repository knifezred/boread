package model

import "time"

type BookTag struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TagName     string    `gorm:"column:tag_name;size:64;not null" json:"tagName"`
	UsageCount  uint32    `gorm:"column:usage_count;default:0" json:"usageCount"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime:milli" json:"createTime"`
}

func (BookTag) TableName() string { return "book_tag" }