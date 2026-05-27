package dto

import "boread/internal/model"

// RoleRequest 创建/更新角色
type RoleRequest struct {
	RoleName  string             `json:"roleName" binding:"required,max=64"`
	RoleCode  string             `json:"roleCode" binding:"required,max=64"`
	RoleDesc  string             `json:"roleDesc"`
	DataScope model.DataScope    `json:"dataScope" binding:"required,oneof=1 2 3 4 5"`
	SortOrder int                `json:"sortOrder"`
	Status    model.EnableStatus `json:"status"`
	DeptIDs   []uint64           `json:"deptIds"` // data_scope=2 时的自定义部门
}

// RoleSearch 角色列表搜索 (对齐前端 Api.SystemManage.RoleSearchParams)
type RoleSearch struct {
	PageRequest
	RoleName string             `json:"roleName"`
	RoleCode string             `json:"roleCode"`
	Status   model.EnableStatus `json:"status"`
}

// RoleMenuRequest 授权菜单
type RoleMenuRequest struct {
	MenuIDs []uint64 `json:"menuIds"`
}

// RoleButtonRequest 授权按钮
type RoleButtonRequest struct {
	ButtonIDs []uint64 `json:"buttonIds"`
}

// RoleBrief 全量角色 (下拉用, 对齐前端 Api.SystemManage.AllRole)
type RoleBrief struct {
	ID       uint64 `json:"id"`
	RoleName string `json:"roleName"`
	RoleCode string `json:"roleCode"`
}
