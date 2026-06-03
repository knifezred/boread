package service

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
	jwtPkg "boread/pkg/jwt"
)

const (
	maxPwdErrorCount = 5
	lockDuration     = 15 * time.Minute
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserDisabled    = errors.New("user is disabled")
	ErrUserLocked      = errors.New("user is locked")
)

// AuthService 认证服务
type AuthService struct {
	userRepo        *repository.SysUserRepository
	chapterRuleRepo *repository.BookChapterRuleRepository
	db              *gorm.DB
}

func NewAuthService(userRepo *repository.SysUserRepository, chapterRuleRepo *repository.BookChapterRuleRepository, db *gorm.DB) *AuthService {
	return &AuthService{userRepo: userRepo, chapterRuleRepo: chapterRuleRepo, db: db}
}

// Login 登录: 校验密码 + 风控 + 记录日志 + 签发 token
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest, ip, ua string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByUserName(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.writeLoginLog(ctx, nil, req.Username, ip, ua, model.LoginResultFail, "user not found")
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.Status == model.StatusDisabled {
		s.writeLoginLog(ctx, &user.ID, user.UserName, ip, ua, model.LoginResultFail, "user disabled")
		return nil, ErrUserDisabled
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		s.writeLoginLog(ctx, &user.ID, user.UserName, ip, ua, model.LoginResultFail, "user locked")
		return nil, ErrUserLocked
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		_ = s.userRepo.IncrErrorCount(ctx, user.ID)
		if user.PwdErrorCount+1 >= maxPwdErrorCount {
			lockUntil := time.Now().Add(lockDuration)
			_ = s.userRepo.LockUser(ctx, user.ID, lockUntil)
		}
		s.writeLoginLog(ctx, &user.ID, user.UserName, ip, ua, model.LoginResultFail, "wrong password")
		return nil, ErrInvalidPassword
	}

	now := time.Now()
	if err := s.userRepo.UpdateLoginSuccess(ctx, user.ID, ip, now); err != nil {
		return nil, err
	}

	token, expiresAt, err := jwtPkg.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return nil, err
	}
	// 生成refresh token
	refreshToken, refreshExpiresAt, err := jwtPkg.GenerateRefreshToken(user.ID, user.UserName)
	if err != nil {
		return nil, err
	}

	s.writeLoginLog(ctx, &user.ID, user.UserName, ip, ua, model.LoginResultSuccess, "ok")

	// 登录成功后初始化用户章节识别规则（如果尚未初始化）
	s.initUserChapterRules(ctx, user.ID)

	return &dto.LoginResponse{
		Token:            token,
		RefreshToken:     refreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

// GetUserInfo 获取登录用户的角色 + 按钮码集合 (供前端权限渲染)
func (s *AuthService) GetUserInfo(ctx context.Context, userID uint64) (*dto.UserInfoResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles, err := s.userRepo.GetRoleCodesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roleIDs, err := s.userRepo.GetRoleIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	buttons, err := s.userRepo.GetButtonCodesByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	return &dto.UserInfoResponse{
		UserID:   strconv.FormatUint(user.ID, 10),
		UserName: user.UserName,
		Roles:    roles,
		Buttons:  buttons,
	}, nil
}

// GetButtons 获取用户的按钮码集合 (用于鉴权中间件缓存)
func (s *AuthService) GetButtons(ctx context.Context, userID uint64) ([]string, error) {
	roleIDs, err := s.userRepo.GetRoleIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetButtonCodesByRoleIDs(ctx, roleIDs)
}

// GetUserMenuTree 获取用户菜单树 (角色聚合并构建树形)
func (s *AuthService) GetUserMenuTree(ctx context.Context, userID uint64) (*dto.MenuResponse, error) {
	roleIDs, err := s.userRepo.GetRoleIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	menus, err := s.userRepo.GetMenusByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	routes := buildMenuTree(menus)

	// 找到第一个非目录的菜单作为首页路由名
	home := "home"
	for _, m := range menus {
		if m.MenuType == model.MenuTypeMenu && !m.HideInMenu {
			home = m.RouteName
			break
		}
	}

	return &dto.MenuResponse{
		Routes: routes,
		Home:   home,
	}, nil
}

// buildMenuTree 把菜单列表构建为前端期望的树形结构
func buildMenuTree(menus []model.SysMenu) []dto.MenuRoute {
	// 按 parent_id 分组
	childrenMap := make(map[uint64][]model.SysMenu)
	for _, m := range menus {
		childrenMap[m.ParentID] = append(childrenMap[m.ParentID], m)
	}
	// 每组按 sort_order 升序
	for k := range childrenMap {
		sort.Slice(childrenMap[k], func(i, j int) bool {
			return childrenMap[k][i].SortOrder < childrenMap[k][j].SortOrder
		})
	}
	return buildLevel(0, childrenMap)
}

func buildLevel(parentID uint64, childrenMap map[uint64][]model.SysMenu) []dto.MenuRoute {
	siblings := childrenMap[parentID]
	if len(siblings) == 0 {
		return []dto.MenuRoute{} // 返回空数组而非nil
	}
	out := make([]dto.MenuRoute, 0, len(siblings))
	for _, m := range siblings {
		r := dto.MenuRoute{
			Name: m.RouteName,
			Path: m.RoutePath,
			Meta: dto.MenuMeta{
				Title:           m.MenuName,
				Order:           m.SortOrder,
				HideInMenu:      m.HideInMenu,
				KeepAlive:       m.KeepAlive,
				Constant:        m.Constant,
				MultiTab:        m.MultiTab,
				FixedIndexInTab: m.FixedIndexInTab,
			},
		}
		if m.Component != nil {
			r.Component = *m.Component
		}
		if m.Icon != nil {
			switch m.IconType {
			case model.IconTypeLocal:
				r.Meta.LocalIcon = *m.Icon
			default:
				r.Meta.Icon = *m.Icon
			}
		}
		if m.I18nKey != nil {
			r.Meta.I18nKey = *m.I18nKey
		}
		if m.Href != nil {
			r.Meta.Href = *m.Href
		}
		if m.ActiveMenu != nil {
			r.Meta.ActiveMenu = *m.ActiveMenu
		}
		children := buildLevel(m.ID, childrenMap)
		if len(children) > 0 {
			r.Children = children
			// 目录类型默认重定向到第一个子菜单
			if m.MenuType == model.MenuTypeDir && len(children) > 0 {
				r.Redirect = children[0].Path
			}
		}
		out = append(out, r)
	}
	return out
}

// initUserChapterRules 初始化用户章节识别规则
// 检查用户是否已有自定义规则，如果没有则复制系统默认规则
func (s *AuthService) initUserChapterRules(ctx context.Context, userID uint64) {
	// 检查用户是否已有规则
	existing, err := s.chapterRuleRepo.ListByUserID(ctx, userID)
	if err != nil || len(existing) > 0 {
		return // 已有规则或查询出错，不做操作
	}

	// 获取系统默认规则
	systemRules, err := s.chapterRuleRepo.ListSystemDefaults(ctx)
	if err != nil || len(systemRules) == 0 {
		return
	}

	// 逐条复制为用户自定义规则
	for _, sr := range systemRules {
		cr := &model.BookChapterRule{
			RuleName:      sr.RuleName,
			RuleType:      model.RuleTypeCustom,
			UserID:        &userID,
			TitlePattern:  sr.TitlePattern,
			GroupPattern:  sr.GroupPattern,
			MinChapterLen: sr.MinChapterLen,
			MaxChapterLen: sr.MaxChapterLen,
			SortOrder:     sr.SortOrder,
			Description:   sr.Description,
			Status:        sr.Status,
		}
		cr.CreateBy = &userID
		cr.UpdateBy = &userID
		_ = s.chapterRuleRepo.Create(ctx, cr)
	}
}

// writeLoginLog 写登录日志, 失败不影响主流程
func (s *AuthService) writeLoginLog(ctx context.Context, userID *uint64, userName, ip, ua string, result model.LoginResult, msg string) {
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
