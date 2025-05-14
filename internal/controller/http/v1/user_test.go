package v1

import (
	"KnowledgeHub/internal/models"
	"KnowledgeHub/internal/repo/mocks"
	"KnowledgeHub/internal/services"
	"KnowledgeHub/pkg/logger"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_GetUser(t *testing.T) {
	// Налаштування Gin у тестовому режимі
	gin.SetMode(gin.TestMode)

	// Створення мок-репозиторію та додавання тестових даних
	mockRepo := mocks.NewRepository()
	mockRepo.AddUser(&models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	})

	// Створення сервісу з мок-репозиторієм
	userService := services.NewUserService(mockRepo.User())

	// Створення логера
	l := logger.New("debug")

	// Створення роутера
	r := gin.New()

	// Створення контролера
	userHandler := NewUserHandler(userService, l)

	// Реєстрація маршруту
	r.GET("/users/:id", userHandler.GetUser)

	// Тестові випадки
	tests := []struct {
		name         string
		url          string
		wantStatus   int
		wantResponse map[string]interface{}
	}{
		{
			name:       "Valid user",
			url:        "/users/1",
			wantStatus: http.StatusOK,
			wantResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"id":       float64(1), // JSON перетворює int на float64
					"username": "testuser",
					"email":    "test@example.com",
				},
			},
		},
		{
			name:       "User not found",
			url:        "/users/999",
			wantStatus: http.StatusNotFound,
			wantResponse: map[string]interface{}{
				"error": "user not found",
			},
		},
	}

	// Виконання тестів
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Створення запиту
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			// Виконання запиту
			r.ServeHTTP(w, req)

			// Перевірка статус-коду
			if w.Code != tt.wantStatus {
				t.Errorf("Status code = %d, want %d", w.Code, tt.wantStatus)
			}

			// Перевірка відповіді
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}

			// Порівняння відповіді з очікуваною
			// Тут можна використовувати більш складну логіку порівняння
			if tt.wantStatus == http.StatusOK {
				if data, ok := response["data"].(map[string]interface{}); ok {
					expectedData := tt.wantResponse["data"].(map[string]interface{})
					if data["id"] != expectedData["id"] ||
						data["username"] != expectedData["username"] ||
						data["email"] != expectedData["email"] {
						t.Errorf("Response body = %v, want %v", data, expectedData)
					}
				} else {
					t.Errorf("Response does not contain data field")
				}
			}
		})
	}
}
