package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

type SysSettingRepository struct {
	db *gorm.DB
}

func NewSysSettingRepository(db *gorm.DB) *SysSettingRepository {
	return &SysSettingRepository{db: db}
}

func (r *SysSettingRepository) Create(ctx context.Context, m *model.SysSetting) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SysSettingRepository) Update(ctx context.Context, m *model.SysSetting) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *SysSettingRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysSetting{}, id).Error
}

func (r *SysSettingRepository) GetByID(ctx context.Context, id uint64) (*model.SysSetting, error) {
	var m model.SysSetting
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysSettingRepository) GetByCategoryKey(ctx context.Context, category, key string) (*model.SysSetting, error) {
	var m model.SysSetting
	if err := r.db.WithContext(ctx).
		Where("category = ? AND `key` = ?", category, key).
		First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysSettingRepository) Page(ctx context.Context, s *dto.SettingSearch) ([]model.SysSetting, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.SysSetting{})
	if s.Category != "" {
		tx = tx.Where("category = ?", s.Category)
	}
	if s.Keyword != "" {
		tx = tx.Where("`key` LIKE ?", "%"+s.Keyword+"%")
	}
	if s.Status != "" {
		tx = tx.Where("status = ?", s.Status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.SysSetting
	if err := tx.Order("category ASC, id ASC").Offset(s.Offset()).Limit(s.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListCategories 获取所有不重复的分类
func (r *SysSettingRepository) ListCategories(ctx context.Context) ([]string, error) {
	var cats []string
	if err := r.db.WithContext(ctx).
		Model(&model.SysSetting{}).
		Where("status = '1'").
		Distinct("category").
		Order("category ASC").
		Pluck("category", &cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

// ListByCategory 获取指定分类下所有配置
func (r *SysSettingRepository) ListByCategory(ctx context.Context, category string) ([]model.SysSetting, error) {
	var rows []model.SysSetting
	if err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Order("`key` ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
