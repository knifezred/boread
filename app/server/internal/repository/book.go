package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(ctx context.Context, m *model.Book) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookRepository) Update(ctx context.Context, m *model.Book) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Book{}, id).Error
}

func (r *BookRepository) GetByID(ctx context.Context, id uint64) (*model.Book, error) {
	var m model.Book
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *BookRepository) Page(ctx context.Context, req *dto.BookSearch) ([]model.Book, int64, error) {
	var rows []model.Book
	tx := r.db.WithContext(ctx).Model(&model.Book{})
	if req.Title != "" {
		tx = tx.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Author != "" {
		tx = tx.Where("author LIKE ?", "%"+req.Author+"%")
	}
	if len(req.CategoryIDs) > 0 {
		tx = tx.Where("category_id IN ?", req.CategoryIDs)
	}
	if req.Status != "" {
		tx = tx.Where("status = ?", req.Status)
	}
	if req.Visibility != "" {
		tx = tx.Where("visibility = ?", req.Visibility)
	}
	if req.SerialStatus != "" {
		tx = tx.Where("serial_status = ?", req.SerialStatus)
	}
	if req.TagID != nil && *req.TagID > 0 {
		tx = tx.Where("id IN (SELECT book_id FROM book_tag_rel WHERE tag_id = ?)", *req.TagID)
	}
	if req.MinWords != nil && *req.MinWords > 0 {
		tx = tx.Where("total_words >= ?", *req.MinWords*10000)
	}
	if req.MaxWords != nil && *req.MaxWords > 0 {
		tx = tx.Where("total_words <= ?", *req.MaxWords*10000)
	}
	if req.UpdateTimeFrom != nil && *req.UpdateTimeFrom != "" {
		tx = tx.Where("update_time >= ?", *req.UpdateTimeFrom)
	}
	if req.UpdateTimeTo != nil && *req.UpdateTimeTo != "" {
		tx = tx.Where("update_time <= ?", *req.UpdateTimeTo)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("id DESC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// FindByTitleAndAuthor 按书名和作者精确查找
func (r *BookRepository) FindByTitleAndAuthor(ctx context.Context, title, author string) (*model.Book, error) {
	var m model.Book
	err := r.db.WithContext(ctx).Where("title = ? AND author = ?", title, author).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByIDs 批量查询 (用于标签关联)
func (r *BookRepository) ListByIDs(ctx context.Context, ids []uint64) ([]model.Book, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []model.Book
	if err := r.db.WithContext(ctx).Model(&model.Book{}).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ======================== BookTagRel ========================

type BookTagRelRepository struct {
	db *gorm.DB
}

func NewBookTagRelRepository(db *gorm.DB) *BookTagRelRepository {
	return &BookTagRelRepository{db: db}
}

func (r *BookTagRelRepository) GetTagIDsByBookID(ctx context.Context, bookID uint64) ([]uint64, error) {
	var ids []uint64
	if err := r.db.WithContext(ctx).Model(&model.BookTagRel{}).
		Where("book_id = ?", bookID).
		Pluck("tag_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *BookTagRelRepository) DeleteByBookID(ctx context.Context, bookID uint64) error {
	return r.db.WithContext(ctx).Where("book_id = ?", bookID).Delete(&model.BookTagRel{}).Error
}

func (r *BookTagRelRepository) BatchCreate(ctx context.Context, rels []model.BookTagRel) error {
	if len(rels) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&rels).Error
}

// GetTagsByBookIDs 批量查多本书的标签关联
func (r *BookTagRelRepository) GetTagsByBookIDs(ctx context.Context, bookIDs []uint64) (map[uint64][]uint64, error) {
	if len(bookIDs) == 0 {
		return nil, nil
	}
	var rels []model.BookTagRel
	if err := r.db.WithContext(ctx).Model(&model.BookTagRel{}).
		Where("book_id IN ?", bookIDs).Find(&rels).Error; err != nil {
		return nil, err
	}
	result := make(map[uint64][]uint64, len(bookIDs))
	for _, rel := range rels {
		result[rel.BookID] = append(result[rel.BookID], rel.TagID)
	}
	return result, nil
}

// ListTagsByBookID 获取书籍的标签详情
func (r *BookTagRelRepository) ListTagsByBookID(ctx context.Context, bookID uint64) ([]model.BookTag, error) {
	var tags []model.BookTag
	if err := r.db.WithContext(ctx).
		Select("book_tag.id, book_tag.tag_name").
		Joins("JOIN book_tag_rel ON book_tag_rel.tag_id = book_tag.id").
		Where("book_tag_rel.book_id = ?", bookID).
		Where("book_tag.deleted_at IS NULL").
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
