package repository

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, id).Error
}

func (r *BaseRepository[T]) List(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	db := r.DB.WithContext(ctx).Model(new(T))
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}