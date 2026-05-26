package repository

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
)

type SysLogRepository struct {
	db *gorm.DB
}

func NewSysLogRepository(db *gorm.DB) *SysLogRepository {
	return &SysLogRepository{db: db}
}

// PageLogin 登录日志分页
func (r *SysLogRepository) PageLogin(ctx context.Context, s *dto.LoginLogSearch) ([]model.SysLoginLog, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.SysLoginLog{})
	if s.UserName != "" {
		tx = tx.Where("user_name LIKE ?", "%"+s.UserName+"%")
	}
	if s.LoginIP != "" {
		tx = tx.Where("login_ip LIKE ?", "%"+s.LoginIP+"%")
	}
	if s.LoginType != "" {
		tx = tx.Where("login_type = ?", s.LoginType)
	}
	if s.LoginResult != "" {
		tx = tx.Where("login_result = ?", s.LoginResult)
	}
	if s.StartTime != "" {
		tx = tx.Where("login_time >= ?", s.StartTime)
	}
	if s.EndTime != "" {
		tx = tx.Where("login_time <= ?", s.EndTime)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.SysLoginLog
	if err := tx.Order("login_time DESC").Offset(s.Offset()).Limit(s.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// PageOperation 操作日志分页
func (r *SysLogRepository) PageOperation(ctx context.Context, s *dto.OperationLogSearch) ([]model.SysOperationLog, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.SysOperationLog{})
	if s.UserName != "" {
		tx = tx.Where("user_name LIKE ?", "%"+s.UserName+"%")
	}
	if s.Module != "" {
		tx = tx.Where("module = ?", s.Module)
	}
	if s.Action != "" {
		tx = tx.Where("action = ?", s.Action)
	}
	if s.ClientIP != "" {
		tx = tx.Where("client_ip LIKE ?", "%"+s.ClientIP+"%")
	}
	if s.StartTime != "" {
		tx = tx.Where("operate_time >= ?", s.StartTime)
	}
	if s.EndTime != "" {
		tx = tx.Where("operate_time <= ?", s.EndTime)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.SysOperationLog
	if err := tx.Order("operate_time DESC").Offset(s.Offset()).Limit(s.Size).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
