package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

type UserService struct {
	repo *repository.SysUserRepository
}

func NewUserService(repo *repository.SysUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req *dto.UserCreateRequest, opID uint64) (*model.SysUser, error) {
	if _, err := s.repo.GetByUserName(ctx, req.UserName); err == nil {
		return nil, code.ErrUserExists
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	status := req.Status
	if status == "" {
		status = model.StatusEnabled
	}
	now := time.Now()
	m := &model.SysUser{
		DeptID:       req.DeptID,
		UserName:     req.UserName,
		Password:     string(hashed),
		PwdUpdatedAt: &now,
		NickName:     req.NickName,
		UserGender:   req.UserGender,
		Status:       status,
	}
	if req.UserPhone != "" {
		m.UserPhone = &req.UserPhone
	}
	if req.UserEmail != "" {
		m.UserEmail = &req.UserEmail
	}
	if req.Avatar != "" {
		m.Avatar = &req.Avatar
	}
	m.CreateBy = &opID
	m.UpdateBy = &opID
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	if len(req.RoleIDs) > 0 {
		if err := s.repo.ReplaceRoles(ctx, m.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (s *UserService) Update(ctx context.Context, id uint64, req *dto.UserUpdateRequest, opID uint64) (*model.SysUser, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	m.DeptID = req.DeptID
	m.NickName = req.NickName
	m.UserGender = req.UserGender
	if req.UserPhone != "" {
		m.UserPhone = &req.UserPhone
	}
	if req.UserEmail != "" {
		m.UserEmail = &req.UserEmail
	}
	if req.Avatar != "" {
		m.Avatar = &req.Avatar
	}
	if req.Status != "" {
		m.Status = req.Status
	}
	m.UpdateBy = &opID
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	if req.RoleIDs != nil {
		if err := s.repo.ReplaceRoles(ctx, m.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (s *UserService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) BatchDelete(ctx context.Context, ids []uint64) error {
	return s.repo.BatchDelete(ctx, ids)
}

func (s *UserService) GetByID(ctx context.Context, id uint64) (*model.SysUser, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Page(ctx context.Context, req *dto.UserSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	// 批量补 roles
	vos := make([]dto.UserVO, 0, len(rows))
	for _, u := range rows {
		roles, _ := s.repo.GetRoleCodesByUserID(ctx, u.ID)
		vos = append(vos, toUserVO(&u, roles))
	}
	return dto.NewPageResponse(vos, total, &req.PageRequest), nil
}

// ResetPassword 重置密码 (清除锁定状态)
func (s *UserService) ResetPassword(ctx context.Context, id uint64, newPwd string, opID uint64) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdateFields(ctx, id, map[string]any{
		"password":        string(hashed),
		"pwd_updated_at":  time.Now(),
		"pwd_error_count": 0,
		"locked_until":    nil,
		"update_by":       opID,
	})
}

func toUserVO(u *model.SysUser, roles []string) dto.UserVO {
	if roles == nil {
		roles = []string{}
	}
	return dto.UserVO{
		ID:         u.ID,
		DeptID:     u.DeptID,
		UserName:   u.UserName,
		NickName:   u.NickName,
		UserGender: u.UserGender,
		UserPhone:  u.UserPhone,
		UserEmail:  u.UserEmail,
		Avatar:     u.Avatar,
		Status:     u.Status,
		UserRoles:  roles,
	}
}
