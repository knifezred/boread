package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 后台登录
// @Summary  登录
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    body  body      dto.LoginRequest  true  "登录参数"
// @Success  200   {object}  response.Response{data=dto.LoginResponse}
// @Router   /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, code.ParamInvalid, err.Error())
		return
	}

	resp, err := h.authService.Login(
		c.Request.Context(), &req, c.ClientIP(), c.Request.UserAgent(),
	)
	if err != nil {
		switch {
		case errors.Is(err, code.ErrUserNotFound), errors.Is(err, code.ErrInvalidPassword):
			response.Error(c, code.AuthFailed, "用户名或密码错误")
		case errors.Is(err, code.ErrUserDisabled):
			response.Error(c, code.UserDisabled, "账号已禁用")
		case errors.Is(err, code.ErrUserLocked):
			response.Error(c, code.UserLocked, "账号已锁定, 请稍后再试")
		default:
			response.Error(c, code.ServerError, "登录失败")
		}
		return
	}

	response.Success(c, resp)
}

// GetUserInfo 当前登录用户信息
// @Summary   当前用户信息
// @Tags      auth
// @Security  BearerAuth
// @Produce   json
// @Success   200  {object}  response.Response{data=dto.UserInfoResponse}
// @Router    /api/auth/userInfo [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, code.AuthFailed, "unauthorized")
		return
	}

	info, err := h.authService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, info)
}

// GetUserMenu 当前用户菜单树
// @Summary   当前用户菜单树
// @Tags      auth
// @Security  BearerAuth
// @Produce   json
// @Success   200  {object}  response.Response{data=dto.MenuResponse}
// @Router    /api/auth/menu [get]
func (h *AuthHandler) GetUserMenu(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, code.AuthFailed, "unauthorized")
		return
	}

	tree, err := h.authService.GetUserMenuTree(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, tree)
}

// GetButtons 当前用户按钮码集合
// @Summary   当前用户按钮码集合
// @Tags      auth
// @Security  BearerAuth
// @Produce   json
// @Success   200  {object}  response.Response{data=[]string}
// @Router    /api/auth/buttons [get]
func (h *AuthHandler) GetButtons(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, code.AuthFailed, "unauthorized")
		return
	}

	codes, err := h.authService.GetButtons(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, code.ServerError, err.Error())
		return
	}
	response.Success(c, codes)
}

func getUserID(c *gin.Context) uint64 {
	v, ok := c.Get("user_id")
	if !ok {
		return 0
	}
	switch x := v.(type) {
	case uint64:
		return x
	case uint:
		return uint64(x)
	case int64:
		return uint64(x)
	case string:
		id, _ := strconv.ParseUint(x, 10, 64)
		return id
	}
	return 0
}
