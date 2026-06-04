package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/model"
)

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

// ======================== ReaderReadEventRepository ========================

type ReaderReadEventRepository struct {
	db *gorm.DB
}

func NewReaderReadEventRepository(db *gorm.DB) *ReaderReadEventRepository {
	return &ReaderReadEventRepository{db: db}
}

// Create 追加一条阅读事件
func (r *ReaderReadEventRepository) Create(ctx context.Context, m *model.ReaderReadEvent) error {
	return r.db.WithContext(ctx).Create(m).Error
}
