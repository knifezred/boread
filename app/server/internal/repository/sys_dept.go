package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/model"
)

type SysDeptRepository struct {
	db *gorm.DB
}

func NewSysDeptRepository(db *gorm.DB) *SysDeptRepository {
	return &SysDeptRepository{db: db}
}

func (r *SysDeptRepository) Create(ctx context.Context, m *model.SysDept) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SysDeptRepository) Update(ctx context.Context, m *model.SysDept) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *SysDeptRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.SysDept{}, id).Error
}

func (r *SysDeptRepository) GetByID(ctx context.Context, id uint64) (*model.SysDept, error) {
	var m model.SysDept
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *SysDeptRepository) GetByCode(ctx context.Context, code string) (*model.SysDept, error) {
	var m model.SysDept
	if err := r.db.WithContext(ctx).Where("dept_code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListAll 返回所有未删除部门 (按 sort_order)
func (r *SysDeptRepository) ListAll(ctx context.Context, name, code, status string) ([]model.SysDept, error) {
	var rows []model.SysDept
	tx := r.db.WithContext(ctx).Model(&model.SysDept{})
	if name != "" {
		tx = tx.Where("dept_name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		tx = tx.Where("dept_code LIKE ?", "%"+code+"%")
	}
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	if err := tx.Order("parent_id ASC, sort_order ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// HasChildren 是否存在子部门
func (r *SysDeptRepository) HasChildren(ctx context.Context, id uint64) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.SysDept{}).Where("parent_id = ?", id).Count(&n).Error
	return n > 0, err
}

// HasUsers 是否还有用户归属
func (r *SysDeptRepository) HasUsers(ctx context.Context, id uint64) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.SysUser{}).Where("dept_id = ?", id).Count(&n).Error
	return n > 0, err
}
