package dto

// ==================== 阅读统计 ====================

// ReadEventDailyResponse 每日阅读统计响应
type ReadEventDailyResponse struct {
	StatDate     string `json:"statDate"`
	DurationSec  uint64 `json:"durationSec"`
	WordCount    uint64 `json:"wordCount"`
	ChapterCount uint64 `json:"chapterCount"`
	BookCount    uint64 `json:"bookCount"`
	SessionCount uint64 `json:"sessionCount"`
}

// ReadEventBookResponse 每本书阅读统计响应
type ReadEventBookResponse struct {
	BookID       uint64 `json:"bookId"`
	BookTitle    string `json:"bookTitle"`
	DurationSec  uint64 `json:"durationSec"`
	WordCount    uint64 `json:"wordCount"`
	ChapterCount uint64 `json:"chapterCount"`
}

// ReadEventTotalResponse 总阅读统计响应
type ReadEventTotalResponse struct {
	DurationSec  uint64 `json:"durationSec"`
	WordCount    uint64 `json:"wordCount"`
	BookCount    uint64 `json:"bookCount"`
	SessionCount uint64 `json:"sessionCount"`
}

// ReadStatsQuery 阅读统计查询
type ReadStatsQuery struct {
	StartDate string  `json:"startDate"` // 起始日期 2006-01-02, 空表示不限
	EndDate   string  `json:"endDate"`   // 结束日期 2006-01-02, 空表示不限
	BookID    *uint64 `json:"bookId"`    // 指定书籍(可选), 用于查单书统计
}
