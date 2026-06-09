package dto

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

// LoginResponse 登录响应 (对齐前端 Api.Auth.LoginToken)
type LoginResponse struct {
	Token            string `json:"token"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresAt        int64  `json:"expiresAt"`
	RefreshExpiresAt int64  `json:"refreshExpiresAt"`
}

// UserInfoResponse 当前用户信息 (对齐前端 Api.Auth.UserInfo)
type UserInfoResponse struct {
	UserID   string   `json:"userId"`
	UserName string   `json:"userName"`
	NickName string   `json:"nickName"`
	Roles    []string `json:"roles"`
	Buttons  []string `json:"buttons"`
}

// MenuRoute 前端 ElegantConstRoute 子集
type MenuRoute struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Component string      `json:"component"`
	Redirect  string      `json:"redirect"`
	Meta      MenuMeta    `json:"meta"`
	Children  []MenuRoute `json:"children"`
}

// MenuMeta 路由 meta
type MenuMeta struct {
	Title           string   `json:"title"`
	I18nKey         string   `json:"i18nKey"`
	Icon            string   `json:"icon"`
	LocalIcon       string   `json:"localIcon"`
	Order           int      `json:"order"`
	HideInMenu      bool     `json:"hideInMenu"`
	KeepAlive       bool     `json:"keepAlive"`
	Constant        bool     `json:"constant"`
	Href            string   `json:"href"`
	ActiveMenu      string   `json:"activeMenu"`
	MultiTab        bool     `json:"multiTab"`
	FixedIndexInTab *int     `json:"fixedIndexInTab"`
	Roles           []string `json:"roles"`
}

// MenuResponse /api/auth/menu 响应 (对齐前端 Api.Route.MenuRoute 结构)
type MenuResponse struct {
	Routes []MenuRoute `json:"routes"`
	Home   string      `json:"home"`
}
