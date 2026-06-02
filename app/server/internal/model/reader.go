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

// ReaderReadProgress 阅读进度表 (reader_read_progress)
type ReaderReadProgress struct {
	BaseModel
	ReaderID     uint64    `gorm:"column:reader_id;not null;uniqueIndex:uk_reader_book" json:"readerId"`
	BookID       uint64    `gorm:"column:book_id;not null;uniqueIndex:uk_reader_book" json:"bookId"`
	FileID       *uint64   `gorm:"column:file_id" json:"fileId"`
	ChapterID    uint64    `gorm:"column:chapter_id;not null" json:"chapterId"`
	ChapterNo    uint32    `gorm:"column:chapter_no;not null" json:"chapterNo"`
	Position     uint32    `gorm:"column:position;not null;default:0" json:"position"`
	Percent      float64   `gorm:"column:percent;type:decimal(5,2);not null;default:0.00" json:"percent"`
	ReadDuration uint32    `gorm:"column:read_duration;not null;default:0" json:"readDuration"`
	LastReadTime time.Time `gorm:"column:last_read_time;autoUpdateTime:milli" json:"lastReadTime"`
}

func (ReaderReadProgress) TableName() string { return "reader_read_progress" }

// 锚定引用
var (
	_ = ReaderBookshelf{}
	_ = ReaderReadProgress{}
	_ = time.Time{}
)
