package repository

import (
	"context"

	"boread/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository[model.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: *NewBaseRepository[model.User](db),
	}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ListWithKeyword(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := r.DB.WithContext(ctx).Model(&model.User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("username LIKE ? OR nickname LIKE ?", like, like)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}