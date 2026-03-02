package service

import (
	"errors"
	"swiftly/backend/internal/user"
	"swiftly/backend/internal/user/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req user.CreateUserRequest) (*user.User, error) {
	if req.Email == "" || req.FullName == "" {
		return nil, errors.New("email and full name are required")
	}

	u := &user.User{
		Email:    req.Email,
		FullName: req.FullName,
	}

	err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) GetUserByID(id string) (*user.User, error) {
	return s.repo.GetByID(id)
}
