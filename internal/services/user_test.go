package services

import (
	"testing"

	"KnowledgeHub/internal/models"
	"KnowledgeHub/internal/repo/mocks"
)

func TestUserService_GetUser(t *testing.T) {
	// Підготовка
	mockRepo := mocks.NewRepository()

	mockRepo.AddUser(&models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	})

	service := NewUserService(mockRepo.User())

	// Тестові випадки
	tests := []struct {
		name     string
		userID   uint
		wantUser *models.User
		wantErr  bool
	}{
		{
			name:   "Existing user",
			userID: 1,
			wantUser: &models.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name:     "Non-existing user",
			userID:   2,
			wantUser: nil,
			wantErr:  false, // або true, залежно від вашої логіки
		},
	}

	// Виконання тестів
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := service.GetUser(tt.userID)

			// Перевірка помилки
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Перевірка результату
			if tt.wantUser == nil && gotUser != nil {
				t.Errorf("UserService.GetUser() = %v, want nil", gotUser)
			} else if tt.wantUser != nil {
				if gotUser == nil {
					t.Errorf("UserService.GetUser() = nil, want %v", tt.wantUser)
				} else if gotUser.ID != tt.wantUser.ID ||
					gotUser.Username != tt.wantUser.Username ||
					gotUser.Email != tt.wantUser.Email {
					t.Errorf("UserService.GetUser() = %v, want %v", gotUser, tt.wantUser)
				}
			}
		})
	}
}
