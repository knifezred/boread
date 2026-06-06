package model

// SysSetting 系统配置 (sys_setting)
type SysSetting struct {
	BaseModel
	Category    string       `gorm:"column:category;size:64;not null" json:"category"`
	Key         string       `gorm:"column:key;size:128;not null" json:"key"`
	Value       string       `gorm:"column:value;type:longtext;not null" json:"value"`
	ValueType   string       `gorm:"column:value_type;size:16;default:string" json:"valueType"`
	Description *string      `gorm:"column:description;size:255" json:"description"`
	Editable    bool         `gorm:"column:editable;default:1" json:"editable"`
	IsSystem    bool         `gorm:"column:is_system;default:0" json:"isSystem"`
	Status      EnableStatus `gorm:"column:status;type:char(1);default:'1'" json:"status"`
}

func (SysSetting) TableName() string { return "sys_setting" }
