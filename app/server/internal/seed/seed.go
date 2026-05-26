package seed

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"boread/internal/model"
)

// Run 幂等初始化系统数据
// 多次执行不会产生重复数据 (按 dept_code/role_code/user_name/route_name 判存)
func Run(ctx context.Context, db *gorm.DB) error {
	dept, err := upsertRootDept(ctx, db)
	if err != nil {
		return fmt.Errorf("seed root dept: %w", err)
	}

	role, err := upsertSuperRole(ctx, db)
	if err != nil {
		return fmt.Errorf("seed super role: %w", err)
	}

	admin, err := upsertAdminUser(ctx, db, dept.ID)
	if err != nil {
		return fmt.Errorf("seed admin user: %w", err)
	}

	if err := ensureUserRole(ctx, db, admin.ID, role.ID); err != nil {
		return fmt.Errorf("seed user role: %w", err)
	}

	menus, buttons, err := upsertMenus(ctx, db)
	if err != nil {
		return fmt.Errorf("seed menus: %w", err)
	}

	if err := ensureRoleMenus(ctx, db, role.ID, menus); err != nil {
		return fmt.Errorf("seed role menus: %w", err)
	}

	if err := ensureRoleButtons(ctx, db, role.ID, buttons); err != nil {
		return fmt.Errorf("seed role buttons: %w", err)
	}

	return nil
}

func upsertRootDept(ctx context.Context, db *gorm.DB) (*model.SysDept, error) {
	var d model.SysDept
	err := db.WithContext(ctx).Where("dept_code = ?", "ROOT").First(&d).Error
	if err == nil {
		return &d, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	d = model.SysDept{
		ParentID:  0,
		Ancestors: "0",
		DeptName:  "总部",
		DeptCode:  "ROOT",
		SortOrder: 0,
		Status:    model.StatusEnabled,
	}
	if err := db.WithContext(ctx).Create(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func upsertSuperRole(ctx context.Context, db *gorm.DB) (*model.SysRole, error) {
	var r model.SysRole
	err := db.WithContext(ctx).Where("role_code = ?", "SUPER_ADMIN").First(&r).Error
	if err == nil {
		return &r, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	desc := "系统超级管理员, 拥有全部权限"
	r = model.SysRole{
		RoleName:  "超级管理员",
		RoleCode:  "SUPER_ADMIN",
		RoleDesc:  &desc,
		DataScope: model.DataScopeAll,
		IsSystem:  true,
		SortOrder: 0,
		Status:    model.StatusEnabled,
	}
	if err := db.WithContext(ctx).Create(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func upsertAdminUser(ctx context.Context, db *gorm.DB, deptID uint64) (*model.SysUser, error) {
	var u model.SysUser
	err := db.WithContext(ctx).Where("user_name = ?", "admin").First(&u).Error
	if err == nil {
		return &u, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u = model.SysUser{
		DeptID:   &deptID,
		UserName: "admin",
		Password: string(hashed),
		NickName: "超级管理员",
		Status:   model.StatusEnabled,
	}
	if err := db.WithContext(ctx).Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func ensureUserRole(ctx context.Context, db *gorm.DB, userID, roleID uint64) error {
	var count int64
	if err := db.WithContext(ctx).Model(&model.SysUserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return db.WithContext(ctx).Create(&model.SysUserRole{UserID: userID, RoleID: roleID}).Error
}

// menuSpec 描述一个菜单 + 它的按钮
type menuSpec struct {
	RouteName string
	MenuName  string
	RoutePath string
	Component string
	Icon      string
	I18nKey   string
	MenuType  model.MenuType
	SortOrder int
	Buttons   []buttonSpec
	Children  []menuSpec
}

type buttonSpec struct {
	Code string
	Desc string
}

// menuTree 完整菜单种子定义
func menuTree() []menuSpec {
	return []menuSpec{
		{
			RouteName: "home",
			MenuName:  "首页",
			RoutePath: "/home",
			Component: "layout.base$view.home",
			Icon:      "mdi:home",
			I18nKey:   "route.home",
			MenuType:  model.MenuTypeMenu,
			SortOrder: 1,
		},
		{
			RouteName: "manage",
			MenuName:  "系统管理",
			RoutePath: "/manage",
			Component: "layout.base",
			Icon:      "carbon:cloud-service-management",
			I18nKey:   "route.manage",
			MenuType:  model.MenuTypeDir,
			SortOrder: 9,
			Children: []menuSpec{
				{
					RouteName: "manage_dept",
					MenuName:  "部门管理",
					RoutePath: "/manage/dept",
					Component: "view.manage_dept",
					Icon:      "mingcute:department-line",
					I18nKey:   "route.manage_dept",
					MenuType:  model.MenuTypeMenu,
					SortOrder: 1,
					Buttons: []buttonSpec{
						{Code: "dept:create", Desc: "新增部门"},
						{Code: "dept:update", Desc: "编辑部门"},
						{Code: "dept:delete", Desc: "删除部门"},
					},
				},
				{
					RouteName: "manage_menu",
					MenuName:  "菜单管理",
					RoutePath: "/manage/menu",
					Component: "view.manage_menu",
					Icon:      "material-symbols:route",
					I18nKey:   "route.manage_menu",
					MenuType:  model.MenuTypeMenu,
					SortOrder: 2,
					Buttons: []buttonSpec{
						{Code: "menu:create", Desc: "新增菜单"},
						{Code: "menu:update", Desc: "编辑菜单"},
						{Code: "menu:delete", Desc: "删除菜单"},
					},
				},
				{
					RouteName: "manage_role",
					MenuName:  "角色管理",
					RoutePath: "/manage/role",
					Component: "view.manage_role",
					Icon:      "carbon:user-role",
					I18nKey:   "route.manage_role",
					MenuType:  model.MenuTypeMenu,
					SortOrder: 3,
					Buttons: []buttonSpec{
						{Code: "role:create", Desc: "新增角色"},
						{Code: "role:update", Desc: "编辑角色"},
						{Code: "role:delete", Desc: "删除角色"},
						{Code: "role:grant", Desc: "授权菜单/按钮"},
					},
				},
				{
					RouteName: "manage_user",
					MenuName:  "用户管理",
					RoutePath: "/manage/user",
					Component: "view.manage_user",
					Icon:      "ic:round-manage-accounts",
					I18nKey:   "route.manage_user",
					MenuType:  model.MenuTypeMenu,
					SortOrder: 4,
					Buttons: []buttonSpec{
						{Code: "user:create", Desc: "新增用户"},
						{Code: "user:update", Desc: "编辑用户"},
						{Code: "user:delete", Desc: "删除用户"},
						{Code: "user:reset_pwd", Desc: "重置密码"},
					},
				},
				{
					RouteName: "manage_dict",
					MenuName:  "字典管理",
					RoutePath: "/manage/dict",
					Component: "view.manage_dict",
					Icon:      "mdi:book-cog",
					I18nKey:   "route.manage_dict",
					MenuType:  model.MenuTypeMenu,
					SortOrder: 5,
					Buttons: []buttonSpec{
						{Code: "dict:create", Desc: "新增字典"},
						{Code: "dict:update", Desc: "编辑字典"},
						{Code: "dict:delete", Desc: "删除字典"},
					},
				},
				{
					RouteName: "manage_log",
					MenuName:  "日志管理",
					RoutePath: "/manage/log",
					Component: "view.manage_log",
					Icon:      "carbon:cloud-logging",
					I18nKey:   "route.manage_log",
					MenuType:  model.MenuTypeMenu,
					SortOrder: 6,
					Buttons: []buttonSpec{
						{Code: "log:export", Desc: "导出日志"},
					},
				},
			},
		},
	}
}

// upsertMenus 写菜单树, 返回 (所有菜单 id, 所有按钮 id)
func upsertMenus(ctx context.Context, db *gorm.DB) ([]uint64, []uint64, error) {
	var allMenus []uint64
	var allButtons []uint64

	var walk func(parentID uint64, specs []menuSpec) error
	walk = func(parentID uint64, specs []menuSpec) error {
		for _, s := range specs {
			menu, err := upsertMenu(ctx, db, parentID, s)
			if err != nil {
				return err
			}
			allMenus = append(allMenus, menu.ID)

			for _, b := range s.Buttons {
				btn, err := upsertButton(ctx, db, menu.ID, b)
				if err != nil {
					return err
				}
				allButtons = append(allButtons, btn.ID)
			}
			if err := walk(menu.ID, s.Children); err != nil {
				return err
			}
		}
		return nil
	}

	if err := walk(0, menuTree()); err != nil {
		return nil, nil, err
	}
	return allMenus, allButtons, nil
}

func upsertMenu(ctx context.Context, db *gorm.DB, parentID uint64, s menuSpec) (*model.SysMenu, error) {
	var m model.SysMenu
	err := db.WithContext(ctx).Where("route_name = ?", s.RouteName).First(&m).Error
	if err == nil {
		return &m, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	component := s.Component
	icon := s.Icon
	i18n := s.I18nKey
	m = model.SysMenu{
		ParentID:  parentID,
		MenuType:  s.MenuType,
		MenuName:  s.MenuName,
		RouteName: s.RouteName,
		RoutePath: s.RoutePath,
		Component: &component,
		Icon:      &icon,
		IconType:  model.IconTypeIconify,
		I18nKey:   &i18n,
		SortOrder: s.SortOrder,
		IsSystem:  true,
		Status:    model.StatusEnabled,
	}
	if err := db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func upsertButton(ctx context.Context, db *gorm.DB, menuID uint64, b buttonSpec) (*model.SysMenuButton, error) {
	var btn model.SysMenuButton
	err := db.WithContext(ctx).
		Where("menu_id = ? AND button_code = ?", menuID, b.Code).
		First(&btn).Error
	if err == nil {
		return &btn, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	desc := b.Desc
	btn = model.SysMenuButton{
		MenuID:     menuID,
		ButtonCode: b.Code,
		ButtonDesc: &desc,
	}
	if err := db.WithContext(ctx).Create(&btn).Error; err != nil {
		return nil, err
	}
	return &btn, nil
}

func ensureRoleMenus(ctx context.Context, db *gorm.DB, roleID uint64, menuIDs []uint64) error {
	for _, mid := range menuIDs {
		var count int64
		if err := db.WithContext(ctx).Model(&model.SysRoleMenu{}).
			Where("role_id = ? AND menu_id = ?", roleID, mid).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if err := db.WithContext(ctx).Create(&model.SysRoleMenu{RoleID: roleID, MenuID: mid}).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureRoleButtons(ctx context.Context, db *gorm.DB, roleID uint64, buttonIDs []uint64) error {
	for _, bid := range buttonIDs {
		var count int64
		if err := db.WithContext(ctx).Model(&model.SysRoleButton{}).
			Where("role_id = ? AND button_id = ?", roleID, bid).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if err := db.WithContext(ctx).Create(&model.SysRoleButton{RoleID: roleID, ButtonID: bid}).Error; err != nil {
			return err
		}
	}
	return nil
}
