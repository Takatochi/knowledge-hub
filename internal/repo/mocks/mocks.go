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

func (m *Mocks) User() repo.UserRepository {
	if m.mockUserRepository != nil {
		return m.mockUserRepository
	}

	m.mockUserRepository = &MockUserRepository{
		store: m,
	}

	return m.mockUserRepository
}
