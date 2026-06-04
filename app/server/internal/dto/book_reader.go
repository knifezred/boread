package dto

import "boread/internal/model"

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

// ==================== 阅读事件 ====================

// ReadEventRequest 上报阅读事件请求
type ReadEventRequest struct {
	BookID      uint64 `json:"bookId" binding:"required"`
	ChapterID   uint64 `json:"chapterId" binding:"required"`
	SessionID   string `json:"sessionId"`                      // 服务端生成 UUID, 客户端可传空
	DurationSec uint32 `json:"durationSec" binding:"required"` // 本次心跳区间阅读时长(秒)
	WordCount   uint32 `json:"wordCount"`                      // 本次心跳区间滚动字数
	DeviceType  string `json:"deviceType"`                     // 默认 web
}
