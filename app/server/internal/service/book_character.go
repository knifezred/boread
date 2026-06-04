package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"
)

var (
	ErrCharacterNotFound    = errors.New("角色不存在")
	ErrCharacterRelNotFound = errors.New("角色关系不存在")
	ErrBookExists           = errors.New("书籍不存在")
)

// ======================== BookCharacterService ========================

type BookCharacterService struct {
	db       *gorm.DB
	charRepo *repository.BookCharacterRepository
	bookRepo *repository.BookRepository
}

func NewBookCharacterService(db *gorm.DB, charRepo *repository.BookCharacterRepository, bookRepo *repository.BookRepository) *BookCharacterService {
	return &BookCharacterService{db: db, charRepo: charRepo, bookRepo: bookRepo}
}

func (s *BookCharacterService) Create(ctx context.Context, userID uint64, req *dto.CharacterRequest) (*dto.CharacterResponse, error) {
	if _, err := s.bookRepo.GetByID(ctx, req.BookID); err != nil {
		return nil, ErrBookExists
	}
	m := &model.BookCharacter{
		BookID:    req.BookID,
		Name:      req.Name,
		Alias:     req.Alias,
		RoleType:  req.RoleType,
		Avatar:    req.Avatar,
		Intro:     req.Intro,
		Extra:     req.Extra,
		SortOrder: req.SortOrder,
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.charRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return &dto.CharacterResponse{BookCharacter: *m}, nil
}

func (s *BookCharacterService) Update(ctx context.Context, userID, charID uint64, req *dto.CharacterRequest) (*dto.CharacterResponse, error) {
	m, err := s.charRepo.GetByID(ctx, charID)
	if err != nil {
		return nil, ErrCharacterNotFound
	}
	m.Name = req.Name
	m.Alias = req.Alias
	m.RoleType = req.RoleType
	m.Avatar = req.Avatar
	m.Intro = req.Intro
	m.Extra = req.Extra
	m.SortOrder = req.SortOrder
	m.UpdateBy = &userID
	if err := s.charRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return &dto.CharacterResponse{BookCharacter: *m}, nil
}

func (s *BookCharacterService) Delete(ctx context.Context, charID uint64) error {
	if _, err := s.charRepo.GetByID(ctx, charID); err != nil {
		return ErrCharacterNotFound
	}
	return s.charRepo.Delete(ctx, charID)
}

func (s *BookCharacterService) GetByID(ctx context.Context, charID uint64) (*dto.CharacterResponse, error) {
	m, err := s.charRepo.GetByID(ctx, charID)
	if err != nil {
		return nil, ErrCharacterNotFound
	}
	return &dto.CharacterResponse{BookCharacter: *m}, nil
}

func (s *BookCharacterService) Page(ctx context.Context, req *dto.CharacterSearch) (*dto.PageResponse, error) {
	req.Normalize()
	rows, total, err := s.charRepo.PageByBook(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.CharacterResponse, len(rows))
	for i, r := range rows {
		resp[i] = dto.CharacterResponse{BookCharacter: r}
	}
	return dto.NewPageResponse(resp, total, &req.PageRequest), nil
}

func (s *BookCharacterService) ListByBook(ctx context.Context, bookID uint64) ([]dto.CharacterResponse, error) {
	rows, err := s.charRepo.ListByBook(ctx, bookID)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.CharacterResponse, len(rows))
	for i, r := range rows {
		resp[i] = dto.CharacterResponse{BookCharacter: r}
	}
	return resp, nil
}

// ======================== BookCharacterRelService ========================

type BookCharacterRelService struct {
	db       *gorm.DB
	relRepo  *repository.BookCharacterRelRepository
	charRepo *repository.BookCharacterRepository
}

func NewBookCharacterRelService(db *gorm.DB, relRepo *repository.BookCharacterRelRepository, charRepo *repository.BookCharacterRepository) *BookCharacterRelService {
	return &BookCharacterRelService{db: db, relRepo: relRepo, charRepo: charRepo}
}

func (s *BookCharacterRelService) Create(ctx context.Context, userID uint64, req *dto.CharacterRelRequest) (*dto.CharacterRelResponse, error) {
	if _, err := s.charRepo.GetByID(ctx, req.CharacterAID); err != nil {
		return nil, ErrCharacterNotFound
	}
	if _, err := s.charRepo.GetByID(ctx, req.CharacterBID); err != nil {
		return nil, ErrCharacterNotFound
	}
	m := &model.BookCharacterRel{
		BookID:       req.BookID,
		CharacterAID: req.CharacterAID,
		CharacterBID: req.CharacterBID,
		RelationType: req.RelationType,
		RelationDesc: req.RelationDesc,
		SortOrder:    req.SortOrder,
	}
	m.CreateBy = &userID
	m.UpdateBy = &userID
	if err := s.relRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return s.buildRelResponse(ctx, m)
}

func (s *BookCharacterRelService) Delete(ctx context.Context, relID uint64) error {
	if _, err := s.relRepo.GetByID(ctx, relID); err != nil {
		return ErrCharacterRelNotFound
	}
	return s.relRepo.Delete(ctx, relID)
}

func (s *BookCharacterRelService) GetByID(ctx context.Context, relID uint64) (*dto.CharacterRelResponse, error) {
	m, err := s.relRepo.GetByID(ctx, relID)
	if err != nil {
		return nil, ErrCharacterRelNotFound
	}
	return s.buildRelResponse(ctx, m)
}

func (s *BookCharacterRelService) ListByCharacter(ctx context.Context, characterID uint64) ([]dto.CharacterRelResponse, error) {
	rows, err := s.relRepo.ListByCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}
	return s.buildRelResponses(ctx, rows)
}

func (s *BookCharacterRelService) ListByBook(ctx context.Context, bookID uint64) ([]dto.CharacterRelResponse, error) {
	rows, err := s.relRepo.ListByBook(ctx, bookID)
	if err != nil {
		return nil, err
	}
	return s.buildRelResponses(ctx, rows)
}

func (s *BookCharacterRelService) buildRelResponse(ctx context.Context, m *model.BookCharacterRel) (*dto.CharacterRelResponse, error) {
	resp := &dto.CharacterRelResponse{BookCharacterRel: *m}
	if a, err := s.charRepo.GetByID(ctx, m.CharacterAID); err == nil {
		resp.CharacterAName = &a.Name
	}
	if b, err := s.charRepo.GetByID(ctx, m.CharacterBID); err == nil {
		resp.CharacterBName = &b.Name
	}
	return resp, nil
}

func (s *BookCharacterRelService) buildRelResponses(ctx context.Context, rows []model.BookCharacterRel) ([]dto.CharacterRelResponse, error) {
	nameMap := make(map[uint64]string)
	for _, r := range rows {
		nameMap[r.CharacterAID] = ""
		nameMap[r.CharacterBID] = ""
	}
	for id := range nameMap {
		if c, err := s.charRepo.GetByID(ctx, id); err == nil {
			nameMap[id] = c.Name
		}
	}
	resp := make([]dto.CharacterRelResponse, len(rows))
	for i, r := range rows {
		resp[i] = dto.CharacterRelResponse{BookCharacterRel: r}
		if n, ok := nameMap[r.CharacterAID]; ok {
			resp[i].CharacterAName = &n
		}
		if n, ok := nameMap[r.CharacterBID]; ok {
			resp[i].CharacterBName = &n
		}
	}
	return resp, nil
}
