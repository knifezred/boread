package service

import (
	"context"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

type MenuService struct {
	repo *repository.SysMenuRepository
}

func NewMenuService(repo *repository.SysMenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) Create(ctx context.Context, req *dto.MenuRequest, opID uint64) (*model.SysMenu, error) {
	if _, err := s.repo.GetByRouteName(ctx, req.RouteName); err == nil {
		return nil, code.ErrMenuRouteExists
	}
	status := req.Status
	if status == "" {
		status = model.StatusEnabled
	}
	iconType := req.IconType
	if iconType == "" {
		iconType = model.IconTypeIconify
	}
	m := &model.SysMenu{
		ParentID:        req.ParentID,
		MenuType:        req.MenuType,
		MenuName:        req.MenuName,
		RouteName:       req.RouteName,
		RoutePath:       req.RoutePath,
		IconType:        iconType,
		KeepAlive:       req.KeepAlive,
		Constant:        req.Constant,
		SortOrder:       req.SortOrder,
		HideInMenu:      req.HideInMenu,
		MultiTab:        req.MultiTab,
		FixedIndexInTab: req.FixedIndexInTab,
		Status:          status,
	}
	// 兼容处理query字段：只有是map类型才赋值，数组/空值赋值为nil
	if q, ok := req.Query.(map[string]any); ok && q != nil {
		m.Query = model.JSONMap(q)
	} else {
		m.Query = nil
	}
	if req.Component != "" {
		m.Component = &req.Component
	}
	if req.Icon != "" {
		m.Icon = &req.Icon
	}
	if req.I18nKey != "" {
		m.I18nKey = &req.I18nKey
	}
	if req.Href != "" {
		m.Href = &req.Href
	}
	if req.ActiveMenu != "" {
		m.ActiveMenu = &req.ActiveMenu
	}
	m.CreateBy = &opID
	m.UpdateBy = &opID
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MenuService) Update(ctx context.Context, id uint64, req *dto.MenuRequest, opID uint64) (*model.SysMenu, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m.IsSystem && req.RouteName != m.RouteName {
		return nil, code.ErrMenuSystem
	}
	if req.RouteName != m.RouteName {
		if _, e := s.repo.GetByRouteName(ctx, req.RouteName); e == nil {
			return nil, code.ErrMenuRouteExists
		}
	}
	m.ParentID = req.ParentID
	m.MenuType = req.MenuType
	m.MenuName = req.MenuName
	m.RouteName = req.RouteName
	m.RoutePath = req.RoutePath
	m.KeepAlive = req.KeepAlive
	m.Constant = req.Constant
	m.SortOrder = req.SortOrder
	m.HideInMenu = req.HideInMenu
	m.MultiTab = req.MultiTab
	m.FixedIndexInTab = req.FixedIndexInTab
	// 兼容处理query字段：只有是map类型才赋值，数组/空值赋值为nil
	if q, ok := req.Query.(map[string]any); ok && q != nil {
		m.Query = model.JSONMap(q)
	} else {
		m.Query = nil
	}
	if req.IconType != "" {
		m.IconType = req.IconType
	}
	if req.Status != "" {
		m.Status = req.Status
	}
	if req.Component != "" {
		m.Component = &req.Component
	}
	if req.Icon != "" {
		m.Icon = &req.Icon
	}
	if req.I18nKey != "" {
		m.I18nKey = &req.I18nKey
	}
	if req.Href != "" {
		m.Href = &req.Href
	}
	if req.ActiveMenu != "" {
		m.ActiveMenu = &req.ActiveMenu
	}
	m.UpdateBy = &opID
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	if req.Buttons != nil {
		if err := s.repo.ReplaceMenuButtons(ctx, id, req.Buttons); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (s *MenuService) Delete(ctx context.Context, id uint64) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if m.IsSystem {
		return code.ErrMenuSystem
	}
	if has, err := s.repo.HasChildren(ctx, id); err != nil {
		return err
	} else if has {
		return code.ErrMenuHasChildren
	}
	return s.repo.Delete(ctx, id)
}

func (s *MenuService) GetByID(ctx context.Context, id uint64) (*model.SysMenu, error) {
	return s.repo.GetByID(ctx, id)
}

// Tree 全量菜单树 (含按钮)
func (s *MenuService) Tree(ctx context.Context) ([]*dto.MenuNode, error) {
	menus, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	buttons, err := s.repo.ListAllButtons(ctx)
	if err != nil {
		return nil, err
	}
	buttonMap := make(map[uint64][]*model.SysMenuButton)
	for i := range buttons {
		b := buttons[i]
		buttonMap[b.MenuID] = append(buttonMap[b.MenuID], &b)
	}
	nodeMap := make(map[uint64]*dto.MenuNode, len(menus))
	for i := range menus {
		m := menus[i]
		btns := buttonMap[m.ID]
		if btns == nil {
			btns = []*model.SysMenuButton{}
		}
		nodeMap[m.ID] = &dto.MenuNode{SysMenu: m, Buttons: btns, Children: []*dto.MenuNode{}}
	}
	var roots []*dto.MenuNode
	// 按原始查询顺序遍历菜单（已排序），保证子节点顺序正确
	for _, m := range menus {
		node := nodeMap[m.ID]
		if node.ParentID == 0 {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}
	return roots, nil
}

// Page 菜单分页列表 (树形结构, 基于顶级菜单分页)
func (s *MenuService) Page(ctx context.Context, req *dto.MenuSearch) (*dto.PageResponse, error) {
	req.Normalize()
	// 1. 分页查询顶级菜单（已按 sort_order desc 排序）
	topMenus, total, err := s.repo.Page(ctx, req.MenuName, string(req.Status), req.Current, req.Size)
	if err != nil {
		return nil, err
	}
	if len(topMenus) == 0 {
		return dto.NewPageResponse([]*dto.MenuNode{}, total, &req.PageRequest), nil
	}
	// 2. 查询所有按钮按菜单分组
	buttons, err := s.repo.ListAllButtons(ctx)
	if err != nil {
		return nil, err
	}
	buttonMap := make(map[uint64][]*model.SysMenuButton)
	for i := range buttons {
		b := buttons[i]
		buttonMap[b.MenuID] = append(buttonMap[b.MenuID], &b)
	}
	// 3. 初始化节点map & 根节点列表（按查询顺序保持排序）
	nodeMap := make(map[uint64]*dto.MenuNode, len(topMenus))
	roots := make([]*dto.MenuNode, 0, len(topMenus))
	for _, m := range topMenus {
		btns := buttonMap[m.ID]
		if btns == nil {
			btns = []*model.SysMenuButton{}
		}
		node := &dto.MenuNode{
			SysMenu:  m,
			Buttons:  btns,
			Children: []*dto.MenuNode{},
		}
		nodeMap[m.ID] = node
		roots = append(roots, node)
	}
	// 4. 层级查询所有子菜单，按查询顺序直接挂载（子菜单查询已按 sort_order desc 排序）
	currentParentIDs := make([]uint64, 0, len(topMenus))
	for _, m := range topMenus {
		currentParentIDs = append(currentParentIDs, m.ID)
	}
	for i := 0; i < 10; i++ { // 限制最多10级，避免循环嵌套
		if len(currentParentIDs) == 0 {
			break
		}
		children, err := s.repo.ListByParentIDs(ctx, currentParentIDs) // 返回的子菜单已按 sort_order desc 排序
		if err != nil {
			return nil, err
		}
		if len(children) == 0 {
			break
		}
		// 按查询顺序挂载子节点，保持排序
		nextParentIDs := make([]uint64, 0, len(children))
		for _, m := range children {
			btns := buttonMap[m.ID]
			if btns == nil {
				btns = []*model.SysMenuButton{}
			}
			node := &dto.MenuNode{
				SysMenu:  m,
				Buttons:  btns,
				Children: []*dto.MenuNode{},
			}
			nodeMap[m.ID] = node
			// 挂载到父节点
			if parent, ok := nodeMap[m.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
			nextParentIDs = append(nextParentIDs, m.ID)
		}
		currentParentIDs = nextParentIDs
	}
	return dto.NewPageResponse(roots, total, &req.PageRequest), nil
}

// === 按钮 ===

func (s *MenuService) CreateButton(ctx context.Context, req *dto.MenuButtonRequest) (*model.SysMenuButton, error) {
	b := &model.SysMenuButton{
		MenuID:     req.MenuID,
		ButtonCode: req.ButtonCode,
	}
	if req.ButtonDesc != "" {
		b.ButtonDesc = &req.ButtonDesc
	}
	if err := s.repo.CreateButton(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *MenuService) DeleteButton(ctx context.Context, id uint64) error {
	return s.repo.DeleteButton(ctx, id)
}

func (s *MenuService) ListButtonsByMenu(ctx context.Context, menuID uint64) ([]model.SysMenuButton, error) {
	return s.repo.ListButtonsByMenu(ctx, menuID)
}
