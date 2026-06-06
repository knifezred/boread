package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

// SysUserRepository 后台用户数据访问层
type SysUserRepository struct {
	db *gorm.DB
}

func NewSysUserRepository(db *gorm.DB) *SysUserRepository {
	return &SysUserRepository{db: db}
}

// GetByUserName 按用户名查 (含已禁用, 业务层判断状态)
func (r *SysUserRepository) GetByUserName(ctx context.Context, userName string) (*model.SysUser, error) {
	var u model.SysUser
	if err := r.db.WithContext(ctx).Where("user_name = ?", userName).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByID 按 id 查
func (r *SysUserRepository) GetByID(ctx context.Context, id uint64) (*model.SysUser, error) {
	var u model.SysUser
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateLoginSuccess 登录成功: 清错误计数 + 记录登录信息
func (r *SysUserRepository) UpdateLoginSuccess(ctx context.Context, id uint64, ip string, t any) error {
	return r.db.WithContext(ctx).Model(&model.SysUser{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"pwd_error_count": 0,
			"locked_until":    nil,
			"last_login_ip":   ip,
			"last_login_time": t,
		}).Error
}

// IncrErrorCount 密码错误自增 (并发安全)
func (r *SysUserRepository) IncrErrorCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.SysUser{}).
		Where("id = ?", id).
		UpdateColumn("pwd_error_count", gorm.Expr("pwd_error_count + 1")).
		Error
}

// LockUser 锁定账号到指定时间
func (r *SysUserRepository) LockUser(ctx context.Context, id uint64, until any) error {
	return r.db.WithContext(ctx).Model(&model.SysUser{}).
		Where("id = ?", id).
		Update("locked_until", until).Error
}

// GetRoleCodesByUserID 查用户的所有启用角色 code
func (r *SysUserRepository) GetRoleCodesByUserID(ctx context.Context, userID uint64) ([]string, error) {
	var codes []string
	err := r.db.WithContext(ctx).
		Table("sys_user_role AS ur").
		Select("r.role_code").
		Joins("JOIN sys_role AS r ON r.id = ur.role_id AND r.deleted_at IS NULL AND r.status = '1'").
		Where("ur.user_id = ?", userID).
		Scan(&codes).Error
	return codes, err
}

// GetRoleIDsByUserID 查用户的所有启用角色 id
func (r *SysUserRepository) GetRoleIDsByUserID(ctx context.Context, userID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).
		Table("sys_user_role AS ur").
		Select("r.id").
		Joins("JOIN sys_role AS r ON r.id = ur.role_id AND r.deleted_at IS NULL AND r.status = '1'").
		Where("ur.user_id = ?", userID).
		Scan(&ids).Error
	return ids, err
}

// GetButtonCodesByRoleIDs 查角色集合下的所有按钮 code (去重)
func (r *SysUserRepository) GetButtonCodesByRoleIDs(ctx context.Context, roleIDs []uint64) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}
	var codes []string
	err := r.db.WithContext(ctx).
		Table("sys_role_button AS rb").
		Distinct("b.button_code").
		Joins("JOIN sys_menu_button AS b ON b.id = rb.button_id AND b.deleted_at IS NULL").
		Where("rb.role_id IN ?", roleIDs).
		Scan(&codes).Error
	return codes, err
}

// GetMenusByRoleIDs 查角色集合下所有可见菜单 (去重, 按 sort_order)
func (r *SysUserRepository) GetMenusByRoleIDs(ctx context.Context, roleIDs []uint64) ([]model.SysMenu, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var menus []model.SysMenu
	err := r.db.WithContext(ctx).
		Table("sys_menu AS m").
		Distinct("m.*").
		Joins("JOIN sys_role_menu AS rm ON rm.menu_id = m.id").
		Where("rm.role_id IN ?", roleIDs).
		Where("m.status = '1' AND m.deleted_at IS NULL").
		Order("m.parent_id ASC, m.sort_order ASC").
		Scan(&menus).Error
	return menus, err
}

// GetConstantMenus 获取所有常量路由 (无需鉴权, 全员可见)
func (r *SysUserRepository) GetConstantMenus(ctx context.Context) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.WithContext(ctx).
		Where("constant = 1 AND status = '1'").
		Order("parent_id ASC, sort_order ASC").
		Find(&menus).Error
	return menus, err
}

// === CRUD 扩展 ===

func (r *SysUserRepository) Create(ctx context.Context, m *model.SysUser) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SysUserRepository) Update(ctx context.Context, m *model.SysUser) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *SysUserRepository) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return r.db.WithContext(ctx).Model(&model.SysUser{}).Where("id = ?", id).Updates(fields).Error
}

func (r *SysUserRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysUser{}, id).Error
}

func (r *SysUserRepository) BatchDelete(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Delete(&model.SysUser{}, "id IN ?", ids).Error
}

func (r *SysUserRepository) Page(ctx context.Context, s *dto.UserSearch) ([]model.SysUser, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.SysUser{})
	if s.UserName != "" {
		tx = tx.Where("user_name LIKE ?", "%"+s.UserName+"%")
	}
	if s.NickName != "" {
		tx = tx.Where("nick_name LIKE ?", "%"+s.NickName+"%")
	}
	if s.UserPhone != "" {
		tx = tx.Where("user_phone LIKE ?", "%"+s.UserPhone+"%")
	}
	if s.UserEmail != "" {
		tx = tx.Where("user_email LIKE ?", "%"+s.UserEmail+"%")
	}
	if s.UserGender != "" {
		tx = tx.Where("user_gender = ?", s.UserGender)
	}
	if s.Status != "" {
		tx = tx.Where("status = ?", s.Status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.SysUser
	if err := tx.Order("id DESC").Offset(s.Offset()).Limit(s.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ReplaceRoles 替换用户的角色集合 (事务)
func (r *SysUserRepository) ReplaceRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.SysUserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}
		rows := make([]model.SysUserRole, 0, len(roleIDs))
		for _, rid := range roleIDs {
			rows = append(rows, model.SysUserRole{UserID: userID, RoleID: rid})
		}
		return tx.Create(&rows).Error
	})
}
