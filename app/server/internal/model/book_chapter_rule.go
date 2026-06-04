package model

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
