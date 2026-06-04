package dto

import "boread/internal/model"

// ==================== 章节 ====================

// ChapterSearch 章节搜索
type ChapterSearch struct {
	PageRequest
	BookID    *uint64 `json:"bookId"`
	FileID    *uint64 `json:"fileId"`
	ChapterNo *uint32 `json:"chapterNo"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
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

// ==================== 章节管理 ====================

// ChapterTitleUpdateRequest 单章标题更新请求
type ChapterTitleUpdateRequest struct {
	Title string `json:"title" binding:"required"`
}

// ChapterTitleBatchRequest 批量标题更新请求
type ChapterTitleBatchRequest struct {
	IDs   []uint64 `json:"ids" binding:"required,min=1"`
	Title string   `json:"title" binding:"required"`
}

// ChapterStatusBatchRequest 批量状态更新请求
type ChapterStatusBatchRequest struct {
	IDs    []uint64 `json:"ids" binding:"required,min=1"`
	Status string   `json:"status" binding:"required,oneof=1 2 3"`
}

// ChapterMergeRequest 合并章节请求
type ChapterMergeRequest struct {
	BookID    uint64   `json:"bookId" binding:"required"`
	TargetID  uint64   `json:"targetId" binding:"required"`
	SourceIDs []uint64 `json:"sourceIds" binding:"required,min=1"`
}

// ChapterFormatRequest 格式化章节编号请求
type ChapterFormatRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// ChapterContentSaveRequest 保存章节内容请求
type ChapterContentSaveRequest struct {
	BookID    uint64 `json:"bookId" binding:"required"`
	ChapterID uint64 `json:"chapterId" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
