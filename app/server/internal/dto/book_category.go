package dto

import "boread/internal/model"

type CategoryRequest struct {
	ParentID     uint64             `json:"parentId"`
	CategoryName string             `json:"categoryName" binding:"required,max=64"`
	CategoryCode string             `json:"categoryCode" binding:"required,max=64"`
	SortOrder    int                `json:"sortOrder"`
	Status       model.EnableStatus `json:"status"`
}

type CategoryNode struct {
	ID           uint64             `json:"id"`
	ParentID     uint64             `json:"parentId"`
	CategoryName string             `json:"categoryName"`
	CategoryCode string             `json:"categoryCode"`
	SortOrder    int                `json:"sortOrder"`
	Status       model.EnableStatus `json:"status"`
	Children     []*CategoryNode    `json:"children"`
}

type CategorySearch struct {
	PageRequest
	CategoryName string             `json:"categoryName"`
	CategoryCode string             `json:"categoryCode"`
	Status       model.EnableStatus `json:"status"`
}
