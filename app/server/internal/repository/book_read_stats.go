package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
)

// ReaderReadStatsRepository 阅读统计仓库 (只读聚合查询)
type ReaderReadStatsRepository struct {
	db *gorm.DB
}

func NewReaderReadStatsRepository(db *gorm.DB) *ReaderReadStatsRepository {
	return &ReaderReadStatsRepository{db: db}
}

// SumDailyByReader 按 reader_id + 日期区间聚合每日数据
func (r *ReaderReadStatsRepository) SumDailyByReader(ctx context.Context, readerID uint64, startDate, endDate string) ([]dto.ReadEventDailyResponse, error) {
	sql := `SELECT event_date,
		SUM(duration_sec) AS duration_sec,
		SUM(word_count) AS word_count,
		COUNT(DISTINCT chapter_id) AS chapter_count,
		COUNT(DISTINCT book_id) AS book_count,
		COUNT(DISTINCT session_id) AS session_count
	FROM reader_read_event
	WHERE reader_id = ? AND event_date BETWEEN ? AND ?
	GROUP BY event_date
	ORDER BY event_date DESC`
	var rows []dto.ReadEventDailyResponse
	if err := r.db.WithContext(ctx).Raw(sql, readerID, startDate, endDate).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// SumBookByReader 按 reader_id + 日期区间聚合每本书数据
func (r *ReaderReadStatsRepository) SumBookByReader(ctx context.Context, readerID uint64, startDate, endDate string) ([]dto.ReadEventBookResponse, error) {
	sql := `SELECT e.book_id, b.title AS book_title,
		SUM(e.duration_sec) AS duration_sec,
		SUM(e.word_count) AS word_count,
		COUNT(DISTINCT e.chapter_id) AS chapter_count
	FROM reader_read_event e
	LEFT JOIN book b ON b.id = e.book_id
	WHERE e.reader_id = ? AND e.event_date BETWEEN ? AND ?
	GROUP BY e.book_id, b.title
	ORDER BY duration_sec DESC`
	var rows []dto.ReadEventBookResponse
	if err := r.db.WithContext(ctx).Raw(sql, readerID, startDate, endDate).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// SumBookStatsByReaderAndBook 查询某本书在日期区间内的聚合
func (r *ReaderReadStatsRepository) SumBookStatsByReaderAndBook(ctx context.Context, readerID, bookID uint64, startDate, endDate string) (*dto.ReadEventBookResponse, error) {
	sql := `SELECT e.book_id, b.title AS book_title,
		SUM(e.duration_sec) AS duration_sec,
		SUM(e.word_count) AS word_count,
		COUNT(DISTINCT e.chapter_id) AS chapter_count
	FROM reader_read_event e
	LEFT JOIN book b ON b.id = e.book_id
	WHERE e.reader_id = ? AND e.book_id = ? AND e.event_date BETWEEN ? AND ?
	GROUP BY e.book_id, b.title`
	var row dto.ReadEventBookResponse
	if err := r.db.WithContext(ctx).Raw(sql, readerID, bookID, startDate, endDate).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// SumTotalByReader 查询某个读者全部统计（不限日期）
func (r *ReaderReadStatsRepository) SumTotalByReader(ctx context.Context, readerID uint64) (*dto.ReadEventTotalResponse, error) {
	sql := `SELECT
		COALESCE(SUM(duration_sec), 0) AS duration_sec,
		COALESCE(SUM(word_count), 0) AS word_count,
		COALESCE(COUNT(DISTINCT book_id), 0) AS book_count,
		COALESCE(COUNT(DISTINCT session_id), 0) AS session_count
	FROM reader_read_event
	WHERE reader_id = ?`
	var row dto.ReadEventTotalResponse
	if err := r.db.WithContext(ctx).Raw(sql, readerID).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
