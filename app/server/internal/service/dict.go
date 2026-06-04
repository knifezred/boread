package service

import (
	"context"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

type DictService struct {
	repo *repository.SysDictRepository
}

func NewDictService(repo *repository.SysDictRepository) *DictService {
	return &DictService{repo: repo}
}

func (s *DictService) Create(ctx context.Context, req *dto.DictRequest, opID uint64) (*model.SysDict, error) {
	if _, err := s.repo.GetByCode(ctx, req.DictCode); err == nil {
		return nil, code.ErrDictCodeExists
	}
	status := req.Status
	if status == "" {
		status = model.StatusEnabled
	}
	m := &model.SysDict{
		DictName: req.DictName,
		DictCode: req.DictCode,
		Status:   status,
	}
	if req.DictDesc != "" {
		m.DictDesc = &req.DictDesc
	}
	m.CreateBy = &opID
	m.UpdateBy = &opID
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *DictService) Update(ctx context.Context, id uint64, req *dto.DictRequest, opID uint64) (*model.SysDict, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m.IsSystem && req.DictCode != m.DictCode {
		return nil, code.ErrDictSystem
	}
	if req.DictCode != m.DictCode {
		if _, e := s.repo.GetByCode(ctx, req.DictCode); e == nil {
			return nil, code.ErrDictCodeExists
		}
	}
	m.DictName = req.DictName
	m.DictCode = req.DictCode
	if req.Status != "" {
		m.Status = req.Status
	}
	if req.DictDesc != "" {
		m.DictDesc = &req.DictDesc
	}
	m.UpdateBy = &opID
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *DictService) Delete(ctx context.Context, id uint64) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if m.IsSystem {
		return code.ErrDictSystem
	}
	return s.repo.Delete(ctx, id)
}

func (s *DictService) GetByID(ctx context.Context, id uint64) (*model.SysDict, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DictService) Page(ctx context.Context, req *dto.DictSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	return dto.NewPageResponse(rows, total, &req.PageRequest), nil
}

// === 字典项 ===

func (s *DictService) CreateItem(ctx context.Context, req *dto.DictItemRequest, opID uint64) (*model.SysDictItem, error) {
	status := req.Status
	if status == "" {
		status = model.StatusEnabled
	}
	m := &model.SysDictItem{
		DictID:    req.DictID,
		ItemLabel: req.ItemLabel,
		ItemValue: req.ItemValue,
		SortOrder: req.SortOrder,
		Status:    status,
	}
	if req.ItemDesc != "" {
		m.ItemDesc = &req.ItemDesc
	}
	m.CreateBy = &opID
	m.UpdateBy = &opID
	if err := s.repo.CreateItem(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *DictService) UpdateItem(ctx context.Context, id uint64, req *dto.DictItemRequest, opID uint64) (*model.SysDictItem, error) {
	m, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	m.ItemLabel = req.ItemLabel
	m.ItemValue = req.ItemValue
	m.SortOrder = req.SortOrder
	if req.Status != "" {
		m.Status = req.Status
	}
	if req.ItemDesc != "" {
		m.ItemDesc = &req.ItemDesc
	}
	m.UpdateBy = &opID
	if err := s.repo.UpdateItem(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *DictService) DeleteItem(ctx context.Context, id uint64) error {
	return s.repo.DeleteItem(ctx, id)
}

func (s *DictService) ListItemsByDictID(ctx context.Context, dictID uint64) ([]model.SysDictItem, error) {
	return s.repo.ListItemsByDictID(ctx, dictID)
}

func (s *DictService) ListItemsByCode(ctx context.Context, code string) ([]model.SysDictItem, error) {
	return s.repo.ListItemsByDictCode(ctx, code)
}
