package mocks

import (
	"KnowledgeHub/internal/models"
)

// MockUserRepository реалізує інтерфейс UserRepository для тестування
type MockUserRepository struct {
	store *Mocks
}

func (m *MockUserRepository) CreateUser() error {
	return nil
}

func (m *MockUserRepository) GetUserByID(id uint) (*models.User, error) {
	user, exists := m.store.users[id]
	if !exists {
		return nil, nil // або повернути помилку "не знайдено"
	}
	return user, nil
}

func (m *MockUserRepository) UpdateUser() error {
	return nil
}

func (m *MockUserRepository) DeleteUser(id int) error {
	delete(m.store.users, uint(id))
	return nil
}
