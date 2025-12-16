package repository

import (
	"moonshine/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: DB}
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

