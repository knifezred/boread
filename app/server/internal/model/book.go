package model

// SerialStatus 连载状态
type SerialStatus string

const (
	SerialOngoing  SerialStatus = "1" // 连载中
	SerialFinished SerialStatus = "2" // 已完结
	SerialStopped  SerialStatus = "3" // 断更
)

// Visibility 可见性
type Visibility string

const (
	VisibilityPublic  Visibility = "1" // 公开
	VisibilityPrivate Visibility = "2" // 仅自己
	VisibilityDept    Visibility = "3" // 部门内
)

// AggregateStatus 聚合状态
type AggregateStatus string

const (
	AggregateSingle  AggregateStatus = "1" // 单文件(无需聚合)
	AggregateMerging AggregateStatus = "2" // 多文件聚合中
	AggregateDone    AggregateStatus = "3" // 聚合完成
)

// BookStatus 上架状态
type BookStatus string

const (
	BookListed    BookStatus = "1" // 已上架
	BookUnlisted  BookStatus = "2" // 下架
	BookReviewing BookStatus = "3" // 审核中
	BookRejected  BookStatus = "4" // 审核拒绝
)

type Book struct {
	BaseModel
	Title           string          `gorm:"column:title;size:255;not null" json:"title"`
	Author          string          `gorm:"column:author;size:128;not null;default:''" json:"author"`
	Cover           *string         `gorm:"column:cover;size:512" json:"cover,omitempty"`
	Intro           *string         `gorm:"column:intro;type:text" json:"intro,omitempty"`
	CategoryID      *uint64         `gorm:"column:category_id;index" json:"categoryId,omitempty"`
	Language        string          `gorm:"column:language;size:16;not null;default:'zh-CN'" json:"language"`
	SerialStatus    SerialStatus    `gorm:"column:serial_status;type:char(1);not null;default:'1'" json:"serialStatus"`
	Visibility      Visibility      `gorm:"column:visibility;type:char(1);not null;default:'1'" json:"visibility"`
	PrimaryFileID   *uint64         `gorm:"column:primary_file_id" json:"primaryFileId,omitempty"`
	TotalChapters   uint32          `gorm:"column:total_chapters;not null;default:0" json:"totalChapters"`
	TotalWords      uint32          `gorm:"column:total_words;not null;default:0" json:"totalWords"`
	AggregateStatus AggregateStatus `gorm:"column:aggregate_status;type:char(1);not null;default:'1'" json:"aggregateStatus"`
	AvgRating       float64         `gorm:"column:avg_rating;type:decimal(2,1);not null;default:0.0" json:"avgRating"`
	RatingCount     uint32          `gorm:"column:rating_count;not null;default:0" json:"ratingCount"`
	OwnerID         uint64          `gorm:"column:owner_id;not null" json:"ownerId"`
	DeptID          *uint64         `gorm:"column:dept_id;index" json:"deptId,omitempty"`
	Status          BookStatus      `gorm:"column:status;type:char(1);not null;default:'1'" json:"status"`
}

func (Book) TableName() string { return "book" }

type BookTagRel struct {
	BaseModel
	BookID uint64 `gorm:"column:book_id;not null;uniqueIndex:uk_book_tag" json:"bookId"`
	TagID  uint64 `gorm:"column:tag_id;not null;uniqueIndex:uk_book_tag" json:"tagId"`
}

func (BookTagRel) TableName() string { return "book_tag_rel" }
