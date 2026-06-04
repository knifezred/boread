package model

import "time"

// ReaderBookshelf 读者书架表 (reader_bookshelf)
type ReaderBookshelf struct {
	BaseModel
	ReaderID     uint64     `gorm:"column:reader_id;not null;uniqueIndex:uk_reader_book" json:"readerId"`
	BookID       uint64     `gorm:"column:book_id;not null;uniqueIndex:uk_reader_book" json:"bookId"`
	GroupName    string     `gorm:"column:group_name;size:64;not null;default:'默认'" json:"groupName"`
	IsTop        bool       `gorm:"column:is_top;not null;default:0" json:"isTop"`
	LastReadTime *time.Time `gorm:"column:last_read_time" json:"lastReadTime"`
	AddTime      time.Time  `gorm:"column:add_time;autoCreateTime:milli" json:"addTime"`
}

func (ReaderBookshelf) TableName() string { return "reader_bookshelf" }

// 锚定引用
var _ = ReaderBookshelf{}
