package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

var (
	ErrBookNotFound    = errors.New("书籍不存在")
	ErrBookTagInvalid  = errors.New("标签无效")
	ErrCategoryInvalid = errors.New("分类不存在")
)

type BookService struct {
	db           *gorm.DB
	bookRepo     *repository.BookRepository
	tagRelRepo   *repository.BookTagRelRepository
	categoryRepo *repository.BookCategoryRepository
	tagRepo      *repository.BookTagRepository
}

func NewBookService(
	db *gorm.DB,
	bookRepo *repository.BookRepository,
	tagRelRepo *repository.BookTagRelRepository,
	categoryRepo *repository.BookCategoryRepository,
	tagRepo *repository.BookTagRepository,
) *BookService {
	return &BookService{
		db:           db,
		bookRepo:     bookRepo,
		tagRelRepo:   tagRelRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

func (s *BookService) Create(ctx context.Context, req *dto.BookRequest, userID uint64) (*dto.BookResponse, error) {
	if req.CategoryID != nil && *req.CategoryID > 0 {
		if _, err := s.categoryRepo.GetByID(ctx, *req.CategoryID); err != nil {
			return nil, ErrCategoryInvalid
		}
	}
	if len(req.TagIDs) > 0 {
		req.TagIDs = uniqueUint64Slice(req.TagIDs)
		for _, tid := range req.TagIDs {
			if _, err := s.tagRepo.GetByID(ctx, tid); err != nil {
				return nil, fmt.Errorf("%w: tag_id=%d", ErrBookTagInvalid, tid)
			}
		}
	}
	serialStatus := model.SerialStatus(req.SerialStatus)
	if serialStatus == "" {
		serialStatus = model.SerialOngoing
	}
	visibility := model.Visibility(req.Visibility)
	if visibility == "" {
		visibility = model.VisibilityPublic
	}
	language := req.Language
	if language == "" {
		language = "zh-CN"
	}

	m := &model.Book{
		Title:        req.Title,
		Author:       req.Author,
		Cover:        req.Cover,
		Intro:        req.Intro,
		CategoryID:   req.CategoryID,
		Language:     language,
		SerialStatus: serialStatus,
		Visibility:   visibility,
		OwnerID:      userID,
		Status:       model.BookReviewing,
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		if len(req.TagIDs) > 0 {
			rels := make([]model.BookTagRel, len(req.TagIDs))
			for i, tagID := range req.TagIDs {
				rels[i] = model.BookTagRel{
					BookID: m.ID,
					TagID:  tagID,
				}
				rels[i].CreateBy = &userID
				rels[i].UpdateBy = &userID
			}
			if err := tx.Create(&rels).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.BookTag{}).
				Where("id IN ?", req.TagIDs).
				UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.BookResponse{Book: *m}, nil
}

func (s *BookService) Update(ctx context.Context, id uint64, req *dto.BookRequest, userID uint64) (*dto.BookResponse, error) {
	m, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBookNotFound
	}
	if req.CategoryID != nil && *req.CategoryID > 0 {
		if _, err := s.categoryRepo.GetByID(ctx, *req.CategoryID); err != nil {
			return nil, ErrCategoryInvalid
		}
	}
	if len(req.TagIDs) > 0 {
		req.TagIDs = uniqueUint64Slice(req.TagIDs)
		for _, tid := range req.TagIDs {
			if _, err := s.tagRepo.GetByID(ctx, tid); err != nil {
				return nil, fmt.Errorf("%w: tag_id=%d", ErrBookTagInvalid, tid)
			}
		}
	}

	m.Title = req.Title
	m.Author = req.Author
	m.Cover = req.Cover
	m.Intro = req.Intro
	m.CategoryID = req.CategoryID
	if req.Language != "" {
		m.Language = req.Language
	}
	m.SerialStatus = model.SerialStatus(req.SerialStatus)
	m.Visibility = model.Visibility(req.Visibility)
	m.UpdateBy = &userID

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		// 取 old 和 new 标签 ID 的差集，只处理变更的部分
		var oldTagIDs []uint64
		if err := tx.Model(&model.BookTagRel{}).Where("book_id = ?", id).Pluck("tag_id", &oldTagIDs).Error; err != nil {
			return err
		}
		oldSet := sliceToSet(oldTagIDs)
		newSet := sliceToSet(req.TagIDs)
		var toDelete, toAdd []uint64
		for _, tid := range oldTagIDs {
			if _, exists := newSet[tid]; !exists {
				toDelete = append(toDelete, tid)
			}
		}
		for _, tid := range req.TagIDs {
			if _, exists := oldSet[tid]; !exists {
				toAdd = append(toAdd, tid)
			}
		}
		// 删除不再关联的标签
		if len(toDelete) > 0 {
			if err := tx.Where("book_id = ? AND tag_id IN ?", id, toDelete).Delete(&model.BookTagRel{}).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.BookTag{}).
				Where("id IN ?", toDelete).
				UpdateColumn("usage_count", gorm.Expr("usage_count - 1")).Error; err != nil {
				return err
			}
		}
		// 新增关联的标签
		if len(toAdd) > 0 {
			rels := make([]model.BookTagRel, len(toAdd))
			for i, tagID := range toAdd {
				rels[i] = model.BookTagRel{
					BookID: id,
					TagID:  tagID,
				}
				rels[i].CreateBy = &userID
				rels[i].UpdateBy = &userID
			}
			if err := tx.Create(&rels).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.BookTag{}).
				Where("id IN ?", toAdd).
				UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.BookResponse{Book: *m}, nil
}

func (s *BookService) Delete(ctx context.Context, id uint64) error {
	if _, err := s.bookRepo.GetByID(ctx, id); err != nil {
		return ErrBookNotFound
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Book{}, id).Error; err != nil {
			return err
		}
		var tagIDs []uint64
		if err := tx.Model(&model.BookTagRel{}).Where("book_id = ?", id).Pluck("tag_id", &tagIDs).Error; err != nil {
			return err
		}
		if err := tx.Where("book_id = ?", id).Delete(&model.BookTagRel{}).Error; err != nil {
			return err
		}
		if len(tagIDs) > 0 {
			if err := tx.Model(&model.BookTag{}).
				Where("id IN ?", tagIDs).
				UpdateColumn("usage_count", gorm.Expr("usage_count - 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BookService) GetByID(ctx context.Context, id uint64) (*dto.BookResponse, error) {
	m, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBookNotFound
	}
	tagIDs, err := s.tagRelRepo.GetTagIDsByBookID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := &dto.BookResponse{Book: *m}
	if len(tagIDs) > 0 {
		resp.TagIDs = tagIDs
	}
	if m.CategoryID != nil && *m.CategoryID > 0 {
		cat, err := s.categoryRepo.GetByID(ctx, *m.CategoryID)
		if err == nil {
			resp.CategoryName = cat.CategoryName
		}
	}
	return resp, nil
}

func (s *BookService) Page(ctx context.Context, req *dto.BookSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.bookRepo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return dto.NewPageResponse([]dto.BookResponse{}, total, &req.PageRequest), nil
	}
	bookIDs := make([]uint64, len(rows))
	for i, r := range rows {
		bookIDs[i] = r.ID
	}
	tagMap, err := s.tagRelRepo.GetTagsByBookIDs(ctx, bookIDs)
	if err != nil {
		return nil, err
	}
	// 收集所有 categoryID 批量查分类名
	var catIDs []uint64
	for _, r := range rows {
		if r.CategoryID != nil && *r.CategoryID > 0 {
			catIDs = append(catIDs, *r.CategoryID)
		}
	}
	catMap, err := s.categoryRepo.GetMapByIDs(ctx, catIDs)
	if err != nil {
		return nil, err
	}
	records := make([]dto.BookResponse, len(rows))
	for i, r := range rows {
		resp := dto.BookResponse{Book: r}
		if ids, ok := tagMap[r.ID]; ok {
			resp.TagIDs = ids
		}
		if r.CategoryID != nil {
			resp.CategoryName = catMap[*r.CategoryID]
		}
		records[i] = resp
	}
	return dto.NewPageResponse(records, total, &req.PageRequest), nil
}

func (s *BookService) UpdateStatus(ctx context.Context, id uint64, status string, userID uint64) error {
	m, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return ErrBookNotFound
	}
	m.Status = model.BookStatus(status)
	m.UpdateBy = &userID
	return s.bookRepo.Update(ctx, m)
}

/** 对 uint64 切片去重，保持原顺序 */
func uniqueUint64Slice(input []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(input))
	result := make([]uint64, 0, len(input))
	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

/** 将 uint64 切片转为 set 用于 O(1) 查询 */
func sliceToSet(input []uint64) map[uint64]struct{} {
	s := make(map[uint64]struct{}, len(input))
	for _, v := range input {
		s[v] = struct{}{}
	}
	return s
}
