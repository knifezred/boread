package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

// BookReaderService 读者阅读服务 (阅读进度 + 阅读事件)
type BookReaderService struct {
	db            *gorm.DB
	bookshelfRepo *repository.ReaderBookshelfRepository
	progressRepo  *repository.ReaderReadProgressRepository
	eventRepo     *repository.ReaderReadEventRepository
	bookRepo      *repository.BookRepository
	chapterRepo   *repository.BookChapterRepository
}

func NewBookReaderService(
	db *gorm.DB,
	bookshelfRepo *repository.ReaderBookshelfRepository,
	progressRepo *repository.ReaderReadProgressRepository,
	eventRepo *repository.ReaderReadEventRepository,
	bookRepo *repository.BookRepository,
	chapterRepo *repository.BookChapterRepository,
) *BookReaderService {
	return &BookReaderService{
		db:            db,
		bookshelfRepo: bookshelfRepo,
		progressRepo:  progressRepo,
		eventRepo:     eventRepo,
		bookRepo:      bookRepo,
		chapterRepo:   chapterRepo,
	}
}

// ==================== 阅读进度 ====================

// ReportProgress 上报阅读进度
func (s *BookReaderService) ReportProgress(ctx context.Context, userID uint64, bookID uint64, req *dto.ReadProgressRequest) (*dto.ReadProgressResponse, error) {
	// 校验书籍存在
	if _, err := s.bookRepo.GetByID(ctx, bookID); err != nil {
		return nil, code.ErrBookNotExist
	}

	m := &model.ReaderReadProgress{
		ReaderID:     userID,
		BookID:       bookID,
		FileID:       req.FileID,
		ChapterID:    req.ChapterID,
		ChapterNo:    req.ChapterNo,
		Position:     req.Position,
		Percent:      req.Percent,
		ReadDuration: req.ReadDuration,
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID

	if err := s.progressRepo.UpsertProgress(ctx, m); err != nil {
		return nil, err
	}

	// 更新书架的最后阅读时间
	now := time.Now()
	_ = s.bookshelfRepo.UpdateLastReadTime(ctx, userID, bookID, now)

	// 如果不在书架中，自动加入
	if _, err := s.bookshelfRepo.GetByReaderAndBook(ctx, userID, bookID); err != nil {
		_ = s.bookshelfRepo.Create(ctx, &model.ReaderBookshelf{
			ReaderID:     userID,
			BookID:       bookID,
			GroupName:    "默认",
			IsTop:        false,
			LastReadTime: &now,
			AddTime:      now,
		})
	}

	resp := &dto.ReadProgressResponse{}
	// 重新查询获取最新数据
	if p, err := s.progressRepo.GetByReaderAndBook(ctx, userID, bookID); err == nil {
		resp.ReaderReadProgress = *p
	}
	return resp, nil
}

// GetProgress 获取阅读进度
func (s *BookReaderService) GetProgress(ctx context.Context, userID uint64, bookID uint64) (*dto.ReadProgressResponse, error) {
	m, err := s.progressRepo.GetByReaderAndBook(ctx, userID, bookID)
	if err != nil {
		return nil, code.ErrProgressNotFound
	}
	return &dto.ReadProgressResponse{ReaderReadProgress: *m}, nil
}

// ==================== 阅读事件 ====================

// ReportEvent 上报阅读事件
func (s *BookReaderService) ReportEvent(ctx context.Context, readerID uint64, req *dto.ReadEventRequest) error {
	if _, err := s.bookRepo.GetByID(ctx, req.BookID); err != nil {
		return code.ErrBookNotExist
	}
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	deviceType := req.DeviceType
	if deviceType == "" {
		deviceType = "web"
	}
	now := time.Now()
	eventDate := now.Format("2006-01-02")
	m := &model.ReaderReadEvent{
		ReaderID:    readerID,
		BookID:      req.BookID,
		ChapterID:   req.ChapterID,
		SessionID:   sessionID,
		DurationSec: req.DurationSec,
		WordCount:   req.WordCount,
		EventDate:   eventDate,
		EventTime:   now,
		DeviceType:  deviceType,
	}
	return s.eventRepo.Create(ctx, m)
}
