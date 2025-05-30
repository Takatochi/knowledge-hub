package services

import (
	"KnowledgeHub/config"
	"testing"
	"time"
)

func getTestConfig() *config.Config {
	return &config.Config{
		JWT: config.JWT{
			Secret:           "test_secret_key_for_testing_purposes_only",
			AccessTokenTTL:   900,
			RefreshTokenTTL:  604800,
			SigningAlgorithm: "HS256",
		},
	}
}

func TestJWTService_GenerateTokenPair(t *testing.T) {
	cfg := getTestConfig()
	jwtService := NewJWTService(cfg)

	userID := uint(1)
	username := "testuser"
	email := "test@example.com"

	tokenPair, err := jwtService.GenerateTokenPair(userID, username, email)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tokenPair == nil {
		t.Fatal("Expected token pair, got nil")
	}

	if tokenPair.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if tokenPair.RefreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}

	if tokenPair.ExpiresAt == 0 {
		t.Error("Expected expires at timestamp, got 0")
	}

	// Перевіряємо, що expires at в майбутньому
	if tokenPair.ExpiresAt <= time.Now().Unix() {
		t.Error("Expected expires at to be in the future")
	}
}

func TestJWTService_ValidateAccessToken(t *testing.T) {
	cfg := getTestConfig()
	jwtService := NewJWTService(cfg)

	userID := uint(1)
	username := "testuser"
	email := "test@example.com"

	// Генеруємо токен
	tokenPair, err := jwtService.GenerateTokenPair(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Валідуємо access токен
	claims, err := jwtService.ValidateAccessToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("Expected username %s, got %s", username, claims.Username)
	}

	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}

	if claims.Subject != "access_token" {
		t.Errorf("Expected subject 'access_token', got %s", claims.Subject)
	}
}

func TestJWTService_ValidateAccessToken_InvalidToken(t *testing.T) {
	cfg := getTestConfig()
	jwtService := NewJWTService(cfg)

	// Тестуємо з невалідним токеном
	_, err := jwtService.ValidateAccessToken("invalid_token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	if err != ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTService_ValidateRefreshToken(t *testing.T) {
	cfg := getTestConfig()
	jwtService := NewJWTService(cfg)

	userID := uint(1)
	username := "testuser"
	email := "test@example.com"

	// Генеруємо токен
	tokenPair, err := jwtService.GenerateTokenPair(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Валідуємо refresh токен
	claims, err := jwtService.ValidateRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if claims.Subject != "refresh_token" {
		t.Errorf("Expected subject 'refresh_token', got %s", claims.Subject)
	}

	if claims.Issuer != "KnowledgeHub" {
		t.Errorf("Expected issuer 'KnowledgeHub', got %s", claims.Issuer)
	}
}

func TestJWTService_RefreshTokens(t *testing.T) {
	cfg := getTestConfig()
	jwtService := NewJWTService(cfg)

	userID := uint(1)
	username := "testuser"
	email := "test@example.com"

	// Генеруємо початкову пару токенів
	originalTokenPair, err := jwtService.GenerateTokenPair(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Додаємо невелику затримку для забезпечення різних timestamp
	time.Sleep(time.Millisecond)

	// Оновлюємо токени
	newTokenPair, err := jwtService.RefreshTokens(originalTokenPair.RefreshToken, userID, username, email)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if newTokenPair.AccessToken == originalTokenPair.AccessToken {
		t.Error("Expected new access token to be different from original")
	}

	if newTokenPair.RefreshToken == originalTokenPair.RefreshToken {
		t.Error("Expected new refresh token to be different from original")
	}

	// Перевіряємо, що новий access токен валідний
	claims, err := jwtService.ValidateAccessToken(newTokenPair.AccessToken)
	if err != nil {
		t.Fatalf("New access token should be valid: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, claims.UserID)
	}
}

func TestJWTService_ExtractTokenFromHeader(t *testing.T) {
	cfg := getTestConfig()
	jwtService := NewJWTService(cfg)

	tests := []struct {
		name        string
		authHeader  string
		expectedErr error
		expectToken bool
	}{
		{
			name:        "Valid Bearer token",
			authHeader:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedErr: nil,
			expectToken: true,
		},
		{
			name:        "Empty header",
			authHeader:  "",
			expectedErr: ErrInvalidToken,
			expectToken: false,
		},
		{
			name:        "Invalid format - no Bearer",
			authHeader:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedErr: ErrInvalidToken,
			expectToken: false,
		},
		{
			name:        "Invalid format - wrong prefix",
			authHeader:  "Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedErr: ErrInvalidToken,
			expectToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwtService.ExtractTokenFromHeader(tt.authHeader)

			if tt.expectedErr != nil {
				if err != tt.expectedErr {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}

			if tt.expectToken && token == "" {
				t.Error("Expected token, got empty string")
			}

			if !tt.expectToken && token != "" {
				t.Error("Expected empty token, got non-empty string")
			}
		})
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	// Створюємо конфігурацію з дуже коротким TTL
	cfg := &config.Config{
		JWT: config.JWT{
			Secret:           "test_secret_key_for_testing_purposes_only",
			AccessTokenTTL:   -1, // Негативний TTL для створення вже прострочених токенів
			RefreshTokenTTL:  604800,
			SigningAlgorithm: "HS256",
		},
	}

	jwtService := NewJWTService(cfg)

	userID := uint(1)
	username := "testuser"
	email := "test@example.com"

	// Генеруємо токен (він буде вже прострочений)
	tokenPair, err := jwtService.GenerateTokenPair(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Спробуємо валідувати прострочений токен
	_, err = jwtService.ValidateAccessToken(tokenPair.AccessToken)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}

	if err != ErrExpiredToken {
		t.Errorf("Expected ErrExpiredToken, got %v", err)
	}
}
