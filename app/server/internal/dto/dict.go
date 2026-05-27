package dto

import "boread/internal/model"

// DictRequest 字典分类请求
type DictRequest struct {
	DictName string             `json:"dictName" binding:"required,max=64"`
	DictCode string             `json:"dictCode" binding:"required,max=64"`
	DictDesc string             `json:"dictDesc"`
	Status   model.EnableStatus `json:"status"`
}

// DictSearch 字典分页搜索
type DictSearch struct {
	PageRequest
	DictName string             `json:"dictName"`
	DictCode string             `json:"dictCode"`
	Status   model.EnableStatus `json:"status"`
}

// DictItemRequest 字典项请求
type DictItemRequest struct {
	DictID    uint64             `json:"dictId" binding:"required"`
	ItemLabel string             `json:"itemLabel" binding:"required,max=128"`
	ItemValue string             `json:"itemValue" binding:"required,max=128"`
	ItemDesc  string             `json:"itemDesc"`
	SortOrder int                `json:"sortOrder"`
	Status    model.EnableStatus `json:"status"`
}

// DictItemSearch 字典项查询
type DictItemSearch struct {
	DictID uint64 `form:"dictId"`
}
