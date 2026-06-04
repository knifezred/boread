package model

import (
	"time"

	"gorm.io/gorm"
)

// ======================== ReaderNote ========================

// ReaderNote 读者笔记/划线表 (reader_note)
type ReaderNote struct {
	BaseModel
	ReaderID       uint64  `gorm:"column:reader_id;not null" json:"readerId"`
	BookID         uint64  `gorm:"column:book_id;not null" json:"bookId"`
	ChapterID      *uint64 `gorm:"column:chapter_id" json:"chapterId"`
	NoteType       string  `gorm:"column:note_type;type:char(1);not null;default:'1'" json:"noteType"`
	SelectedText   *string `gorm:"column:selected_text;type:text" json:"selectedText"`
	StartOffset    *uint32 `gorm:"column:start_offset" json:"startOffset"`
	EndOffset      *uint32 `gorm:"column:end_offset" json:"endOffset"`
	HighlightColor *string `gorm:"column:highlight_color;size:16" json:"highlightColor"`
	Content        *string `gorm:"column:content;type:text" json:"content"`
	Visibility     string  `gorm:"column:visibility;type:char(1);not null;default:'2'" json:"visibility"`
}

func (ReaderNote) TableName() string { return "reader_note" }

// ======================== BookReview ========================

// BookReview 整本书评表 (book_review)
type BookReview struct {
	BaseModel
	BookID   uint64  `gorm:"column:book_id;not null" json:"bookId"`
	ReaderID uint64  `gorm:"column:reader_id;not null" json:"readerId"`
	Rating   *uint8  `gorm:"column:rating" json:"rating"`
	Title    *string `gorm:"column:title;size:255" json:"title"`
	Content  string  `gorm:"column:content;type:text;not null" json:"content"`
	OwnerID  uint64  `gorm:"column:owner_id;not null" json:"ownerId"`
	DeptID   *uint64 `gorm:"column:dept_id" json:"deptId"`
	Status   string  `gorm:"column:status;type:char(1);not null;default:'1'" json:"status"`
}

func (BookReview) TableName() string { return "book_review" }

// ======================== BookChapterComment ========================

// BookChapterComment 章节评论表 (book_chapter_comment)
type BookChapterComment struct {
	BaseModel
	BookID    uint64  `gorm:"column:book_id;not null" json:"bookId"`
	ChapterID uint64  `gorm:"column:chapter_id;not null" json:"chapterId"`
	ReaderID  uint64  `gorm:"column:reader_id;not null" json:"readerId"`
	ParentID  uint64  `gorm:"column:parent_id;not null;default:0" json:"parentId"`
	ReplyToID *uint64 `gorm:"column:reply_to_id" json:"replyToId"`
	Content   string  `gorm:"column:content;type:text;not null" json:"content"`
	OwnerID   uint64  `gorm:"column:owner_id;not null" json:"ownerId"`
	DeptID    *uint64 `gorm:"column:dept_id" json:"deptId"`
	Status    string  `gorm:"column:status;type:char(1);not null;default:'1'" json:"status"`
}

func (BookChapterComment) TableName() string { return "book_chapter_comment" }

// ======================== ReaderLike ========================

// ReaderLike 点赞表 (reader_like)
// 不继承 BaseModel：无 create_by/update_by，软删除即取消点赞
type ReaderLike struct {
	ID         uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ReaderID   uint64         `gorm:"column:reader_id;not null;uniqueIndex:uk_reader_target" json:"readerId"`
	TargetType string         `gorm:"column:target_type;type:char(1);not null;uniqueIndex:uk_reader_target" json:"targetType"`
	TargetID   uint64         `gorm:"column:target_id;not null;uniqueIndex:uk_reader_target" json:"targetId"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime:milli" json:"createTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (ReaderLike) TableName() string { return "reader_like" }

// 锚定引用
var (
	_ = ReaderNote{}
	_ = BookReview{}
	_ = BookChapterComment{}
	_ = ReaderLike{}
	_ = time.Time{}
	_ = gorm.DeletedAt{}
)
