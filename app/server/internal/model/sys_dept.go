package model

// SysDept 部门表 (sys_dept)
type SysDept struct {
	BaseModel
	ParentID  uint64       `gorm:"column:parent_id;default:0;index" json:"parentId"`
	Ancestors string       `gorm:"column:ancestors;size:512;default:''" json:"ancestors"`
	DeptName  string       `gorm:"column:dept_name;size:64;not null" json:"deptName"`
	DeptCode  string       `gorm:"column:dept_code;size:64;not null" json:"deptCode"`
	Leader    *string      `gorm:"column:leader;size:64" json:"leader"`
	SortOrder int          `gorm:"column:sort_order;default:0" json:"sortOrder"`
	Status    EnableStatus `gorm:"column:status;type:char(1);default:'1'" json:"status"`
}

func (SysDept) TableName() string { return "sys_dept" }
