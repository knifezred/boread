package dto

import "boread/internal/model"

// UserCreateRequest 创建用户
type UserCreateRequest struct {
	DeptID     *uint64            `json:"deptId"`
	UserName   string             `json:"userName" binding:"required,min=3,max=64"`
	Password   string             `json:"password" binding:"required,min=6,max=64"`
	NickName   string             `json:"nickName"`
	UserGender *string            `json:"userGender" binding:"omitempty,oneof=1 2"`
	UserPhone  string             `json:"userPhone"`
	UserEmail  string             `json:"userEmail" binding:"omitempty,email"`
	Avatar     string             `json:"avatar"`
	Status     model.EnableStatus `json:"status"`
	RoleIDs    []uint64           `json:"roleIds"`
}

// UserUpdateRequest 更新用户 (不含密码)
type UserUpdateRequest struct {
	DeptID     *uint64            `json:"deptId"`
	NickName   string             `json:"nickName"`
	UserGender *string            `json:"userGender" binding:"omitempty,oneof=1 2"`
	UserPhone  string             `json:"userPhone"`
	UserEmail  string             `json:"userEmail" binding:"omitempty,email"`
	Avatar     string             `json:"avatar"`
	Status     model.EnableStatus `json:"status"`
	RoleIDs    []uint64           `json:"roleIds"`
}

// UserSearch 用户分页搜索
type UserSearch struct {
	PageRequest
	UserName   string             `json:"userName"`
	NickName   string             `json:"nickName"`
	UserPhone  string             `json:"userPhone"`
	UserEmail  string             `json:"userEmail"`
	UserGender string             `json:"userGender"`
	Status     model.EnableStatus `json:"status"`
}

// UserResetPwdRequest 重置密码
type UserResetPwdRequest struct {
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// UserVO 用户输出 (含 roleIds, 对齐前端 Api.SystemManage.User)
type UserVO struct {
	ID         uint64             `json:"id"`
	DeptID     *uint64            `json:"deptId"`
	UserName   string             `json:"userName"`
	NickName   string             `json:"nickName"`
	UserGender *string            `json:"userGender"`
	UserPhone  *string            `json:"userPhone"`
	UserEmail  *string            `json:"userEmail"`
	Avatar     *string            `json:"avatar"`
	Status     model.EnableStatus `json:"status"`
	UserRoles  []string           `json:"userRoles"`
}
