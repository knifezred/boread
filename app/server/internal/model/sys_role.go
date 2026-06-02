package model

// DataScope 数据权限范围
type DataScope string

const (
	DataScopeAll        DataScope = "1" // 全部
	DataScopeCustom     DataScope = "2" // 自定义部门
	DataScopeDept       DataScope = "3" // 本部门
	DataScopeDeptAndSub DataScope = "4" // 本部门及子部门
	DataScopeSelf       DataScope = "5" // 仅本人
)

// SysRole 角色表 (sys_role)
type SysRole struct {
	BaseModel
	RoleName  string       `gorm:"column:role_name;size:64;not null" json:"roleName"`
	RoleCode  string       `gorm:"column:role_code;size:64;not null" json:"roleCode"`
	RoleDesc  *string      `gorm:"column:role_desc;size:255" json:"roleDesc"`
	DataScope DataScope    `gorm:"column:data_scope;type:char(1);default:'5'" json:"dataScope"`
	IsSystem  bool         `gorm:"column:is_system;default:0" json:"isSystem"`
	SortOrder int          `gorm:"column:sort_order;default:0" json:"sortOrder"`
	Status    EnableStatus `gorm:"column:status;type:char(1);default:'1'" json:"status"`
}

func (SysRole) TableName() string { return "sys_role" }

// SysRoleDept 角色-自定义部门关联 (sys_role_dept)
type SysRoleDept struct {
	ID     uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RoleID uint64 `gorm:"column:role_id;not null" json:"roleId"`
	DeptID uint64 `gorm:"column:dept_id;not null" json:"deptId"`
}

func (SysRoleDept) TableName() string { return "sys_role_dept" }
