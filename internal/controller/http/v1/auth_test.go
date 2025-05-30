package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"KnowledgeHub/config"
	"KnowledgeHub/internal/controller/http/middleware"
	"KnowledgeHub/internal/repo/mocks"
	"KnowledgeHub/internal/services"
	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
)

func getTestAuthHandler() (*AuthHandler, *services.JWTService) {
	cfg := &config.Config{
		JWT: config.JWT{
			Secret:           "test_secret_key_for_testing_purposes_only",
			AccessTokenTTL:   900,
			RefreshTokenTTL:  604800,
			SigningAlgorithm: "HS256",
		},
	}

	jwtService := services.NewJWTService(cfg)
	mockRepo := mocks.NewRepository()
	userService := services.NewUserService(mockRepo.User())
	logger := logger.New("debug")

	return NewAuthHandler(jwtService, userService, logger), jwtService
}

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, _ := getTestAuthHandler()

	router := gin.New()
	router.POST("/auth/login", authHandler.Login)

	loginReq := LoginRequest{
		Username: "admin",
		Password: "password",
	}

	jsonData, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if response.RefreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}

	if response.User.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", response.User.Username)
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, _ := getTestAuthHandler()

	router := gin.New()
	router.POST("/auth/login", authHandler.Login)

	loginReq := LoginRequest{
		Username: "admin",
		Password: "wrongpassword",
	}

	jsonData, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, _ := getTestAuthHandler()

	router := gin.New()
	router.POST("/auth/login", authHandler.Login)

	// Відправляємо невалідний JSON
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, _ := getTestAuthHandler()

	router := gin.New()
	router.POST("/auth/register", authHandler.Register)

	registerReq := RegisterRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.User.Username != "newuser" {
		t.Errorf("Expected username 'newuser', got '%s'", response.User.Username)
	}

	if response.User.Email != "newuser@example.com" {
		t.Errorf("Expected email 'newuser@example.com', got '%s'", response.User.Email)
	}
}

func TestAuthHandler_Register_ExistingUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, _ := getTestAuthHandler()

	router := gin.New()
	router.POST("/auth/register", authHandler.Register)

	registerReq := RegisterRequest{
		Username: "admin", // Існуючий користувач
		Email:    "admin@example.com",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, jwtService := getTestAuthHandler()

	// Генеруємо токени
	tokenPair, err := jwtService.GenerateTokenPair(1, "admin", "admin@example.com")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Додаємо затримку для забезпечення різних timestamp
	time.Sleep(time.Millisecond)

	router := gin.New()
	router.POST("/auth/refresh", authHandler.RefreshToken)

	refreshReq := RefreshRequest{
		RefreshToken: tokenPair.RefreshToken,
	}

	jsonData, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("Expected new access token, got empty string")
	}

	if response.AccessToken == tokenPair.AccessToken {
		t.Error("Expected new access token to be different from original")
	}
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, _ := getTestAuthHandler()

	router := gin.New()
	router.POST("/auth/refresh", authHandler.RefreshToken)

	refreshReq := RefreshRequest{
		RefreshToken: "invalid_token",
	}

	jsonData, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, jwtService := getTestAuthHandler()
	logger := logger.New("debug")

	// Генеруємо токен
	tokenPair, err := jwtService.GenerateTokenPair(1, "admin", "admin@example.com")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	router := gin.New()
	router.Use(middleware.JWTAuthMiddleware(jwtService, logger))
	router.POST("/auth/logout", authHandler.Logout)

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authHandler, jwtService := getTestAuthHandler()
	logger := logger.New("debug")

	// Генеруємо токен
	tokenPair, err := jwtService.GenerateTokenPair(1, "admin", "admin@example.com")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	router := gin.New()
	router.Use(middleware.JWTAuthMiddleware(jwtService, logger))
	router.GET("/auth/me", authHandler.Me)

	req := httptest.NewRequest("GET", "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response UserInfo
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", response.Username)
	}

	if response.Email != "admin@example.com" {
		t.Errorf("Expected email 'admin@example.com', got '%s'", response.Email)
	}
}
