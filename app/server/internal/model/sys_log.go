package model

import "time"

// LoginUserType 登录用户类型
type LoginUserType string

const (
	LoginUserTypeAdmin  LoginUserType = "1"
	LoginUserTypeReader LoginUserType = "2"
)

// LoginType 登录/登出
type LoginType string

const (
	LoginTypeLogin  LoginType = "1"
	LoginTypeLogout LoginType = "2"
)

// LoginResult 登录结果
type LoginResult string

const (
	LoginResultSuccess LoginResult = "1"
	LoginResultFail    LoginResult = "2"
)

// SysLoginLog 登录日志 (sys_login_log)
type SysLoginLog struct {
	ID          uint64        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserType    LoginUserType `gorm:"column:user_type;type:char(1);default:'1'" json:"userType"`
	UserID      *uint64       `gorm:"column:user_id" json:"userId,omitempty"`
	UserName    string        `gorm:"column:user_name;size:64;not null" json:"userName"`
	LoginIP     *string       `gorm:"column:login_ip;size:64" json:"loginIp,omitempty"`
	UserAgent   *string       `gorm:"column:user_agent;type:text" json:"userAgent,omitempty"`
	LoginType   LoginType     `gorm:"column:login_type;type:char(1);default:'1'" json:"loginType"`
	LoginResult LoginResult   `gorm:"column:login_result;type:char(1);not null" json:"loginResult"`
	Message     *string       `gorm:"column:message;size:255" json:"message,omitempty"`
	LoginTime   time.Time     `gorm:"column:login_time;autoCreateTime:milli" json:"loginTime"`
}

func (SysLoginLog) TableName() string { return "sys_login_log" }

// SysOperationLog 操作日志 (sys_operation_log)
type SysOperationLog struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"column:user_id;not null" json:"userId"`
	UserName      string    `gorm:"column:user_name;size:64;not null" json:"userName"`
	Module        string    `gorm:"column:module;size:64;not null" json:"module"`
	Action        string    `gorm:"column:action;size:32;not null" json:"action"`
	TargetID      *uint64   `gorm:"column:target_id" json:"targetId,omitempty"`
	TargetName    *string   `gorm:"column:target_name;size:128" json:"targetName,omitempty"`
	RequestURL    *string   `gorm:"column:request_url;size:255" json:"requestUrl,omitempty"`
	RequestMethod *string   `gorm:"column:request_method;size:16" json:"requestMethod,omitempty"`
	RequestBody   *string   `gorm:"column:request_body;type:text" json:"requestBody,omitempty"`
	ResponseCode  *int      `gorm:"column:response_code" json:"responseCode,omitempty"`
	ClientIP      *string   `gorm:"column:client_ip;size:64" json:"clientIp,omitempty"`
	UserAgent     *string   `gorm:"column:user_agent;type:text" json:"userAgent,omitempty"`
	CostMs        uint32    `gorm:"column:cost_ms;default:0" json:"costMs"`
	OperateTime   time.Time `gorm:"column:operate_time;autoCreateTime:milli" json:"operateTime"`
}

func (SysOperationLog) TableName() string { return "sys_operation_log" }
