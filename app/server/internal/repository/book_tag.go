package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

type BookTagRepository struct {
	db *gorm.DB
}

func NewBookTagRepository(db *gorm.DB) *BookTagRepository {
	return &BookTagRepository{db: db}
}

func (r *BookTagRepository) Create(ctx context.Context, m *model.BookTag) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookTagRepository) Update(ctx context.Context, m *model.BookTag) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookTagRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BookTag{}, id).Error
}

func (r *BookTagRepository) GetByID(ctx context.Context, id uint64) (*model.BookTag, error) {
	var m model.BookTag
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookTagRepository) GetByName(ctx context.Context, name string) (*model.BookTag, error) {
	var m model.BookTag
	if err := r.db.WithContext(ctx).Where("tag_name = ?", name).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookTagRepository) Page(ctx context.Context, req *dto.TagSearch) ([]model.BookTag, int64, error) {
	var rows []model.BookTag
	tx := r.db.WithContext(ctx).Model(&model.BookTag{})
	if req.TagName != "" {
		tx = tx.Where("tag_name LIKE ?", "%"+req.TagName+"%")
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("usage_count DESC, id ASC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}