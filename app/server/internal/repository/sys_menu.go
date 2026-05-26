package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/model"
)

type SysMenuRepository struct {
	db *gorm.DB
}

func NewSysMenuRepository(db *gorm.DB) *SysMenuRepository {
	return &SysMenuRepository{db: db}
}

func (r *SysMenuRepository) Create(ctx context.Context, m *model.SysMenu) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SysMenuRepository) Update(ctx context.Context, m *model.SysMenu) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *SysMenuRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysMenu{}, id).Error
}

func (r *SysMenuRepository) GetByID(ctx context.Context, id uint64) (*model.SysMenu, error) {
	var m model.SysMenu
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysMenuRepository) GetByRouteName(ctx context.Context, routeName string) (*model.SysMenu, error) {
	var m model.SysMenu
	if err := r.db.WithContext(ctx).Where("route_name = ?", routeName).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysMenuRepository) ListAll(ctx context.Context) ([]model.SysMenu, error) {
	var rows []model.SysMenu
	err := r.db.WithContext(ctx).Order("parent_id ASC, sort_order ASC, id ASC").Find(&rows).Error
	return rows, err
}

// HasChildren 检查是否有子菜单
func (r *SysMenuRepository) HasChildren(ctx context.Context, id uint64) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.SysMenu{}).Where("parent_id = ?", id).Count(&n).Error
	return n > 0, err
}

// === 按钮 ===

func (r *SysMenuRepository) CreateButton(ctx context.Context, b *model.SysMenuButton) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *SysMenuRepository) DeleteButton(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysMenuButton{}, id).Error
}

func (r *SysMenuRepository) ListButtonsByMenu(ctx context.Context, menuID uint64) ([]model.SysMenuButton, error) {
	var rows []model.SysMenuButton
	err := r.db.WithContext(ctx).Where("menu_id = ?", menuID).Find(&rows).Error
	return rows, err
}

// ListAllButtons 获取所有按钮 (用于全量菜单树展示)
func (r *SysMenuRepository) ListAllButtons(ctx context.Context) ([]model.SysMenuButton, error) {
	var rows []model.SysMenuButton
	err := r.db.WithContext(ctx).Find(&rows).Error
	return rows, err
}

// Page 菜单分页查询 (只查顶级菜单 parent_id=0)
func (r *SysMenuRepository) Page(ctx context.Context, name string, status string, page, size int) ([]model.SysMenu, int64, error) {
	var rows []model.SysMenu
	var total int64
	db := r.db.WithContext(ctx).Model(&model.SysMenu{}).Where("parent_id = 0") // 分页基于顶级菜单
	if name != "" {
		db = db.Where("menu_name like ?", "%"+name+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Offset((page - 1) * size).Limit(size).Order("sort_order asc, id asc").Find(&rows).Error
	return rows, total, err
}

// ListByParentIDs 批量查询父ID下的所有子菜单
func (r *SysMenuRepository) ListByParentIDs(ctx context.Context, parentIDs []uint64) ([]model.SysMenu, error) {
	var rows []model.SysMenu
	err := r.db.WithContext(ctx).Where("parent_id IN ?", parentIDs).Order("sort_order asc, id asc").Find(&rows).Error
	return rows, err
}

// DB 获取DB实例，供service扩展查询用
func (r *SysMenuRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}
