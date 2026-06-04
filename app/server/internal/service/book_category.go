package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"gorm.io/gorm"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

type BookCategoryService struct {
	repo *repository.BookCategoryRepository
}

func NewBookCategoryService(repo *repository.BookCategoryRepository) *BookCategoryService {
	return &BookCategoryService{repo: repo}
}

func (s *BookCategoryService) Create(ctx context.Context, req *dto.CategoryRequest, userID uint64) (*model.BookCategory, error) {
	if _, err := s.repo.GetByCode(ctx, req.CategoryCode); err == nil {
		return nil, code.ErrCategoryCodeExists
	}

	ancestors := "0"
	if req.ParentID > 0 {
		parent, err := s.repo.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, code.ErrCategoryParentNotFound
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
		Description:  req.Description,
		SortOrder:    req.SortOrder,
		Status:       status,
	}
	if req.IsHot != nil {
		m.IsHot = *req.IsHot
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
			return nil, code.ErrCategoryCodeExists
		}
	}

	if req.ParentID != m.ParentID {
		ancestors := "0"
		if req.ParentID > 0 {
			parent, e := s.repo.GetByID(ctx, req.ParentID)
			if e != nil {
				return nil, code.ErrCategoryParentNotFound
			}
			ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
		}
		m.Ancestors = ancestors
		m.ParentID = req.ParentID
	}

	m.CategoryName = req.CategoryName
	m.CategoryCode = req.CategoryCode
	m.Description = req.Description
	m.SortOrder = req.SortOrder
	if req.IsHot != nil {
		m.IsHot = *req.IsHot
	}
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
		return code.ErrCategoryHasChildren
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
			Description:  r.Description,
			SortOrder:    r.SortOrder,
			IsHot:        r.IsHot,
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
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].SortOrder < roots[j].SortOrder
	})
	for _, n := range nodeMap {
		sort.Slice(n.Children, func(i, j int) bool {
			return n.Children[i].SortOrder < n.Children[j].SortOrder
		})
	}
	return roots
}

// Page 分类分页列表 (树形, 基于顶级分类分页 + 递归加载子级)
func (s *BookCategoryService) Page(ctx context.Context, req *dto.CategorySearch) (*dto.PageResponse, error) {
	req.Normalize()
	topRows, total, err := s.repo.PageTop(ctx, req.CategoryName, req.CategoryCode, req.ParentID, req.IsHot, string(req.Status), req.Current, req.Size)
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
			Description:  r.Description,
			SortOrder:    r.SortOrder,
			IsHot:        r.IsHot,
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
			Description:  c.Description,
			SortOrder:    c.SortOrder,
			IsHot:        c.IsHot,
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

// GetHotCategories 获取所有已启用的热门分类
func (s *BookCategoryService) GetHotCategories(ctx context.Context) ([]dto.HotCategoryItem, error) {
	rows, err := s.repo.ListHot(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]dto.HotCategoryItem, len(rows))
	for i, r := range rows {
		items[i] = dto.HotCategoryItem{
			ID:           r.ID,
			CategoryName: r.CategoryName,
			CategoryCode: r.CategoryCode,
		}
	}
	return items, nil
}
