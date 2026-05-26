package dto

// PageRequest 通用分页+搜索请求
type PageRequest struct {
	Current int    `json:"current" form:"current"`
	Size    int    `json:"size" form:"size"`
	Keyword string `json:"keyword" form:"keyword"`
}

// Normalize 兜底默认值
func (p *PageRequest) Normalize() {
	if p.Current <= 0 {
		p.Current = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Size > 200 {
		p.Size = 200
	}
}

// Offset 分页偏移
func (p *PageRequest) Offset() int {
	return (p.Current - 1) * p.Size
}

// PageResponse 通用分页响应 (对齐前端 Common.PaginatingQueryRecord)
type PageResponse struct {
	Records any   `json:"records"`
	Current int   `json:"current"`
	Size    int   `json:"size"`
	Total   int64 `json:"total"`
}

// NewPageResponse 构造分页响应
func NewPageResponse(records any, total int64, req *PageRequest) *PageResponse {
	return &PageResponse{
		Records: records,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}
}

// IDRequest 通用 id 请求体
type IDRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

// IDsRequest 通用 ids 请求体 (批量操作)
type IDsRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}
