package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

// ======================== ReaderNoteRepository ========================

type ReaderNoteRepository struct {
	db *gorm.DB
}

func NewReaderNoteRepository(db *gorm.DB) *ReaderNoteRepository {
	return &ReaderNoteRepository{db: db}
}

func (r *ReaderNoteRepository) Create(ctx context.Context, m *model.ReaderNote) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ReaderNoteRepository) Update(ctx context.Context, m *model.ReaderNote) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *ReaderNoteRepository) Delete(ctx context.Context, id, readerID uint64) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND reader_id = ?", id, readerID).
		Delete(&model.ReaderNote{}).Error
}

func (r *ReaderNoteRepository) GetByID(ctx context.Context, id uint64) (*model.ReaderNote, error) {
	var m model.ReaderNote
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// PageByReader 分页查询读者的笔记，可选按 bookID 过滤
func (r *ReaderNoteRepository) PageByReader(ctx context.Context, readerID uint64, bookID *uint64, req *dto.PageRequest) ([]model.ReaderNote, int64, error) {
	var rows []model.ReaderNote
	tx := r.db.WithContext(ctx).Model(&model.ReaderNote{}).
		Where("reader_id = ?", readerID)
	if bookID != nil && *bookID > 0 {
		tx = tx.Where("book_id = ?", *bookID)
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

// PageByBook 分页查询某本书的公开笔记
func (r *ReaderNoteRepository) PageByBook(ctx context.Context, bookID uint64, noteType string, req *dto.PageRequest) ([]model.ReaderNote, int64, error) {
	var rows []model.ReaderNote
	tx := r.db.WithContext(ctx).Model(&model.ReaderNote{}).
		Where("book_id = ? AND visibility = '1'", bookID)
	if noteType != "" {
		tx = tx.Where("note_type = ?", noteType)
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

// ======================== BookReviewRepository ========================

type BookReviewRepository struct {
	db *gorm.DB
}

func NewBookReviewRepository(db *gorm.DB) *BookReviewRepository {
	return &BookReviewRepository{db: db}
}

func (r *BookReviewRepository) Create(ctx context.Context, m *model.BookReview) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookReviewRepository) Update(ctx context.Context, m *model.BookReview) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BookReviewRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.BookReview{}, id).Error
}

func (r *BookReviewRepository) GetByID(ctx context.Context, id uint64) (*model.BookReview, error) {
	var m model.BookReview
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// PageByBook 分页查询某本书的书评，按状态过滤
func (r *BookReviewRepository) PageByBook(ctx context.Context, bookID uint64, status string, req *dto.PageRequest) ([]model.BookReview, int64, error) {
	var rows []model.BookReview
	tx := r.db.WithContext(ctx).Model(&model.BookReview{}).
		Where("book_id = ?", bookID)
	if status != "" {
		tx = tx.Where("status = ?", status)
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

// PageByReader 分页查询读者的书评
func (r *BookReviewRepository) PageByReader(ctx context.Context, readerID uint64, req *dto.PageRequest) ([]model.BookReview, int64, error) {
	var rows []model.BookReview
	tx := r.db.WithContext(ctx).Model(&model.BookReview{}).
		Where("reader_id = ?", readerID)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("id DESC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ======================== BookChapterCommentRepository ========================

type BookChapterCommentRepository struct {
	db *gorm.DB
}

func NewBookChapterCommentRepository(db *gorm.DB) *BookChapterCommentRepository {
	return &BookChapterCommentRepository{db: db}
}

func (r *BookChapterCommentRepository) Create(ctx context.Context, m *model.BookChapterComment) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *BookChapterCommentRepository) Delete(ctx context.Context, id, readerID uint64) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND reader_id = ?", id, readerID).
		Delete(&model.BookChapterComment{}).Error
}

func (r *BookChapterCommentRepository) GetByID(ctx context.Context, id uint64) (*model.BookChapterComment, error) {
	var m model.BookChapterComment
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// PageByChapter 分页查询章节评论
// parentID=0 查顶层评论，>0 查子回复
func (r *BookChapterCommentRepository) PageByChapter(ctx context.Context, chapterID uint64, parentID uint64, req *dto.PageRequest) ([]model.BookChapterComment, int64, error) {
	var rows []model.BookChapterComment
	tx := r.db.WithContext(ctx).Model(&model.BookChapterComment{}).
		Where("chapter_id = ? AND parent_id = ?", chapterID, parentID)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := tx.Order("id ASC").Offset(req.Offset()).Limit(req.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ======================== ReaderLikeRepository ========================

type ReaderLikeRepository struct {
	db *gorm.DB
}

func NewReaderLikeRepository(db *gorm.DB) *ReaderLikeRepository {
	return &ReaderLikeRepository{db: db}
}

// Toggle 切换点赞状态：已赞则取消（软删除），未赞则创建
func (r *ReaderLikeRepository) Toggle(ctx context.Context, readerID uint64, targetType string, targetID uint64) (bool, error) {
	var m model.ReaderLike
	err := r.db.WithContext(ctx).
		Unscoped().
		Where("reader_id = ? AND target_type = ? AND target_id = ?", readerID, targetType, targetID).
		First(&m).Error

	if err != nil {
		// 不存在，创建点赞
		like := &model.ReaderLike{
			ReaderID:   readerID,
			TargetType: targetType,
			TargetID:   targetID,
		}
		if err := r.db.WithContext(ctx).Create(like).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	// 存在：切换 deleted_at
	if m.DeletedAt.Valid {
		// 已取消点赞 → 恢复
		if err := r.db.WithContext(ctx).Unscoped().
			Model(&m).
			UpdateColumn("deleted_at", nil).
			Error; err != nil {
			return false, err
		}
		return true, nil
	}

	// 已点赞 → 取消
	if err := r.db.WithContext(ctx).Delete(&m).Error; err != nil {
		return false, err
	}
	return false, nil
}

// CountByTarget 查询点赞总数
func (r *ReaderLikeRepository) CountByTarget(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ReaderLike{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error
	return count, err
}

// Exists 检查是否已点赞
func (r *ReaderLikeRepository) Exists(ctx context.Context, readerID uint64, targetType string, targetID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ReaderLike{}).
		Where("reader_id = ? AND target_type = ? AND target_id = ?", readerID, targetType, targetID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// BatchCheckExists 批量查询是否已点赞
func (r *ReaderLikeRepository) BatchCheckExists(ctx context.Context, readerID uint64, targets []struct {
	TargetType string
	TargetID   uint64
}) (map[string]bool, error) {
	if len(targets) == 0 {
		return nil, nil
	}
	result := make(map[string]bool, len(targets))
	for _, t := range targets {
		ok, err := r.Exists(ctx, readerID, t.TargetType, t.TargetID)
		if err != nil {
			return nil, err
		}
		key := t.TargetType + ":" + fmt.Sprint(t.TargetID)
		result[key] = ok
	}
	return result, nil
}
