package service

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/code"
	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

const (
	DefaultChapterLen = 100
	MaxChapterLen     = 100000
)

// BookChapterRuleService 章节识别规则管理服务
type BookChapterRuleService struct {
	db                 *gorm.DB
	chapterRuleRepo    *repository.BookChapterRuleRepository
	chapterRuleRelRepo *repository.BookChapterRuleRelRepository
}

func NewBookChapterRuleService(
	db *gorm.DB,
	chapterRuleRepo *repository.BookChapterRuleRepository,
	chapterRuleRelRepo *repository.BookChapterRuleRelRepository,
) *BookChapterRuleService {
	return &BookChapterRuleService{
		db:                 db,
		chapterRuleRepo:    chapterRuleRepo,
		chapterRuleRelRepo: chapterRuleRelRepo,
	}
}

// ==================== 章节识别规则 CRUD ====================

func (s *BookChapterRuleService) CreateChapterRule(ctx context.Context, req *dto.ChapterRuleRequest, userID uint64) (*model.BookChapterRule, error) {
	// 如果是系统默认规则，UserID 设为 nil；否则设为当前用户
	userIDPtr := &userID
	if req.RuleType == string(model.RuleTypeSystem) {
		userIDPtr = nil
	}
	m := &model.BookChapterRule{
		RuleName:      req.RuleName,
		RuleType:      model.RuleType(req.RuleType),
		UserID:        userIDPtr,
		TitlePattern:  req.TitlePattern,
		GroupPattern:  req.GroupPattern,
		MinChapterLen: req.MinChapterLen,
		MaxChapterLen: req.MaxChapterLen,
		SortOrder:     req.SortOrder,
		Description:   req.Description,
		Status:        model.EnableStatus(req.Status),
	}
	if m.RuleType == "" {
		m.RuleType = model.RuleTypeCustom
	}
	if m.Status == "" {
		m.Status = model.StatusEnabled
	}
	if m.MinChapterLen == 0 {
		m.MinChapterLen = DefaultChapterLen
	}
	if m.MaxChapterLen == 0 {
		m.MaxChapterLen = MaxChapterLen
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.chapterRuleRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookChapterRuleService) UpdateChapterRule(ctx context.Context, id uint64, req *dto.ChapterRuleRequest, userID uint64) (*model.BookChapterRule, error) {
	m, err := s.chapterRuleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, code.ErrRuleNotFound
	}
	m.RuleName = req.RuleName
	m.RuleType = model.RuleType(req.RuleType)
	// 如果是系统默认规则，UserID 设为 nil
	if req.RuleType == string(model.RuleTypeSystem) {
		m.UserID = nil
	}
	m.TitlePattern = req.TitlePattern
	m.GroupPattern = req.GroupPattern
	m.MinChapterLen = req.MinChapterLen
	m.MaxChapterLen = req.MaxChapterLen
	m.SortOrder = req.SortOrder
	m.Description = req.Description
	if req.Status != "" {
		m.Status = model.EnableStatus(req.Status)
	}
	m.UpdateBy = &userID
	if err := s.chapterRuleRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *BookChapterRuleService) DeleteChapterRule(ctx context.Context, id uint64) error {
	if _, err := s.chapterRuleRepo.GetByID(ctx, id); err != nil {
		return code.ErrRuleNotFound
	}
	return s.chapterRuleRepo.Delete(ctx, id)
}

func (s *BookChapterRuleService) GetChapterRuleByID(ctx context.Context, id uint64) (*model.BookChapterRule, error) {
	return s.chapterRuleRepo.GetByID(ctx, id)
}

func (s *BookChapterRuleService) PageChapterRule(ctx context.Context, req *dto.ChapterRuleSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.chapterRuleRepo.Page(ctx, req)
	if err != nil {
		return nil, err
	}
	records := make([]dto.ChapterRuleResponse, len(rows))
	for i, r := range rows {
		records[i] = dto.ChapterRuleResponse{BookChapterRule: r}
	}
	return dto.NewPageResponse(records, total, &req.PageRequest), nil
}

// ==================== 章节规则绑定 ====================

// BindChapterRule 为书籍绑定固定章节识别规则
func (s *BookChapterRuleService) BindChapterRule(ctx context.Context, req *dto.ChapterRuleBindRequest, userID uint64) (*dto.ChapterRuleBindResponse, error) {
	// 验证规则存在
	_, err := s.chapterRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		return nil, code.ErrRuleNotFound
	}

	// 查找已有绑定，存在则直接修改 rule_id，不存在则新建
	var rel *model.BookChapterRuleRel
	existing, err := s.chapterRuleRelRepo.GetByBookAndReader(ctx, req.BookID, userID)
	if err == nil && existing != nil {
		existing.RuleID = req.RuleID
		if err := s.db.WithContext(ctx).Save(existing).Error; err != nil {
			return nil, err
		}
		rel = existing
	} else {
		rel = &model.BookChapterRuleRel{
			BookID:   req.BookID,
			ReaderID: userID,
			RuleID:   req.RuleID,
		}
		rel.CreateBy = &userID
		rel.UpdateBy = &userID
		if err := s.chapterRuleRelRepo.Create(ctx, rel); err != nil {
			return nil, err
		}
	}

	// 获取规则名称
	rule, _ := s.chapterRuleRepo.GetByID(ctx, req.RuleID)
	ruleName := ""
	if rule != nil {
		ruleName = rule.RuleName
	}

	return &dto.ChapterRuleBindResponse{
		ID:         rel.ID,
		BookID:     rel.BookID,
		ReaderID:   rel.ReaderID,
		RuleID:     rel.RuleID,
		RuleName:   ruleName,
		CreateTime: rel.CreateTime.Format("2006-01-02 15:04:05"),
	}, nil
}

// UnbindChapterRule 解绑书籍的章节识别规则
func (s *BookChapterRuleService) UnbindChapterRule(ctx context.Context, bookID, userID uint64) error {
	return s.chapterRuleRelRepo.DeleteByBookAndReader(ctx, bookID, userID)
}

// GetBoundChapterRule 获取书籍绑定的章节识别规则
func (s *BookChapterRuleService) GetBoundChapterRule(ctx context.Context, bookID, userID uint64) (*dto.ChapterRuleBindResponse, error) {
	rel, err := s.chapterRuleRelRepo.GetByBookAndReader(ctx, bookID, userID)
	if err != nil {
		return nil, nil // 未绑定返回 nil
	}

	rule, _ := s.chapterRuleRepo.GetByID(ctx, rel.RuleID)
	ruleName := ""
	if rule != nil {
		ruleName = rule.RuleName
	}

	return &dto.ChapterRuleBindResponse{
		ID:         rel.ID,
		BookID:     rel.BookID,
		ReaderID:   rel.ReaderID,
		RuleID:     rel.RuleID,
		RuleName:   ruleName,
		CreateTime: rel.CreateTime.Format("2006-01-02 15:04:05"),
	}, nil
}
