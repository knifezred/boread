package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
	jwtPkg "boread/pkg/jwt"
)

// UgreenAuthService 绿联认证服务
type UgreenAuthService struct {
	userRepo *repository.SysUserRepository
	roleRepo *repository.SysRoleRepository
	db       *gorm.DB
}

func NewUgreenAuthService(userRepo *repository.SysUserRepository, roleRepo *repository.SysRoleRepository, db *gorm.DB) *UgreenAuthService {
	return &UgreenAuthService{userRepo: userRepo, roleRepo: roleRepo, db: db}
}

// LoginOrRegister 绿联用户登录或自动注册
// ugreenUserID: 绿联平台用户ID
// ugreenUserName: 绿联平台用户名
// ugreenUserType: 绿联平台用户类型 (admin/users)
func (s *UgreenAuthService) LoginOrRegister(ctx context.Context, ugreenUserID, ugreenUserName, ugreenUserType string) (*dto.UgreenLoginResponse, error) {
	if ugreenUserID == "" {
		return nil, code.ErrUgreenAuthFailed
	}

	// 1. 按 ugreen_user_id 查找已有映射用户
	user, err := s.FindByUgreenID(ctx, ugreenUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 2. 不存在则自动创建本地用户
	if user == nil {
		user, err = s.createUgreenUser(ctx, ugreenUserID, ugreenUserName, ugreenUserType)
		if err != nil {
			return nil, err
		}
	}

	// 3. 生成 JWT token
	now := time.Now()
	_ = s.userRepo.UpdateLoginSuccess(ctx, user.ID, "", now)

	token, expiresAt, err := jwtPkg.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExpiresAt, err := jwtPkg.GenerateRefreshToken(user.ID, user.UserName)
	if err != nil {
		return nil, err
	}

	return &dto.UgreenLoginResponse{
		Token:            token,
		RefreshToken:     refreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

// FindByUgreenID 按绿联用户ID查找本地用户
func (s *UgreenAuthService) FindByUgreenID(ctx context.Context, ugreenUserID string) (*model.SysUser, error) {
	var u model.SysUser
	err := s.db.WithContext(ctx).Where("ugreen_user_id = ?", ugreenUserID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// createUgreenUser 创建绿联映射的本地用户
func (s *UgreenAuthService) createUgreenUser(ctx context.Context, ugreenUserID, ugreenUserName, ugreenUserType string) (*model.SysUser, error) {
	var user model.SysUser

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 生成唯一用户名
		userName := fmt.Sprintf("ug_%s", ugreenUserID)
		if len(userName) > 64 {
			userName = userName[:64]
		}

		nickName := ugreenUserName
		if nickName == "" {
			nickName = "绿联用户"
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user = model.SysUser{
			UserName:     userName,
			Password:     string(hashed), // 绿联用户密码默认123456
			NickName:     nickName,
			UgreenUserID: &ugreenUserID,
			Status:       model.StatusEnabled,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 分配默认角色，根据ugreenUserType判断，admin映射管理员，users映射用户
		var defaultRole *model.SysRole
		var roleErr error
		if ugreenUserType == "admin" {
			defaultRole, roleErr = s.roleRepo.GetByCode(ctx, "SUPER_ADMIN")
		} else {
			defaultRole, roleErr = s.roleRepo.GetByCode(ctx, "USERS")
		}
		if roleErr == nil && defaultRole != nil {
			ur := model.SysUserRole{
				UserID: user.ID,
				RoleID: defaultRole.ID,
			}
			if err := tx.Create(&ur).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}
