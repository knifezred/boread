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
	ErrDeptCodeExists  = errors.New("部门编码已存在")
	ErrDeptHasChildren = errors.New("存在子部门, 不能删除")
	ErrDeptHasUsers    = errors.New("部门下有用户, 不能删除")
	ErrParentNotFound  = errors.New("父部门不存在")
)

type DeptService struct {
	repo *repository.SysDeptRepository
}

func NewDeptService(repo *repository.SysDeptRepository) *DeptService {
	return &DeptService{repo: repo}
}

// Create 新建部门
func (s *DeptService) Create(ctx context.Context, req *dto.DeptRequest, userID uint64) (*model.SysDept, error) {
	if _, err := s.repo.GetByCode(ctx, req.DeptCode); err == nil {
		return nil, ErrDeptCodeExists
	}

	ancestors := "0"
	if req.ParentID > 0 {
		parent, err := s.repo.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, ErrParentNotFound
		}
		ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
	}

	status := req.Status
	if status == "" {
		status = model.StatusEnabled
	}
	m := &model.SysDept{
		ParentID:  req.ParentID,
		Ancestors: ancestors,
		DeptName:  req.DeptName,
		DeptCode:  req.DeptCode,
		SortOrder: req.SortOrder,
		Status:    status,
	}
	if req.Leader != "" {
		m.Leader = &req.Leader
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// Update 更新部门
func (s *DeptService) Update(ctx context.Context, id uint64, req *dto.DeptRequest, userID uint64) (*model.SysDept, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 编码变更需查重
	if req.DeptCode != m.DeptCode {
		if _, e := s.repo.GetByCode(ctx, req.DeptCode); e == nil {
			return nil, ErrDeptCodeExists
		}
	}

	// 父级变更需重算 ancestors
	if req.ParentID != m.ParentID {
		ancestors := "0"
		if req.ParentID > 0 {
			parent, e := s.repo.GetByID(ctx, req.ParentID)
			if e != nil {
				return nil, ErrParentNotFound
			}
			ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
		}
		m.Ancestors = ancestors
		m.ParentID = req.ParentID
	}

	m.DeptName = req.DeptName
	m.DeptCode = req.DeptCode
	m.SortOrder = req.SortOrder
	if req.Status != "" {
		m.Status = req.Status
	}
	if req.Leader != "" {
		m.Leader = &req.Leader
	}
	m.UpdateBy = &userID

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// Delete 删除部门 (检查子部门 + 用户)
func (s *DeptService) Delete(ctx context.Context, id uint64) error {
	if has, err := s.repo.HasChildren(ctx, id); err != nil {
		return err
	} else if has {
		return ErrDeptHasChildren
	}
	if has, err := s.repo.HasUsers(ctx, id); err != nil {
		return err
	} else if has {
		return ErrDeptHasUsers
	}
	return s.repo.Delete(ctx, id)
}

// GetByID 获取部门详情
func (s *DeptService) GetByID(ctx context.Context, id uint64) (*model.SysDept, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("dept %d not found", id)
		}
		return nil, err
	}
	return m, nil
}

// Tree 部门树 (返回根节点列表)
func (s *DeptService) Tree(ctx context.Context, search *dto.DeptSearch) ([]*dto.DeptNode, error) {
	rows, err := s.repo.ListAll(ctx, search.DeptName, search.DeptCode, string(search.Status))
	if err != nil {
		return nil, err
	}
	return buildDeptTree(rows), nil
}

func buildDeptTree(rows []model.SysDept) []*dto.DeptNode {
	nodeMap := make(map[uint64]*dto.DeptNode, len(rows))
	for _, r := range rows {
		node := &dto.DeptNode{
			ID:        r.ID,
			ParentID:  r.ParentID,
			DeptName:  r.DeptName,
			DeptCode:  r.DeptCode,
			SortOrder: r.SortOrder,
			Status:    r.Status,
			Children:  []*dto.DeptNode{}, // 初始化空数组
		}
		if r.Leader != nil {
			node.Leader = *r.Leader
		}
		nodeMap[r.ID] = node
	}
	roots := make([]*dto.DeptNode, 0) // 初始化为空数组而非nil
	for _, n := range nodeMap {
		if parent, ok := nodeMap[n.ParentID]; ok {
			parent.Children = append(parent.Children, n)
		} else {
			roots = append(roots, n)
		}
	}
	return roots
}

// Page 部门分页列表 (树形结构, 基于顶级部门分页)
func (s *DeptService) Page(ctx context.Context, req *dto.DeptSearch) (*dto.PageResponse, error) {
	req.Normalize()
	// 1. 分页查询顶级部门（已按 sort_order asc 排序）
	topDepts, total, err := s.repo.PageTop(ctx, req.DeptName, string(req.Status), req.Current, req.Size)
	if err != nil {
		return nil, err
	}
	if len(topDepts) == 0 {
		return dto.NewPageResponse([]*dto.DeptNode{}, total, &req.PageRequest), nil
	}
	// 2. 初始化节点 map & 根节点列表（按查询顺序保持排序）
	nodeMap := make(map[uint64]*dto.DeptNode, len(topDepts))
	roots := make([]*dto.DeptNode, 0, len(topDepts))
	for _, r := range topDepts {
		leader := ""
		if r.Leader != nil {
			leader = *r.Leader
		}
		node := &dto.DeptNode{
			ID:        r.ID,
			ParentID:  r.ParentID,
			DeptName:  r.DeptName,
			DeptCode:  r.DeptCode,
			Leader:    leader,
			SortOrder: r.SortOrder,
			Status:    r.Status,
			Children:  []*dto.DeptNode{},
		}
		nodeMap[r.ID] = node
		roots = append(roots, node)
	}
	// 3. 层级查询所有子部门，按查询顺序直接挂载（子部门查询已按 sort_order asc 排序）
	currentParentIDs := make([]uint64, 0, len(topDepts))
	for _, r := range topDepts {
		currentParentIDs = append(currentParentIDs, r.ID)
	}
	for i := 0; i < 10; i++ { // 限制最多10级，避免循环嵌套
		if len(currentParentIDs) == 0 {
			break
		}
		children, err := s.repo.ListByParentIDs(ctx, currentParentIDs)
		if err != nil {
			return nil, err
		}
		if len(children) == 0 {
			break
		}
		nextParentIDs := make([]uint64, 0, len(children))
		for _, r := range children {
			leader := ""
			if r.Leader != nil {
				leader = *r.Leader
			}
			node := &dto.DeptNode{
				ID:        r.ID,
				ParentID:  r.ParentID,
				DeptName:  r.DeptName,
				DeptCode:  r.DeptCode,
				Leader:    leader,
				SortOrder: r.SortOrder,
				Status:    r.Status,
				Children:  []*dto.DeptNode{},
			}
			nodeMap[r.ID] = node
			if parent, ok := nodeMap[r.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
			nextParentIDs = append(nextParentIDs, r.ID)
		}
		currentParentIDs = nextParentIDs
	}
	return dto.NewPageResponse(roots, total, &req.PageRequest), nil
}
