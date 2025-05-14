package services

import (
	"KnowledgeHub/internal/models"
	"KnowledgeHub/internal/repo"
)

type UserService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (uc *UserService) GetUser(id uint) (*models.User, error) {
	// Бізнес-логіка
	return uc.userRepo.GetUserByID(id)
}

// Інші методи
