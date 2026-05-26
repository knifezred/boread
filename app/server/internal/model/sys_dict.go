package model

// SysDict 字典分类 (sys_dict)
type SysDict struct {
	BaseModel
	DictName string       `gorm:"column:dict_name;size:64;not null" json:"dictName"`
	DictCode string       `gorm:"column:dict_code;size:64;not null" json:"dictCode"`
	DictDesc *string      `gorm:"column:dict_desc;size:255" json:"dictDesc,omitempty"`
	IsSystem bool         `gorm:"column:is_system;default:0" json:"isSystem"`
	Status   EnableStatus `gorm:"column:status;type:char(1);default:'1'" json:"status"`
}

func (SysDict) TableName() string { return "sys_dict" }

// SysDictItem 字典项 (sys_dict_item)
type SysDictItem struct {
	BaseModel
	DictID    uint64       `gorm:"column:dict_id;not null;index" json:"dictId"`
	ItemLabel string       `gorm:"column:item_label;size:128;not null" json:"itemLabel"`
	ItemValue string       `gorm:"column:item_value;size:128;not null" json:"itemValue"`
	ItemDesc  *string      `gorm:"column:item_desc;size:255" json:"itemDesc,omitempty"`
	SortOrder int          `gorm:"column:sort_order;default:0" json:"sortOrder"`
	Status    EnableStatus `gorm:"column:status;type:char(1);default:'1'" json:"status"`
}

func (SysDictItem) TableName() string { return "sys_dict_item" }
