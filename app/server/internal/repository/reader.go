package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

// ======================== ReaderBookshelfRepository ========================

type ReaderBookshelfRepository struct {
	db *gorm.DB
}

func NewReaderBookshelfRepository(db *gorm.DB) *ReaderBookshelfRepository {
	return &ReaderBookshelfRepository{db: db}
}

func (r *ReaderBookshelfRepository) Create(ctx context.Context, m *model.ReaderBookshelf) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ReaderBookshelfRepository) Update(ctx context.Context, m *model.ReaderBookshelf) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *ReaderBookshelfRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.ReaderBookshelf{}, id).Error
}

func (r *ReaderBookshelfRepository) GetByID(ctx context.Context, id uint64) (*model.ReaderBookshelf, error) {
	var m model.ReaderBookshelf
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ReaderBookshelfRepository) GetByReaderAndBook(ctx context.Context, readerID, bookID uint64) (*model.ReaderBookshelf, error) {
	var m model.ReaderBookshelf
	err := r.db.WithContext(ctx).
		Where("reader_id = ? AND book_id = ?", readerID, bookID).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ReaderBookshelfRepository) DeleteByReaderAndBook(ctx context.Context, readerID, bookID uint64) error {
	return r.db.WithContext(ctx).
		Where("reader_id = ? AND book_id = ?", readerID, bookID).
		Delete(&model.ReaderBookshelf{}).Error
}

// PageByReader 分页查询书架，按最后阅读时间降序排序，置顶优先
func (r *ReaderBookshelfRepository) PageByReader(ctx context.Context, readerID uint64, req *dto.BookshelfSearch) ([]model.ReaderBookshelf, int64, error) {
	var rows []model.ReaderBookshelf
	tx := r.db.WithContext(ctx).Model(&model.ReaderBookshelf{}).
		Where("reader_id = ?", readerID)

	if req.GroupName != "" {
		tx = tx.Where("group_name = ?", req.GroupName)
	}
	if req.Keyword != "" {
		// 关联 book 表按书名搜索
		tx = tx.Joins("JOIN book ON book.id = reader_bookshelf.book_id").
			Where("book.title LIKE ?", "%"+req.Keyword+"%")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := tx.
		Order("is_top DESC, COALESCE(last_read_time, '1970-01-01') DESC, add_time DESC").
		Offset(req.Offset()).
		Limit(req.Size).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListGroupsByReader 获取读者的所有分组及其书籍数量
func (r *ReaderBookshelfRepository) ListGroupsByReader(ctx context.Context, readerID uint64) ([]dto.BookshelfGroupItem, error) {
	type groupCount struct {
		GroupName string `gorm:"column:group_name"`
		BookCount int64  `gorm:"column:cnt"`
	}
	var items []groupCount
	if err := r.db.WithContext(ctx).
		Model(&model.ReaderBookshelf{}).
		Select("group_name, COUNT(*) AS cnt").
		Where("reader_id = ?", readerID).
		Group("group_name").
		Order("cnt DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	result := make([]dto.BookshelfGroupItem, len(items))
	for i, item := range items {
		result[i] = dto.BookshelfGroupItem{
			GroupName: item.GroupName,
			BookCount: item.BookCount,
		}
	}
	return result, nil
}

// UpdateIsTop 更新置顶状态
func (r *ReaderBookshelfRepository) UpdateIsTop(ctx context.Context, id uint64, isTop bool) error {
	return r.db.WithContext(ctx).Model(&model.ReaderBookshelf{}).
		Where("id = ?", id).
		UpdateColumn("is_top", isTop).Error
}

// UpdateGroupName 更新分组名
func (r *ReaderBookshelfRepository) UpdateGroupName(ctx context.Context, id uint64, groupName string) error {
	return r.db.WithContext(ctx).Model(&model.ReaderBookshelf{}).
		Where("id = ?", id).
		UpdateColumn("group_name", groupName).Error
}

// UpdateLastReadTime 更新最后阅读时间
func (r *ReaderBookshelfRepository) UpdateLastReadTime(ctx context.Context, readerID, bookID uint64, t time.Time) error {
	return r.db.WithContext(ctx).Model(&model.ReaderBookshelf{}).
		Where("reader_id = ? AND book_id = ?", readerID, bookID).
		UpdateColumn("last_read_time", t).Error
}

// ======================== ReaderReadProgressRepository ========================

type ReaderReadProgressRepository struct {
	db *gorm.DB
}

func NewReaderReadProgressRepository(db *gorm.DB) *ReaderReadProgressRepository {
	return &ReaderReadProgressRepository{db: db}
}

// GetByReaderAndBook 获取某本书的阅读进度
func (r *ReaderReadProgressRepository) GetByReaderAndBook(ctx context.Context, readerID, bookID uint64) (*model.ReaderReadProgress, error) {
	var m model.ReaderReadProgress
	err := r.db.WithContext(ctx).
		Where("reader_id = ? AND book_id = ?", readerID, bookID).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// BatchGetByReader 批量获取一个读者的多本书进度
func (r *ReaderReadProgressRepository) BatchGetByReader(ctx context.Context, readerID uint64, bookIDs []uint64) (map[uint64]*model.ReaderReadProgress, error) {
	if len(bookIDs) == 0 {
		return nil, nil
	}
	var rows []model.ReaderReadProgress
	if err := r.db.WithContext(ctx).
		Where("reader_id = ? AND book_id IN ?", readerID, bookIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[uint64]*model.ReaderReadProgress, len(rows))
	for i := range rows {
		result[rows[i].BookID] = &rows[i]
	}
	return result, nil
}

// UpsertProgress 直接使用 SQL 实现 MySQL ON DUPLICATE KEY UPDATE
func (r *ReaderReadProgressRepository) UpsertProgress(ctx context.Context, m *model.ReaderReadProgress) error {
	sql := `INSERT INTO reader_read_progress
		(reader_id, book_id, file_id, chapter_id, chapter_no, position, percent, read_duration, last_read_time, create_time, create_by, update_by, update_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
		file_id = VALUES(file_id),
		chapter_id = VALUES(chapter_id),
		chapter_no = VALUES(chapter_no),
		position = VALUES(position),
		percent = VALUES(percent),
		read_duration = read_duration + VALUES(read_duration),
		last_read_time = NOW(),
		update_time = NOW()`
	return r.db.WithContext(ctx).Exec(sql,
		m.ReaderID, m.BookID, m.FileID, m.ChapterID, m.ChapterNo,
		m.Position, m.Percent, m.ReadDuration,
		m.CreateBy, m.UpdateBy,
	).Error
}
