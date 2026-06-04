package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

type BookTagService struct {
	repo *repository.BookTagRepository
}

func NewBookTagService(repo *repository.BookTagRepository) *BookTagService {
	return &BookTagService{repo: repo}
}

func (s *BookTagService) Create(ctx context.Context, req *dto.TagRequest, userID uint64) (*model.BookTag, error) {
	if _, err := s.repo.GetByName(ctx, req.TagName); err == nil {
		return nil, code.ErrTagNameExists
	}
	m := &model.BookTag{TagName: req.TagName, Description: req.Description}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookTagService) Update(ctx context.Context, id uint64, req *dto.TagRequest, userID uint64) (*model.BookTag, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.TagName != m.TagName {
		if _, e := s.repo.GetByName(ctx, req.TagName); e == nil {
			return nil, code.ErrTagNameExists
		}
	}
	m.TagName = req.TagName
	m.Description = req.Description
	m.UpdateBy = &userID
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookTagService) Delete(ctx context.Context, id uint64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *BookTagService) GetByID(ctx context.Context, id uint64) (*model.BookTag, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return m, nil
}

func (s *BookTagService) Page(ctx context.Context, req *dto.TagSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	return dto.NewPageResponse(rows, total, &req.PageRequest), nil
}
