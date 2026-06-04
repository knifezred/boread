package service

import (
	"context"

	"boread/internal/dto"
	"boread/internal/repository"
)

// BookReadStatsService 阅读统计服务 (只读聚合查询)
type BookReadStatsService struct {
	statsRepo *repository.ReaderReadStatsRepository
	bookRepo  *repository.BookRepository
}

func NewBookReadStatsService(
	statsRepo *repository.ReaderReadStatsRepository,
	bookRepo *repository.BookRepository,
) *BookReadStatsService {
	return &BookReadStatsService{
		statsRepo: statsRepo,
		bookRepo:  bookRepo,
	}
}

// GetDailyStats 按日聚合
func (s *BookReadStatsService) GetDailyStats(ctx context.Context, readerID uint64, req *dto.ReadStatsQuery) ([]dto.ReadEventDailyResponse, error) {
	startDate, endDate := normalizeDateRange(req.StartDate, req.EndDate)
	return s.statsRepo.SumDailyByReader(ctx, readerID, startDate, endDate)
}

// GetBookStats 按书聚合
func (s *BookReadStatsService) GetBookStats(ctx context.Context, readerID uint64, req *dto.ReadStatsQuery) ([]dto.ReadEventBookResponse, error) {
	startDate, endDate := normalizeDateRange(req.StartDate, req.EndDate)
	return s.statsRepo.SumBookByReader(ctx, readerID, startDate, endDate)
}

// GetBookDetail 单本书明细统计
func (s *BookReadStatsService) GetBookDetail(ctx context.Context, readerID, bookID uint64, req *dto.ReadStatsQuery) (*dto.ReadEventBookResponse, error) {
	startDate, endDate := normalizeDateRange(req.StartDate, req.EndDate)
	return s.statsRepo.SumBookStatsByReaderAndBook(ctx, readerID, bookID, startDate, endDate)
}

// GetTotalStats 总统计
func (s *BookReadStatsService) GetTotalStats(ctx context.Context, readerID uint64) (*dto.ReadEventTotalResponse, error) {
	return s.statsRepo.SumTotalByReader(ctx, readerID)
}

func normalizeDateRange(startDate, endDate string) (string, string) {
	if startDate == "" {
		startDate = "2000-01-01"
	}
	if endDate == "" {
		endDate = "2099-12-31"
	}
	return startDate, endDate
}
