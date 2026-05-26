package service

import (
	"context"
	"errors"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

var (
	ErrRoleCodeExists    = errors.New("角色编码已存在")
	ErrRoleSystem        = errors.New("系统内置角色不可操作")
	ErrRoleHasUsers      = errors.New("角色下还有用户, 不能删除")
)

type RoleService struct {
	repo *repository.SysRoleRepository
}

func NewRoleService(repo *repository.SysRoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) Create(ctx context.Context, req *dto.RoleRequest, userID uint64) (*model.SysRole, error) {
	if _, err := s.repo.GetByCode(ctx, req.RoleCode); err == nil {
		return nil, ErrRoleCodeExists
	}
	status := req.Status
	if status == "" {
		status = model.StatusEnabled
	}
	m := &model.SysRole{
		RoleName:  req.RoleName,
		RoleCode:  req.RoleCode,
		DataScope: req.DataScope,
		SortOrder: req.SortOrder,
		Status:    status,
	}
	if req.RoleDesc != "" {
		m.RoleDesc = &req.RoleDesc
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	if m.DataScope == model.DataScopeCustom {
		_ = s.repo.ReplaceDepts(ctx, m.ID, req.DeptIDs)
	}
	return m, nil
}

func (s *RoleService) Update(ctx context.Context, id uint64, req *dto.RoleRequest, userID uint64) (*model.SysRole, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m.IsSystem && req.RoleCode != m.RoleCode {
		return nil, ErrRoleSystem // 系统角色不允许改 code
	}
	if req.RoleCode != m.RoleCode {
		if _, e := s.repo.GetByCode(ctx, req.RoleCode); e == nil {
			return nil, ErrRoleCodeExists
		}
	}
	m.RoleName = req.RoleName
	m.RoleCode = req.RoleCode
	m.DataScope = req.DataScope
	m.SortOrder = req.SortOrder
	if req.Status != "" {
		m.Status = req.Status
	}
	if req.RoleDesc != "" {
		m.RoleDesc = &req.RoleDesc
	}
	m.UpdateBy = &userID
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	if m.DataScope == model.DataScopeCustom {
		_ = s.repo.ReplaceDepts(ctx, m.ID, req.DeptIDs)
	}
	return m, nil
}

func (s *RoleService) Delete(ctx context.Context, id uint64) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if m.IsSystem {
		return ErrRoleSystem
	}
	if has, err := s.repo.HasUsers(ctx, id); err != nil {
		return err
	} else if has {
		return ErrRoleHasUsers
	}
	return s.repo.Delete(ctx, id)
}

func (s *RoleService) GetByID(ctx context.Context, id uint64) (*model.SysRole, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RoleService) Page(ctx context.Context, req *dto.RoleSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	return dto.NewPageResponse(rows, total, &req.PageRequest), nil
}

func (s *RoleService) AllBrief(ctx context.Context) ([]dto.RoleBrief, error) {
	rows, err := s.repo.AllBrief(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RoleBrief, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.RoleBrief{ID: r.ID, RoleName: r.RoleName, RoleCode: r.RoleCode})
	}
	return out, nil
}

// GrantMenus 授权菜单
func (s *RoleService) GrantMenus(ctx context.Context, roleID uint64, menuIDs []uint64) error {
	if _, err := s.repo.GetByID(ctx, roleID); err != nil {
		return err
	}
	return s.repo.ReplaceMenus(ctx, roleID, menuIDs)
}

// GrantButtons 授权按钮
func (s *RoleService) GrantButtons(ctx context.Context, roleID uint64, buttonIDs []uint64) error {
	if _, err := s.repo.GetByID(ctx, roleID); err != nil {
		return err
	}
	return s.repo.ReplaceButtons(ctx, roleID, buttonIDs)
}

func (s *RoleService) GetMenuIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	return s.repo.GetMenuIDs(ctx, roleID)
}

func (s *RoleService) GetButtonIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	return s.repo.GetButtonIDs(ctx, roleID)
}
