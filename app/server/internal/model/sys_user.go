package model

import "time"

// SysUser 后台用户表 (sys_user)
type SysUser struct {
	BaseModel
	DeptID        *uint64      `gorm:"column:dept_id;index" json:"deptId"`
	UserName      string       `gorm:"column:user_name;size:64;not null" json:"userName"`
	Password      string       `gorm:"column:password;size:128;not null" json:"-"`
	PwdUpdatedAt  *time.Time   `gorm:"column:pwd_updated_at" json:"pwdUpdatedAt"`
	PwdErrorCount uint16       `gorm:"column:pwd_error_count;default:0" json:"-"`
	LockedUntil   *time.Time   `gorm:"column:locked_until" json:"lockedUntil"`
	UserGender    *string      `gorm:"column:user_gender;type:char(1)" json:"userGender"`
	NickName      string       `gorm:"column:nick_name;size:64;default:''" json:"nickName"`
	UserPhone     *string      `gorm:"column:user_phone;size:20;index" json:"userPhone"`
	UserEmail     *string      `gorm:"column:user_email;size:128;index" json:"userEmail"`
	Avatar        *string      `gorm:"column:avatar;size:255" json:"avatar"`
	UgreenUserID  *string      `gorm:"column:ugreen_user_id;size:64;index" json:"ugreenUserId"`
	LastLoginTime *time.Time   `gorm:"column:last_login_time" json:"lastLoginTime"`
	LastLoginIP   *string      `gorm:"column:last_login_ip;size:64" json:"lastLoginIp"`
	Status        EnableStatus `gorm:"column:status;type:char(1);default:'1'" json:"status"`
	Version       uint64       `gorm:"column:version;default:0" json:"-"`
}

func (SysUser) TableName() string { return "sys_user" }

// SysUserRole 用户-角色关联 (sys_user_role)
type SysUserRole struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint64    `gorm:"column:user_id;not null" json:"userId"`
	RoleID     uint64    `gorm:"column:role_id;not null" json:"roleId"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli" json:"createTime"`
}

func (SysUserRole) TableName() string { return "sys_user_role" }
