package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// BaseModel 业务表通用字段
// 字段名与 SQL 列对齐 (sys_user / sys_role 等), 而非 GORM 默认命名
type BaseModel struct {
	ID         uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreateBy   *uint64        `gorm:"column:create_by" json:"createBy,omitempty"`
	CreateTime time.Time      `gorm:"column:create_time;autoCreateTime:milli" json:"createTime"`
	UpdateBy   *uint64        `gorm:"column:update_by" json:"updateBy,omitempty"`
	UpdateTime time.Time      `gorm:"column:update_time;autoUpdateTime:milli" json:"updateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// EnableStatus 通用启用/禁用枚举 (与前端 Common.EnableStatus 对齐)
type EnableStatus string

const (
	StatusEnabled  EnableStatus = "1"
	StatusDisabled EnableStatus = "2"
)

// JSONMap 用于 sys_menu.query 等 JSON 字段
type JSONMap map[string]any

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for JSONMap: %T", value)
	}
	return json.Unmarshal(bytes, j)
}
