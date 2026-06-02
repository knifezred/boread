package dto

import "boread/internal/model"

// ==================== 书架 ====================

// BookshelfRequest 添加到书架请求
type BookshelfRequest struct {
	BookID    uint64 `json:"bookId" binding:"required"`
	GroupName string `json:"groupName"` // 可选，默认"默认"
}

// BookshelfUpdateRequest 更新书架请求
type BookshelfUpdateRequest struct {
	GroupName *string `json:"groupName"` // 修改分组
	IsTop     *bool   `json:"isTop"`     // 修改置顶状态
}

// BookshelfSearch 书架分页搜索
type BookshelfSearch struct {
	PageRequest
	GroupName string `json:"groupName"` // 按分组筛选
	Keyword   string `json:"keyword"`   // 按书名关键词搜索
}

// BookshelfResponse 书架响应 (包含书籍基本信息+阅读进度)
type BookshelfResponse struct {
	model.ReaderBookshelf
	BookTitle     string  `json:"bookTitle"`
	BookAuthor    string  `json:"bookAuthor"`
	BookCover     *string `json:"bookCover"`
	TotalChapters uint32  `json:"totalChapters"`
	TotalWords    uint32  `json:"totalWords"`
	ChapterID     *uint64 `json:"chapterId"`
	ChapterNo     *uint32 `json:"chapterNo"`
	Position      *uint32 `json:"position"`
	ReadPercent   float64 `json:"readPercent"`
	ReadDuration  uint32  `json:"readDuration"`
	LastReadTime  *string `json:"lastReadTime"`
}

// BookshelfGroupItem 分组项
type BookshelfGroupItem struct {
	GroupName string `json:"groupName"`
	BookCount int64  `json:"bookCount"`
}

// ==================== 阅读进度 ====================

// ReadProgressRequest 上报阅读进度请求
type ReadProgressRequest struct {
	FileID       *uint64 `json:"fileId"` // 文件切换时传入
	ChapterID    uint64  `json:"chapterId" binding:"required"`
	ChapterNo    uint32  `json:"chapterNo" binding:"required"`
	Position     uint32  `json:"position"`
	Percent      float64 `json:"percent"`
	ReadDuration uint32  `json:"readDuration"` // 本次增加的阅读时长(秒)
}

// ReadProgressResponse 阅读进度响应
type ReadProgressResponse struct {
	model.ReaderReadProgress
}

// BookshelfPageResponse 书架分页响应
type BookshelfPageResponse struct {
	Records []BookshelfResponse `json:"records"`
	Current int                 `json:"current"`
	Size    int                 `json:"size"`
	Total   int64               `json:"total"`
}
