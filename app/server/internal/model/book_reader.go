package model

import "time"

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

// ReaderReadEvent 阅读事件原子表 (reader_read_event)
// 纯追加日志, 不嵌入 BaseModel (无 create_by/update_by/deleted_at)
type ReaderReadEvent struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ReaderID    uint64    `gorm:"column:reader_id;not null" json:"readerId"`
	BookID      uint64    `gorm:"column:book_id;not null" json:"bookId"`
	ChapterID   uint64    `gorm:"column:chapter_id;not null" json:"chapterId"`
	SessionID   string    `gorm:"column:session_id;size:36;not null" json:"sessionId"`
	DurationSec uint32    `gorm:"column:duration_sec;not null;default:0" json:"durationSec"`
	WordCount   uint32    `gorm:"column:word_count;not null;default:0" json:"wordCount"`
	EventDate   string    `gorm:"column:event_date;type:date;not null" json:"eventDate"`
	EventTime   time.Time `gorm:"column:event_time;autoCreateTime:milli" json:"eventTime"`
	DeviceType  string    `gorm:"column:device_type;size:32;not null;default:'web'" json:"deviceType"`
}

func (ReaderReadEvent) TableName() string { return "reader_read_event" }

// 锚定引用
var (
	_ = ReaderReadProgress{}
	_ = ReaderReadEvent{}
	_ = time.Time{}
)
