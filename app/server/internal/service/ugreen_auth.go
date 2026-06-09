package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
// ip, ua: 客户端 IP 和 User-Agent，用于登录日志
func (s *UgreenAuthService) LoginOrRegister(ctx context.Context, ugreenUserID, ugreenUserName, ugreenUserType, ip, ua string) (*dto.UgreenLoginResponse, error) {
	if ugreenUserID == "" {
		return nil, code.ErrUgreenAuthFailed
	}

	// 1. 按 ugreen_user_id 查找已有映射用户
	user, err := s.FindByUgreenID(ctx, ugreenUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.writeLoginLog(ctx, nil, ugreenUserName, ip, ua, model.LoginResultFail, "query failed: "+err.Error())
		return nil, err
	}

	// 2. 不存在则自动创建本地用户
	if user == nil {
		user, err = s.createUgreenUser(ctx, ugreenUserID, ugreenUserName, ugreenUserType)
		if err != nil {
			s.writeLoginLog(ctx, nil, ugreenUserName, ip, ua, model.LoginResultFail, "create user failed: "+err.Error())
			return nil, err
		}
	}

	// 3. 生成 JWT token
	now := time.Now()
	_ = s.userRepo.UpdateLoginSuccess(ctx, user.ID, "", now)

	token, expiresAt, err := jwtPkg.GenerateToken(user.ID, user.UserName)
	if err != nil {
		s.writeLoginLog(ctx, &user.ID, user.UserName, ip, ua, model.LoginResultFail, "generate token failed")
		return nil, err
	}
	refreshToken, refreshExpiresAt, err := jwtPkg.GenerateRefreshToken(user.ID, user.UserName)
	if err != nil {
		s.writeLoginLog(ctx, &user.ID, user.UserName, ip, ua, model.LoginResultFail, "generate refresh token failed")
		return nil, err
	}

	s.writeLoginLog(ctx, &user.ID, user.UserName, ip, ua, model.LoginResultSuccess, "ugreen login")

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

		// 绿联用户只能通过网关授权登录，设置随机不可猜测密码禁用密码登录
		randomPwd := generateRandomPassword(32)
		hashed, err := bcrypt.GenerateFromPassword([]byte(randomPwd), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user = model.SysUser{
			UserName:     userName,
			Password:     string(hashed), // 随机密码，无人知晓
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

// writeLoginLog 写绿联登录日志，失败不影响主流程
func (s *UgreenAuthService) writeLoginLog(ctx context.Context, userID *uint64, userName, ip, ua string, result model.LoginResult, msg string) {
	log := &model.SysLoginLog{
		UserType:    model.LoginUserTypeAdmin,
		UserID:      userID,
		UserName:    userName,
		LoginIP:     &ip,
		UserAgent:   &ua,
		LoginType:   model.LoginTypeLogin,
		LoginResult: result,
		Message:     &msg,
	}
	_ = s.db.WithContext(ctx).Create(log).Error
}

// generateRandomPassword 生成指定长度的随机十六进制密码
func generateRandomPassword(length int) string {
	buf := make([]byte, (length+1)/2)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)[:length]
}
