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
	ErrCategoryCodeExists     = errors.New("分类编码已存在")
	ErrCategoryHasChildren    = errors.New("存在子分类, 不能删除")
	ErrCategoryParentNotFound = errors.New("父分类不存在")
)

type BookCategoryService struct {
	repo *repository.BookCategoryRepository
}

func NewBookCategoryService(repo *repository.BookCategoryRepository) *BookCategoryService {
	return &BookCategoryService{repo: repo}
}

func (s *BookCategoryService) Create(ctx context.Context, req *dto.CategoryRequest, userID uint64) (*model.BookCategory, error) {
	if _, err := s.repo.GetByCode(ctx, req.CategoryCode); err == nil {
		return nil, ErrCategoryCodeExists
	}

	ancestors := "0"
	if req.ParentID > 0 {
		parent, err := s.repo.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, ErrCategoryParentNotFound
		}
		ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
	}

	status := req.Status
	if status == "" {
		status = model.StatusEnabled
	}
	m := &model.BookCategory{
		ParentID:     req.ParentID,
		Ancestors:    ancestors,
		CategoryName: req.CategoryName,
		CategoryCode: req.CategoryCode,
		SortOrder:    req.SortOrder,
		Status:       status,
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookCategoryService) Update(ctx context.Context, id uint64, req *dto.CategoryRequest, userID uint64) (*model.BookCategory, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CategoryCode != m.CategoryCode {
		if _, e := s.repo.GetByCode(ctx, req.CategoryCode); e == nil {
			return nil, ErrCategoryCodeExists
		}
	}

	if req.ParentID != m.ParentID {
		ancestors := "0"
		if req.ParentID > 0 {
			parent, e := s.repo.GetByID(ctx, req.ParentID)
			if e != nil {
				return nil, ErrCategoryParentNotFound
			}
			ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
		}
		m.Ancestors = ancestors
		m.ParentID = req.ParentID
	}

	m.CategoryName = req.CategoryName
	m.CategoryCode = req.CategoryCode
	m.SortOrder = req.SortOrder
	if req.Status != "" {
		m.Status = req.Status
	}
	m.UpdateBy = &userID

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookCategoryService) Delete(ctx context.Context, id uint64) error {
	if has, err := s.repo.HasChildren(ctx, id); err != nil {
		return err
	} else if has {
		return ErrCategoryHasChildren
	}
	return s.repo.Delete(ctx, id)
}

func (s *BookCategoryService) GetByID(ctx context.Context, id uint64) (*model.BookCategory, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category %d not found", id)
		}
		return nil, err
	}
	return m, nil
}

func (s *BookCategoryService) Tree(ctx context.Context) ([]*dto.CategoryNode, error) {
	rows, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildCategoryTree(rows), nil
}

func buildCategoryTree(rows []model.BookCategory) []*dto.CategoryNode {
	nodeMap := make(map[uint64]*dto.CategoryNode, len(rows))
	for _, r := range rows {
		node := &dto.CategoryNode{
			ID:           r.ID,
			ParentID:     r.ParentID,
			CategoryName: r.CategoryName,
			CategoryCode: r.CategoryCode,
			SortOrder:    r.SortOrder,
			Status:       r.Status,
			Children:     []*dto.CategoryNode{},
		}
		nodeMap[r.ID] = node
	}
	roots := make([]*dto.CategoryNode, 0)
	for _, n := range nodeMap {
		if parent, ok := nodeMap[n.ParentID]; ok {
			parent.Children = append(parent.Children, n)
		} else {
			roots = append(roots, n)
		}
	}
	return roots
}

// Page 分类分页列表 (树形, 基于顶级分类分页 + 递归加载子级)
func (s *BookCategoryService) Page(ctx context.Context, req *dto.CategorySearch) (*dto.PageResponse, error) {
	req.Normalize()
	topRows, total, err := s.repo.PageTop(ctx, req.CategoryName, req.CategoryCode, string(req.Status), req.Current, req.Size)
	if err != nil {
		return nil, err
	}
	if len(topRows) == 0 {
		return dto.NewPageResponse([]*dto.CategoryNode{}, total, &req.PageRequest), nil
	}
	nodeMap := make(map[uint64]*dto.CategoryNode, len(topRows))
	roots := make([]*dto.CategoryNode, 0, len(topRows))
	for _, r := range topRows {
		n := &dto.CategoryNode{
			ID:           r.ID,
			ParentID:     r.ParentID,
			CategoryName: r.CategoryName,
			CategoryCode: r.CategoryCode,
			SortOrder:    r.SortOrder,
			Status:       r.Status,
			Children:     []*dto.CategoryNode{},
		}
		nodeMap[r.ID] = n
		roots = append(roots, n)
	}
	s.loadChildrenRecursively(ctx, nodeMap)
	return dto.NewPageResponse(roots, total, &req.PageRequest), nil
}

func (s *BookCategoryService) loadChildrenRecursively(ctx context.Context, nodeMap map[uint64]*dto.CategoryNode) {
	ids := make([]uint64, 0, len(nodeMap))
	for id := range nodeMap {
		ids = append(ids, id)
	}
	children, err := s.repo.ListByIDs(ctx, ids)
	if err != nil || len(children) == 0 {
		return
	}
	childMap := make(map[uint64]*dto.CategoryNode, len(children))
	for _, c := range children {
		n := &dto.CategoryNode{
			ID:           c.ID,
			ParentID:     c.ParentID,
			CategoryName: c.CategoryName,
			CategoryCode: c.CategoryCode,
			SortOrder:    c.SortOrder,
			Status:       c.Status,
			Children:     []*dto.CategoryNode{},
		}
		childMap[c.ID] = n
		if parent, ok := nodeMap[c.ParentID]; ok {
			parent.Children = append(parent.Children, n)
		}
	}
	if len(childMap) > 0 {
		s.loadChildrenRecursively(ctx, childMap)
	}
}
