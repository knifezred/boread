package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

// ======================== BookCharacterRepository ========================

type BookCharacterRepository struct {
	db *gorm.DB
}

func NewBookCharacterRepository(db *gorm.DB) *BookCharacterRepository {
	return &BookCharacterRepository{db: db}
}

func (r *BookCharacterRepository) Create(ctx context.Context, m *model.BookCharacter) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookCharacterRepository) Update(ctx context.Context, m *model.BookCharacter) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookCharacterRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BookCharacter{}, id).Error
}

func (r *BookCharacterRepository) GetByID(ctx context.Context, id uint64) (*model.BookCharacter, error) {
	var m model.BookCharacter
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// PageByBook 分页查询某本书的角色，支持 keyword(模糊name) 和 role_type 筛选
func (r *BookCharacterRepository) PageByBook(ctx context.Context, req *dto.CharacterSearch) ([]model.BookCharacter, int64, error) {
	var rows []model.BookCharacter
	tx := r.db.WithContext(ctx).Model(&model.BookCharacter{}).Where("book_id = ?", req.BookID)
	if req.RoleType != "" {
		tx = tx.Where("role_type = ?", req.RoleType)
	}
	if req.Keyword != "" {
		tx = tx.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("sort_order ASC, id ASC").
		Offset(req.Offset()).
		Limit(req.Size).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListByBook 列出某本书全部角色，按 sort_order 排序
func (r *BookCharacterRepository) ListByBook(ctx context.Context, bookID uint64) ([]model.BookCharacter, error) {
	var rows []model.BookCharacter
	err := r.db.WithContext(ctx).Model(&model.BookCharacter{}).
		Where("book_id = ?", bookID).
		Order("sort_order ASC, id ASC").
		Find(&rows).Error
	return rows, err
}

// ======================== BookCharacterRelRepository ========================

type BookCharacterRelRepository struct {
	db *gorm.DB
}

func NewBookCharacterRelRepository(db *gorm.DB) *BookCharacterRelRepository {
	return &BookCharacterRelRepository{db: db}
}

func (r *BookCharacterRelRepository) Create(ctx context.Context, m *model.BookCharacterRel) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookCharacterRelRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BookCharacterRel{}, id).Error
}

func (r *BookCharacterRelRepository) GetByID(ctx context.Context, id uint64) (*model.BookCharacterRel, error) {
	var m model.BookCharacterRel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByCharacter 查询某个角色的所有关系
func (r *BookCharacterRelRepository) ListByCharacter(ctx context.Context, characterID uint64) ([]model.BookCharacterRel, error) {
	var rows []model.BookCharacterRel
	err := r.db.WithContext(ctx).Model(&model.BookCharacterRel{}).
		Where("character_a_id = ? OR character_b_id = ?", characterID, characterID).
		Order("sort_order ASC, id ASC").
		Find(&rows).Error
	return rows, err
}

// ListByBook 查询某本书的所有关系
func (r *BookCharacterRelRepository) ListByBook(ctx context.Context, bookID uint64) ([]model.BookCharacterRel, error) {
	var rows []model.BookCharacterRel
	err := r.db.WithContext(ctx).Model(&model.BookCharacterRel{}).
		Where("book_id = ?", bookID).
		Order("sort_order ASC, id ASC").
		Find(&rows).Error
	return rows, err
}
