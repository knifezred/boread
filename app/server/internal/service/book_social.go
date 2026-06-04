package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

type BookSocialService struct {
	db          *gorm.DB
	noteRepo    *repository.ReaderNoteRepository
	reviewRepo  *repository.BookReviewRepository
	commentRepo *repository.BookChapterCommentRepository
	likeRepo    *repository.ReaderLikeRepository
	bookRepo    *repository.BookRepository
	chapterRepo *repository.BookChapterRepository
	userRepo    *repository.SysUserRepository
}

func NewBookSocialService(
	db *gorm.DB,
	noteRepo *repository.ReaderNoteRepository,
	reviewRepo *repository.BookReviewRepository,
	commentRepo *repository.BookChapterCommentRepository,
	likeRepo *repository.ReaderLikeRepository,
	bookRepo *repository.BookRepository,
	chapterRepo *repository.BookChapterRepository,
	userRepo *repository.SysUserRepository,
) *BookSocialService {
	return &BookSocialService{
		db:          db,
		noteRepo:    noteRepo,
		reviewRepo:  reviewRepo,
		commentRepo: commentRepo,
		likeRepo:    likeRepo,
		bookRepo:    bookRepo,
		chapterRepo: chapterRepo,
		userRepo:    userRepo,
	}
}

// ======================== 笔记/划线 ========================

func (s *BookSocialService) CreateNote(ctx context.Context, readerID uint64, req *dto.NoteRequest) (*dto.NoteResponse, error) {
	if _, err := s.bookRepo.GetByID(ctx, req.BookID); err != nil {
		return nil, code.ErrBookNotExists
	}
	visibility := req.Visibility
	if visibility == "" {
		visibility = "2"
	}
	m := &model.ReaderNote{
		ReaderID:       readerID,
		BookID:         req.BookID,
		ChapterID:      req.ChapterID,
		NoteType:       req.NoteType,
		SelectedText:   req.SelectedText,
		StartOffset:    req.StartOffset,
		EndOffset:      req.EndOffset,
		HighlightColor: req.HighlightColor,
		Content:        req.Content,
		Visibility:     visibility,
	}
	m.CreateBy = &readerID
	m.UpdateBy = &readerID

	if err := s.noteRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return s.buildNoteResponse(ctx, m, readerID)
}

func (s *BookSocialService) UpdateNote(ctx context.Context, readerID, noteID uint64, req *dto.NoteRequest) (*dto.NoteResponse, error) {
	m, err := s.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, code.ErrNoteNotFound
	}
	if m.ReaderID != readerID {
		return nil, code.ErrNoteNotOwner
	}
	// 可更新字段
	m.NoteType = req.NoteType
	m.SelectedText = req.SelectedText
	m.StartOffset = req.StartOffset
	m.EndOffset = req.EndOffset
	m.HighlightColor = req.HighlightColor
	m.Content = req.Content
	if req.Visibility != "" {
		m.Visibility = req.Visibility
	}
	m.UpdateBy = &readerID

	if err := s.noteRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return s.buildNoteResponse(ctx, m, readerID)
}

func (s *BookSocialService) DeleteNote(ctx context.Context, readerID, noteID uint64) error {
	return s.noteRepo.Delete(ctx, noteID, readerID)
}

func (s *BookSocialService) GetNote(ctx context.Context, noteID, readerID uint64) (*dto.NoteResponse, error) {
	m, err := s.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, code.ErrNoteNotFound
	}
	return s.buildNoteResponse(ctx, m, readerID)
}

func (s *BookSocialService) PageNote(ctx context.Context, readerID uint64, req *dto.NoteSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.noteRepo.PageByReader(ctx, readerID, req.BookID, &req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.NoteResponse, len(rows))
	for i, r := range rows {
		br, _ := s.buildNoteResponse(ctx, &r, readerID)
		if br != nil {
			resp[i] = *br
		}
	}
	return dto.NewPageResponse(resp, total, &req.PageRequest), nil
}

func (s *BookSocialService) ListNotesByBook(ctx context.Context, bookID uint64, noteType string, req *dto.PageRequest) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.noteRepo.PageByBook(ctx, bookID, noteType, req)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.NoteResponse, len(rows))
	for i, r := range rows {
		br, _ := s.buildNoteResponse(ctx, &r, 0)
		if br != nil {
			resp[i] = *br
		}
	}
	return dto.NewPageResponse(resp, total, req), nil
}

// ======================== 书评 ========================

func (s *BookSocialService) CreateReview(ctx context.Context, readerID uint64, req *dto.ReviewRequest) (*dto.ReviewResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, req.BookID)
	if err != nil {
		return nil, code.ErrBookNotExists
	}
	m := &model.BookReview{
		BookID:   req.BookID,
		ReaderID: readerID,
		Rating:   req.Rating,
		Title:    req.Title,
		Content:  req.Content,
		OwnerID:  readerID,
		DeptID:   book.DeptID,
		Status:   "3", // 默认审核中
	}
	m.CreateBy = &readerID
	m.UpdateBy = &readerID

	if err := s.reviewRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return s.buildReviewResponse(ctx, m, readerID)
}

func (s *BookSocialService) UpdateReview(ctx context.Context, readerID, reviewID uint64, req *dto.ReviewRequest) (*dto.ReviewResponse, error) {
	m, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, code.ErrReviewNotFound
	}
	if m.ReaderID != readerID {
		return nil, code.ErrNoteNotOwner // 复用: 无权修改他人书评
	}
	m.Rating = req.Rating
	m.Title = req.Title
	m.Content = req.Content
	m.Status = "3" // 修改后重新审核
	m.UpdateBy = &readerID

	if err := s.reviewRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return s.buildReviewResponse(ctx, m, readerID)
}

func (s *BookSocialService) DeleteReview(ctx context.Context, readerID, reviewID uint64) error {
	m, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return code.ErrReviewNotFound
	}
	if m.ReaderID != readerID {
		return code.ErrNoteNotOwner
	}
	return s.reviewRepo.Delete(ctx, reviewID)
}

func (s *BookSocialService) GetReview(ctx context.Context, reviewID, readerID uint64) (*dto.ReviewResponse, error) {
	m, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, code.ErrReviewNotFound
	}
	return s.buildReviewResponse(ctx, m, readerID)
}

func (s *BookSocialService) PageReview(ctx context.Context, req *dto.ReviewSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.reviewRepo.PageByBook(ctx, req.BookID, req.Status, &req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.ReviewResponse, len(rows))
	for i, r := range rows {
		br, _ := s.buildReviewResponse(ctx, &r, 0)
		if br != nil {
			resp[i] = *br
		}
	}
	return dto.NewPageResponse(resp, total, &req.PageRequest), nil
}

// ======================== 章节评论 ========================

func (s *BookSocialService) CreateComment(ctx context.Context, readerID uint64, req *dto.CommentRequest) (*dto.CommentResponse, error) {
	if _, err := s.chapterRepo.GetByID(ctx, req.ChapterID); err != nil {
		return nil, code.ErrChapterNotExists
	}
	if req.ParentID > 0 {
		if _, err := s.commentRepo.GetByID(ctx, req.ParentID); err != nil {
			return nil, code.ErrParentCommentNotExist
		}
	}
	m := &model.BookChapterComment{
		BookID:    req.BookID,
		ChapterID: req.ChapterID,
		ReaderID:  readerID,
		ParentID:  req.ParentID,
		ReplyToID: req.ReplyToID,
		Content:   req.Content,
		OwnerID:   readerID,
		Status:    "1",
	}
	m.CreateBy = &readerID
	m.UpdateBy = &readerID

	if err := s.commentRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return s.buildCommentResponse(ctx, m, readerID)
}

func (s *BookSocialService) DeleteComment(ctx context.Context, readerID, commentID uint64) error {
	return s.commentRepo.Delete(ctx, commentID, readerID)
}

func (s *BookSocialService) GetComment(ctx context.Context, commentID, readerID uint64) (*dto.CommentResponse, error) {
	m, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, code.ErrCommentNotFound
	}
	return s.buildCommentResponse(ctx, m, readerID)
}

func (s *BookSocialService) PageComment(ctx context.Context, req *dto.CommentSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.commentRepo.PageByChapter(ctx, req.ChapterID, req.ParentID, &req.PageRequest)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.CommentResponse, len(rows))
	for i, r := range rows {
		br, _ := s.buildCommentResponse(ctx, &r, 0)
		if br != nil {
			resp[i] = *br
		}
	}
	return dto.NewPageResponse(resp, total, &req.PageRequest), nil
}

// ======================== 点赞 ========================

func (s *BookSocialService) ToggleLike(ctx context.Context, readerID uint64, req *dto.LikeRequest) (*dto.LikeResponse, error) {
	liked, err := s.likeRepo.Toggle(ctx, readerID, req.TargetType, req.TargetID)
	if err != nil {
		return nil, err
	}
	count, err := s.likeRepo.CountByTarget(ctx, req.TargetType, req.TargetID)
	if err != nil {
		return nil, err
	}
	return &dto.LikeResponse{Liked: liked, Count: count}, nil
}

func (s *BookSocialService) GetLikeStatus(ctx context.Context, readerID uint64, req *dto.LikeStatusRequest) (map[string]bool, error) {
	targets := make([]struct {
		TargetType string
		TargetID   uint64
	}, len(req.Targets))
	for i, t := range req.Targets {
		targets[i] = struct {
			TargetType string
			TargetID   uint64
		}{TargetType: t.TargetType, TargetID: t.TargetID}
	}
	return s.likeRepo.BatchCheckExists(ctx, readerID, targets)
}

func (s *BookSocialService) CountLikes(ctx context.Context, targetType string, targetID uint64) (int64, error) {
	return s.likeRepo.CountByTarget(ctx, targetType, targetID)
}

// ======================== 辅助方法 ========================

func (s *BookSocialService) buildNoteResponse(ctx context.Context, m *model.ReaderNote, currentReaderID uint64) (*dto.NoteResponse, error) {
	resp := &dto.NoteResponse{ReaderNote: *m}
	// 查询读者昵称
	if user, err := s.userRepo.GetByID(ctx, m.ReaderID); err == nil {
		nick := user.NickName
		resp.ReaderName = &nick
	}
	// 查询书籍标题
	if book, err := s.bookRepo.GetByID(ctx, m.BookID); err == nil {
		resp.BookTitle = &book.Title
	}
	// 查询章节标题
	if m.ChapterID != nil {
		if ch, err := s.chapterRepo.GetByID(ctx, *m.ChapterID); err == nil {
			resp.ChapterTitle = &ch.Title
		}
	}
	// 点赞信息
	count, _ := s.likeRepo.CountByTarget(ctx, "3", m.ID)
	resp.LikeCount = count
	if currentReaderID > 0 {
		liked, _ := s.likeRepo.Exists(ctx, currentReaderID, "3", m.ID)
		resp.Liked = liked
	}
	return resp, nil
}

func (s *BookSocialService) buildReviewResponse(ctx context.Context, m *model.BookReview, currentReaderID uint64) (*dto.ReviewResponse, error) {
	resp := &dto.ReviewResponse{BookReview: *m}
	if user, err := s.userRepo.GetByID(ctx, m.ReaderID); err == nil {
		nick := user.NickName
		resp.ReaderName = &nick
	}
	if book, err := s.bookRepo.GetByID(ctx, m.BookID); err == nil {
		resp.BookTitle = &book.Title
	}
	count, _ := s.likeRepo.CountByTarget(ctx, "1", m.ID)
	resp.LikeCount = count
	if currentReaderID > 0 {
		liked, _ := s.likeRepo.Exists(ctx, currentReaderID, "1", m.ID)
		resp.Liked = liked
	}
	return resp, nil
}

func (s *BookSocialService) buildCommentResponse(ctx context.Context, m *model.BookChapterComment, currentReaderID uint64) (*dto.CommentResponse, error) {
	resp := &dto.CommentResponse{BookChapterComment: *m}
	if user, err := s.userRepo.GetByID(ctx, m.ReaderID); err == nil {
		nick := user.NickName
		resp.ReaderName = &nick
	}
	if m.ReplyToID != nil {
		if replyUser, err := s.userRepo.GetByID(ctx, *m.ReplyToID); err == nil {
			name := replyUser.NickName
			resp.ReplyToName = &name
		}
	}
	count, _ := s.likeRepo.CountByTarget(ctx, "2", m.ID)
	resp.LikeCount = count
	if currentReaderID > 0 {
		liked, _ := s.likeRepo.Exists(ctx, currentReaderID, "2", m.ID)
		resp.Liked = liked
	}
	return resp, nil
}

// 锚定引用，确保 swag 等工具能解析
var _ = time.Time{}
