package dto

import "boread/internal/model"

// ==================== 章节识别规则 ====================

// ChapterRuleRequest 规则请求
type ChapterRuleRequest struct {
	RuleName      string  `json:"ruleName" binding:"required,max=64"`
	RuleType      string  `json:"ruleType" binding:"required,oneof=1 2"`
	TitlePattern  string  `json:"titlePattern" binding:"required,max=512"`
	GroupPattern  *string `json:"groupPattern"`
	MinChapterLen uint32  `json:"minChapterLen"`
	MaxChapterLen uint32  `json:"maxChapterLen"`
	SortOrder     int     `json:"sortOrder"`
	Description   *string `json:"description"`
	Status        string  `json:"status"`
}

// ChapterRuleSearch 规则搜索
type ChapterRuleSearch struct {
	PageRequest
	RuleName string  `json:"ruleName"`
	RuleType string  `json:"ruleType"`
	UserID   *uint64 `json:"userId"`
	Status   string  `json:"status"`
}

// ChapterRuleResponse 规则响应
type ChapterRuleResponse struct {
	model.BookChapterRule
}

// ==================== 章节规则绑定 ====================

// ChapterRuleBindRequest 绑定规则到书籍请求
type ChapterRuleBindRequest struct {
	BookID uint64 `json:"bookId" binding:"required"`
	RuleID uint64 `json:"ruleId" binding:"required"`
}

// ChapterRuleBindResponse 绑定规则响应
type ChapterRuleBindResponse struct {
	ID         uint64 `json:"id"`
	BookID     uint64 `json:"bookId"`
	ReaderID   uint64 `json:"readerId"`
	RuleID     uint64 `json:"ruleId"`
	RuleName   string `json:"ruleName"`
	CreateTime string `json:"createTime"`
}
