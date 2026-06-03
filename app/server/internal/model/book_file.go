package model

import "gorm.io/gorm"

// FileSourceType 文件来源
type FileSourceType string

const (
	FileSourceUserUpload  FileSourceType = "1" // 用户上传
	FileSourceAdminUpload FileSourceType = "2" // 管理员上传
	FileSourceLocalScan   FileSourceType = "3" // 本地扫描
)

// FileStatus 文件解析状态
type FileStatus string

const (
	FilePending    FileStatus = "1" // 待处理
	FileProcessing FileStatus = "2" // 处理中
	FileSuccess    FileStatus = "3" // 解析成功
	FileFailed     FileStatus = "4" // 解析失败
)

// BookFile 书籍物理文件
type BookFile struct {
	BaseModel
	BookID         uint64         `gorm:"column:book_id;not null" json:"bookId"`
	OriginalName   string         `gorm:"column:original_name;size:255;not null" json:"originalName"`
	SourceType     FileSourceType `gorm:"column:source_type;type:char(1);not null;default:'1'" json:"sourceType"`
	SourceFormat   *string        `gorm:"column:source_format;size:16" json:"sourceFormat"`
	SourceFileURL  *string        `gorm:"column:source_file_url;size:512" json:"sourceFileUrl"`
	ContentPath    *string        `gorm:"column:content_path;size:512" json:"contentPath"`
	ContentSize    uint64         `gorm:"column:content_size;not null;default:0" json:"contentSize"`
	ContentMD5     *string        `gorm:"column:content_md5;size:32" json:"contentMd5"`
	ContentCharset string         `gorm:"column:content_charset;size:16;not null;default:'utf-8'" json:"contentCharset"`
	ContentVersion uint32         `gorm:"column:content_version;not null;default:1" json:"contentVersion"`
	ChapterCount   uint32         `gorm:"column:chapter_count;not null;default:0" json:"chapterCount"`
	IsPrimary      bool           `gorm:"column:is_primary;not null;default:0" json:"isPrimary"`
	FileStatus     FileStatus     `gorm:"column:file_status;type:char(1);not null;default:'1'" json:"fileStatus"`
	ParseMessage   *string        `gorm:"column:parse_message;size:512" json:"parseMessage"`
}

func (BookFile) TableName() string { return "book_file" }

// ChapterStatus 章节状态
type ChapterStatus string

const (
	ChapterPublished ChapterStatus = "1" // 发布
	ChapterDraft     ChapterStatus = "2" // 草稿
	ChapterRemoved   ChapterStatus = "3" // 下架
)

// BookChapter 章节索引
type BookChapter struct {
	BaseModel
	BookID      uint64        `gorm:"column:book_id;not null" json:"bookId"`
	FileID      uint64        `gorm:"column:file_id;not null" json:"fileId"`
	VolumeNo    *uint32       `gorm:"column:volume_no" json:"volumeNo"`
	VolumeTitle *string       `gorm:"column:volume_title;size:255" json:"volumeTitle"`
	ChapterNo   uint32        `gorm:"column:chapter_no;not null" json:"chapterNo"`
	Title       string        `gorm:"column:title;size:255;not null" json:"title"`
	ByteOffset  uint64        `gorm:"column:byte_offset;not null" json:"byteOffset"`
	ByteLength  uint32        `gorm:"column:byte_length;not null" json:"byteLength"`
	WordCount   uint32        `gorm:"column:word_count;not null;default:0" json:"wordCount"`
	IsVip       bool          `gorm:"column:is_vip;not null;default:0" json:"isVip"`
	Status      ChapterStatus `gorm:"column:status;type:char(1);not null;default:'1'" json:"status"`
}

func (BookChapter) TableName() string { return "book_chapter" }

// ParseStatus 解析状态
type ParseStatus string

const (
	ParsePending    ParseStatus = "1" // 待解析
	ParseProcessing ParseStatus = "2" // 解析中
	ParseSuccess    ParseStatus = "3" // 解析成功
	ParseFailed     ParseStatus = "4" // 解析失败
)

// BookUpload 上传/解析任务
type BookUpload struct {
	BaseModel
	BookID       *uint64     `gorm:"column:book_id" json:"bookId"`
	OriginalName string      `gorm:"column:original_name;size:255;not null" json:"originalName"`
	FilePath     string      `gorm:"column:file_path;size:512;not null" json:"filePath"`
	FileSize     uint64      `gorm:"column:file_size;not null;default:0" json:"fileSize"`
	FileMD5      *string     `gorm:"column:file_md5;size:32" json:"fileMd5"`
	SourceFormat *string     `gorm:"column:source_format;size:16" json:"sourceFormat"`
	ParseStatus  ParseStatus `gorm:"column:parse_status;type:char(1);not null;default:'1'" json:"parseStatus"`
	ParseMessage *string     `gorm:"column:parse_message;size:512" json:"parseMessage"`
	ChapterCount *uint32     `gorm:"column:chapter_count" json:"chapterCount"`
}

func (BookUpload) TableName() string { return "book_upload" }

// RuleType 规则类型
type RuleType string

const (
	RuleTypeSystem RuleType = "1" // 系统默认规则
	RuleTypeCustom RuleType = "2" // 用户自定义规则
)

// BookChapterRule 章节识别规则
type BookChapterRule struct {
	BaseModel
	RuleName      string       `gorm:"column:rule_name;size:64;not null" json:"ruleName"`
	RuleType      RuleType     `gorm:"column:rule_type;type:char(1);not null;default:'1'" json:"ruleType"`
	UserID        *uint64      `gorm:"column:user_id" json:"userId"`
	TitlePattern  string       `gorm:"column:title_pattern;size:512;not null" json:"titlePattern"`
	GroupPattern  *string      `gorm:"column:group_pattern;size:512" json:"groupPattern"`
	MinChapterLen uint32       `gorm:"column:min_chapter_len;not null;default:100" json:"minChapterLen"`
	MaxChapterLen uint32       `gorm:"column:max_chapter_len;not null;default:100000" json:"maxChapterLen"`
	SortOrder     int          `gorm:"column:sort_order;not null;default:0" json:"sortOrder"`
	Description   *string      `gorm:"column:description;size:255" json:"description"`
	Status        EnableStatus `gorm:"column:status;type:char(1);not null;default:'1'" json:"status"`
}

func (BookChapterRule) TableName() string { return "book_chapter_rule" }

// BookChapterRuleRel 书籍章节规则关联表
type BookChapterRuleRel struct {
	BaseModel
	BookID   uint64 `gorm:"column:book_id;not null" json:"bookId"`
	ReaderID uint64 `gorm:"column:reader_id;not null" json:"readerId"`
	RuleID   uint64 `gorm:"column:rule_id;not null" json:"ruleId"`
}

func (BookChapterRuleRel) TableName() string { return "book_chapter_rule_rel" }

// FilterMatchType 匹配方式
type FilterMatchType string

const (
	FilterKeyword FilterMatchType = "1" // 关键词
	FilterRegex   FilterMatchType = "2" // 正则
)

// FilterAction 过滤动作
type FilterAction string

const (
	FilterReplace    FilterAction = "1" // 替换
	FilterBlock      FilterAction = "2" // 拦截整章
	FilterMarkReview FilterAction = "3" // 标记审核
)

// FilterApplyStage 应用阶段
type FilterApplyStage string

const (
	FilterStageInput  FilterApplyStage = "1" // 入库时(解析阶段)
	FilterStageOutput FilterApplyStage = "2" // 出库时(读章节时)
)

// FilterSeverity 严重程度
type FilterSeverity string

const (
	FilterLow    FilterSeverity = "1"
	FilterMedium FilterSeverity = "2"
	FilterHigh   FilterSeverity = "3"
)

// BookContentFilterRule 内容净化规则
type BookContentFilterRule struct {
	BaseModel
	RuleName    string           `gorm:"column:rule_name;size:64;not null" json:"ruleName"`
	MatchType   FilterMatchType  `gorm:"column:match_type;type:char(1);not null;default:'1'" json:"matchType"`
	Pattern     string           `gorm:"column:pattern;size:512;not null" json:"pattern"`
	Action      FilterAction     `gorm:"column:action;type:char(1);not null;default:'1'" json:"action"`
	Replacement string           `gorm:"column:replacement;size:255;not null;default:'***'" json:"replacement"`
	ApplyStage  FilterApplyStage `gorm:"column:apply_stage;type:char(1);not null;default:'1'" json:"applyStage"`
	Category    *string          `gorm:"column:category;size:32" json:"category"`
	Severity    FilterSeverity   `gorm:"column:severity;type:char(1);not null;default:'1'" json:"severity"`
	Description *string          `gorm:"column:description;size:255" json:"description"`
	Status      EnableStatus     `gorm:"column:status;type:char(1);not null;default:'1'" json:"status"`
}

func (BookContentFilterRule) TableName() string { return "book_content_filter_rule" }

// 锚定引用让 swag 能解析
var (
	_ = BookFile{}
	_ = BookChapter{}
	_ = BookUpload{}
	_ = BookChapterRule{}
	_ = BookChapterRuleRel{}
	_ = BookContentFilterRule{}
	_ = gorm.DeletedAt{}
)
