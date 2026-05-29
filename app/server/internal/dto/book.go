package dto

import "boread/internal/model"

// BookRequest 新增/编辑书籍请求
type BookRequest struct {
	Title        string   `json:"title" binding:"required,max=255"`
	Author       string   `json:"author" binding:"max=128"`
	Cover        *string  `json:"cover"`
	Intro        *string  `json:"intro"`
	CategoryID   *uint64  `json:"categoryId"`
	Language     string   `json:"language" binding:"max=16"`
	SerialStatus string   `json:"serialStatus" binding:"oneof=1 2 3"`
	Visibility   string   `json:"visibility" binding:"oneof=1 2 3"`
	TagIDs       []uint64 `json:"tagIds"`
}

// BookUpdateStatusRequest 更新上架状态请求
type BookUpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=1 2 3 4"`
}

// BookSearch 书籍分页搜索
type BookSearch struct {
	PageRequest
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	CategoryID   *uint64  `json:"categoryId"`
	CategoryIDs  []uint64 // 展开后的分类ID列表（含自身+子分类），由 service 层填充
	Status       string   `json:"status"`
	Visibility   string   `json:"visibility"`
	SerialStatus string   `json:"serialStatus"`
	TagID        *uint64  `json:"tagId"`
}

// BookResponse 书籍响应 (包含标签列表和分类名称)
type BookResponse struct {
	model.Book
	TagIDs       []uint64   `json:"tagIds"`
	Tags         []TagBrief `json:"tags,omitempty"`
	CategoryName string     `json:"categoryName,omitempty"`
}

// TagBrief 标签简要信息
type TagBrief struct {
	ID      uint64 `json:"id"`
	TagName string `json:"tagName"`
}
