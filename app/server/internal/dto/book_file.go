package dto

import "boread/internal/model"

// ==================== 文件上传 ====================

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	BookID *uint64 `form:"bookId"` // 可选，关联已有书籍
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	UploadID         uint64  `json:"uploadId"`
	OriginalName     string  `json:"originalName"`
	FileSize         uint64  `json:"fileSize"`
	SourceFormat     *string `json:"sourceFormat"`
	SuggestedTitle   string  `json:"suggestedTitle"`
	SuggestedAuthor  string  `json:"suggestedAuthor"`
	MatchedBookID    *uint64 `json:"matchedBookId"`
	MatchedBookTitle string  `json:"matchedBookTitle"`
}

// ConfirmImportRequest 确认入库请求
type ConfirmImportRequest struct {
	UploadID uint64 `json:"uploadId" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Author   string `json:"author"`
}

// ConfirmImportResponse 确认入库响应
type ConfirmImportResponse struct {
	UploadID     uint64  `json:"uploadId"`
	BookID       uint64  `json:"bookId"`
	BookTitle    string  `json:"bookTitle"`
	BookAuthor   string  `json:"bookAuthor"`
	ChapterCount uint32  `json:"chapterCount"`
	ParseStatus  string  `json:"parseStatus"`
	ParseMessage *string `json:"parseMessage"`
}

// ==================== 上传任务 ====================

// UploadSearch 上传任务搜索
type UploadSearch struct {
	PageRequest
	OriginalName string  `json:"originalName"`
	ParseStatus  string  `json:"parseStatus"`
	BookID       *uint64 `json:"bookId"`
}

// UploadResponse 上传任务响应
type UploadResponse struct {
	model.BookUpload
}

// ==================== 文件管理 ====================

// FileSearch 文件搜索
type FileSearch struct {
	PageRequest
	BookID     *uint64 `json:"bookId"`
	FileStatus string  `json:"fileStatus"`
	SourceType string  `json:"sourceType"`
}

// FileResponse 文件响应
type FileResponse struct {
	model.BookFile
}

// ==================== 章节 ====================

// ChapterSearch 章节搜索
type ChapterSearch struct {
	PageRequest
	BookID    *uint64 `json:"bookId"`
	FileID    *uint64 `json:"fileId"`
	ChapterNo *uint32 `json:"chapterNo"`
}

// ChapterListRequest 章节列表请求（不分页）
type ChapterListRequest struct {
	BookID uint64 `json:"bookId"` // 书籍ID
}

// ChapterResponse 章节响应
type ChapterResponse struct {
	model.BookChapter
}

// ChapterContentResponse 章节内容响应（含文本）
type ChapterContentResponse struct {
	model.BookChapter
	Content string `json:"content"`
}

// ==================== 重新识别章节 ====================

// ReParseRequest 重新识别章节请求
type ReParseRequest struct {
	BookID uint64 `json:"bookId" binding:"required"`
}

// ReParseResponse 重新识别章节响应
type ReParseResponse struct {
	BookID     uint64 `json:"bookId"`
	BookTitle  string `json:"bookTitle"`
	OldCount   uint32 `json:"oldCount"`
	NewCount   uint32 `json:"newCount"`
	TotalWords uint32 `json:"totalWords"`
}

// ==================== 章节识别规则 ====================

// ChapterRuleRequest 规则请求
type ChapterRuleRequest struct {
	RuleName      string  `json:"ruleName" binding:"required,max=64"`
	ScopeType     string  `json:"scopeType" binding:"required,oneof=1 2"`
	BookID        *uint64 `json:"bookId"`
	Pattern       string  `json:"pattern" binding:"required,max=512"`
	TitleGroup    int     `json:"titleGroup"`
	MinChapterLen uint32  `json:"minChapterLen"`
	MaxChapterLen uint32  `json:"maxChapterLen"`
	Priority      int     `json:"priority"`
	Description   *string `json:"description"`
	Status        string  `json:"status"`
}

// ChapterRuleSearch 规则搜索
type ChapterRuleSearch struct {
	PageRequest
	RuleName  string  `json:"ruleName"`
	ScopeType string  `json:"scopeType"`
	BookID    *uint64 `json:"bookId"`
	Status    string  `json:"status"`
}

// ChapterRuleResponse 规则响应
type ChapterRuleResponse struct {
	model.BookChapterRule
}

// ==================== 内容净化规则 ====================

// FilterRuleRequest 过滤规则请求
type FilterRuleRequest struct {
	RuleName    string  `json:"ruleName" binding:"required,max=64"`
	MatchType   string  `json:"matchType" binding:"required,oneof=1 2"`
	Pattern     string  `json:"pattern" binding:"required,max=512"`
	Action      string  `json:"action" binding:"required,oneof=1 2 3"`
	Replacement string  `json:"replacement"`
	ApplyStage  string  `json:"applyStage" binding:"required,oneof=1 2"`
	Category    *string `json:"category"`
	Severity    string  `json:"severity"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}

// FilterRuleSearch 过滤规则搜索
type FilterRuleSearch struct {
	PageRequest
	RuleName   string `json:"ruleName"`
	ApplyStage string `json:"applyStage"`
	Category   string `json:"category"`
	Status     string `json:"status"`
}

// FilterRuleResponse 过滤规则响应
type FilterRuleResponse struct {
	model.BookContentFilterRule
}

// ==================== 扫描入库 ====================

// ScanResult 单文件扫描结果
type ScanResult struct {
	UploadID     uint64  `json:"uploadId"`
	OriginalName string  `json:"originalName"`
	FileSize     uint64  `json:"fileSize"`
	ParseStatus  string  `json:"parseStatus"`
	ParseMessage *string `json:"parseMessage"`
	BookID       *uint64 `json:"bookId"`
	ChapterCount *uint32 `json:"chapterCount"`
}

// ScanAllResponse 批量扫描响应
type ScanAllResponse struct {
	Results []ScanResult `json:"results"`
	Success int          `json:"success"`
	Failed  int          `json:"failed"`
}

// ScanPathRequest 扫描本地路径请求
type ScanPathRequest struct {
	Path string `json:"path" binding:"required"`
}

// ScanPathResponse 扫描路径响应
type ScanPathResponse struct {
	Total    int          `json:"total"`
	Imported int          `json:"imported"`
	Failed   int          `json:"failed"`
	Results  []ScanResult `json:"results"`
}
