package dto

import "boread/internal/model"

// MenuRequest 创建/更新菜单
type MenuRequest struct {
	ParentID        uint64             `json:"parentId"`
	MenuType        model.MenuType     `json:"menuType" binding:"required,oneof=1 2"`
	MenuName        string             `json:"menuName" binding:"required,max=64"`
	RouteName       string             `json:"routeName" binding:"required,max=128"`
	RoutePath       string             `json:"routePath" binding:"required,max=255"`
	Component       string             `json:"component"`
	Icon            string             `json:"icon"`
	IconType        model.IconType     `json:"iconType"`
	I18nKey         string             `json:"i18nKey"`
	KeepAlive       bool               `json:"keepAlive"`
	Constant        bool               `json:"constant"`
	SortOrder       int                `json:"order"`
	Href            string             `json:"href"`
	HideInMenu      bool               `json:"hideInMenu"`
	ActiveMenu      string             `json:"activeMenu"`
	MultiTab        bool               `json:"multiTab"`
	FixedIndexInTab *int               `json:"fixedIndexInTab"`
	Query           any                `json:"query"`
	Status          model.EnableStatus `json:"status"`
	Buttons         []MenuButtonItem   `json:"buttons"`
}

// MenuButtonItem 菜单中的按钮项 (前端编辑菜单时内嵌提交)
type MenuButtonItem struct {
	ID         uint64 `json:"id"`
	ButtonCode string `json:"buttonCode"`
	ButtonDesc string `json:"buttonDesc"`
}

// MenuSearch 菜单分页搜索 (对齐前端 Api.SystemManage.MenuSearchParams)
type MenuSearch struct {
	PageRequest
	MenuName string             `json:"menuName"`
	Status   model.EnableStatus `json:"status"`
}

// MenuButtonRequest 菜单按钮请求
type MenuButtonRequest struct {
	MenuID     uint64 `json:"menuId" binding:"required"`
	ButtonCode string `json:"buttonCode" binding:"required,max=64"`
	ButtonDesc string `json:"buttonDesc"`
}

// MenuNode 菜单节点 (树/列表通用)
type MenuNode struct {
	model.SysMenu
	Children []*MenuNode            `json:"children"`
	Buttons  []*model.SysMenuButton `json:"buttons"`
}
