package model

// MenuType 菜单类型
type MenuType string

const (
	MenuTypeDir  MenuType = "1" // 目录
	MenuTypeMenu MenuType = "2" // 菜单
)

// IconType 图标类型
type IconType string

const (
	IconTypeIconify IconType = "1"
	IconTypeLocal   IconType = "2"
)

// SysMenu 菜单表 (sys_menu)
type SysMenu struct {
	BaseModel
	ParentID        uint64       `gorm:"column:parent_id;default:0;index" json:"parentId"`
	MenuType        MenuType     `gorm:"column:menu_type;type:char(1);not null" json:"menuType"`
	MenuName        string       `gorm:"column:menu_name;size:64;not null" json:"menuName"`
	RouteName       string       `gorm:"column:route_name;size:128;not null" json:"routeName"`
	RoutePath       string       `gorm:"column:route_path;size:255;not null" json:"routePath"`
	Component       *string      `gorm:"column:component;size:255" json:"component"`
	Icon            *string      `gorm:"column:icon;size:128" json:"icon"`
	IconType        IconType     `gorm:"column:icon_type;type:char(1);default:'1'" json:"iconType"`
	I18nKey         *string      `gorm:"column:i18n_key;size:255" json:"i18nKey"`
	KeepAlive       bool         `gorm:"column:keep_alive;default:0" json:"keepAlive"`
	Constant        bool         `gorm:"column:constant;default:0" json:"constant"`
	SortOrder       int          `gorm:"column:sort_order;default:0" json:"order"`
	Href            *string      `gorm:"column:href;size:255" json:"href"`
	HideInMenu      bool         `gorm:"column:hide_in_menu;default:0" json:"hideInMenu"`
	ActiveMenu      *string      `gorm:"column:active_menu;size:128" json:"activeMenu"`
	MultiTab        bool         `gorm:"column:multi_tab;default:0" json:"multiTab"`
	FixedIndexInTab *int         `gorm:"column:fixed_index_in_tab" json:"fixedIndexInTab"`
	Query           JSONMap      `gorm:"column:query;type:json" json:"query"`
	IsSystem        bool         `gorm:"column:is_system;default:0" json:"isSystem"`
	Status          EnableStatus `gorm:"column:status;type:char(1);default:'1'" json:"status"`
}

func (SysMenu) TableName() string { return "sys_menu" }
