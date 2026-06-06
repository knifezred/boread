package service

import (
	"context"
	"time"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

type SettingService struct {
	repo *repository.SysSettingRepository
}

func NewSettingService(repo *repository.SysSettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) Create(ctx context.Context, req *dto.SettingRequest, opID uint64) (*model.SysSetting, error) {
	if _, err := s.repo.GetByCategoryKey(ctx, req.Category, req.Key); err == nil {
		return nil, code.ErrSettingKeyExists
	}
	status := req.Status
	if status == "" {
		status = model.StatusEnabled
	}
	valueType := req.ValueType
	if valueType == "" {
		valueType = "string"
	}
	m := &model.SysSetting{
		Category:  req.Category,
		Key:       req.Key,
		Value:     req.Value,
		ValueType: valueType,
		Status:    status,
	}
	if req.Description != "" {
		m.Description = &req.Description
	}
	if req.Editable != nil {
		m.Editable = *req.Editable
	}
	m.CreateBy = &opID
	m.UpdateBy = &opID
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *SettingService) Update(ctx context.Context, id uint64, req *dto.SettingUpdateRequest, opID uint64) (*model.SysSetting, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !m.Editable {
		return nil, code.ErrSettingNotEditable
	}
	m.Value = req.Value
	m.ValueType = req.ValueType
	if req.Description != "" {
		m.Description = &req.Description
	}
	if req.Status != "" {
		m.Status = req.Status
	}
	m.UpdateBy = &opID
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *SettingService) Delete(ctx context.Context, id uint64) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if m.IsSystem {
		return code.ErrSettingSystem
	}
	return s.repo.Delete(ctx, id)
}

func (s *SettingService) GetByID(ctx context.Context, id uint64) (*dto.SettingVO, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toSettingVO(m), nil
}

func (s *SettingService) Page(ctx context.Context, req *dto.SettingSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	vos := make([]dto.SettingVO, len(rows))
	for i, r := range rows {
		vos[i] = *toSettingVO(&r)
	}
	return dto.NewPageResponse(vos, total, &req.PageRequest), nil
}

func (s *SettingService) ListCategories(ctx context.Context) ([]string, error) {
	return s.repo.ListCategories(ctx)
}

// GetByCategory 获取指定分类下所有配置 (key→value map)
func (s *SettingService) GetByCategory(ctx context.Context, category string) (map[string]dto.SettingVO, error) {
	rows, err := s.repo.ListByCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	result := make(map[string]dto.SettingVO, len(rows))
	for _, r := range rows {
		result[r.Key] = *toSettingVO(&r)
	}
	return result, nil
}

// BatchSave 批量保存配置 (按 category+key upsert)
func (s *SettingService) BatchSave(ctx context.Context, req *dto.BatchSaveRequest, opID uint64) (*dto.BatchSaveResult, error) {
	var created, updated int
	for _, item := range req.Items {
		existing, err := s.repo.GetByCategoryKey(ctx, req.Category, item.Key)
		if err != nil {
			// 不存在 → 创建
			m := &model.SysSetting{
				Category:  req.Category,
				Key:       item.Key,
				Value:     item.Value,
				ValueType: item.ValueType,
				Status:    model.StatusEnabled,
				Editable:  true,
			}
			m.CreateBy = &opID
			m.UpdateBy = &opID
			if err := s.repo.Create(ctx, m); err != nil {
				return nil, err
			}
			created++
		} else {
			// 存在 → 更新
			existing.Value = item.Value
			existing.ValueType = item.ValueType
			existing.UpdateBy = &opID
			if err := s.repo.Update(ctx, existing); err != nil {
				return nil, err
			}
			updated++
		}
	}
	return &dto.BatchSaveResult{Created: created, Updated: updated}, nil
}

func toSettingVO(m *model.SysSetting) *dto.SettingVO {
	vo := &dto.SettingVO{
		ID:        m.ID,
		Category:  m.Category,
		Key:       m.Key,
		Value:     m.Value,
		ValueType: m.ValueType,
		Editable:  m.Editable,
		IsSystem:  m.IsSystem,
		Status:    m.Status,
	}
	if m.Description != nil {
		vo.Description = m.Description
	}
	vo.CreateTime = m.CreateTime.Format(time.DateTime)
	vo.UpdateTime = m.UpdateTime.Format(time.DateTime)
	return vo
}
