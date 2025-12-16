package service

import (
	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) CreateUser(user *domain.User) (*domain.User, error) {
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepo.FindByUsername(username)
}

func CreateUser(user *domain.User) (*domain.User, error) {
	service := NewUserService()
	return service.CreateUser(user)
}
