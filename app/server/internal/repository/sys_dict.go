package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

type SysDictRepository struct {
	db *gorm.DB
}

func NewSysDictRepository(db *gorm.DB) *SysDictRepository {
	return &SysDictRepository{db: db}
}

func (r *SysDictRepository) Create(ctx context.Context, m *model.SysDict) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SysDictRepository) Update(ctx context.Context, m *model.SysDict) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *SysDictRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysDict{}, id).Error
}

func (r *SysDictRepository) GetByID(ctx context.Context, id uint64) (*model.SysDict, error) {
	var m model.SysDict
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysDictRepository) GetByCode(ctx context.Context, code string) (*model.SysDict, error) {
	var m model.SysDict
	if err := r.db.WithContext(ctx).Where("dict_code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysDictRepository) Page(ctx context.Context, s *dto.DictSearch) ([]model.SysDict, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.SysDict{})
	if s.DictName != "" {
		tx = tx.Where("dict_name LIKE ?", "%"+s.DictName+"%")
	}
	if s.DictCode != "" {
		tx = tx.Where("dict_code LIKE ?", "%"+s.DictCode+"%")
	}
	if s.Status != "" {
		tx = tx.Where("status = ?", s.Status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.SysDict
	if err := tx.Order("id DESC").Offset(s.Offset()).Limit(s.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// === 字典项 ===

func (r *SysDictRepository) CreateItem(ctx context.Context, m *model.SysDictItem) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SysDictRepository) UpdateItem(ctx context.Context, m *model.SysDictItem) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *SysDictRepository) DeleteItem(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysDictItem{}, id).Error
}

func (r *SysDictRepository) GetItemByID(ctx context.Context, id uint64) (*model.SysDictItem, error) {
	var m model.SysDictItem
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysDictRepository) ListItemsByDictID(ctx context.Context, dictID uint64) ([]model.SysDictItem, error) {
	var rows []model.SysDictItem
	err := r.db.WithContext(ctx).Where("dict_id = ?", dictID).Order("sort_order ASC").Find(&rows).Error
	return rows, err
}

// ListItemsByDictCode 按 dict_code 拉项 (前端高频)
func (r *SysDictRepository) ListItemsByDictCode(ctx context.Context, code string) ([]model.SysDictItem, error) {
	d, err := r.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.ListItemsByDictID(ctx, d.ID)
}
