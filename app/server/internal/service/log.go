package service

import (
	"context"

	"boread/internal/dto"
	"boread/internal/repository"
)

type LogService struct {
	repo *repository.SysLogRepository
}

func NewLogService(repo *repository.SysLogRepository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) PageLogin(ctx context.Context, req *dto.LoginLogSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.repo.PageLogin(ctx, req)
	if err != nil {
		return nil, err
	}
	return dto.NewPageResponse(rows, total, &req.PageRequest), nil
}

func (s *LogService) PageOperation(ctx context.Context, req *dto.OperationLogSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.repo.PageOperation(ctx, req)
	if err != nil {
		return nil, err
	}
	return dto.NewPageResponse(rows, total, &req.PageRequest), nil
}
