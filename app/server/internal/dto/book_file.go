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

// ==================== 笔记/划线 ====================

// NoteRequest 创建/更新笔记请求
type NoteRequest struct {
	BookID         uint64  `json:"bookId" binding:"required"`
	ChapterID      *uint64 `json:"chapterId"`
	NoteType       string  `json:"noteType" binding:"required,oneof=1 2 3"`
	SelectedText   *string `json:"selectedText"`
	StartOffset    *uint32 `json:"startOffset"`
	EndOffset      *uint32 `json:"endOffset"`
	HighlightColor *string `json:"highlightColor"`
	Content        *string `json:"content"`
	Visibility     string  `json:"visibility" binding:"oneof=1 2"`
}

// NoteResponse 笔记响应
type NoteResponse struct {
	model.ReaderNote
	ReaderName   *string `json:"readerName,omitempty"`
	BookTitle    *string `json:"bookTitle,omitempty"`
	ChapterTitle *string `json:"chapterTitle,omitempty"`
	LikeCount    int64   `json:"likeCount"`
	Liked        bool    `json:"liked"`
}

// NoteSearch 笔记搜索
type NoteSearch struct {
	PageRequest
	BookID    *uint64 `json:"bookId"`
	ChapterID *uint64 `json:"chapterId"`
	NoteType  string  `json:"noteType"`
}

// ==================== 书评 ====================

// ReviewRequest 创建/更新书评请求
type ReviewRequest struct {
	BookID  uint64  `json:"bookId" binding:"required"`
	Rating  *uint8  `json:"rating" binding:"omitempty,min=1,max=5"`
	Title   *string `json:"title"`
	Content string  `json:"content" binding:"required"`
}

// ReviewResponse 书评响应
type ReviewResponse struct {
	model.BookReview
	ReaderName *string `json:"readerName,omitempty"`
	BookTitle  *string `json:"bookTitle,omitempty"`
	LikeCount  int64   `json:"likeCount"`
	Liked      bool    `json:"liked"`
}

// ReviewSearch 书评搜索
type ReviewSearch struct {
	PageRequest
	BookID uint64 `json:"bookId"`
	Status string `json:"status"`
}

// ==================== 章节评论 ====================

// CommentRequest 创建章节评论请求
type CommentRequest struct {
	BookID    uint64  `json:"bookId" binding:"required"`
	ChapterID uint64  `json:"chapterId" binding:"required"`
	ParentID  uint64  `json:"parentId"`
	ReplyToID *uint64 `json:"replyToId"`
	Content   string  `json:"content" binding:"required"`
}

// CommentResponse 章节评论响应
type CommentResponse struct {
	model.BookChapterComment
	ReaderName  *string           `json:"readerName,omitempty"`
	ReplyToName *string           `json:"replyToName,omitempty"`
	LikeCount   int64             `json:"likeCount"`
	Liked       bool              `json:"liked"`
	Replies     []CommentResponse `json:"replies,omitempty"`
}

// CommentSearch 章节评论搜索
type CommentSearch struct {
	PageRequest
	ChapterID uint64 `json:"chapterId"`
	ParentID  uint64 `json:"parentId"`
}

// ==================== 点赞 ====================

// LikeRequest 点赞请求
type LikeRequest struct {
	TargetType string `json:"targetType" binding:"required,oneof=1 2 3"`
	TargetID   uint64 `json:"targetId" binding:"required"`
}

// LikeResponse 点赞响应
type LikeResponse struct {
	Liked bool  `json:"liked"`
	Count int64 `json:"count"`
}

// LikeStatusRequest 批量查询点赞状态请求
type LikeStatusRequest struct {
	Targets []LikeTarget `json:"targets" binding:"required,min=1"`
}

// LikeTarget 点赞目标
type LikeTarget struct {
	TargetType string `json:"targetType"`
	TargetID   uint64 `json:"targetId"`
}
