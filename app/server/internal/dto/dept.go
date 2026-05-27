package dto

import "boread/internal/model"

// DeptRequest 创建/更新部门
type DeptRequest struct {
	ParentID  uint64             `json:"parentId"`
	DeptName  string             `json:"deptName" binding:"required,max=64"`
	DeptCode  string             `json:"deptCode" binding:"required,max=64"`
	Leader    string             `json:"leader"`
	SortOrder int                `json:"sortOrder"`
	Status    model.EnableStatus `json:"status"`
}

// DeptSearch 部门分页搜索 (对齐前端 Api.SystemManage.DeptSearchParams)
type DeptSearch struct {
	PageRequest
	DeptName string             `json:"deptName"`
	DeptCode string             `json:"deptCode"`
	Status   model.EnableStatus `json:"status"`
}

// DeptNode 部门树节点
type DeptNode struct {
	ID        uint64             `json:"id"`
	ParentID  uint64             `json:"parentId"`
	DeptName  string             `json:"deptName"`
	DeptCode  string             `json:"deptCode"`
	Leader    string             `json:"leader"`
	SortOrder int                `json:"sortOrder"`
	Status    model.EnableStatus `json:"status"`
	Children  []*DeptNode        `json:"children"`
}
