package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

type SysRoleRepository struct {
	db *gorm.DB
}

func NewSysRoleRepository(db *gorm.DB) *SysRoleRepository {
	return &SysRoleRepository{db: db}
}

func (r *SysRoleRepository) Create(ctx context.Context, m *model.SysRole) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SysRoleRepository) Update(ctx context.Context, m *model.SysRole) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *SysRoleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysRole{}, id).Error
}

func (r *SysRoleRepository) GetByID(ctx context.Context, id uint64) (*model.SysRole, error) {
	var m model.SysRole
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysRoleRepository) GetByCode(ctx context.Context, code string) (*model.SysRole, error) {
	var m model.SysRole
	if err := r.db.WithContext(ctx).Where("role_code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysRoleRepository) Page(ctx context.Context, s *dto.RoleSearch) ([]model.SysRole, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.SysRole{})
	if s.RoleName != "" {
		tx = tx.Where("role_name LIKE ?", "%"+s.RoleName+"%")
	}
	if s.RoleCode != "" {
		tx = tx.Where("role_code LIKE ?", "%"+s.RoleCode+"%")
	}
	if s.Status != "" {
		tx = tx.Where("status = ?", s.Status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.SysRole
	if err := tx.Order("sort_order ASC, id ASC").
		Offset(s.Offset()).Limit(s.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *SysRoleRepository) AllBrief(ctx context.Context) ([]model.SysRole, error) {
	var rows []model.SysRole
	err := r.db.WithContext(ctx).Where("status = '1'").Order("sort_order ASC").Find(&rows).Error
	return rows, err
}

// ReplaceMenus 替换角色的菜单授权 (事务)
func (r *SysRoleRepository) ReplaceMenus(ctx context.Context, roleID uint64, menuIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		rows := make([]model.SysRoleMenu, 0, len(menuIDs))
		for _, mid := range menuIDs {
			rows = append(rows, model.SysRoleMenu{RoleID: roleID, MenuID: mid})
		}
		return tx.Create(&rows).Error
	})
}

// ReplaceButtons 替换角色的按钮授权
func (r *SysRoleRepository) ReplaceButtons(ctx context.Context, roleID uint64, buttonIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRoleButton{}).Error; err != nil {
			return err
		}
		if len(buttonIDs) == 0 {
			return nil
		}
		rows := make([]model.SysRoleButton, 0, len(buttonIDs))
		for _, bid := range buttonIDs {
			rows = append(rows, model.SysRoleButton{RoleID: roleID, ButtonID: bid})
		}
		return tx.Create(&rows).Error
	})
}

// ReplaceDepts 替换角色数据权限部门 (data_scope=2)
func (r *SysRoleRepository) ReplaceDepts(ctx context.Context, roleID uint64, deptIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRoleDept{}).Error; err != nil {
			return err
		}
		if len(deptIDs) == 0 {
			return nil
		}
		rows := make([]model.SysRoleDept, 0, len(deptIDs))
		for _, did := range deptIDs {
			rows = append(rows, model.SysRoleDept{RoleID: roleID, DeptID: did})
		}
		return tx.Create(&rows).Error
	})
}

// GetMenuIDs 角色已授权的菜单 ids
func (r *SysRoleRepository) GetMenuIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).Model(&model.SysRoleMenu{}).
		Where("role_id = ?", roleID).Pluck("menu_id", &ids).Error
	return ids, err
}

// GetButtonIDs 角色已授权的按钮 ids
func (r *SysRoleRepository) GetButtonIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).Model(&model.SysRoleButton{}).
		Where("role_id = ?", roleID).Pluck("button_id", &ids).Error
	return ids, err
}

// HasUsers 是否还有用户挂在这个角色上
func (r *SysRoleRepository) HasUsers(ctx context.Context, roleID uint64) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.SysUserRole{}).
		Where("role_id = ?", roleID).Count(&n).Error
	return n > 0, err
}
