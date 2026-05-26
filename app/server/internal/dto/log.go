package dto

import "boread/internal/model"

// LoginLogSearch 登录日志分页搜索
type LoginLogSearch struct {
	PageRequest
	UserName    string             `form:"userName"`
	LoginIP     string             `form:"loginIp"`
	LoginType   model.LoginType    `form:"loginType"`
	LoginResult model.LoginResult  `form:"loginResult"`
	StartTime   string             `form:"startTime"` // YYYY-MM-DD HH:MM:SS
	EndTime     string             `form:"endTime"`
}

// OperationLogSearch 操作日志分页搜索
type OperationLogSearch struct {
	PageRequest
	UserName  string `form:"userName"`
	Module    string `form:"module"`
	Action    string `form:"action"`
	ClientIP  string `form:"clientIp"`
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}
