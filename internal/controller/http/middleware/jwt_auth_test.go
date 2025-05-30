package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"KnowledgeHub/config"
	"KnowledgeHub/internal/services"
	"KnowledgeHub/pkg/logger"

	"github.com/gin-gonic/gin"
)

func getTestJWTService() *services.JWTService {
	cfg := &config.Config{
		JWT: config.JWT{
			Secret:           "test_secret_key_for_testing_purposes_only",
			AccessTokenTTL:   900,
			RefreshTokenTTL:  604800,
			SigningAlgorithm: "HS256",
		},
	}
	return services.NewJWTService(cfg)
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtService := getTestJWTService()
	logger := logger.New("debug")

	// Генеруємо валідний токен
	tokenPair, err := jwtService.GenerateTokenPair(1, "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Створюємо роутер з middleware
	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtService, logger))
	router.GET("/protected", func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// Створюємо запит з валідним токеном
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Перевіряємо результат
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestJWTAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtService := getTestJWTService()
	logger := logger.New("debug")

	// Створюємо роутер з middleware
	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtService, logger))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Створюємо запит без токена
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Перевіряємо результат
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWTAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtService := getTestJWTService()
	logger := logger.New("debug")

	// Створюємо роутер з middleware
	router := gin.New()
	router.Use(JWTAuthMiddleware(jwtService, logger))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Тестуємо різні невалідні формати
	testCases := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "invalid_token"},
		{"Wrong prefix", "Basic invalid_token"},
		{"Empty token", "Bearer "},
		{"Invalid token", "Bearer invalid_token"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tc.header)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected status %d, got %d for case %s", http.StatusUnauthorized, w.Code, tc.name)
			}
		})
	}
}

func TestOptionalJWTAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtService := getTestJWTService()
	logger := logger.New("debug")

	// Створюємо роутер з опціональним middleware
	router := gin.New()
	router.Use(OptionalJWTAuthMiddleware(jwtService, logger))
	router.GET("/optional", func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if exists {
			c.JSON(http.StatusOK, gin.H{"user_id": userID, "authenticated": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"authenticated": false})
		}
	})

	// Створюємо запит без токена
	req := httptest.NewRequest("GET", "/optional", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Перевіряємо, що запит пройшов успішно
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOptionalJWTAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtService := getTestJWTService()
	logger := logger.New("debug")

	// Генеруємо валідний токен
	tokenPair, err := jwtService.GenerateTokenPair(1, "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Створюємо роутер з опціональним middleware
	router := gin.New()
	router.Use(OptionalJWTAuthMiddleware(jwtService, logger))
	router.GET("/optional", func(c *gin.Context) {
		userID, exists := GetUserIDFromContext(c)
		if exists {
			c.JSON(http.StatusOK, gin.H{"user_id": userID, "authenticated": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"authenticated": false})
		}
	})

	// Створюємо запит з валідним токеном
	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Перевіряємо результат
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Створюємо контекст з user_id
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("user_id", uint(123))

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		t.Error("Expected user_id to exist in context")
	}
	if userID != 123 {
		t.Errorf("Expected user_id 123, got %d", userID)
	}

	// Тестуємо контекст без user_id
	ctx2, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, exists2 := GetUserIDFromContext(ctx2)
	if exists2 {
		t.Error("Expected user_id to not exist in empty context")
	}
}

func TestGetUsernameFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Створюємо контекст з username
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("username", "testuser")

	username, exists := GetUsernameFromContext(ctx)
	if !exists {
		t.Error("Expected username to exist in context")
	}
	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}

	// Тестуємо контекст без username
	ctx2, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, exists2 := GetUsernameFromContext(ctx2)
	if exists2 {
		t.Error("Expected username to not exist in empty context")
	}
}

func TestGetEmailFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Створюємо контекст з email
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("email", "test@example.com")

	email, exists := GetEmailFromContext(ctx)
	if !exists {
		t.Error("Expected email to exist in context")
	}
	if email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", email)
	}

	// Тестуємо контекст без email
	ctx2, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, exists2 := GetEmailFromContext(ctx2)
	if exists2 {
		t.Error("Expected email to not exist in empty context")
	}
}
