package service

import (
	"context"
	"errors"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	_, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context, req *dto.PageRequest) ([]model.User, int64, error) {
	return s.userRepo.ListWithKeyword(ctx, req.Page, req.PageSize, req.Keyword)
}