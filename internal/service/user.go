package service

import (
	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

func (s *UserService) Create(user *domain.User) error {
	return s.repo.Create(user)
}

func (s *UserService) GetByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) GetByUsername(username string) (*domain.User, error) {
	return s.repo.FindByUsername(username)
}
