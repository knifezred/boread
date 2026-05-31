package dto

import "boread/internal/model"

type CategoryRequest struct {
	ParentID     uint64             `json:"parentId"`
	CategoryName string             `json:"categoryName" binding:"required,max=64"`
	CategoryCode string             `json:"categoryCode" binding:"required,max=64"`
	Description  string             `json:"description"`
	SortOrder    int                `json:"sortOrder"`
	IsHot        *bool              `json:"isHot"`
	Status       model.EnableStatus `json:"status"`
}

type CategoryNode struct {
	ID           uint64             `json:"id"`
	ParentID     uint64             `json:"parentId"`
	CategoryName string             `json:"categoryName"`
	CategoryCode string             `json:"categoryCode"`
	Description  string             `json:"description"`
	SortOrder    int                `json:"sortOrder"`
	IsHot        bool               `json:"isHot"`
	Status       model.EnableStatus `json:"status"`
	Children     []*CategoryNode    `json:"children"`
}

type CategorySearch struct {
	PageRequest
	CategoryName string             `json:"categoryName"`
	CategoryCode string             `json:"categoryCode"`
	ParentID     uint64             `json:"parentId"`
	IsHot        *bool              `json:"isHot"`
	Status       model.EnableStatus `json:"status"`
}

// HotCategoryItem 热门分类响应项
type HotCategoryItem struct {
	ID           uint64 `json:"id"`
	CategoryName string `json:"categoryName"`
	CategoryCode string `json:"categoryCode"`
}
