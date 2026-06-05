package dto

// UgreenLoginRequest 绿联登录请求
type UgreenLoginRequest struct {
	UgreenToken string `json:"ugreenToken"`
}

// UgreenLoginResponse 绿联登录响应
type UgreenLoginResponse struct {
	Token            string `json:"token"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresAt        int64  `json:"expiresAt"`
	RefreshExpiresAt int64  `json:"refreshExpiresAt"`
}

// UgreenProfileResponse 绿联用户信息响应
type UgreenProfileResponse struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	UserType string `json:"userType"`
	IsNew    bool   `json:"isNew"`
}
