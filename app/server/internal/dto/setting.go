package dto

import "boread/internal/model"

// SettingRequest 系统配置请求
type SettingRequest struct {
	Category    string             `json:"category" binding:"required,max=64"`
	Key         string             `json:"key" binding:"required,max=128"`
	Value       string             `json:"value" binding:"required"`
	ValueType   string             `json:"valueType" binding:"required,oneof=string number boolean json array"`
	Description string             `json:"description"`
	Editable    *bool              `json:"editable"`
	Status      model.EnableStatus `json:"status"`
}

// SettingUpdateRequest 编辑配置请求（key 不可改，category 不可改）
type SettingUpdateRequest struct {
	Value       string             `json:"value" binding:"required"`
	ValueType   string             `json:"valueType" binding:"required,oneof=string number boolean json array"`
	Description string             `json:"description"`
	Status      model.EnableStatus `json:"status"`
}

// SettingSearch 配置分页搜索
type SettingSearch struct {
	PageRequest
	Category string `json:"category"`
	Keyword  string `json:"keyword"` // 按 key 模糊
	Status   string `json:"status"`
}

// SettingVO 配置输出
type SettingVO struct {
	ID          uint64             `json:"id"`
	Category    string             `json:"category"`
	Key         string             `json:"key"`
	Value       string             `json:"value"`
	ValueType   string             `json:"valueType"`
	Description *string            `json:"description"`
	Editable    bool               `json:"editable"`
	IsSystem    bool               `json:"isSystem"`
	Status      model.EnableStatus `json:"status"`
	CreateTime  string             `json:"createTime"`
	UpdateTime  string             `json:"updateTime"`
}

// BatchSaveRequest 批量保存配置请求
type BatchSaveRequest struct {
	Category string          `json:"category" binding:"required,max=64"`
	Items    []BatchSaveItem `json:"items" binding:"required,dive"`
}

// BatchSaveItem 批量保存单项
type BatchSaveItem struct {
	Key       string `json:"key" binding:"required,max=128"`
	Value     string `json:"value"`
	ValueType string `json:"valueType" binding:"required,oneof=string number boolean json array"`
}

// BatchSaveResult 批量保存结果
type BatchSaveResult struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
}

// SettingCategoryItem 分类下的配置项
type SettingCategoryItem struct {
	Key         string  `json:"key"`
	Value       string  `json:"value"`
	ValueType   string  `json:"valueType"`
	Description *string `json:"description"`
	Editable    bool    `json:"editable"`
	IsSystem    bool    `json:"isSystem"`
}

// SettingCategoryGroup 按分类分组的配置
type SettingCategoryGroup struct {
	Category string                `json:"category"`
	Settings []SettingCategoryItem `json:"settings"`
}
