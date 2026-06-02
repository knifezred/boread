package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

var (
	ErrBookshelfNotFound  = errors.New("书架记录不存在")
	ErrBookNotExist       = errors.New("书籍不存在")
	ErrAlreadyInBookshelf = errors.New("该书已在书架中")
	ErrProgressNotFound   = errors.New("阅读进度不存在")
)

type ReaderService struct {
	db            *gorm.DB
	bookshelfRepo *repository.ReaderBookshelfRepository
	progressRepo  *repository.ReaderReadProgressRepository
	bookRepo      *repository.BookRepository
	chapterRepo   *repository.BookChapterRepository
}

func NewReaderService(
	db *gorm.DB,
	bookshelfRepo *repository.ReaderBookshelfRepository,
	progressRepo *repository.ReaderReadProgressRepository,
	bookRepo *repository.BookRepository,
	chapterRepo *repository.BookChapterRepository,
) *ReaderService {
	return &ReaderService{
		db:            db,
		bookshelfRepo: bookshelfRepo,
		progressRepo:  progressRepo,
		bookRepo:      bookRepo,
		chapterRepo:   chapterRepo,
	}
}

// AddToBookshelf 添加到书架
func (s *ReaderService) AddToBookshelf(ctx context.Context, userID uint64, req *dto.BookshelfRequest) (*dto.BookshelfResponse, error) {
	// 校验书籍存在
	if _, err := s.bookRepo.GetByID(ctx, req.BookID); err != nil {
		return nil, ErrBookNotExist
	}
	// 检查是否已在书架
	existing, err := s.bookshelfRepo.GetByReaderAndBook(ctx, userID, req.BookID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyInBookshelf
	}
	groupName := req.GroupName
	if groupName == "" {
		groupName = "默认"
	}
	now := time.Now()
	m := &model.ReaderBookshelf{
		ReaderID:     userID,
		BookID:       req.BookID,
		GroupName:    groupName,
		IsTop:        false,
		LastReadTime: nil,
		AddTime:      now,
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID

	if err := s.bookshelfRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return s.buildBookshelfResponse(ctx, m)
}

// RemoveFromBookshelf 从书架移除
func (s *ReaderService) RemoveFromBookshelf(ctx context.Context, userID uint64, bookID uint64) error {
	_, err := s.bookshelfRepo.GetByReaderAndBook(ctx, userID, bookID)
	if err != nil {
		return ErrBookshelfNotFound
	}
	return s.bookshelfRepo.DeleteByReaderAndBook(ctx, userID, bookID)
}

// UpdateBookshelf 更新书架 (分组/置顶)
func (s *ReaderService) UpdateBookshelf(ctx context.Context, userID uint64, bookID uint64, req *dto.BookshelfUpdateRequest) (*dto.BookshelfResponse, error) {
	m, err := s.bookshelfRepo.GetByReaderAndBook(ctx, userID, bookID)
	if err != nil {
		return nil, ErrBookshelfNotFound
	}
	if req.GroupName != nil {
		m.GroupName = *req.GroupName
	}
	if req.IsTop != nil {
		m.IsTop = *req.IsTop
	}
	m.UpdateBy = &userID

	if err := s.bookshelfRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return s.buildBookshelfResponse(ctx, m)
}

// GetBookshelfPage 分页获取书架列表
func (s *ReaderService) GetBookshelfPage(ctx context.Context, userID uint64, req *dto.BookshelfSearch) (*dto.BookshelfPageResponse, error) {
	req.Normalize()

	rows, total, err := s.bookshelfRepo.PageByReader(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &dto.BookshelfPageResponse{
			Records: []dto.BookshelfResponse{},
			Current: req.Current,
			Size:    req.Size,
			Total:   0,
		}, nil
	}

	// 批量构建响应
	responses := make([]dto.BookshelfResponse, len(rows))
	for i, r := range rows {
		resp, err := s.buildBookshelfResponse(ctx, &r)
		if err != nil {
			// 如果书已被删除，仍然返回书架记录但基本信息为空
			emptyResp := &dto.BookshelfResponse{ReaderBookshelf: r}
			responses[i] = *emptyResp
			continue
		}
		responses[i] = *resp
	}

	// 批量填充阅读进度
	bookIDs := make([]uint64, len(rows))
	for i, r := range rows {
		bookIDs[i] = r.BookID
	}
	progressMap, _ := s.progressRepo.BatchGetByReader(ctx, userID, bookIDs)
	for i := range responses {
		if p, ok := progressMap[responses[i].BookID]; ok {
			responses[i].ChapterID = &p.ChapterID
			responses[i].ChapterNo = &p.ChapterNo
			responses[i].Position = &p.Position
			responses[i].ReadPercent = p.Percent
			responses[i].ReadDuration = p.ReadDuration
			lastRead := p.LastReadTime.Format("2006-01-02 15:04:05")
			responses[i].LastReadTime = &lastRead
		}
	}

	return &dto.BookshelfPageResponse{
		Records: responses,
		Current: req.Current,
		Size:    req.Size,
		Total:   total,
	}, nil
}

// ListGroups 获取分组列表
func (s *ReaderService) ListGroups(ctx context.Context, userID uint64) ([]dto.BookshelfGroupItem, error) {
	groups, err := s.bookshelfRepo.ListGroupsByReader(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 如果没有分组，返回默认分组
	if len(groups) == 0 {
		groups = []dto.BookshelfGroupItem{
			{GroupName: "默认", BookCount: 0},
		}
	}
	return groups, nil
}

// ReportProgress 上报阅读进度
func (s *ReaderService) ReportProgress(ctx context.Context, userID uint64, bookID uint64, req *dto.ReadProgressRequest) (*dto.ReadProgressResponse, error) {
	// 校验书籍存在
	if _, err := s.bookRepo.GetByID(ctx, bookID); err != nil {
		return nil, ErrBookNotExist
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
func (s *ReaderService) GetProgress(ctx context.Context, userID uint64, bookID uint64) (*dto.ReadProgressResponse, error) {
	m, err := s.progressRepo.GetByReaderAndBook(ctx, userID, bookID)
	if err != nil {
		return nil, ErrProgressNotFound
	}
	return &dto.ReadProgressResponse{ReaderReadProgress: *m}, nil
}

// buildBookshelfResponse 构建书架响应（填充书籍基本信息）
func (s *ReaderService) buildBookshelfResponse(ctx context.Context, m *model.ReaderBookshelf) (*dto.BookshelfResponse, error) {
	resp := &dto.BookshelfResponse{ReaderBookshelf: *m}

	book, err := s.bookRepo.GetByID(ctx, m.BookID)
	if err != nil {
		return nil, err
	}
	resp.BookTitle = book.Title
	resp.BookAuthor = book.Author
	resp.BookCover = book.Cover
	resp.TotalChapters = book.TotalChapters
	resp.TotalWords = book.TotalWords

	// 填充阅读进度
	if p, err := s.progressRepo.GetByReaderAndBook(ctx, m.ReaderID, m.BookID); err == nil {
		resp.ChapterID = &p.ChapterID
		resp.ChapterNo = &p.ChapterNo
		resp.Position = &p.Position
		resp.ReadPercent = p.Percent
		resp.ReadDuration = p.ReadDuration
		lastRead := p.LastReadTime.Format("2006-01-02 15:04:05")
		resp.LastReadTime = &lastRead
	}

	return resp, nil
}
