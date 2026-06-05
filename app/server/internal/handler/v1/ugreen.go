package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/service"
	"boread/pkg/response"
)

// UgreenHandler 绿联认证处理器
type UgreenHandler struct {
	ugreenAuthService *service.UgreenAuthService
}

func NewUgreenHandler(ugreenAuthService *service.UgreenAuthService) *UgreenHandler {
	return &UgreenHandler{ugreenAuthService: ugreenAuthService}
}

// Profile 获取绿联用户信息（登录前确认用）
// @Summary  获取绿联用户信息
// @Tags     auth
// @Accept   json
// @Produce  json
// @Success  200  {object}  response.Response{data=dto.UgreenProfileResponse}
// @Router   /api/auth/ugreen-profile [get]
func (h *UgreenHandler) Profile(c *gin.Context) {
	ugreenUserID := c.GetHeader("Ugreen-User-ID")
	ugreenUserName := c.GetHeader("Ugreen-User-Name")
	ugreenUserType := c.GetHeader("Ugreen-User-Type")

	if ugreenUserID == "" {
		response.Error(c, code.UgreenAuthFailed, "missing ugreen user id, not in ugos gateway environment")
		return
	}

	// 查询本地是否已有映射用户
	user, err := h.ugreenAuthService.FindByUgreenID(c.Request.Context(), ugreenUserID)
	isNew := err != nil // 不存在则为新用户

	resp := dto.UgreenProfileResponse{
		UserID:   ugreenUserID,
		UserName: ugreenUserName,
		UserType: ugreenUserType,
		IsNew:    isNew,
	}

	if user != nil && user.NickName != "" {
		resp.UserName = user.NickName
	}

	response.Success(c, resp)
}

// Login 绿联登录
// @Summary  绿联NAS登录
// @Tags     auth
// @Accept   json
// @Produce  json
// @Success  200  {object}  response.Response{data=dto.UgreenLoginResponse}
// @Router   /api/auth/ugreen-login [post]
//
// 该接口需要运行在UGOS系统网关之后，网关会在请求头中注入:
//   - Ugreen-User-ID:    用户ID
//   - Ugreen-User-Name:  用户名
//   - Ugreen-User-Type:  用户身份 (admin/users)
func (h *UgreenHandler) Login(c *gin.Context) {
	ugreenUserID := c.GetHeader("Ugreen-User-ID")
	ugreenUserName := c.GetHeader("Ugreen-User-Name")
	ugreenUserType := c.GetHeader("Ugreen-User-Type")

	if ugreenUserID == "" {
		response.Error(c, code.UgreenAuthFailed, "missing ugreen user id, not in ugos gateway environment")
		return
	}

	resp, err := h.ugreenAuthService.LoginOrRegister(
		c.Request.Context(), ugreenUserID, ugreenUserName, ugreenUserType,
	)
	if err != nil {
		if errors.Is(err, code.ErrUgreenAuthFailed) {
			response.Error(c, code.UgreenAuthFailed, "ugreen auth failed")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, resp)
}
