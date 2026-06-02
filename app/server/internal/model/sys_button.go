package model

import (
	"time"

	"gorm.io/gorm"
)

// SysMenuButton 菜单按钮表 (sys_menu_button)
type SysMenuButton struct {
	ID         uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MenuID     uint64         `gorm:"column:menu_id;not null" json:"menuId"`
	ButtonCode string         `gorm:"column:button_code;size:64;not null" json:"buttonCode"`
	ButtonDesc *string        `gorm:"column:button_desc;size:255" json:"buttonDesc"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime:milli" json:"createTime"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime:milli" json:"updateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (SysMenuButton) TableName() string { return "sys_menu_button" }

// SysRoleMenu 角色-菜单关联 (sys_role_menu)
type SysRoleMenu struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID     uint64    `gorm:"column:role_id;not null" json:"roleId"`
	MenuID     uint64    `gorm:"column:menu_id;not null" json:"menuId"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli" json:"createTime"`
}

func (SysRoleMenu) TableName() string { return "sys_role_menu" }

// SysRoleButton 角色-按钮关联 (sys_role_button)
type SysRoleButton struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID     uint64    `gorm:"column:role_id;not null" json:"roleId"`
	ButtonID   uint64    `gorm:"column:button_id;not null" json:"buttonId"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime:milli" json:"createTime"`
}

func (SysRoleButton) TableName() string { return "sys_role_button" }
