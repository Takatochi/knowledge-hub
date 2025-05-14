package mocks

import (
	"KnowledgeHub/internal/models"
	"KnowledgeHub/internal/repo"
)

type Mocks struct {
	users              map[uint]*models.User
	mockUserRepository *MockUserRepository
}

func NewRepository() *Mocks {
	return &Mocks{
		users: make(map[uint]*models.User),
	}
}

// Допоміжні методи для тестування
func (m *Mocks) AddUser(user *models.User) {
	m.users[user.ID] = user
}

func (r *Mocks) User() repo.UserRepository {
	if r.mockUserRepository != nil {
		return r.mockUserRepository
	}

	r.mockUserRepository = &MockUserRepository{
		store: r,
	}

	return r.mockUserRepository
}
