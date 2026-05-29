package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/model"
)

type BookCategoryRepository struct {
	db *gorm.DB
}

func NewBookCategoryRepository(db *gorm.DB) *BookCategoryRepository {
	return &BookCategoryRepository{db: db}
}

func (r *BookCategoryRepository) Create(ctx context.Context, m *model.BookCategory) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookCategoryRepository) Update(ctx context.Context, m *model.BookCategory) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookCategoryRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BookCategory{}, id).Error
}

func (r *BookCategoryRepository) GetByID(ctx context.Context, id uint64) (*model.BookCategory, error) {
	var m model.BookCategory
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookCategoryRepository) GetByCode(ctx context.Context, code string) (*model.BookCategory, error) {
	var m model.BookCategory
	if err := r.db.WithContext(ctx).Where("category_code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookCategoryRepository) ListAll(ctx context.Context) ([]model.BookCategory, error) {
	var rows []model.BookCategory
	if err := r.db.WithContext(ctx).Model(&model.BookCategory{}).Order("parent_id ASC, sort_order ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *BookCategoryRepository) HasChildren(ctx context.Context, id uint64) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.BookCategory{}).Where("parent_id = ?", id).Count(&n).Error
	return n > 0, err
}

func (r *BookCategoryRepository) PageTop(ctx context.Context, name, code, status string, current, size int) ([]model.BookCategory, int64, error) {
	var rows []model.BookCategory
	tx := r.db.WithContext(ctx).Model(&model.BookCategory{}).Where("parent_id = 0")
	if name != "" {
		tx = tx.Where("category_name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		tx = tx.Where("category_code LIKE ?", "%"+code+"%")
	}
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("sort_order ASC").Offset((current - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *BookCategoryRepository) ListByIDs(ctx context.Context, ids []uint64) ([]model.BookCategory, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []model.BookCategory
	if err := r.db.WithContext(ctx).Model(&model.BookCategory{}).Where("parent_id IN ?", ids).Order("sort_order ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetMapByIDs 根据主键 ID 批量查询分类，返回 id → CategoryName 映射
func (r *BookCategoryRepository) GetMapByIDs(ctx context.Context, ids []uint64) (map[uint64]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []model.BookCategory
	if err := r.db.WithContext(ctx).Model(&model.BookCategory{}).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[uint64]string, len(rows))
	for _, c := range rows {
		m[c.ID] = c.CategoryName
	}
	return m, nil
}