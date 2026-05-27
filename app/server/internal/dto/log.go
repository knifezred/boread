package dto

import "boread/internal/model"

// LoginLogSearch 登录日志分页搜索
type LoginLogSearch struct {
	PageRequest
	UserName    string             `json:"userName"`
	LoginIP     string             `json:"loginIp"`
	LoginType   model.LoginType    `json:"loginType"`
	LoginResult model.LoginResult  `json:"loginResult"`
	StartTime   string             `json:"startTime"` // YYYY-MM-DD HH:MM:SS
	EndTime     string             `json:"endTime"`
}

// OperationLogSearch 操作日志分页搜索
type OperationLogSearch struct {
	PageRequest
	UserName  string `json:"userName"`
	Module    string `json:"module"`
	Action    string `json:"action"`
	ClientIP  string `json:"clientIp"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
